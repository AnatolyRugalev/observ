package coffeshop

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"log"
)

var meter = global.Meter("coffeeshop")

var metrics = coffeeMetrics{
	cups: promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "observ",
		Subsystem: "coffee_shop",
		Name:      "cups",
	}, []string{"kind"}),
	extras: promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "observ",
		Subsystem: "coffee_shop",
		Name:      "extras",
	}, []string{"extra"}),
	latency: promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "observ",
		Subsystem: "coffee_shop",
		Name:      "latency",
	}, []string{"kind"}),
	cupsCounter: lo.Must(meter.Int64Counter("cups")),
	latencyHist: lo.Must(meter.Float64Histogram("latency")),
}

type coffeeMetrics struct {
	cups    *prometheus.CounterVec
	extras  *prometheus.CounterVec
	latency *prometheus.HistogramVec

	cupsCounter instrument.Int64Counter
	latencyHist instrument.Float64Histogram
}

type State int

const (
	StateClosed = State(iota)
	StateOpen
	StateFinalizing
)

type CoffeeShop struct {
	timeScale int

	stateMu  sync.Mutex
	state    State
	baristas chan order
}

type CoffeeKind string

const (
	CoffeeEspresso  = CoffeeKind("espresso")
	CoffeeLungo     = CoffeeKind("lungo")
	CoffeeCafeLatte = CoffeeKind("cafe-latte")
	CoffeeCappucino = CoffeeKind("cappuccino")
	CofeeFlatWhite  = CoffeeKind("flat-white")
)

type Extra string

const (
	ExtraShot          = Extra("shot")
	ExtraCinnamon      = Extra("cinnamon")
	ExtraBaileys       = Extra("baileys")
	ExtraCondencedMilk = Extra("condenced-milk")
)

type Coffee struct {
	Kind   CoffeeKind
	Extras []Extra
}

func (c Coffee) Duration() time.Duration {
	var d time.Duration
	switch c.Kind {
	case CoffeeEspresso:
		d = time.Minute
	case CoffeeLungo:
		d = time.Minute + 30*time.Second
	case CofeeFlatWhite:
		d = 2 * time.Minute
	case CoffeeCappucino:
		d = 2 * time.Minute
	case CoffeeCafeLatte:
		d = 2*time.Minute + 30*time.Second
	}
	for _, extra := range c.Extras {
		switch extra {
		case ExtraShot:
			d += time.Minute
		case ExtraCinnamon:
			d += 10 * time.Second
		case ExtraBaileys:
			d += 20 * time.Second
		case ExtraCondencedMilk:
			panic("barista started laughing uncontrollably and died")
		}
	}
	return d
}

type order struct {
	// todo: allow trace injection
	coffee Coffee
	to     chan Coffee
}

func NewCoffeeShop(timeScale int) *CoffeeShop {
	return &CoffeeShop{
		timeScale: timeScale,
	}
}

// PlaceOrder allocates one barista to make caller a coffee. The duration of making
// depends on order complexity
func (c *CoffeeShop) PlaceOrder(ctx context.Context, coffee Coffee) (<-chan Coffee, error) {
	c.stateMu.Lock()
	state := c.state
	c.stateMu.Unlock()
	logrus.Info("Placing order")
	switch state {
	case StateClosed:
		return nil, fmt.Errorf("sorry, coffeeshop is closed")
	case StateFinalizing:
		return nil, fmt.Errorf("sorry, baristas are making their last drinks. please come back tomorrow")
	}
	order := order{
		coffee: coffee,
		to:     make(chan Coffee),
	}
	select {
	case <-ctx.Done():
		close(order.to)
		return nil, fmt.Errorf("no free baristas: please come later: %w", ctx.Err())
	case c.baristas <- order:
		// order accepted
	}
	return order.to, nil
}

// Open opens a coffeeshop with desired context
func (c *CoffeeShop) Open(ctx context.Context, baristasNum int) error {
	if baristasNum < 1 {
		return fmt.Errorf("you can't make coffee with %d baristas", baristasNum)
	}
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	switch c.state {
	case StateOpen:
		return fmt.Errorf("coffeeshop is already open")
	case StateFinalizing:
		return fmt.Errorf("previous barisas are still inside")
	}
	c.state = StateOpen
	c.baristas = make(chan order, 1)
	wg := sync.WaitGroup{}
	for i := 0; i < baristasNum; i++ {
		wg.Add(1)
		go c.baristaLoop(&wg)
	}
	go func(wg *sync.WaitGroup) {
		<-ctx.Done()
		c.stateMu.Lock()
		c.state = StateFinalizing
		// baristas should make coffee for commited orders
		close(c.baristas)
		c.stateMu.Unlock()
		wg.Wait()
		c.stateMu.Lock()
		c.state = StateClosed
		c.baristas = nil
		c.stateMu.Unlock()
	}(&wg)
	return nil
}

func (c *CoffeeShop) baristaLoop(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		if p := recover(); p != nil {
			// TODO: funerals
		}
	}()
	for order := range c.baristas {
		start := time.Now()
		time.Sleep(order.coffee.Duration() / time.Duration(c.timeScale)) // sleeping at work, bastards?
		log.Printf("Sleepy barista woke up")
		latency := time.Since(start)
		metrics.cups.WithLabelValues(string(order.coffee.Kind)).Inc()
		metrics.cupsCounter.Add(context.Background(), 1, attribute.String("kind", string(order.coffee.Kind)))
		metrics.latency.WithLabelValues(string(order.coffee.Kind)).Observe(latency.Seconds())
		metrics.latencyHist.Record(context.Background(), latency.Seconds(), attribute.String("kind", string(order.coffee.Kind)))
		for _, extra := range order.coffee.Extras {
			metrics.extras.WithLabelValues(string(extra)).Inc()
		}
		order.to <- order.coffee
	}
}
