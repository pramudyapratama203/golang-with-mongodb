# Go API - Penyimpanan Data In-Memory (Slice & Map)

Repositori ini berisi implementasi API CRUD menggunakan bahasa Go, framework Gin, dan penyimpanan data sementara (in-memory) menggunakan struktur data map dan slice. Cocok untuk pemula yang ingin memahami dasar backend API development.

## Struktur Proyek
```
go-api-in-memory/
├── main.go               // Titik masuk utama
├── models/
│   └── book.go           // Definisi struct Book
├── handlers/
│   └── book_handlers.go  // Handler HTTP CRUD
└── services/
    └── book_service.go   // Logika CRUD in-memory (map)
```

## 1. models/book.go
```go
package models

type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    ISBN   string `json:"isbn"`
    Year   int    `json:"year"`
}
```

## 2. services/book_service.go
```go
package services

import (
    "fmt"
    "sync"
    "go-api-in-memory/models"
)

var booksData = make(map[int]models.Book)
var bookMutex sync.Mutex
var nextID = 1

func InitBook() {
    bookMutex.Lock()
    defer bookMutex.Unlock()
    booksData[1] = models.Book{ID: 1, Title: "1984", Author: "George Orwell", ISBN: "9780451524935", Year: 1949}
    booksData[2] = models.Book{ID: 2, Title: "To Kill a Mockingbird", Author: "Harper Lee", ISBN: "9780061120084", Year: 1960}
    nextID = 3
}

func GetAllBooks() []models.Book {
    bookMutex.Lock()
    defer bookMutex.Unlock()
    var allBooks []models.Book
    for _, book := range booksData {
        allBooks = append(allBooks, book)
    }
    return allBooks
}

func GetBookById(id int) (models.Book, bool) {
    bookMutex.Lock()
    defer bookMutex.Unlock()
    book, found := booksData[id]
    return book, found
}

func CreateBook(newBook models.Book) models.Book {
    bookMutex.Lock()
    defer bookMutex.Unlock()
    newBook.ID = nextID
    nextID++
    booksData[newBook.ID] = newBook
    return newBook
}

func UpdateBook(id int, updatedBook models.Book) (models.Book, bool) {
    bookMutex.Lock()
    defer bookMutex.Unlock()
    book, found := booksData[id]
    if !found {
        return models.Book{}, false
    }
    book.Title = updatedBook.Title
    book.Author = updatedBook.Author
    book.ISBN = updatedBook.ISBN
    book.Year = updatedBook.Year
    booksData[id] = book
    return book, true
}

func DeleteBook(id int) bool {
    bookMutex.Lock()
    defer bookMutex.Unlock()
    _, found := booksData[id]
    if !found {
        return false
    }
    delete(booksData, id)
    return true
}

func DeleteAllBooks() {
    bookMutex.Lock()
    defer bookMutex.Unlock()
    booksData = make(map[int]models.Book)
    nextID = 1
}
```

## 3. handlers/book_handlers.go
```go
package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "go-api-in-memory/models"
    "go-api-in-memory/services"
)

func sendErrorResponse(c *gin.Context, status int, msg string) {
    c.JSON(status, gin.H{"error": msg})
}

func GetAllBooks(c *gin.Context) {
    books := services.GetAllBooks()
    c.JSON(http.StatusOK, books)
}

func GetBookById(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        sendErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
        return
    }
    book, found := services.GetBookById(id)
    if !found {
        sendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan")
        return
    }
    c.JSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
    var book models.Book
    if err := c.BindJSON(&book); err != nil {
        sendErrorResponse(c, http.StatusBadRequest, "Format JSON salah")
        return
    }
    if book.Title == "" {
        sendErrorResponse(c, http.StatusBadRequest, "Judul tidak boleh kosong")
        return
    }
    created := services.CreateBook(book)
    c.JSON(http.StatusCreated, created)
}

func UpdateBook(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        sendErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
        return
    }
    var book models.Book
    if err := c.BindJSON(&book); err != nil {
        sendErrorResponse(c, http.StatusBadRequest, "Format JSON salah")
        return
    }
    updated, ok := services.UpdateBook(id, book)
    if !ok {
        sendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan")
        return
    }
    c.JSON(http.StatusOK, updated)
}

func DeleteBook(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        sendErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
        return
    }
    if !services.DeleteBook(id) {
        sendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan")
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Buku dihapus"})
}

func DeleteAllBooks(c *gin.Context) {
    services.DeleteAllBooks()
    c.JSON(http.StatusOK, gin.H{"message": "Semua buku dihapus"})
}
```

## 4. main.go
```go
package main

import (
    "fmt"
    "log"
    "github.com/gin-gonic/gin"
    "go-api-in-memory/handlers"
    "go-api-in-memory/services"
)

func main() {
    gin.SetMode(gin.ReleaseMode)
    router := gin.Default()
    services.InitBook()

    router.GET("/books", handlers.GetAllBooks)
    router.GET("/books/:id", handlers.GetBookById)
    router.POST("/books", handlers.CreateBook)
    router.PUT("/books/:id", handlers.UpdateBook)
    router.DELETE("/books/:id", handlers.DeleteBook)
    router.DELETE("/books", handlers.DeleteAllBooks)

    fmt.Println("Server berjalan di http://localhost:8080")
    log.Fatal(router.Run(":8080"))
}
```

## Cara Menjalankan
```bash
go mod init go-api-in-memory
go get github.com/gin-gonic/gin
go mod tidy
go run main.go
```

Uji API dengan Postman / cURL:
- `GET /books`
- `GET /books/:id`
- `POST /books`
- `PUT /books/:id`
- `DELETE /books/:id`
- `DELETE /books`
