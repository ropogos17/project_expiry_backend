package api

//
//import (
//	"database/sql"
//	"net/http"
//	db "project_expiry_backend/db/sqlc"
//
//	"github.com/gin-gonic/gin"
//)
//
//type createShareRequest struct {
//	StorageID    int32 `json:"storage_id" binding:"required,min=1,max=1000000"`
//	TargetUserID int32 `json:"target_user_id" binding:"required,min=1,max=1000000"`
//	CanShare     bool  `json:"can_share" binding:"required"`
//}
//
//func (server *Server) createShare(c *gin.Context) {
//	var req createShareRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//
//	authUser, status, err := server.getAuthenticatedUser(c)
//	if err != nil {
//		c.JSON(status, errorResponse(err))
//		return
//	}
//
//	arg := db.CreateShareParams{
//		SourceUserID: authUser.ID,
//		StorageID:    req.StorageID,
//		TargetUserID: req.TargetUserID,
//		CanShare:     req.CanShare,
//	}
//	share, err := server.store.CreateShare(c.Request.Context(), arg)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//	c.JSON(http.StatusCreated, share)
//}
//
//type getShareRequest struct {
//	ID int32 `uri:"id" binding:"required,min=1,max=1000000"`
//}
//
//func (server *Server) getShare(c *gin.Context) {
//	var req getShareRequest
//	if err := c.ShouldBindUri(&req); err != nil {
//		c.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//
//	/*	authUser, status, err := server.getAuthenticatedUser(c)
//		if err != nil {
//			c.JSON(status, errorResponse(err))
//			return
//		}*/
//
//	share, err := server.store.GetShareByID(c.Request.Context(), req.ID)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			c.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}
//		c.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//	c.JSON(http.StatusOK, share)
//}
//
//type listShareRequest struct {
//	PageID   int32 `form:"page_id" binding:"required,min=1,max=1000000"`
//	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
//}
//
//func (server *Server) listShare(c *gin.Context) {
//	var req listShareRequest
//	if err := c.ShouldBindQuery(&req); err != nil {
//		c.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//
//	arg := db.ListSharesParams{
//		Limit:  req.PageSize,
//		Offset: (req.PageID - 1) * req.PageSize,
//	}
//
//	shares, err := server.store.ListShares(c.Request.Context(), arg)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//	if shares == nil {
//		shares = []db.Share{}
//	}
//	c.JSON(http.StatusOK, shares)
//}
//
//type ListSharesByUserRequest struct {
//	Limit        int32 `form:"limit" binding:"required,min=1,max=1000000"`
//	Offset       int32 `form:"offset" binding:"required,min=1,max=100"`
//	SourceUserID int32 `form:"source_user_id" binding:"required,min=1,max=1000000"`
//}
//
//func (server *Server) ListSharesByUser(c *gin.Context) {
//	var req ListSharesByUserRequest
//	if err := c.ShouldBindQuery(&req); err != nil {
//		c.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//
//}
