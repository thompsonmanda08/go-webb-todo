package main

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/thompsonmanda08/go-webb-todo/database"
	"github.com/thompsonmanda08/go-webb-todo/todo"
	"github.com/thompsonmanda08/go-webb-todo/user"

	jwtware "github.com/gofiber/contrib/jwt"
)

func InitDB() error {
	var err error

	database.DBConn, err = gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")

	}

	fmt.Print("Database connected successfully!\nMigrating tables... \n\n")

	database.DBConn.AutoMigrate(&todo.Todo{})
	database.DBConn.AutoMigrate(&user.User{})
	fmt.Print("Database migrated!")
	return nil

}

func handleSetupRoutes(app *fiber.App) {

	public := app.Group("/api/v1")
	private := app.Group("/api/v1")

	// JWT Middleware
	private.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret-key")},
	}))

	// PUBLIC ROUTES
	public.Post("/login", user.HandleLogin)
	public.Post("/register", user.HandleRegistration)

	// PRIVATE ROUTES
	private.Get("/users", user.GetAllUsers)
	private.Get("/todos", todo.GetTodos)
	private.Get("/todo/:id", todo.GetTodo)
	private.Post("/todo", todo.NewTodo)
	private.Delete("/todo/:id", todo.DeleteTodo)
	private.Patch("/todo/:id", todo.UpdateTodo)

	// For testing purposes, we can add this route here to fetch all todos.
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Public Route Access Successful",
			"status":  fiber.StatusOK,
		})
	})
}

func main() {
	app := fiber.New()

	if err := InitDB(); err != nil {
		panic(err)
	}

	handleSetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

	fmt.Println("Server running on http://localhost:8080...")

}
