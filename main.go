package main

import (
	"fmt"
	"github.com/mahdimehrabi/grpc-base-currency/data"
	"github.com/mahdimehrabi/grpc-base-currency/proto/currency"
	"github.com/mahdimehrabi/grpc-base-currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {

	rates, err := data.NewRates()
	if err != nil {
		fmt.Println("Unable to generate rates", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()
	cs := server.NewCurrencyServer(rates)

	currency.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		fmt.Println("unable to create listener", "error", err)
		os.Exit(1)
	}
	err = gs.Serve(l)
	panic(err)
}
