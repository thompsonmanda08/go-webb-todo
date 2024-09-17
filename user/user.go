package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thompsonmanda08/go-webb-todo/database"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleLogin(c *fiber.Ctx) error {

	loginReq := new(LoginRequest)
	user := new(User)

	// PARSE REQUEST BODY
	if err := c.BodyParser(loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// VALIDATE REQUEST BODY - NONE EMPTY
	if loginReq.Email == "" || loginReq.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
		})
	}

	// FIND THE USER WITH MATCHING EMAIL
	if err := database.DBConn.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User does not exist",
			"status":  fiber.StatusNotFound,
			"data":    nil,
		})
	}

	// CHECK HASHED THE PASSWORD
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
		})
	}

	/// GENERATE JWT TOKEN AND LOG USER IN
	token, expiry, err := CreateJWT(*user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(), // FAILED TO GENERATE TOKEN
			"status":  fiber.StatusInternalServerError,
			"data":    nil,
		})
	}

	sanitizedUser := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	// SEND RESPONSE WITH AUTHENTICATED USER
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusAccepted,
		"message": "Login Successful",
		"data": fiber.Map{
			"token":  token,
			"expiry": expiry,
			"user":   sanitizedUser,
		},
	})
}

func HandleRegistration(c *fiber.Ctx) error {

	newUser := new(User)

	// PARSE REQUEST BODY
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// VALIDATE REQUEST BODY - NONE EMPTY
	if newUser.Email == "" || newUser.Password == "" || newUser.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "All fields are required",
		})
	}

	// HASH THE PASSWORD
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	// SET THE HASHED PASSWORD
	newUser.Password = string(hash)

	// INSERT USER INTO DATABASE
	if err := database.DBConn.Create(newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
			"status":  fiber.StatusInternalServerError,
		})
	}

	/// GENERATE JWT TOKEN AND LOG USER IN
	token, expiry, err := CreateJWT(*newUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(), // FAILED TO GENERATE TOKEN
			"status":  fiber.StatusInternalServerError,
		})
	}

	// SEND RESPONSE WITH AUTHENTICATED USER
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusAccepted,
		"message": "Registration Successful",
		"data": fiber.Map{
			"token":  token,
			"expiry": expiry,
			"user":   newUser,
		},
	})
}

func GetAllUsers(c *fiber.Ctx) error {

	// GET DATABASE CONNECTION
	db := database.DBConn
	var users []User

	// GET ALL USERS
	result := db.Find(&users)

	// Fetch all users from the database

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error fetching users"})
	}

	// Create a new slice to hold sanitized users
	var sanitizedUsers []UserResponse

	for _, newUser := range users {
		// Add the sanitized version without password and other unwanted fields
		sanitizedUser := UserResponse{
			ID:    newUser.ID,
			Name:  newUser.Name,
			Email: newUser.Email,
		}
		sanitizedUsers = append(sanitizedUsers, sanitizedUser)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "",
		"data":    sanitizedUsers,
		"status":  fiber.StatusOK,
	})
}

func CreateJWT(user User) (string, int64, error) {
	expires := time.Now().Add(time.Minute + 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["expires"] = expires

	t, err := token.SignedString([]byte("secret-key"))

	if err != nil {
		return "", 0, err
	}

	return t, expires, nil

}
