package server

import (
	"context"
	"fmt"
	"github.com/mahdimehrabi/grpc-base-currency/data"
	"github.com/mahdimehrabi/grpc-base-currency/proto/currency"
	"io"
	"time"
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

func (s CurrencyServer) SubscribeRates(src currency.Currency_SubscribeRatesServer) error {
	go func() {
		for {
			rr, err := src.Recv()
			if err == io.EOF {
				fmt.Println("client has closed connection")
				break
			}
			if err != nil {
				fmt.Println("unable to read from client", err)
				break
			}
			fmt.Println("Handle client request", "request.base", rr.Base, "request.destionation", rr.Destination)
		}
	}()
	for {
		err := src.Send(&currency.RateResponse{Rate: 12.1})
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(5 * time.Second)
	}
}
