package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	db "project_expiry_backend/db/sqlc"
	"project_expiry_backend/token"
)

type createItemRequest struct {
	Name        string    `json:"name" binding:"required"`
	CategoryID  int32     `json:"category_id" binding:"required,min=1"`
	Quantity    *string   `json:"quantity"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
	StoredAt    time.Time `json:"stored_at" binding:"required"`
	Description *string   `json:"description"`
	StorageID   int32     `json:"storage_id" binding:"required,min=1"`
	Shelf       *int32    `json:"shelf"`
	HasImage    bool      `json:"has_image"`
}

type createBulkItemsRequest struct {
	Items []createItemRequest `json:"items" binding:"required,min=1,max=100"`
}

type bulkItemResponse struct {
	Success      bool      `json:"success"`
	CreatedItems []db.Item `json:"created_items,omitempty"`
	Errors       []string  `json:"errors,omitempty"`
	Message      string    `json:"message"`
	CreatedCount int       `json:"created_count"`
	TotalCount   int       `json:"total_count"`
}

func (server *Server) createItem(c *gin.Context) {
	var req createItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify user owns the storage or has access to it
	if !server.hasStorageAccess(c, req.StorageID, user.ID) {
		return // Error response already sent
	}

	item, err := server.createSingleItem(c, req, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, item)
}

func (server *Server) createBulkItems(c *gin.Context) {
	var req createBulkItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Group items by storage for access verification
	storageMap := make(map[int32][]createItemRequest)
	for _, item := range req.Items {
		storageMap[item.StorageID] = append(storageMap[item.StorageID], item)
	}

	// Verify access to all storages first
	for storageID := range storageMap {
		if !server.hasStorageAccessNonResponse(c, storageID, user.ID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "access denied to one or more storages",
			})
			return
		}
	}

	var createdItems []db.Item
	var errorStrings []string
	successCount := 0

	// Create items one by one (since we can't use copyfrom)
	for i, itemReq := range req.Items {
		item, err := server.createSingleItem(c, itemReq, user.ID)
		if err != nil {
			errorStrings = append(errorStrings, "Item "+string(rune(i+1))+": "+err.Error())
			continue
		}
		createdItems = append(createdItems, item)
		successCount++
	}

	response := bulkItemResponse{
		Success:      successCount > 0,
		CreatedItems: createdItems,
		Errors:       errorStrings,
		CreatedCount: successCount,
		TotalCount:   len(req.Items),
	}

	if successCount == len(req.Items) {
		response.Message = "All items created successfully"
		c.JSON(http.StatusOK, response)
	} else if successCount > 0 {
		response.Message = "Some items created successfully with errorStrings"
		c.JSON(http.StatusPartialContent, response)
	} else {
		response.Message = "Failed to create any items"
		c.JSON(http.StatusBadRequest, response)
	}
}

// Helper function to create a single item
func (server *Server) createSingleItem(c *gin.Context, req createItemRequest, userID int32) (db.Item, error) {
	arg := db.CreateItemParams{
		Name:         req.Name,
		CategoryID:   req.CategoryID,
		Description:  sql.NullString{String: "", Valid: false},
		Quantity:     sql.NullString{String: "", Valid: false},
		HasImage:     req.HasImage,
		ExpiresAt:    req.ExpiresAt,
		StorageID:    req.StorageID,
		Shelf:        sql.NullInt32{Int32: 0, Valid: false},
		LastModified: userID,
	}

	if req.Description != nil {
		arg.Description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Quantity != nil {
		arg.Quantity = sql.NullString{String: *req.Quantity, Valid: true}
	}
	if req.Shelf != nil {
		arg.Shelf = sql.NullInt32{Int32: *req.Shelf, Valid: true}
	}

	return server.store.CreateItem(c, arg)
}

// Helper function to check storage access (with response)
func (server *Server) hasStorageAccess(c *gin.Context, storageID, userID int32) bool {
	storage, err := server.store.GetStorage(c, storageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "storage not found"})
			return false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if storage.Owner != userID {
		// Check if user has shared access
		share, err := server.store.GetShareByStorageAndOwner(c, db.GetShareByStorageAndOwnerParams{
			StorageID:    storageID,
			TargetUserID: userID,
		})
		if err != nil || share.ID == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this storage"})
			return false
		}
	}
	return true
}

// Helper function to check storage access (without response)
func (server *Server) hasStorageAccessNonResponse(c *gin.Context, storageID, userID int32) bool {
	storage, err := server.store.GetStorage(c, storageID)
	if err != nil {
		return false
	}

	if storage.Owner != userID {
		// Check if user has shared access
		share, err := server.store.GetShareByStorageAndOwner(c, db.GetShareByStorageAndOwnerParams{
			StorageID:    storageID,
			TargetUserID: userID,
		})
		if err != nil || share.ID == 0 {
			return false
		}
	}
	return true
}

type getItemRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getItem(c *gin.Context) {
	var req getItemRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	item, err := server.store.GetItem(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify user owns the storage or has access to it
	if !server.hasStorageAccess(c, item.StorageID, user.ID) {
		return // Error response already sent
	}

	c.JSON(http.StatusOK, item)
}

type listItemsRequest struct {
	PageID    int32  `json:"page_id" binding:"required,min=1"`
	PageSize  int32  `json:"page_size" binding:"required,min=5,max=50"`
	StorageID *int32 `json:"storage_id"`
}

func (server *Server) listItems(c *gin.Context) {
	var req listItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.StorageID != nil {
		// Verify user owns the storage or has access to it
		if !server.hasStorageAccess(c, *req.StorageID, user.ID) {
			return // Error response already sent
		}

		// List items for specific storage
		arg := db.ListItemsByStorageParams{
			StorageID: *req.StorageID,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}

		items, err := server.store.ListItemsByStorage(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, items)
	} else {
		// List all items for user's storages
		arg := db.ListItemsByUserParams{
			Owner:  user.ID,
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}

		items, err := server.store.ListItemsByUser(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

type listActiveItemsRequest struct {
	PageID    int32  `json:"page_id" binding:"required,min=1"`
	PageSize  int32  `json:"page_size" binding:"required,min=5,max=50"`
	StorageID *int32 `json:"storage_id"`
}

func (server *Server) listActiveItems(c *gin.Context) {
	var req listActiveItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	currentTime := time.Now().UTC()

	if req.StorageID != nil {
		// Verify user owns the storage or has access to it
		if !server.hasStorageAccess(c, *req.StorageID, user.ID) {
			return // Error response already sent
		}

		// List active items for specific storage
		arg := db.ListActiveItemsByStorageWithTimeParams{
			StorageID: *req.StorageID,
			ExpiresAt: currentTime,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}

		items, err := server.store.ListActiveItemsByStorageWithTime(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, items)
	} else {
		// List all active items for user's storages
		arg := db.ListActiveItemsByUserWithTimeParams{
			Owner:     user.ID,
			ExpiresAt: currentTime,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}

		items, err := server.store.ListActiveItemsByUserWithTime(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func (server *Server) listExpiredItems(c *gin.Context) {
	var req listActiveItemsRequest // Reusing the same request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	currentTime := time.Now().UTC()

	if req.StorageID != nil {
		// Verify user owns the storage or has access to it
		if !server.hasStorageAccess(c, *req.StorageID, user.ID) {
			return // Error response already sent
		}

		// List expired items for specific storage
		arg := db.ListExpiredItemsByStorageWithTimeParams{
			StorageID: *req.StorageID,
			ExpiresAt: currentTime,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}

		items, err := server.store.ListExpiredItemsByStorageWithTime(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, items)
	} else {
		// List all expired items for user's storages
		arg := db.ListExpiredItemsByUserWithTimeParams{
			Owner:     user.ID,
			ExpiresAt: currentTime,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}

		items, err := server.store.ListExpiredItemsByUserWithTime(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func (server *Server) listDeletedItems(c *gin.Context) {
	var req listActiveItemsRequest // Reusing the same request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.StorageID != nil {
		// Verify user owns the storage or has access to it
		if !server.hasStorageAccess(c, *req.StorageID, user.ID) {
			return // Error response already sent
		}

		// List expired items for specific storage
		arg := db.ListDeletedItemsByStorageWithTimeParams{
			StorageID: *req.StorageID,
			Limit:     req.PageSize,
			Offset:    (req.PageID - 1) * req.PageSize,
		}

		items, err := server.store.ListDeletedItemsByStorageWithTime(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, items)
	} else {
		// List all expired items for user's storages
		arg := db.ListDeletedItemsByUserWithTimeParams{
			Owner:  user.ID,
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}

		items, err := server.store.ListDeletedItemsByUserWithTime(c, arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

type updateItemUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type updateItemBodyRequest struct {
	Name        *string    `json:"name"`
	CategoryID  *int32     `json:"category_id"`
	Description *string    `json:"description"`
	Quantity    *string    `json:"quantity"`
	HasImage    *bool      `json:"has_image"`
	ExpiresAt   *time.Time `json:"expires_at"`
	StorageID   *int32     `json:"storage_id"`
	Shelf       *int32     `json:"shelf"`
}

func (server *Server) updateItem(c *gin.Context) {
	var uriReq updateItemUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateItemBodyRequest
	if err := c.ShouldBindJSON(&bodyReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get current item
	currentItem, err := server.store.GetItem(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify user owns the storage or has access to it
	if !server.hasStorageAccess(c, currentItem.StorageID, user.ID) {
		return // Error response already sent
	}

	// Build update parameters with current values as defaults
	arg := db.UpdateItemParams{
		ID:           uriReq.ID,
		Name:         currentItem.Name,
		CategoryID:   currentItem.CategoryID,
		Description:  currentItem.Description,
		Quantity:     currentItem.Quantity,
		HasImage:     currentItem.HasImage,
		IsDeleted:    currentItem.IsDeleted,
		StoredAt:     currentItem.StoredAt,
		ExpiresAt:    currentItem.ExpiresAt,
		StorageID:    currentItem.StorageID,
		Shelf:        currentItem.Shelf,
		LastModified: user.ID,
		UpdatedAt:    time.Now(),
	}

	// Update only provided fields
	if bodyReq.Name != nil {
		arg.Name = *bodyReq.Name
	}
	if bodyReq.CategoryID != nil {
		arg.CategoryID = *bodyReq.CategoryID
	}
	if bodyReq.Description != nil {
		arg.Description = sql.NullString{String: *bodyReq.Description, Valid: true}
	}
	if bodyReq.Quantity != nil {
		arg.Quantity = sql.NullString{String: *bodyReq.Quantity, Valid: true}
	}
	if bodyReq.HasImage != nil {
		arg.HasImage = *bodyReq.HasImage
	}
	if bodyReq.ExpiresAt != nil {
		arg.ExpiresAt = *bodyReq.ExpiresAt
	}
	if bodyReq.StorageID != nil {
		// Verify user has access to new storage
		if !server.hasStorageAccess(c, *bodyReq.StorageID, user.ID) {
			return // Error response already sent
		}
		arg.StorageID = *bodyReq.StorageID
	}
	if bodyReq.Shelf != nil {
		arg.Shelf = sql.NullInt32{Int32: *bodyReq.Shelf, Valid: true}
	}

	item, err := server.store.UpdateItem(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, item)
}

type deleteItemRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteItem(c *gin.Context) {
	var req deleteItemRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get current item
	currentItem, err := server.store.GetItem(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Get authenticated user
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(c, authPayload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify user owns the storage or has access to it
	if !server.hasStorageAccess(c, currentItem.StorageID, user.ID) {
		return // Error response already sent
	}

	if err := server.store.DeleteItem(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})
}
