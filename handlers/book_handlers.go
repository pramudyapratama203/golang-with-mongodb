package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"api-book/models"
	"api-book/services"
)

func sendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}	

func GetAllBook(c *gin.Context) {
	books, err := services.GetAllBook()
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("Gagal mengambil buku", err))
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetBookByID(c *gin.Context) {
	idStr := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "ID Buku tidak valid")
		return
	}

	book, err := services.GetBookByID(objectID)
	if err != nil {
		sendError(c, http.StatusNotFound, "Buku tidak ditemukan")
		return
	}
	c.JSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.BindJSON(&newBook); err != nil {
		sendError(c, http.StatusBadRequest, "Format buku tidak valid")
		return
	}

	if newBook.Title == "" || newBook.Author == "" {
		sendError(c, http.StatusNotFound, "Judul dan Penulis buku tidak boleh kosong")
		return
	}

	createBook, err := services.CreateBook(newBook)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("Gagal membuat buku: %v", err))
		return
	}
	c.JSON(http.StatusCreated, createBook)
}

func UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	var updateBook models.Book
	if err := c.BindJSON(&updateBook); err != nil {
		sendError(c, http.StatusBadRequest, "Format buku tidak valid")
		return
	}

	book, err := services.UpdateBook(objectID, updateBook)
	if err != nil {
		sendError(c, http.StatusNotFound, "Buku tidak dtemukan")
		return
	}
	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		sendError(c, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	success, err := services.DeleteBook(objectID)
	if err != nil {
		sendError(c, http.StatusNotFound, "Buku tidak dtemukan")
		return
	}
	c.JSON(http.StatusOK, success)
}