package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	"github.com/jordyf15/thullo-api/middlewares"
	"github.com/jordyf15/thullo-api/token/repository"
	"github.com/jordyf15/thullo-api/token/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbClient    *mongo.Database
	redisClient *redis.Client
	rtdbClient  *db.Client
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

func connectToRTDB() {
	ctx := context.Background()
	conf := &firebase.Config{DatabaseURL: os.Getenv("FIREBASE_DB_URL")}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("Unable to initialize firebase realtime db")
	}

	rtdbClient, err = app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing firebase database client: ", err)
	}
}

func health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	connectToDB()
	connectToRTDB()
	connectToRedis()
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
