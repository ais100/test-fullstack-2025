package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Realname string `json:"real_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Username string `json:"username"`
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

func main() {
	red := newRedisClientEnv()
	ctx := context.Background()

	if err := red.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed connect to Redis: %v", err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var req loginRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		if req.Username == "" || req.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "username and password are required",
			})
		}

		key := "login_" + req.Username
		val, err := red.Get(ctx, key).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to get user",
			})
		}

		var user User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to unmarshal user",
			})
		}

		if !compareSHA1(req.Password, user.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid username or password",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"username": user.Email,
			"realname": user.Realname,
			"email":    user.Email,
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func newRedisClientEnv() *redis.Client {
	redisURL := "localhost:6379"
	if redisURL == "" {
		panic("REDIS_URL is not set")
	}
	return redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})
}

func compareSHA1(password, hash string) bool {
	h := sha1.New()
	h.Write([]byte(password))
	hashed := hex.EncodeToString(h.Sum(nil))
	return hashed == hash
}
