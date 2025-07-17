package main

import (
	"api-book/handlers"
	"api-book/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode);
	router := gin.Default()

	services.IntBook();

	router.GET("/book", handlers.GetAllBooks);
	router.GET("/book/:id", handlers.GetBookById);
	router.POST("/book", handlers.CreateBook);
	router.PUT("/book/:id", handlers.UpdateBook);
	router.DELETE("/book/:id", handlers.DeleteBook);
	router.DELETE("/book", handlers.DeleteAllBook);

	fmt.Print("http://localhost:8080");
	log.Fatal(router.Run(":8080"));
}