package controllers

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Abhijeet6387/Blog/database"
	"github.com/Abhijeet6387/Blog/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)
var secretKey string

func init() {
	// Load .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
	secretKey = os.Getenv("SECRETKEY")
}

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome Home")
}

// Register User
func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil{
		return err
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Name: data["name"],
		Email: data["email"],
		Password: string(hashPassword),
	}
	database.DB.Create(&user)
	return c.JSON(user)
}

// Login User
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil{
		return err
	} 
	// Initialize user
	var user models.User

	// Query in Database
	database.DB.Where("email = ?", data["email"]).First(&user)
	
	// If email is not found
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message":"User not found",
		})
	}
	
	// else compare the hashed password from db and password from body
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	
	// if not matched
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message":"Incorrect password",
		})
	}
	// generate jwt token with secret key	
	claimsMap := jwt.MapClaims{
		"Issuer": strconv.Itoa(int(user.Id)),
		"Email": user.Email,
		"ExpiresAt":time.Now().Add(time.Hour * 24).Unix(),
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsMap)
	token, err := claims.SignedString([]byte(secretKey))
	
	// if error occurs in generating token
	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message":"Unable to login",
		})
	}
	
	// save the token information in cookie, retrieved from front-end
	cookie := fiber.Cookie{
		Name : "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour *24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"message":"Logged In",
		// "data" :user,
		// "token":token,
	})
}

// Get user details
func GetUser(c *fiber.Ctx) error {
	
	// Retrieve the JWT token from the cookie
	cookie := c.Cookies("jwt")
	if cookie == "" {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	
	// Parse and validate the JWT token
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	// Check if the token is valid
	if !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	// Extract the claims from the token
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}
	
	// The user is authorized, and you can access the claims from the token
	user_id:= (*claims)["Issuer"].(string)
	email := (*claims)["Email"].(string)

	// store user information
	var user models.User
	database.DB.Where("Id = ? And email = ?", user_id, email).First(&user)
	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
        return c.JSON(fiber.Map{
            "message": "User not found",
        })
	}
	// return  user information
	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"message":"Authorized",
		"data":user,
	})
}

// Logout user
func Logout(c *fiber.Ctx) error {	
	// clear the cookie by setting an expire time in the past
    cookie := fiber.Cookie{
        Name: "jwt",
        Value: "",
        Expires: time.Now().Add(-time.Hour), // expire the cookie in the past
        HTTPOnly: true,
    }
    c.Cookie(&cookie)

    return c.JSON(fiber.Map{
        "message": "Logged out",
    })
}