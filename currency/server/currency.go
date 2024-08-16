package server

import (
	"context"
	"github.com/MousaZa/product-services/currency/data"
	protos "github.com/MousaZa/product-services/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	log hclog.Logger
	r   *data.ExchangeRates
}

func NewCurrency(l hclog.Logger, r *data.ExchangeRates) *Currency {
	return &Currency{log: l, r: r}
}

func (c *Currency) GetRate(ctx context.Context, request *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", request.GetBase(), "destination", request.GetDestination())
	rate, err := c.r.GetRate(request.GetBase().String(), request.GetDestination().String())
	c.log.Info("rate", rate)
	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Rate: rate}, nil
}
