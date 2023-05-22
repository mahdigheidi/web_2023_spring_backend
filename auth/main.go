package main

import (
	"context"
	pb "web-server-2023.com/auth/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type authenticationServer struct {
	pb.UnimplementedAuthenticationServer
}

func (s *authenticationServer) ReqPG(ctx context.Context, req *pb.ReqPGRequest) (*pb.ReqPGResponse, error) {
	
	return &pb.ReqPGResponse{}, nil
}

func (s *authenticationServer) ReqDHParams(ctx context.Context, req *pb.ReqDHParamsRequest) (*pb.ReqDHParamsResponse, error) {
	return &pb.ReqDHParamsResponse{}, nil
}

func main () {
    lis, err := net.Listen("tcp", ":5052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    authServer := &authenticationServer{}
    pb.RegisterAuthenticationServer(s, authServer)
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}