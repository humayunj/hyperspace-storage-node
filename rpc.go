package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	color "github.com/fatih/color"
	proto "github.com/storage-node-p1/storage-node-proto"
	"google.golang.org/grpc"
)

type RPCServer struct{}

func RunRPCServer() *grpc.Server {

	s := grpc.NewServer()

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	proto.RegisterStorageNodeServer(s, &RPCServer{})
	color.Set(color.FgMagenta)
	log.Println("Listening RPC on 8000")
	color.Unset()
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}

	return s
}

func (s *RPCServer) Ping(ctx context.Context, c *proto.PingRequest) (*proto.PingResponse, error) {
	days := float64(c.TimePeriod) / (60 * 60 * 24)

	mbs := float64(c.FileSize) / (1024 * 1024)

	log.Print("FileSize:", c.FileSize)
	log.Print("Tp:", c.TimePeriod)

	k := days * mbs
	log.Print("K:", k)
	price := new(big.Float)
	price, valid := price.SetString(NC.FeeWeiPerMBPerDay)
	if !valid {
		log.Panic("Failed to parse FeeWeiPerMBPerDay", NC.FeeWeiPerMBPerDay)
	}

	price = price.Mul(price, big.NewFloat((k)))
	log.Print("Price:", price.String())
	priceInt := new(big.Int)

	intVal := fmt.Sprintf("%.0f", price)
	log.Print("intVal:", intVal)

	priceInt.SetString(intVal, 10)
	log.Print("PriceInt:", priceInt.String())

	baseFee := new(big.Int)
	baseFee, val := baseFee.SetString(NC.FeeBase, 10)
	if !val {
		panic("failed to parse FeeBase")
	}
	priceInt = priceInt.Add(priceInt, baseFee)

	return &proto.PingResponse{
		CanStore: true,
		BidPrice: priceInt.String(),
	}, nil
}

func (s *RPCServer) GetStats(ctx context.Context, in *proto.Empty) (*proto.GetStatsResponse, error) {
	color.Set(color.FgYellow)
	log.Println("RPC::GetStats")
	color.Unset()
	return &proto.GetStatsResponse{
		FreeStorage: int32(NC.TotalStorage),
	}, nil
}

func (s *RPCServer) InitTransaction(ctx context.Context, in *proto.InitTransactionRequest) (*proto.InitTransactionResponse, error) {
	color.Set(color.FgYellow)
	log.Println("RPC::InitTransaction")
	color.Unset()
	t := FileTokenParams{
		Bid:           (in.Bid),
		FileSize:      uint64(in.FileSize),
		FileHash:      in.FileHash,
		Timeperiod:    uint64(in.Timeperiod),
		SegmentsCount: uint64(in.SegmentsCount),
	}
	token, err := JFS.CreateFileToken(t)
	if err != nil {
		return nil, err
	}
	out := proto.InitTransactionResponse{
		JWT:       token,
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(), // 2 Minute
		HttpURL:   NC.HttpURL,
	}

	return &out, nil
}
