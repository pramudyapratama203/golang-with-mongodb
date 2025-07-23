# API Buku (Go + Gin + MongoDB)

Proyek ini adalah implementasi dari API CRUD (Create, Read, Update, Delete) untuk entitas Buku. Semua data buku disimpan secara persisten di database MongoDB.

## Fitur API
API ini menyediakan operasi dasar untuk mengelola buku:
1. GET /books: Mengambil daftar semua buku yang tersimpan.
2. GET /books/:id: Mengambil detail buku berdasarkan ID uniknya (dari MongoDB).
3. POST /books: Menambahkan buku baru ke koleksi.
4. PUT /books/:id: Memperbarui detail buku yang sudah ada.
5. DELETE /books/:id: Menghapus buku dari koleksi.
6. DELETE /books: Menghapus semua buku yang ada di koleksi.

## Cara Menjalankan
```bash
go mod init api-book
go get github.com/gin-gonic/gin
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
go get go.mongodb.org/mongo-driver/bson
go mod tidy // Disarankan untuk membersihkan dependensi yang tidak terpakai
go run main.go
```
## Struktur Proyek
```
api-book/
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
// api-book/models/book.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive" // Import untuk tipe ObjectID MongoDB

// Book merepresentasikan satu item buku di database MongoDB.
type Book struct {
	// ID dari MongoDB akan otomatis di-generate.
	// `bson:"_id"`: Memberitahu Go ini adalah field '_id' di MongoDB.
	// `json:"id,omitempty"`: Field ini akan disebut 'id' di JSON, dan 'omitempty' berarti
	// tidak akan muncul di JSON jika nilainya kosong (saat membuat baru).
	ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title  string             `bson:"title" json:"title"`
	Author string             `bson:"author" json:"author"`
	ISBN   string             `bson:"isbn" json:"isbn"`
	Year   int                `bson:"year" json:"year"`
}
```

## 2. services/book_service.go
```go
// api-book/services/book_service.go
package services

import (
	"context" // Penting untuk semua operasi database di Go
	"fmt"
	"log"
	"time"

	"api-book/models" // Import struct Book dari models (sesuaikan dengan nama modulmu)

	"go.mongodb.org/mongo-driver/bson"          // Untuk membangun query BSON
	"go.mongodb.org/mongo-driver/bson/primitive" // Untuk tipe ObjectID
	"go.mongodb.org/mongo-driver/mongo"         // Driver utama MongoDB
	"go.mongodb.org/mongo-driver/mongo/options" // Untuk opsi koneksi
)  

// client adalah instance koneksi ke MongoDB.
var client *mongo.Client

// collection adalah referensi ke koleksi "buku" di database "perpustakaanDigital".
var collection *mongo.Collection

// InitMongoDB: Fungsi ini akan menginisialisasi koneksi ke MongoDB.
func InitMongoDB() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    client, err = mongo.Connect(ctx, clientOptions) // Koneksi ke mongoDB
    if err != nil {
        log.Fatal("Gagal terhubung ke MongoDB:", err)
    }

    err = client.Ping(ctx, nil) // Ping database untuk memastikan koneksi berhasil
    if err != nil {
        log.Fatal("Gagal melakukan ping ke MongoDB:", err)
    }

    fmt.Println("Berhasil terhubung ke MongoDB!")

    // Dapatkan referensi ke database "perpustakaanDigital" dan koleksi "buku"
    collection = client.Database("perpustakaanDigital").Collection("buku")
    fmt.Println("Koleksi 'buku' siap digunakan")
}

// CloseMongoDB
func CloseMongoDB() {
    if client == nil {
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := client.Disconnect(ctx)
    if err != nil {
        log.Fatal("Gagal memutuskan koneksi MongoDB:", err)
    }
    fmt.Println("Koneksi berhasil ditutup")
}

// Get All Books
func GetAllBooks() ([]models.Book, error) { // Nama fungsi disesuaikan
    var books []models.Book
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{}) // Cari semua dokumen
    if err != nil {
        return nil, fmt.Errorf("Gagal mencari buku : %w", err)
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var book models.Book
        if err := cursor.Decode(&book); err != nil { // Menggunakan &book dan jangan lupa decode
            return nil, fmt.Errorf("Gagal decode buku dari cursor : %w", err) // Pesan error diperbaiki
        }
        books = append(books, book)
    }

    if err := cursor.Err(); err != nil {
        return nil, fmt.Errorf("Error cursor mongodb: %w", err)
    }

    return books, nil
}

// Get Book by ID
func GetBookByID(id primitive.ObjectID) (models.Book, error) {
    var book models.Book
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"_id" : id} // Filter untuk mencari berdasarkan _id
    err := collection.FindOne(ctx, filter).Decode(&book) // Menggunakan &book
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return models.Book{}, fmt.Errorf("buku tidak ditemukan")
        }
        return models.Book{}, fmt.Errorf("Gagal mencari buku berdasarkan ID : %w", err)
    }

    return book, nil
}

// Create book
func CreateBook(newBook models.Book) (models.Book, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    newBook.ID = primitive.NilObjectID // Set ID ke nilai kosong agar MongoDB yang generate

    result, err := collection.InsertOne(ctx, newBook) // Masukkan dokumen baru
    if err != nil {
        return models.Book{}, fmt.Errorf("Gagal membuat buku : %w", err)
    }

    newBook.ID = result.InsertedID.(primitive.ObjectID) // Ambil ID yang baru di-generate
    fmt.Printf("Buku baru ditambahkan dengan ID: %s\n", newBook.ID.Hex())
    return newBook, nil
}

// Update book
func UpdateBook(id primitive.ObjectID, updateBook models.Book) (models.Book, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"_id" : id}
    update := bson.M{
        "$set": bson.M{
            "title": updateBook.Title,
            "author": updateBook.Author,
            "isbn": updateBook.ISBN,
            "year": updateBook.Year,
        },
    }

    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return models.Book{}, fmt.Errorf("gagal memperbarui buku: %w", err)
    }

    if result.MatchedCount == 0 { // Tidak ada dokumen yang cocok
        return models.Book{}, fmt.Errorf("buku tidak ditemukan untuk diperbarui")
    }

    return GetBookByID(id) // Ambil dokumen yang sudah diperbarui untuk dikembalikan
}

// Delete book
func DeleteBook(id primitive.ObjectID) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"_id" : id}
    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return false, fmt.Errorf("Gagal menghapus buku: %w", err)
    }

    if result.DeletedCount == 0 {
        return false, fmt.Errorf("Buku tidak ditemukan untuk dihapus")
    }

    fmt.Printf("Buku dengan ID %s dihapus.\n", id.Hex())
    return true, nil
}

// Delete All Books
func DeleteAllBooks() (bool, error) { // Nama fungsi disesuaikan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// collection.DeleteMany(ctx, bson.M{}): Menghapus semua dokumen (filter kosong).
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return false, fmt.Errorf("gagal menghapus semua buku: %w", err)
	}

	fmt.Println("Semua buku berhasil dihapus dari MongoDB.")
	return true, nil

}
```

## 3. handlers/book_handlers.go
```go
// api-book/handlers/book_handlers.go
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive" // Import untuk tipe ObjectID
	// go.mongodb.org/mongo-driver/mongo // Tidak lagi dibutuhkan secara langsung di sini

	"api-book/models"
	"api-book/services"
)

// sendErrorResponse adalah fungsi pembantu untuk mengirim response error JSON yang konsisten.
func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

// --- Handler Functions ---

// GetAllBooks: Handler untuk mengambil semua buku.
func GetAllBooks(c *gin.Context) {
	books, err := services.GetAllBooks()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Gagal mengambil buku: %v", err))
		return
	}
	c.JSON(http.StatusOK, books)
}

// GetBookById: Handler untuk mengambil satu buku berdasarkan ID.
func GetBookById(c *gin.Context) {
	idStr := c.Param("id") // Ambil ID dari URL (masih string)

	// Konversi string ID dari URL menjadi primitive.ObjectID.
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	book, err := services.GetBookById(objectID)
	if err != nil {
		if err.Error() == "buku tidak ditemukan" { // Cek error spesifik dari service
			sendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan")
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Gagal mengambil buku: %v", err))
		}
		return
	}
	c.JSON(http.StatusOK, book)
}

// CreateBook: Handler untuk membuat buku baru.
func CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.BindJSON(&newBook); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Format data buku tidak valid: "+err.Error())
		return
	}

	if newBook.Title == "" || newBook.Author == "" { // Validasi judul/penulis tidak boleh kosong
		sendErrorResponse(c, http.StatusBadRequest, "Judul dan Penulis buku tidak boleh kosong")
		return
	}

	createdBook, err := services.CreateBook(newBook)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Gagal membuat buku: %v", err))
		return
	}
	c.JSON(http.StatusCreated, createdBook)
}

// UpdateBook: Handler untuk memperbarui buku yang sudah ada.
func UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idStr) // Konversi ID string ke ObjectID
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	var updatedData models.Book
	if err := c.BindJSON(&updatedData); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Format data update buku tidak valid: "+err.Error())
		return
	}

	if updatedData.Title == "" || updatedData.Author == "" { // Validasi judul/penulis tidak boleh kosong
		sendErrorResponse(c, http.StatusBadRequest, "Judul dan Penulis buku tidak boleh kosong")
		return
	}

	updatedBook, err := services.UpdateBook(objectID, updatedData)
	if err != nil {
		if err.Error() == "buku tidak ditemukan untuk diperbarui" {
			sendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan")
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Gagal memperbarui buku: %v", err))
		}
		return
	}
	c.JSON(http.StatusOK, updatedBook)
}

// DeleteBook: Handler untuk menghapus buku berdasarkan ID.
func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idStr) // Konversi ID string ke ObjectID
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "ID buku tidak valid")
		return
	}

	success, err := services.DeleteBook(objectID)
	if err != nil {
		if err.Error() == "buku tidak ditemukan untuk dihapus" {
			sendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan")
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Gagal menghapus buku: %v", err))
		}
		return
	}

	if !success { // Fallback, seharusnya sudah ditangani oleh error dari service
		sendErrorResponse(c, http.StatusNotFound, "Buku tidak ditemukan")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dihapus"})
}

// DeleteAllBooks: Handler untuk menghapus semua buku.
func DeleteAllBooks(c *gin.Context) {
	success, err := services.DeleteAllBooks()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Gagal menghapus semua buku: %v", err))
		return
	}
	if !success { // Seharusnya selalu true jika tidak ada error
		sendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus semua buku karena alasan tidak diketahui.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Semua buku berhasil dihapus"})
}
```

## 4. main.go
```go
// api-book/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"api-book/handlers" // Import handlers
	"api-book/services" // Import services
)

func main() {
	// Set Gin ke Release Mode untuk mengurangi log debug di konsol saat produksi.
	gin.SetMode(gin.ReleaseMode)

	// --- Inisialisasi Koneksi MongoDB ---
	services.InitMongoDB() // Panggil fungsi inisialisasi koneksi MongoDB
	// defer services.CloseMongoDB() ini memastikan koneksi ke MongoDB ditutup
	// dengan rapi saat aplikasi Go berhenti (misal Ctrl+C ditekan).
	defer services.CloseMongoDB()

	router := gin.Default() // Inisialisasi router Gin

	// --- Definisi Endpoints API ---
	// Setiap jalur (path) menunjuk ke Pelayan (handler) yang tepat.
	router.GET("/books", handlers.GetAllBooks)
	router.GET("/books/:id", handlers.GetBookById)
	router.POST("/books", handlers.CreateBook)
	router.PUT("/books/:id", handlers.UpdateBook)
	router.DELETE("/books/:id", handlers.DeleteBook)
	router.DELETE("/books", handlers.DeleteAllBooks) // Endpoint untuk hapus semua

	fmt.Println("Server API Buku (MongoDB) berjalan di http://localhost:8080")
	// Mulai server Gin di port 8080. Jika ada error fatal saat startup, program akan berhenti.
	log.Fatal(router.Run(":8080"))
}
```

Uji API dengan Postman / cURL:
- `GET /books`
- `GET /books/:id`
- `POST /books`
- `PUT /books/:id`
- `DELETE /books/:id`
- `DELETE /books`
