package controllers

import (
	"net/http"
	"strconv"

	"github.com/garvpatidar/food-delivery/backend/usecases"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentUsecase usecases.PaymentUsecase
}

func NewPaymentController(u usecases.PaymentUsecase) *PaymentController {
	return &PaymentController{paymentUsecase: u}
}

type InitiatePaymentInput struct {
	OrderID uint    `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Method  string  `json:"method" binding:"required"` // card, upi, netbanking, cod
}

type ConfirmPaymentInput struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

// POST /payments - Initiate a payment for an order
func (c *PaymentController) InitiatePayment(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var input InitiatePaymentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := c.paymentUsecase.InitiatePayment(ctx.Request.Context(), input.OrderID, userID, input.Amount, input.Method)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Payment initiated",
		"payment": payment,
	})
}

// PUT /payments/:id/confirm - Confirm payment (webhook simulation)
func (c *PaymentController) ConfirmPayment(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var input ConfirmPaymentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.paymentUsecase.ConfirmPayment(ctx.Request.Context(), uint(id), input.TransactionID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment confirmed successfully"})
}

// GET /payments/order/:orderID - Get payment details for an order
func (c *PaymentController) GetPaymentByOrder(ctx *gin.Context) {
	orderIDParam := ctx.Param("orderID")
	orderID, err := strconv.ParseUint(orderIDParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	payment, err := c.paymentUsecase.GetPaymentByOrder(ctx.Request.Context(), uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
