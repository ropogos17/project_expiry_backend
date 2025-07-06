package DTO

// ShoppingListItemCreateRequest represents the request body for creating a shopping list item
type ShoppingListItemCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Quantity    *string `json:"quantity"`
}

// ShoppingListItemUpdateRequest represents the request body for updating a shopping list item
type ShoppingListItemUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Quantity    *string `json:"quantity"`
	IsChecked   *bool   `json:"is_checked"`
}
