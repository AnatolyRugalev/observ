package hello

import (
	"github.com/AnatolyRugalev/observ/collectors/logrust"
	"github.com/AnatolyRugalev/observ/collectors/prometheust"
	"github.com/AnatolyRugalev/observ/collectors/stdlogt"
	"github.com/AnatolyRugalev/observ/logq"
	"github.com/AnatolyRugalev/observ/logt"
	"github.com/AnatolyRugalev/observ/metrq"
	"github.com/AnatolyRugalev/observ/metrt"
	"testing"
)

func TestHelloStandardLog(t *testing.T) {
	// logt.Start creates new log tester, and starts capturing immediately
	lt := logt.Start(t, stdlogt.Default())
	//                  ^ here we provide standard Go log collector
	Hello("Sam")
	lt.Collect().Assert().Count(1) // 1 record collected
	Hello("Jack")
	lt.Collect().Assert().Count(2) // 2 records collected

	// lt.Collect() will always return full set of records
	// To scope the records, you can use lt.Start, lt.Finish of lt.Scope:
	scope := lt.Start()
	Hello("Mary")
	Hello("Mary")
	Hello("Mary")
	Hello("Alice")
	records := scope.Finish()
	records.Assert().Count(4)

	// And root lt now has 6 records
	lt.Collect().Assert().Count(6)

	// When the scope is active, all higher level collectors also capture records.
	// You can do multiple levels of nesting. Also, lt.Scope() will return captured records right away:
	lt.Scope(func(lt logt.LogT) {
		// this lt ^ refers to current scope. You can interact with records within the scope using it
		Hello("Andy")
		Hello("Andy")
		lt.Collect().Assert().Count(2)
		Hello("Jerry")
	}).Assert().Count(3) // .Scope() return a set of records, so they can be asserted right away.
}

func TestHelloLogrus(t *testing.T) {
	// Now, let's see how we can work with log record attributes
	// This tester is initialized with Logrus default logger:
	lt := logt.Start(t, logrust.Default())
	Hello("Toly")
	lt.Collect().Assert().Attr("name", "Toly")
	// You can also filter records after they are collected:
	Hello("A")
	Hello("A")
	Hello("A")
	Hello("B")
	Hello("A")
	lt.Collect().Attr("name", "A").Assert().Count(4)
	// And, of course, you can filter records however you want.
	// All filters are just functions:
	lt.Collect(func(v logq.Record) bool {
		return v.Message == "Hello!" && v.Attributes["name"] == "B"
	}).Assert().Count(1)

	// You can also perform aggregation operations on records:
	lt.Collect().Group(logq.ByAttr("name")).Assert().Count(map[string]int{
		"A":    4,
		"B":    1,
		"Toly": 1,
	})
}

func TestHelloPrometheus(t *testing.T) {
	// OK, logging is easy, let's try testing some metrics.
	mt := metrt.New(t,
		// use Prometheus default collector
		prometheust.Default(),
		// only take into account metrics starting from hello_
		// By default, there are a bunch of Go metrics in the registry.
		metrt.WithCollectFilter(metrq.Prefix("hello_")),
	)
	// Same principle applies, but we work with metric measurements instead of log records
	Hello("A")
	Hello("A")
	Hello("A")
	Hello("A")
	Hello("B")
	Hello("B")

	// Total 6 increments:
	mt.Collect().Assert().Sum(6)

	// Now, we expect that "A" had 4 instances and "B" only 2.
	// How do we test for this? Easy:
	mt.Collect().Group(metrq.ByAttr("name")).Assert().Sum(map[string]int64{
		"A": 4,
		"B": 2,
	})

	// Scoping works in a similar way as in logs
	mt.Scope(func(mct metrt.MetrT) {
		Hello("A")
		Hello("B")
	}).Group(metrq.ByAttr("name")).Assert().Sum(map[string]int64{
		"A": 1,
		"B": 1,
	})
	// You can also receive negative counter value from the scope. This indicates that value
	// had decreased.
}
