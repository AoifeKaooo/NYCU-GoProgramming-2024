package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Book struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Pages int    `json:"pages"`
}

var bookshelf = []Book{
	{ID: 1, Name: "Blue Bird", Pages: 500},
}
var nextID = 2

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, bookshelf)
}

func getBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid book ID"})
		return
	}
	for _, book := range bookshelf {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func addBook(c *gin.Context) {
	var newBook Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON data"})
		return
	}
	for _, book := range bookshelf {
		if book.Name == newBook.Name {
			c.JSON(http.StatusConflict, gin.H{"message": "duplicate book name"})
			return
		}
	}
	newBook.ID = nextID
	nextID++
	bookshelf = append(bookshelf, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func deleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid book ID"})
		return
	}
	for i, book := range bookshelf {
		if book.ID == id {
			bookshelf = append(bookshelf[:i], bookshelf[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.Status(http.StatusNoContent) // If book ID not found, return 204
}

func updateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid book ID"})
		return
	}

	var updatedBook Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON data"})
		return
	}

	for i, book := range bookshelf {
		if book.ID == id {
			// Check for duplicate name
			for _, b := range bookshelf {
				if b.Name == updatedBook.Name && b.ID != id {
					c.JSON(http.StatusConflict, gin.H{"message": "duplicate book name"})
					return
				}
			}
			bookshelf[i].Name = updatedBook.Name
			bookshelf[i].Pages = updatedBook.Pages
			c.JSON(http.StatusOK, bookshelf[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func main() {
	r := gin.Default()
	r.GET("/bookshelf", getBooks)
	r.GET("/bookshelf/:id", getBook)
	r.POST("/bookshelf", addBook)
	r.DELETE("/bookshelf/:id", deleteBook)
	r.PUT("/bookshelf/:id", updateBook)

	err := r.Run(":8087")
	if err != nil {
		return
	}
}