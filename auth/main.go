package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	pb "webserver/auth/pb"

	"google.golang.org/grpc"
)

type authenticationServer struct {
	pb.UnimplementedAuthenticationServer
}

func (s *authenticationServer) ReqPG(ctx context.Context, req *pb.ReqPGRequest) (*pb.ReqPGResponse, error) {
	clientNonce := req.Nonce
	clientMessageId := req.MessageId
	return &pb.ReqPGResponse{Nonce: clientNonce, ServerNonce: RandomNonce(20), MessageId: clientMessageId+1, P: 17, G: 3}, nil
}

func (s *authenticationServer) ReqDHParams(ctx context.Context, req *pb.ReqDHParamsRequest) (*pb.ReqDHParamsResponse, error) {
	clinetNonce := req.Nonce
	serverNonce := req.ServerNonce
	clientMessageId := req.MessageId
	publicA := req.A
	b := rand.Intn(100) + 1
	p := int32(17)
	g := int32(3)
	publicB := int32(1)
	modulo := int32(1)
	for i := 1; i <= b; i++ {
		modulo = modulo * publicA
		modulo = modulo % p

		publicB *= g
		publicB %= p
	}
	fmt.Printf("this should be saved in redis %d", modulo)
	return &pb.ReqDHParamsResponse{Nonce: clinetNonce, ServerNonce: serverNonce, MessageId: clientMessageId + 1, B: publicB}, nil
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
