package controllers

import (
	"net/http"
	"strconv"

	"github.com/garvpatidar/food-delivery/backend/models"
	"github.com/garvpatidar/food-delivery/backend/usecases"
	"github.com/gin-gonic/gin"
)

type CartController struct {
	cartUsecase usecases.CartUsecase
}

func NewCartController(u usecases.CartUsecase) *CartController {
	return &CartController{cartUsecase: u}
}

// GET /cart - User ki cart dekho
func (c *CartController) GetCart(ctx *gin.Context) {
	userID := ctx.GetUint("userID") // JWT middleware se milta hai

	cart, err := c.cartUsecase.GetCart(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// POST /cart - Cart mein item add karo
func (c *CartController) AddItem(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var item models.CartItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	cart, err := c.cartUsecase.AddItem(ctx.Request.Context(), userID, item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// DELETE /cart/:foodItemId - Cart se item hatao
func (c *CartController) RemoveItem(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	foodItemID, err := strconv.ParseUint(ctx.Param("foodItemId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food item ID"})
		return
	}

	cart, err := c.cartUsecase.RemoveItem(ctx.Request.Context(), userID, uint(foodItemID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item"})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// DELETE /cart - Poori cart clear karo
func (c *CartController) ClearCart(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	if err := c.cartUsecase.ClearCart(ctx.Request.Context(), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
