package main

import (
	"backend/db"
	_ "backend/docs"
	"backend/routes"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Connect to the database
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.Close()

	// Serve Swagger UI files
	router.Static("/docs", "./docs")
	url := ginSwagger.URL("/docs/swagger.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// CORS for https://foo.com and https://github.com origins, allowing:
	// - GET and POST methods
	// - "Authorization" and "Content-Type" headers
	// - Credentials share
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com", "https://github.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	routes.SetupIngredientsRoutes(router, database)
	routes.SetupRecipeRoutes(router, database)
	routes.SetupRecipeIngredientsRoutes(router, database)
	routes.SetupRecipeStepsRoutes(router, database)

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	fmt.Println("Server is running on http://localhost:8080")
}
