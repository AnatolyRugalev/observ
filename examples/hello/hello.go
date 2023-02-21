package hello

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"log"
)

var counter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "hello_calls",
}, []string{"name"})

func Hello(name string) {
	log.Printf("Hello, %s!", name)
	logrus.WithField("name", name).Info("Hello!")
	counter.WithLabelValues(name).Inc()
}
