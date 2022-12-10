package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/jordyf15/thullo-api/middlewares"
	"github.com/jordyf15/thullo-api/token/repository"
	"github.com/jordyf15/thullo-api/token/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbClient    *mongo.Database
	redisClient *redis.Client
	router      *gin.Engine
)

func connectToDB() {
	connectionURL := os.Getenv("DB_URL")
	clientOptions := options.Client().ApplyURI(connectionURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	dbClient = client.Database(os.Getenv("DB_NAME"))
}

func connectToRedis() {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	fmt.Println(pong, err)
}

func health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func main() {
	router = gin.Default()

	trustedProxies := strings.Split(os.Getenv("TRUSTED_PROXIES"), ",")
	if len(trustedProxies) > 0 {
		router.SetTrustedProxies(trustedProxies)
	}

	if allowedOriginsEnvValue := os.Getenv("ALLOWED_ORIGINS"); len(allowedOriginsEnvValue) > 0 {
		allowedOrigins := strings.Split(allowedOriginsEnvValue, ",")
		config := cors.DefaultConfig()
		config.AllowOrigins = allowedOrigins
		config.AllowHeaders = []string{"Origin", "Authorization"}

		router.Use(cors.New(config))
	}

	tokenRepo := repository.NewTokenRepository(dbClient, redisClient)
	loggerMiddleware := middlewares.NewLoggerMiddleware()
	authMiddleware := middlewares.NewAuthMiddleware(usecase.NewTokenUsecase(tokenRepo))

	if gin.IsDebugging() {
		router.Use(loggerMiddleware.PrintClientIP, loggerMiddleware.PrintHeadersAndFormParams, authMiddleware.AuthenticateJWT)
	} else {
		router.Use(loggerMiddleware.PrintClientIP, authMiddleware.AuthenticateJWT)
	}

	router.MaxMultipartMemory = 10 << 20

	initializeRoutes()
	if os.Getenv("ROUTER_PORT") != "" {
		router.Run(fmt.Sprintf(":%s", os.Getenv("ROUTER_PORT")))
	} else {
		router.Run()
	}
}
