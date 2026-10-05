package models

// CartItem represents a single item in the shopping cart
type CartItem struct {
	FoodItemID   uint    `json:"food_item_id"`
	RestaurantID uint    `json:"restaurant_id"` // Ek cart mein items kis restaurant se hain
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
}

// Cart represents the entire shopping cart for a user
type Cart struct {
	UserID      uint       `json:"user_id"`
	Items       []CartItem `json:"items"`
	TotalAmount float64    `json:"total_amount"`
}
