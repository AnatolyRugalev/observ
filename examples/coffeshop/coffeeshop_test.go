package coffeshop

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/suite"

	"github.com/AnatolyRugalev/observ/prometheus/promtest"
)

type CoffeeShopTestSuite struct {
	suite.Suite

	prom *promtest.Tester
	shop *CoffeeShop
}

func TestCoffeeShopTestSuite(t *testing.T) {
	suite.Run(t, &CoffeeShopTestSuite{})
}

func (s *CoffeeShopTestSuite) SetupTest() {
	// Initialize promtest.Tester
	s.prom = promtest.New(
		// testing.T is necessary, all following options are optional
		s.T(),
		// Provide metrics prefix to simplify testing
		promtest.WithPrefix("observ_coffee_shop_"),
		// When no registry provided, promtest creates a new registry
		promtest.WithRegistry(prometheus.NewRegistry()),
		// Set polling rate for eventual assertion (e.g. wait until metric value becomes equal X)
		promtest.WithPollingRate(100*time.Millisecond),
	)
	// promtest.Tester implements prometheus.Registerer.
	// CoffeeShop registers metrics upon creation
	s.shop = NewCoffeeShop(s.prom, 60)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := s.shop.Open(ctx, 3)
	s.Require().NoError(err)
	orderCtx, orderCancel := context.WithCancel(ctx)
	ch, err := s.shop.PlaceOrder(orderCtx, Coffee{Kind: CoffeeEspresso})
	orderCancel()
	s.Require().NoError(err)
	coffee := <-ch
	s.Require().Equal(Coffee{Kind: CoffeeEspresso}, coffee)

	orderCtx, orderCancel = context.WithCancel(ctx)
	ch, err = s.shop.PlaceOrder(orderCtx, Coffee{Kind: CoffeeCappucino})
	orderCancel()
	s.Require().NoError(err)
	coffee = <-ch
	s.Require().Equal(Coffee{Kind: CoffeeCappucino}, coffee)

	orderCtx, orderCancel = context.WithCancel(ctx)
	ch, err = s.shop.PlaceOrder(orderCtx, Coffee{Kind: CoffeeCappucino})
	orderCancel()
	s.Require().NoError(err)
	coffee = <-ch
	s.Require().Equal(Coffee{Kind: CoffeeCappucino}, coffee)

	s.prom.Assert("cups", "kind", "espresso").Equal(1)
	s.prom.Assert("cups", "kind", "cappucino").Equal(2)

	s.prom.Assert("cups").Group("kind").EqualMap(promtest.M{
		"espresso":  1,
		"cappucino": 2,
	})
}
