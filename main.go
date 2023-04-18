package main

import (
	"fmt"
	"github.com/GarnBarn/common-go/database"
	"github.com/GarnBarn/common-go/httpserver"
	"github.com/GarnBarn/common-go/logger"
	"github.com/GarnBarn/common-go/proto"
	"github.com/GarnBarn/gb-tag-service/config"
	"github.com/GarnBarn/gb-tag-service/handler"
	"github.com/GarnBarn/gb-tag-service/repository"
	"github.com/GarnBarn/gb-tag-service/service"
	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var appConfig config.Config

func init() {
	appConfig = config.Load()
	logger.InitLogger(logger.Config{
		Env: appConfig.Env,
	})

}

func main() {
	// Start DB Connection
	db, err := database.Conn(appConfig.MYSQL_CONNECTION_STRING)
	if err != nil {
		logrus.Panic("Can't connect to db: ", err)
	}

	// Create the required dependentices
	validate := validator.New()

	// Create the repositories
	tagRepository := repository.NewTagRepository(db)

	// Create the services
	tagService := service.NewTagService(tagRepository)

	// Init the handler
	tagHandler := handler.NewTagHandler(*validate, tagService)
	grpcHandler := handler.NewGrpcServer(tagService)

	// Create the http server
	httpServer := httpserver.NewHttpServer()

	httpServer.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// Router
	router := httpServer.Group("/api/v1/tag")
	router.GET("/:id", tagHandler.GetTagById)
	router.GET("/", tagHandler.GetAllTag)
	router.POST("/", tagHandler.CreateTag)
	router.PATCH("/:tagId", tagHandler.UpdateTag)
	router.DELETE("/:tagId", tagHandler.DeleteTag)

	// GRPC
	lis, err := net.Listen("tcp", fmt.Sprint(":", appConfig.GRPC_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterTagServer(s, grpcHandler)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	logrus.Info("Listening and serving HTTP on :", appConfig.HTTP_SERVER_PORT)
	httpServer.Run(fmt.Sprint(":", appConfig.HTTP_SERVER_PORT))
}
