package main

import (
	"context"
	"log"
	"net"

	pb "../proto/biz_pb"

	"google.golang.org/grpc"
)

type businessServer struct {
	pb.UnimplementedBusinessServer
}

func (s *businessServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {

	return &pb.GetUsersResponse{}, nil
}

func (s *businessServer) GetUsersWithSQLInject(ctx context.Context, req *pb.GetUsersWithSQLInjectRequest) (*pb.GetUsersResponse, error) {
	return &pb.GetUsersResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5062")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	bizServer := &businessServer{}
	pb.RegisterBusinessServer(s, bizServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
