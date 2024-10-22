package main

import (
	"fmt"
	"lock/config"
	"lock/database"
	"lock/handlers"
	"lock/repository"
	routes "lock/router"
	"lock/usecase"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig()

	if err != nil {

		fmt.Println("Error on config ", err)
	}

	db, err := database.ConnectDatabase(cfg)

	if err != nil {

		fmt.Println("rror found on connectoin of database ", err)
	}

	ur := repository.NewUserRepo(db)
	uc := usecase.UseruseCase(ur)
	uh := handlers.NewUserHandler(uc)
	router := gin.Default()

	routes.UserRoutes(router.Group("/"), db, uh)

	router.Run()

}
