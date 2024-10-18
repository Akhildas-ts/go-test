package main

import (
	"lock/config"
	"lock/database"
	"lock/handlers"
	"lock/repository"
	routes "lock/router"
	"lock/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig()

	if err != nil {

		log.Fatalf("rror found on config connectoin", err)
	}

	db, err := database.ConnectDatabase(cfg)

	if err != nil {

		log.Fatalf("rror found on connectoin of database ", err)
	}

	ur := repository.NewUserRepo(db)
	uc := usecase.UseruseCase(*ur)
	uh := handlers.NewUserHandler(*uc)
	router := gin.Default()

	routes.UserRoutes(router.Group("/"), db, uh)

	router.Run()

}
