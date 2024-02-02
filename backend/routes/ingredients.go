package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Define routes:
func SetupIngredientsRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/ingredients", getIngredients)
	router.GET("/ingredients/:id", getIngredient)
	router.POST("/ingredients", createIngredient)
	router.PUT("/ingredients/:id", updateIngredient)
	router.DELETE("/ingredients/:id", deleteIngredient)
}

func getIngredients(c *gin.Context) {

}

func getIngredient(c *gin.Context) {

}

func createIngredient(c *gin.Context) {

}

func updateIngredient(c *gin.Context) {

}

func deleteIngredient(c *gin.Context) {

}
