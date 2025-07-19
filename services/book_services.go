package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"api-book/models"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)	

// Instance koneksi ke mongoDB
var client *mongo.Client

// Collection 
var collection *mongo.Collection 

// InitMongoDB 
func InitMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions) // Koneksi ke mongoDB
	if err != nil {
		log.Fatal("Gagal terhubung ke mongoDB:", err)
	}

	err = client.Ping(ctx, nil) // Ping database untuk memastikan koneksi berhasil
	if err != nil {
		log.Fatal("Gagal melakukan ping ke MongoDB:", err)
	}

	fmt.Println("Berhasil terhubung ke mongoDB!")

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

// Get All Book
func GetAllBook() ([]models.Book, error) {
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
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("Gagal encode buku : %w", err)
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
	err := collection.FindOne(ctx, filter).Decode(&book)
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

	newBook.ID = primitive.NilObjectID

	result, err := collection.InsertOne(ctx, newBook)
	if err != nil {
		return models.Book{}, fmt.Errorf("Gagal membuat buku : %w", err)
	}

	newBook.ID = result.InsertedID.(primitive.ObjectID)
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

	return GetBookByID(id)
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