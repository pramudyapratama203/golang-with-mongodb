package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"api-book/handlers"
	"api-book/services"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	services.InitMongoDB()
	defer services.CloseMongoDB()
	
	router := gin.Default()

	router.GET("/books", handlers.GetAllBook)
	router.GET("/books/:id", handlers.GetBookByID)
	router.POST("/books", handlers.CreateBook)
	router.PUT("/books/:id", handlers.UpdateBook)
	router.DELETE("/books/:id", handlers.DeleteBook)

	fmt.Println("http://localhost:8080");
	log.Fatal(router.Run(":8080"));
}