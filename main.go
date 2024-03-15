package main

import (
	"lock/config"
	"lock/database"
	routes "lock/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("error loading the config file")
	}

	db, dberr := database.ConnectDatabase(cfg)

	if dberr != nil {
		log.Fatalf("error loading the config file")

	}
	router := gin.Default()
	routes.UserRouter(router.Group("/"), db)

	err = router.Run("localhost:8080")

	if err != nil {

		log.Fatalf("Local host error %v", err)
	}

}
