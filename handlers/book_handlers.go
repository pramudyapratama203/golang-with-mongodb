package handlers

import (
	"api-book/services"
	"api-book/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error" : message});
}

func GetAllBooks(c *gin.Context) {
	books := services.GetAllBooks();
	c.JSON(http.StatusOK, books);
}

func GetBookById(c *gin.Context) {
	idStr := c.Param("id");

	idInt, err := strconv.Atoi(idStr);
	if err != nil {
		sendError(c, http.StatusBadRequest, "Id harus berupa angka");
		return;
	}

	getBook, success := services.GetBookById(idInt);
	if !success {
		sendError(c, http.StatusBadRequest, "Data tidak ada");
		return;
	}

	c.JSON(http.StatusOK, getBook)
}

func CreateBook(c *gin.Context) {
	var newBook models.Book

	if err := c.BindJSON(&newBook); err != nil {
		sendError(c, http.StatusBadRequest, "Format buku tidak valid");
		return;
	}

	if newBook.Title == "" {
		sendError(c, http.StatusBadRequest, "Judul buku harus diisi");
		return;
	}

	create := services.CreateBook(newBook);
	c.JSON(http.StatusOK, create);
}

func UpdateBook(c *gin.Context) {
	idStr := c.Param("id");
	var updateBook models.Book

	if err := c.BindJSON(&updateBook); err != nil {
		sendError(c, http.StatusBadRequest, "Format update tidak valid");
		return;
	}

	idInt, err := strconv.Atoi(idStr);
	if err != nil {
		sendError(c, http.StatusBadRequest, "Id harus berupa angka");
		return;
	}

	update, success := services.UpdateBook(idInt, updateBook);
	if !success {
		sendError(c, http.StatusBadRequest, "Gagal update buku");
		return;
	}

	c.JSON(http.StatusOK, update);
}

func DeleteBook(c *gin.Context) {
	idStr := c.Param("id");

	idInt, err := strconv.Atoi(idStr);
	if err != nil {
		sendError(c, http.StatusBadRequest, "Id harus berupa angka");
		return;
	}

	delete := services.DeleteBook(idInt);
	if !delete {
		sendError(c, http.StatusBadRequest, "Gagal menghapus buku");
		return;
	}
	
	c.JSON(http.StatusOK, delete);
}

func DeleteAllBook(c *gin.Context) {
	services.DeleteAllBook();
	c.JSON(http.StatusOK, gin.H{"message" : "Semua data berhasil dihapus"});
}