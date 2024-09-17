package todo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thompsonmanda08/go-webb-todo/database"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func GetTodos(c *fiber.Ctx) error {

	// GET DATABASE CONNECTION
	db := database.DBConn
	var todos []Todo

	// GET ALL TODOS
	db.Find(&todos)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "",
		"data":    todos,
		"status":  fiber.StatusOK,
	})
}

func GetTodo(c *fiber.Ctx) error {

	id := c.Params("id")
	db := database.DBConn

	var todo Todo

	db.Find(&todo, id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		// "message": ""
		"data":   todo,
		"status": fiber.StatusOK,
	})
}

// func updateTodo(c *fiber.Ctx) error {
// 	return c.SendString("All todos")
// }

func NewTodo(c *fiber.Ctx) error {
	db := database.DBConn

	todo := new(Todo)

	c.BodyParser(todo)

	if err := c.BodyParser(todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if todo.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "A title is required",
			"data":    nil,
			"status":  fiber.StatusBadRequest,
		})
	}

	db.Create(&todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Todo was not found",
		"data":    todo,
		"status":  fiber.StatusNotFound,
	})
}

func UpdateTodo(c *fiber.Ctx) error {
	db := database.DBConn

	id := c.Params("id")

	var todo Todo

	db.First(&todo, id)

	// FIDN THE EXISTING TODO IN DB
	if err := db.First(&todo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Todo was not found",
			"data":    nil,
			"status":  fiber.StatusNotFound,
		})
	}

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"data":    nil,
			"status":  fiber.StatusBadRequest,
		})
	}

	// Save the updated todo item back to the database
	if err := db.Save(&todo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update todo",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Todo updated successfully",
		"data":    todo,
		"status":  fiber.StatusOK,
	})
}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	var todo Todo

	db.First(&todo, id)

	db.Delete(&todo)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Todo deleted successfully",
		"data":    todo,
		"status":  fiber.StatusOK,
	})

}
