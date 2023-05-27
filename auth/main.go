package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"net"
	"os"

	// "os"
	"time"

	pb "webserver/auth/pb"

	redis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type authenticationServer struct {
	pb.UnimplementedAuthenticationServer
}

var p int32 = 23
var g int32 = 5

// var b int32 = int32(rand.Intn(100) + 1)
var b int32 = 15
var redisClient *redis.Client
var redisCtx context.Context

func (s *authenticationServer) ReqPG(ctx context.Context, req *pb.ReqPGRequest) (*pb.ReqPGResponse, error) {
	clientNonce := req.Nonce
	serverNonce := RandomNonce(20)
	clientMessageId := req.MessageId
	redisKey := clientNonce + serverNonce
	redisVal := fmt.Sprintf("%x", sha1.Sum([]byte(redisKey)))
	fmt.Printf("%s %s\n", redisKey, redisVal)
	err := redisClient.Set(redisCtx, redisKey, redisVal, 20*time.Minute).Err()
	if err != nil {
		log.Printf("trying to set %s as %s failed in %e", redisKey, redisVal, err)
	}

	return &pb.ReqPGResponse{Nonce: clientNonce, ServerNonce: serverNonce, MessageId: clientMessageId + 1, P: p, G: g}, nil
}

func (s *authenticationServer) ReqDHParams(ctx context.Context, req *pb.ReqDHParamsRequest) (*pb.ReqDHParamsResponse, error) {
	clientNonce := req.Nonce
	serverNonce := req.ServerNonce
	clientMessageId := req.MessageId
	publicA := req.A
	publicB := int32(1)
	modulo := int32(1)
	for i := int32(1); i <= b; i++ {
		modulo *= publicA
		modulo %= p

		publicB *= g
		publicB %= p
	}
	log.Printf("this should be saved in redis %d", modulo)
	err := redisClient.Set(redisCtx, clientNonce+serverNonce, modulo, 20*time.Minute).Err()
	if err != nil {
		log.Fatal("could not set auth key for client")
	}

	return &pb.ReqDHParamsResponse{Nonce: clientNonce, ServerNonce: serverNonce, MessageId: clientMessageId + 1, B: publicB}, nil
}

func main() {
	redisCtx = context.TODO()
	// log.Println(os.Getenv("REDIS_PASS"))
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis_cache:6379",      // host:port of the redis server
		Password: os.Getenv("REDIS_PASS"), // no password set
		DB:       0,                       // use default DB
	})
	if err := redisClient.Ping(redisCtx).Err(); err != nil {
		log.Fatal(err)
		return
	}
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
