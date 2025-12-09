package main

import (
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var store = sessions.NewCookieStore([]byte("secret-key")) // Use a real secret key in production!

// Define User model
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func main() {
	// Initialize Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Initialize SQLite database
	db, err := gorm.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Migrate User table
	db.AutoMigrate(&User{})
	
	fmt.Println("Hi from login-go-t")
}
