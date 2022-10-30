package main

import (
	"context"
	"log"
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
	log.Println("Listening RPC on 8000")
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}

	return s
}

func (s *RPCServer) Ping(ctx context.Context, c *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{
		CanStore: true,
		BidPrice: "100",
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
		Bid:           uint64(in.Bid),
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
	}

	return &out, nil
}
