package server

import (
	"context"
	protos "github.com/MousaZa/product-services/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	log hclog.Logger
}

func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{log: l}
}

func (c *Currency) GetRate(ctx context.Context, request *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", request.GetBase(), "destination", request.GetDestination())
	return &protos.RateResponse{Rate: 0.5}, nil
}
