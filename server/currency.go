package server

import (
	"context"
	"fmt"
	"github.com/mahdimehrabi/grpc-base-currency/data"
	"github.com/mahdimehrabi/grpc-base-currency/proto/currency"
)

type CurrencyServer struct {
	rates *data.ExchangeRates
}

func NewCurrencyServer(r *data.ExchangeRates) *CurrencyServer {
	return &CurrencyServer{rates: r}
}

func (s CurrencyServer) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	fmt.Println("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := s.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &currency.RateResponse{
		Rate: rate,
	}, nil
}
