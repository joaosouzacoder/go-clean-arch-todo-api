package main

import (
	"fmt"
	"log"

	envconfig "github.com/caarlos0/env/v10"
	"github.com/gabrielsouzacoder/clean-new/api/routes"
	"github.com/gabrielsouzacoder/clean-new/infrastructure/repository"
	"github.com/gabrielsouzacoder/clean-new/usecase/todo"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Port     string `env:"PORT" envDefault:"8080"`
	DBType   string `env:"DB_TYPE" envDefault:"inmemory"`
	MongoURI string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
}

func main() {
	fmt.Println("[Server] Initializing ...")
	config := loadEnvironment()
	todoRepo := selectDatabase(config)
	todoService := todo.NewService(todoRepo)

	server := NewServer(config)
	server.Run(todoService)
}

func loadEnvironment() Config {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("[Warning] The .env file could not be loaded")
	}

	var config Config
	if err := envconfig.Parse(&config); err != nil {
		log.Fatal(err)
	}

	return config
}

type Server struct {
	port   string
	server *gin.Engine
}

func NewServer(config Config) Server {
	return Server{
		port:   config.Port,
		server: gin.Default(),
	}
}

func (s *Server) Run(todo *todo.Service) {
	router := routes.ConfigRoutes(s.server, todo)

	log.Printf("Server running at port: %v", s.port)
	log.Fatal(router.Run(":" + s.port))
}

func selectDatabase(config Config) todo.Repository {
	var todoRepo todo.Repository

	if config.DBType == "mongo" {
		todoRepo = repository.NewMongoDbRepository(options.Client().ApplyURI(config.MongoURI))
	} else {
		todoRepo = repository.NewInMemoryDatabase()
	}
	return todoRepo
}
