package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Define routes:
func SetupRecipeStepsRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/steps", getSteps)
	router.GET("/steps/:id", getStep)
	router.POST("/steps", createStep)
	router.PUT("/steps/:id", updateStep)
	router.DELETE("/steps/:id", deleteStep)
}

func getSteps(c *gin.Context) {

}

func getStep(c *gin.Context) {

}

func createStep(c *gin.Context) {

}

func updateStep(c *gin.Context) {

}

func deleteStep(c *gin.Context) {

}
