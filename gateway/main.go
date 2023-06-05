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

	docs "webserver/gateway/docs"
	authPb "webserver/gateway/pb/auth"
	bizPb "webserver/gateway/pb/biz"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
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
		if authIdentifier == "" {
			authIdentifier = c.PostForm("auth_id")
			authKey = c.PostForm("auth_key")
		}
		redisSavedKey, err := redisClient.Get(redisCtx, authIdentifier).Result()

		fmt.Println(authIdentifier, authKey, redisSavedKey)

		if err != nil {
			c.Error(errors.New("cannot provide biz services to unauthorized users"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if  redisSavedKey == authKey {
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

var authCtx, _ = context.WithTimeout(context.Background(), 5*time.Second)
var authConn, authErr = grpc.DialContext(authCtx, "authentication:5052", grpc.WithTransportCredentials(insecure.NewCredentials()))

var authClient = authPb.NewAuthenticationClient(authConn)

var bizCtx, _ = context.WithTimeout(context.Background(), 5*time.Second)
var bizConn, bizErr = grpc.DialContext(bizCtx, "business_logic:5062", grpc.WithTransportCredentials(insecure.NewCredentials()))
var bizClient = bizPb.NewBusinessClient(bizConn)

// @BasePath /

// RequestPG godoc
// @Summary request p/g params diffie-hellman
// @Schemes
// @Description this endpoint will request the p and g parameters of the diffie-hellman
// @Tags req_pq
// @Params nonce message_id
// @Accept json
// @Produce json
// @Success 200 {string} Diffie-Hellman params
// @Router /auth/req_pq [get]
func HandleReqPQ(c *gin.Context) {
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
}

// @BasePath /

// RequestDHParams godoc
// @Summary handshake on diffie-hellman parameters between server/client
// @Schemes
// @Description after receiving p/g params from req_pq, communicate the keys
// @Tags req_dh_params
// @Accept json
// @Produce json
// @Success 200 {string} Handshake on diffie-hellman keys
// @Router /auth/req_dh_params [get]
func HandleReqDHParams(c *gin.Context) {
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
}

// @BasePath /

// GetUsers godoc
// @Summary lets you do a query on db with secure functionality and no injection
// @Schemes
// @Description this endpoint will fetch the user with the id provided in the request or if no id is provided will return last 100 rows of db
// @Tags get_users
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /biz/get_users [get]
func HandleGetUsers(c *gin.Context) {
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
}

// @BasePath /

// GetUsersWithSQLInjection godoc
// @Summary enables you to perform a query on the db with injection feature
// @Schemes
// @Description will not perform any security check on the input given by user and will execute the query without any checks
// @Tags get_users_with_sql_injection
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /biz/get_users_with_sql_injection [get]
func HandleGetUsersWithSQLInject(c *gin.Context) {
	if c.MustGet("authenticated") != true {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authenticated",
		})
	} else {
		user_id := c.PostForm("user_id")
		message_id, _ := strconv.Atoi(c.PostForm("message_id"))
		auth_key := c.PostForm("auth_key")
		response, err := bizClient.GetUsersWithSQLInject(context.Background(), &bizPb.GetUsersWithSQLInjectRequest{UserId: user_id, AuthKey: auth_key, MessageId: int32(message_id)})
		c.JSON(200, gin.H{
			"response": response,
			"err":      err,
		})
	}
}

// @title           Web 2023 Backend HW - Gateway service
// @version         1.0
// @description     This is a simple documentation on the gateway service which is a part of backend homework of WebDev course in spring 2023.
// @termsOfService  http://swagger.io/terms/

// @contact.name   MohammadMahdi Gheidi
// @contact.email  gheidimahdi@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:6433
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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

	docs.SwaggerInfo.BasePath = "/"
	if authErr != nil {
		log.Fatalf("failed to connect: %v", authErr)
	}
	defer authConn.Close()
	if bizErr != nil {
		log.Fatalf("failed to connect to business services: %v", bizErr)
	}
	defer bizConn.Close()

	r := gin.Default()
	r.Use(Blacklist())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth services
	auth := r.Group("/auth")
	maxEventsPerSec := 100
	maxBurstSize := 100
	auth.Use(Throttle(maxEventsPerSec, maxBurstSize))
	{
		auth.GET("/req_pq", HandleReqPQ)
		auth.GET("/req_dh_params", HandleReqDHParams)
	}

	// Biz services
	biz := r.Group("/biz")
	biz.Use(isAuthenticated())
	{
		biz.GET("/get_users", HandleGetUsers)
		biz.POST("/get_users_with_sql_inject", HandleGetUsersWithSQLInject)
	}

	r.Run("0.0.0.0:6433") // listen and serve gateway service on port 6433
}
