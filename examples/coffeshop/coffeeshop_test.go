package coffeshop

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/suite"
)

type CoffeeShopTestSuite struct {
	suite.Suite

	shop *CoffeeShop
}

func TestCoffeeShopTestSuite(t *testing.T) {
	suite.Run(t, &CoffeeShopTestSuite{})
}

func (s *CoffeeShopTestSuite) SetupTest() {
	registry := prometheus.NewRegistry()
	s.shop = NewCoffeeShop(registry)
}

func (s *CoffeeShopTestSuite) TestClosed() {
	_, err := s.shop.PlaceOrder(context.Background(), Coffee{Kind: CoffeeEspresso})
	s.Require().Error(err)
	s.Require().Equal("sorry, coffeeshop is closed", err.Error())
}
