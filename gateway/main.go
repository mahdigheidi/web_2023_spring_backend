package main

import (
	"context"
	"log"

	authPb "webserver/gateway/pb/auth"
	bizPb "webserver/gateway/pb/biz"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	authConn, authErr := grpc.Dial("localhost:5052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if authErr != nil {
		log.Fatalf("failed to connect: %v", authErr)
	}
	defer authConn.Close()
	authClient := authPb.NewAuthenticationClient(authConn)

	bizConn, bizErr := grpc.Dial("localhost:5062", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if bizErr != nil {
		log.Fatalf("failed to connect to business services: %v", bizErr)
	}
	defer bizConn.Close()
	bizClient := bizPb.NewBusinessClient(bizConn)

	r := gin.Default()

	// Auth services
	r.GET("/req_pq", func(c *gin.Context) {
		var nonce string
		var message_id int32
		nonce = "ABCDEFGHIJ0123456789"
		message_id = 2
		response, err := authClient.ReqPG(context.Background(), &authPb.ReqPGRequest{Nonce: nonce, MessageId: message_id})
		c.JSON(200, gin.H{
			"response": response,
			"err":      err,
		})
	})

	r.GET("/req_DH_params", func(c *gin.Context) {
		var nonce, server_nonce string
		var message_id, a int32
		nonce = "ABCDEFGHIJ0123456789"
		server_nonce = "ABCDEFGHIJ0123456789"
		message_id = 4
		a = 19
		response, err := authClient.ReqDHParams(context.Background(),
			&authPb.ReqDHParamsRequest{
				Nonce:       nonce,
				ServerNonce: server_nonce,
				MessageId:   message_id,
				A:           a,
			})
		c.JSON(200, gin.H{
			"response": response,
			"err":      err,
		})
	})

	// Biz services
	r.GET("/get_users", func(c *gin.Context) {
		var user_id, message_id int32
		var auth_key string
		response, err := bizClient.GetUsers(context.Background(), &bizPb.GetUsersRequest{UserId: user_id, AuthKey: auth_key, MessageId: message_id})
		c.JSON(200, gin.H{
			"response": response,
			"err":      err,
		})
	})

	r.GET("/get_users_with_sql_inject", func(c *gin.Context) {
		var message_id int32
		var user_id, auth_key string
		response, err := bizClient.GetUsersWithSQLInject(context.Background(), &bizPb.GetUsersWithSQLInjectRequest{UserId: user_id, AuthKey: auth_key, MessageId: message_id})
		c.JSON(200, gin.H{
			"response": response,
			"err":      err,
		})
	})

	r.Run("0.0.0.0:6433") // listen and serve on 0.0.0.0:8080
}
