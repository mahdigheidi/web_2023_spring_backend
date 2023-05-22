package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "webserver/auth/pb"

	"google.golang.org/grpc"
)

type authenticationServer struct {
	pb.UnimplementedAuthenticationServer
}

func (s *authenticationServer) ReqPG(ctx context.Context, req *pb.ReqPGRequest) (*pb.ReqPGResponse, error) {
	fmt.Println("fuck this shit")
	return &pb.ReqPGResponse{}, nil
}

func (s *authenticationServer) ReqDHParams(ctx context.Context, req *pb.ReqDHParamsRequest) (*pb.ReqDHParamsResponse, error) {
	fmt.Println("print dh params")
	return &pb.ReqDHParamsResponse{}, nil
}

func main() {
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
