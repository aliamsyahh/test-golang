package routes

import (
	"test-golang/controllers"
	"test-golang/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.POST("/login", controllers.Login)
	router.GET("/protected", controllers.Protected)

	// Auth Routes
	router.GET("/auth/google/login", controllers.GoogleLogin)
	router.GET("/auth/google/callback", controllers.GoogleCallback)
	router.GET("/auth/logout", controllers.Logout)

	// Protected Route (only accessible with JWT)
	protected := router.Group("/user")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		c.JSON(200, gin.H{"id": userID, "username": username})
	})

	// CRUD Routes for Users
	router.GET("/users", controllers.GetUsers)
	router.GET("/users/:id", controllers.GetUser)
	router.POST("/users", controllers.CreateUser)
	router.PUT("/users/:id", controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)

	// GET Route for random user
	router.GET("/fetchuser", controllers.FetchUser)
	// POST route for the checkout
	router.POST("/checkout", controllers.CheckoutController)
}
