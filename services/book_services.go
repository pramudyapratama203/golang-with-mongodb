package services

import (
	"api-book/models"
	"sync"
)	

var booksData = make(map[int]models.Book)
var bookMuted sync.Mutex
var nextId int

func IntBook() {
	bookMuted.Lock()
	defer bookMuted.Unlock()

	booksData[1] = models.Book{ID: 1, Title: "1984", Author: "George Orwell", ISBN: "10239", Year: 1994}
	nextId = 2
}

func GetAllBooks() []models.Book{
	bookMuted.Lock()
	defer bookMuted.Unlock()

	var allBook []models.Book
	for _, book := range booksData {
		allBook = append(allBook, book)
	}

	return allBook
}

func GetBookById(id int) (models.Book, bool) {
	bookMuted.Lock()
	defer bookMuted.Unlock()

	book, found := booksData[id]
	return book, found
}

func CreateBook(newBook models.Book) models.Book{
	bookMuted.Lock()
	defer bookMuted.Unlock()

	newBook.ID = nextId
	nextId++
	booksData[newBook.ID] = newBook
	return newBook
}

func UpdateBook(id int, updateBook models.Book) (models.Book, bool) {
	bookMuted.Lock()
	defer bookMuted.Unlock()

	update, found := booksData[id]
	if !found {
		return models.Book{}, false
	}

	update.Title = updateBook.Title
	update.Author = updateBook.Author
	update.ISBN = updateBook.ISBN
	update.Year = updateBook.Year
	booksData[id] = update
	return update, true
}

func DeleteBook(id int) bool {
	bookMuted.Lock()
	defer bookMuted.Unlock()

	_, found := booksData[id]
	if !found {
		return false
	}

	delete(booksData, id)
	return true
}

func DeleteAllBook() {
	bookMuted.Lock()
	defer bookMuted.Unlock()

	booksData = make(map[int]models.Book);
	nextId = 1;
}