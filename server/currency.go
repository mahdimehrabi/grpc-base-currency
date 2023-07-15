package server

import (
	"context"
	"fmt"
	"github.com/mahdimehrabi/grpc-base-currency/proto/currency"
)

type CurrencyServer struct {
}

func NewCurrencyServer() *CurrencyServer {
	return &CurrencyServer{}
}

func (s CurrencyServer) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	fmt.Println("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	return &currency.RateResponse{
		Rate: 0.5,
	}, nil
}
