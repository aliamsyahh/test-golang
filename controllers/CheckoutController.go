package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckoutRequest defines the structure for the checkout request
type CheckoutRequest struct {
	VoucherCode string  `json:"voucher_code" binding:"required"`
	TotalPrice  float64 `json:"total_price" binding:"required"`
}

// CheckoutController handles the checkout logic
func CheckoutController(c *gin.Context) {
	var req CheckoutRequest

	// Bind the request JSON data to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Static data for the voucher
	voucherDiscount := 0.50      // 50% discount
	voucherPointsPercent := 0.02 // 2% of the voucher value

	// Static product price (can be dynamic in a real app)
	productPrice := 5000000.0

	// Check if the provided voucher code is valid
	if req.VoucherCode != "50OFF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voucher code"})
		return
	}

	// Calculate the discount amount and the final price after discount
	discountAmount := productPrice * voucherDiscount
	finalPrice := productPrice - discountAmount

	// Calculate reward points (2% of the voucher discount applied)
	rewardPoints := discountAmount * voucherPointsPercent

	// Prepare the response
	response := gin.H{
		"original_price":   productPrice,
		"voucher_discount": fmt.Sprintf("Rp%.2f", discountAmount),
		"final_price":      fmt.Sprintf("Rp%.2f", finalPrice),
		"reward_points":    fmt.Sprintf("%.2f points", rewardPoints),
	}

	// Send the response back to the client
	c.JSON(http.StatusOK, response)
}
