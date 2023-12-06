package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	"github.com/AkifhanIlgaz/random-question-selector/controllers"
	"github.com/AkifhanIlgaz/random-question-selector/middleware"
	"github.com/AkifhanIlgaz/random-question-selector/routes"
	"github.com/AkifhanIlgaz/random-question-selector/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const mongoDBName = "random-question-selector"
const mongoUsersCollection = "users"
const mongoQuestionCollection = "questions"

func main() {
	config, err := cfg.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not read environment variables", err)
	}

	ctx := context.TODO()

	mongoClient, err := connectToMongo(ctx, config.MongoUri)
	if err != nil {
		log.Fatal("Could not connect to mongo", err)
	}

	defer mongoClient.Disconnect(ctx)

	redisClient, err := connectToRedis(ctx, config.RedisUri)
	if err != nil {
		log.Fatal("Could not connect to redis", err)
	}

	defer redisClient.Close()

	value, err := redisClient.Get(ctx, "test").Result()
	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		log.Fatal("error with redis", err)
	}

	userCollection := mongoClient.Database(mongoDBName).Collection(mongoUsersCollection)
	questionCollection := mongoClient.Database(mongoDBName).Collection(mongoQuestionCollection)

	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}
	if _, err := userCollection.Indexes().CreateOne(ctx, index); err != nil {
		log.Fatal("could not create index for email")
	}

	userService := services.NewUserService(ctx, userCollection)
	tokenService := services.NewTokenService(ctx, redisClient, &config)
	questionService := services.NewQuestionService(ctx, questionCollection)

	authController := controllers.NewAuthController(userService, tokenService, &config)
	userController := controllers.NewUserController(userService)
	questionController := controllers.NewQuestionController(questionService)

	userMiddleware := middleware.NewUserMiddleware(userService, tokenService)
	questionMiddleware := middleware.NewQuestionMiddleware()

	authRouteController := routes.NewAuthRouteController(authController, userMiddleware)
	userRouteController := routes.NewUserRouteController(userController, userMiddleware)
	questionRouterController := routes.NewQuestionRouteController(questionController, userMiddleware, questionMiddleware)

	server := gin.Default()

	setCors(server)

	router := server.Group("/api")
	router.GET("/healthchecker", userMiddleware.ExtractUser(), func(ctx *gin.Context) {
		fmt.Println(ctx.Get("currentUser"))

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	authRouteController.AuthRoute(router)
	userRouteController.UserRoute(router)
	questionRouterController.QuestionRoute(router)

	log.Fatal(server.Run(":" + config.Port))
}

func connectToMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	mongoConn := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("connect to mongo: %w", err)
	}

	fmt.Println("MongoDB successfully connected...")
	return client, nil
}

func connectToRedis(ctx context.Context, uri string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: uri,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	err := client.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	fmt.Println("Redis client connected successfully...")

	return client, nil
}

func setCors(server *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
}
