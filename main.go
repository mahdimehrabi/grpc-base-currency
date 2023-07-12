package main

import (
	"currency-implementation/proto/currency"
	"currency-implementation/server"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	gs := grpc.NewServer()
	cs := server.NewCurrencyServer()

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
