package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kikudesuyo/arumonogohan-app/api/usecase"
)

func HandleSuggestRecipe(c *gin.Context) {
	var input usecase.RecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe, err := usecase.SuggestRecipe(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipe": recipe})
}