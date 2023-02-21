package coffeshop

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/AnatolyRugalev/observ/collectors/logrust"
	"github.com/AnatolyRugalev/observ/collectors/otelt"
	"github.com/AnatolyRugalev/observ/collectors/prometheust"
	"github.com/AnatolyRugalev/observ/collectors/stdlogt"
	"github.com/AnatolyRugalev/observ/logq"
	"github.com/AnatolyRugalev/observ/logt"
	"github.com/AnatolyRugalev/observ/logt/logwait"
	"github.com/AnatolyRugalev/observ/metrcollect/metrmulti"
	"github.com/AnatolyRugalev/observ/metrq"
	"github.com/AnatolyRugalev/observ/metrt"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"log"
	"time"
)

type CoffeeShopTestSuite struct {
	suite.Suite

	shop *CoffeeShop
}

func TestCoffeeShopTestSuite(t *testing.T) {
	suite.Run(t, &CoffeeShopTestSuite{})
}

func (s *CoffeeShopTestSuite) SetupTest() {
	// promtest.Tester implements prometheus.Registerer.
	// CoffeeShop registers metrics upon creation
	s.shop = NewCoffeeShop(60)
}

func (s *CoffeeShopTestSuite) TestClosed() {
	_, err := s.shop.PlaceOrder(context.Background(), Coffee{Kind: CoffeeEspresso})
	s.Require().Error(err)
	s.Require().Equal("sorry, coffeeshop is closed", err.Error())
}

func (s *CoffeeShopTestSuite) TestOpenNegativeBaristas() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := s.shop.Open(ctx, -10)
	s.Require().Error(err)
	s.Require().Equal("you can't make coffee with -10 baristas", err.Error())
}

func (s *CoffeeShopTestSuite) TestPlaceOrder() {

	otelCollector := otelt.New()
	mt := metrt.New(s.T(), metrmulti.Multi(prometheust.New(prometheus.DefaultGatherer), otelCollector))
	provider := metric.NewMeterProvider(metric.WithReader(otelCollector))
	global.SetMeterProvider(provider)

	log.SetPrefix("[coffeeshop] ")
	log.SetFlags(log.Lmsgprefix | log.LUTC | log.Llongfile | log.Lmicroseconds | log.Ldate | log.Ltime)
	lt := logt.Start(s.T(), logrust.Default(), stdlogt.Default())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := s.shop.Open(ctx, 3)
	s.Require().NoError(err)
	orderCtx, orderCancel := context.WithCancel(ctx)
	ch, err := s.shop.PlaceOrder(orderCtx, Coffee{Kind: CoffeeEspresso})
	orderCancel()
	s.Require().NoError(err)
	lt.Wait(logwait.For(2 * time.Second)).Require().Count(1)
	coffee := <-ch
	s.Require().Equal(Coffee{Kind: CoffeeEspresso}, coffee)
	lt.Collect().Message("Placing order").Require().Count(1)

	logt.Scope(s.T(), func(lt2 logt.LogT) {
		lt.Scope(func(lt logt.LogT) {
			orderCtx, orderCancel = context.WithCancel(ctx)
			ch, err = s.shop.PlaceOrder(orderCtx, Coffee{Kind: CoffeeCappucino})
			orderCancel()
			s.Require().NoError(err)
			coffee = <-ch
			s.Require().Equal(Coffee{Kind: CoffeeCappucino}, coffee)
		}).Message("Placing order").Require().Count(1)

		lt.Collect().Message("Placing order").Require().Count(2)
		lt2.Collect(logq.Message("Placing order")).Require().Count(1)

		orderCtx, orderCancel = context.WithCancel(ctx)
		delta := mt.Scope(func(mq metrt.MetrT) {
			ch, err = s.shop.PlaceOrder(orderCtx, Coffee{Kind: CoffeeCappucino})
			orderCancel()
			s.Require().NoError(err)
			coffee = <-ch
			s.Require().Equal(Coffee{Kind: CoffeeCappucino}, coffee)
		})
		delta.And(metrq.Name("observ_coffee_shop_cups")).Assert().Sum(1)

		lt.Collect().Message("Placing order").Require().Count(3)

	}, logrust.Default()).Message("Placing order").Require().Count(2)

	mt.Collect().
		Where(metrq.Scope("coffeeshop"), metrq.Name("cups")).
		Group(metrq.ByAttr("kind")).
		Assert().
		Sum(map[string]int64{
			"cappuccino": 2,
			"espresso":   1,
		})
	mt.Collect().
		And(metrq.Name("observ_coffee_shop_cups")).
		Group(metrq.ByAttr("kind")).
		Require().
		Sum(map[string]int64{
			"cappuccino": 2,
			"espresso":   1,
		})
}
