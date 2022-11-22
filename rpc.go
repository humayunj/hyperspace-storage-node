package main

import (
	"context"
	"errors"
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

	priceInt, ok := ComputePrice(c.TimePeriod, c.FileSize)
	if !ok {
		return nil, errors.New("failed to compute bid")
	}
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

	inBid := new(big.Int)
	inBid, ok := inBid.SetString(in.Bid, 10)
	if !ok {
		return nil, errors.New("bid parsing failed")
	}
	if requiredBid, ok := ComputePrice(in.TimeStart-in.TimeEnd, in.FileSize); !ok || requiredBid.Cmp(inBid) == -1 {
		return nil, errors.New("minimum bid amount required " + requiredBid.String())
	}
	t := FileTokenParams{
		Bid:             in.Bid,
		FileSize:        in.FileSize,
		FileHash:        in.FileHash,
		TimeStart:       in.TimeStart,
		TimeEnd:         in.TimeEnd,
		ConcludeTimeout: in.ConcludeTimeout,
		ProveTimeout:    in.ProveTimeout,
		SegmentsCount:   in.SegmentsCount,
		UserAddress:     in.UserAddress,
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
