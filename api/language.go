package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	db "project_expiry_backend/db/sqlc"
)

type createLanguageRequest struct {
	En string `json:"en" binding:"required"`
	Hu string `json:"hu" binding:"required"`
	Ro string `json:"ro" binding:"required"`
	Es string `json:"es" binding:"required"`
	Zh string `json:"zh" binding:"required"`
	Fr string `json:"fr" binding:"required"`
	De string `json:"de" binding:"required"`
	Pt string `json:"pt" binding:"required"`
	Ru string `json:"ru" binding:"required"`
	Ua string `json:"ua" binding:"required"`
	Ar string `json:"ar" binding:"required"`
}

func (server *Server) createLanguage(c *gin.Context) {
	var req createLanguageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateLanguageParams{
		En: req.En,
		Hu: req.Hu,
		Ro: req.Ro,
		Es: req.Es,
		Zh: req.Zh,
		Fr: req.Fr,
		De: req.De,
		Pt: req.Pt,
		Ru: req.Ru,
		Ua: req.Ua,
		Ar: req.Ar,
	}

	lang, err := server.store.CreateLanguage(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, lang)
}

type getLanguageRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getLanguage(c *gin.Context) {
	var req getLanguageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lang, err := server.store.GetLanguage(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, lang)
}

type listLanguagesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

func (server *Server) listLanguages(c *gin.Context) {
	var req listLanguagesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListLanguagesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	langs, err := server.store.ListLanguages(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, langs)
}

type deleteLanguageRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteLanguage(c *gin.Context) {
	var req deleteLanguageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteLanguage(c, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}

type updateLanguageUriRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}
type updateLanguageBodyRequest struct {
	En string `json:"en" binding:"required"`
	Hu string `json:"hu" binding:"required"`
	Ro string `json:"ro" binding:"required"`
	Es string `json:"es" binding:"required"`
	Zh string `json:"zh" binding:"required"`
	Fr string `json:"fr" binding:"required"`
	De string `json:"de" binding:"required"`
	Pt string `json:"pt" binding:"required"`
	Ru string `json:"ru" binding:"required"`
	Ua string `json:"ua" binding:"required"`
	Ar string `json:"ar" binding:"required"`
}

func (server *Server) updateLanguage(c *gin.Context) {
	var uriReq updateLanguageUriRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var bodyReq updateLanguageBodyRequest
	if err := c.ShouldBindJSON(&bodyReq); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	currentLang, err := server.store.GetLanguage(c, uriReq.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "language not found"})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	if currentLang.En == bodyReq.En &&
		currentLang.Hu == bodyReq.Hu &&
		currentLang.Ro == bodyReq.Ro &&
		currentLang.Es == bodyReq.Es &&
		currentLang.Zh == bodyReq.Zh &&
		currentLang.Fr == bodyReq.Fr &&
		currentLang.De == bodyReq.De &&
		currentLang.Pt == bodyReq.Pt &&
		currentLang.Ru == bodyReq.Ru &&
		currentLang.Ua == bodyReq.Ua &&
		currentLang.Ar == bodyReq.Ar {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update the language with an existing language"})
	}

	arg := db.UpdateLanguageParams{
		En: bodyReq.En,
		Hu: bodyReq.Hu,
		Ro: bodyReq.Ro,
		Es: bodyReq.Es,
		Zh: bodyReq.Zh,
		Fr: bodyReq.Fr,
		De: bodyReq.De,
		Pt: bodyReq.Pt,
		Ru: bodyReq.Ru,
		Ua: bodyReq.Ua,
		Ar: bodyReq.Ar,
	}

	updatedLang, err := server.store.UpdateLanguage(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, updatedLang)
}
