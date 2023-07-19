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
	rates         *data.ExchangeRates
	subscriptions map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest
}

func NewCurrencyServer(r *data.ExchangeRates) *CurrencyServer {
	c := &CurrencyServer{rates: r, subscriptions: make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest, 0)}
	go c.handleUpdates()
	return c
}

func (c *CurrencyServer) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		fmt.Println("GO updated rates")

		//loop over subscribed clients
		for k, v := range c.subscriptions {

			//loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					fmt.Println("Unable to get updated rate", "base", rr.GetBase().String(), "Destination", rr.GetDestination().String())
				}
				err = k.Send(&currency.RateResponse{Rate: r, Base: rr.GetBase(), Destination: rr.GetDestination()})
				if err != nil {
					fmt.Println("Unable to send updated rate", "base", rr.GetBase().String(), "Destination", rr.GetDestination().String())
				}
			}
		}
	}
}

func (s CurrencyServer) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	fmt.Println("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := s.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &currency.RateResponse{
		Rate:        rate,
		Base:        rr.GetBase(),
		Destination: rr.GetDestination(),
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
			rrs, ok := s.subscriptions[src]
			if !ok {
				rrs = []*currency.RateRequest{}
			}
			rrs = append(rrs, rr)
			s.subscriptions[src] = rrs
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
