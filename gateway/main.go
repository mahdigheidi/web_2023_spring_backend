package main

import (
	"context"
	"log"
	"strconv"

	authPb "webserver/gateway/pb/auth"
	bizPb "webserver/gateway/pb/biz"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	authConn, authErr := grpc.Dial("authentication:5052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if authErr != nil {
		log.Fatalf("failed to connect: %v", authErr)
	}
	defer authConn.Close()
	authClient := authPb.NewAuthenticationClient(authConn)

	bizConn, bizErr := grpc.Dial("business_logic:5062", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if bizErr != nil {
		log.Fatalf("failed to connect to business services: %v", bizErr)
	}
	defer bizConn.Close()
	bizClient := bizPb.NewBusinessClient(bizConn)

	r := gin.Default()

	// Auth services
	r.GET("/req_pq", func(c *gin.Context) {
		nonce := c.Query("nonce")
		message_id, err := strconv.Atoi(c.Query("message_id"))
		if err != nil {
			c.Error(err)
		} else {
			response, err := authClient.ReqPG(context.Background(), &authPb.ReqPGRequest{Nonce: nonce, MessageId: int32(message_id)})
			c.JSON(200, gin.H{
				"response": response,
				"err":      err,
			})
		}
	})

	r.GET("/req_dh_params", func(c *gin.Context) {
		nonce := c.Query("nonce")
		server_nonce := c.Query("server_nonce")
		message_id, _ := strconv.Atoi(c.Query("message_id"))
		a, _ := strconv.Atoi(c.Query("a"))
		response, err := authClient.ReqDHParams(context.Background(),
			&authPb.ReqDHParamsRequest{
				Nonce:       nonce,
				ServerNonce: server_nonce,
				MessageId:   int32(message_id),
				A:           int32(a),
			})
		c.JSON(200, gin.H{
			"response": response,
			"err":      err,
		})
	})

	// Biz services
	r.GET("/get_users", func(c *gin.Context) {
		user_id, _ := strconv.Atoi(c.Query("user_id"))
		message_id, _ := strconv.Atoi(c.Query("message_id"))
		auth_key := c.Query("auth_key")
		response, err := bizClient.GetUsers(context.Background(), &bizPb.GetUsersRequest{UserId: int32(user_id), AuthKey: auth_key, MessageId: int32(message_id)})
		c.JSON(200, gin.H{
			"response": response,
			"err":      err,
		})
	})

	r.GET("/get_users_with_sql_inject", func(c *gin.Context) {
		user_id := c.Query("user_id")
		message_id, _ := strconv.Atoi(c.Query("message_id"))
		auth_key := c.Query("auth_key")
		response, err := bizClient.GetUsersWithSQLInject(context.Background(), &bizPb.GetUsersWithSQLInjectRequest{UserId: user_id, AuthKey: auth_key, MessageId: int32(message_id)})
		c.JSON(200, gin.H{
			"response": response,
			"err":      err,
		})
	})

	r.Run("0.0.0.0:6433") // listen and serve on 0.0.0.0:8080
}
