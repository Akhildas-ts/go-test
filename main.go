package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	// Set the start date to June 1, 2023
	startDate := time.Date(2024, time.January, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()

	for currentDate := startDate; currentDate.Before(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		for j := 0; j < randInt(1, 10); j++ {
			d := strconv.Itoa(currentDate.Day()) + " days ago"
			writeToFile("file.txt", d)
			runCommand("git", "add", ".")
			runCommand("git", "commit", "--date", currentDate.Format("Mon Jan 2 15:04:05 2006 -0700"), "-m", "commit")
			runCommand("git", "push", "-u", "origin", "main")
		}
	}
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func writeToFile(filename, content string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(content + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error executing command:", err)
	}
}



// package main

// import (
// 	"lock/config"
// 	"lock/database"
// 	routes "lock/router"
// 	"log"

// 	"github.com/gin-gonic/gin"
// )

// func main() {

// 	cfg, err := config.LoadConfig()

// 	if err != nil {
// 		log.Fatalf("error loading the config file")
// 	}

// 	db, dberr := database.ConnectDatabase(cfg)

// 	if dberr != nil {
// 		log.Fatalf("error loading the config file")

// 	}
// 	router := gin.Default()
// 	routes.UserRouter(router.Group("/"), db)

// 	err = router.Run("localhost:8080")

// 	if err != nil {

// 		log.Fatalf("Local host error %v", err)
// 	}

// }
