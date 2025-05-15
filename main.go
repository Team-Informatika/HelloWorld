package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Logger middleware
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		fmt.Printf("[%s] %s %s - %v\n",
			c.Method(),
			c.Path(),
			c.IP(),
			duration,
		)
		return err
	}
}

// Message represents a simple message structure
type Message struct {
	Content string `json:"content" validate:"required,min=1,max=100"`
	Author  string `json:"author" validate:"required,min=1,max=50"`
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// In-memory user store
var users = make(map[string]*User)

func main() {
	// Create new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(ErrorResponse{
				Error:   "Request Error",
				Message: err.Error(),
			})
		},
	})

	// Use logger middleware
	app.Use(Logger())

	// Use recover middleware
	app.Use(recover.New())

	// Add compression
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Add rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Define a route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Add a simple API endpoint
	app.Get("/api/info", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app_name":  "Simple API",
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
			"status":    "running",
		})
	})

	// Add message endpoint
	app.Post("/api/message", func(c *fiber.Ctx) error {
		message := new(Message)

		// Parse request body
		if err := c.BodyParser(message); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate message
		if message.Content == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Message content is required")
		}
		if message.Author == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Author is required")
		}
		if len(message.Content) > 100 {
			return fiber.NewError(fiber.StatusBadRequest, "Message content too long")
		}
		if len(message.Author) > 50 {
			return fiber.NewError(fiber.StatusBadRequest, "Author name too long")
		}

		// Return success response
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Message received successfully",
			"data":      message,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Add user management endpoints
	app.Post("/api/users", func(c *fiber.Ctx) error {
		user := new(User)

		// Parse request body
		if err := c.BodyParser(user); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate user data
		if user.Username == "" || user.Email == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Username and email are required")
		}

		// Generate ID and set creation time
		user.ID = fmt.Sprintf("user_%d", time.Now().Unix())
		user.CreatedAt = time.Now()

		// Store user
		users[user.ID] = user

		return c.Status(fiber.StatusCreated).JSON(user)
	})

	app.Get("/api/users/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		user, exists := users[userID]
		if !exists {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return c.JSON(user)
	})

	app.Get("/api/users", func(c *fiber.Ctx) error {
		userList := make([]*User, 0, len(users))
		for _, user := range users {
			userList = append(userList, user)
		}
		return c.JSON(userList)
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
