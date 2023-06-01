package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	authPb "webserver/gateway/pb/auth"
	bizPb "webserver/gateway/pb/biz"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var redisClient *redis.Client
var redisCtx context.Context

func isAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authIdentifier := c.Query("auth_id")
		authKey := c.Query("auth_key")
		redisSavedKey, err := redisClient.Get(redisCtx, authIdentifier).Result()

		fmt.Println(authIdentifier, authKey, redisSavedKey)

		if err != nil {
			log.Fatal("error while fetching the provided auth key")
			return
		}

		if redisSavedKey == authKey {
			c.Set("authenticated", true)
		} else {
			c.Set("authenticated", false)
		}
	}
}

func Throttle(maxEventsPerSec int, maxBurstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(maxEventsPerSec), maxBurstSize)

	return func(context *gin.Context) {
		if limiter.Allow() {
			context.Next()
			return
		}
		badClient := context.ClientIP()
		redisClient.Set(redisCtx, badClient, "blocked", 24*time.Hour)
		context.Error(errors.New("limit exceeded"))
		context.AbortWithStatus(http.StatusTooManyRequests)
	}
}

func Blacklist() gin.HandlerFunc {
    return func(context *gin.Context) {
		clientIP := context.ClientIP()
		clientStatus, clientErr := redisClient.Get(redisCtx, clientIP).Result()
		if clientErr == nil && clientStatus == "blocked" {
			context.Error(errors.New("this IP will not get any service for an small amount of time"))
			context.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		context.Next()
	}
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
	redisClient.FlushAll(redisCtx)

	authCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	authConn, authErr := grpc.DialContext(authCtx, "authentication:5052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if authErr != nil {
		log.Fatalf("failed to connect: %v", authErr)
	}
	defer authConn.Close()
	authClient := authPb.NewAuthenticationClient(authConn)

	bizCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	bizConn, bizErr := grpc.DialContext(bizCtx, "business_logic:5062", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if bizErr != nil {
		log.Fatalf("failed to connect to business services: %v", bizErr)
	}
	defer bizConn.Close()
	bizClient := bizPb.NewBusinessClient(bizConn)

	r := gin.Default()
	r.Use(Blacklist())

	// Auth services
	auth := r.Group("/auth")
	maxEventsPerSec := 100
	maxBurstSize := 100
	auth.Use(Throttle(maxEventsPerSec, maxBurstSize))
	{
		auth.GET("/req_pq", func(c *gin.Context) {
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

		auth.GET("/req_dh_params", func(c *gin.Context) {
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
	}

	// Biz services
	biz := r.Group("/biz")
	biz.Use(isAuthenticated())
	{
		biz.GET("/get_users", func(c *gin.Context) {
			if c.MustGet("authenticated") != true {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "not authenticated",
				})
			} else {
				user_id, _ := strconv.Atoi(c.Query("user_id"))
				message_id, _ := strconv.Atoi(c.Query("message_id"))
				auth_key := c.Query("auth_key")
				response, err := bizClient.GetUsers(context.Background(), &bizPb.GetUsersRequest{UserId: int32(user_id), AuthKey: auth_key, MessageId: int32(message_id)})
				c.JSON(http.StatusOK, gin.H{
					"response": response,
					"err":      err,
				})
			}
		})

		biz.GET("/get_users_with_sql_inject", func(c *gin.Context) {
			if c.MustGet("authenticated") != true {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "not authenticated",
				})
			} else {
				user_id := c.Query("user_id")
				message_id, _ := strconv.Atoi(c.Query("message_id"))
				auth_key := c.Query("auth_key")
				response, err := bizClient.GetUsersWithSQLInject(context.Background(), &bizPb.GetUsersWithSQLInjectRequest{UserId: user_id, AuthKey: auth_key, MessageId: int32(message_id)})
				c.JSON(200, gin.H{
					"response": response,
					"err":      err,
				})
			}
		})
	}

	r.Run("0.0.0.0:6433") // listen and serve on 0.0.0.0:8080
}
