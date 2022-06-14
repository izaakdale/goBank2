package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/izaakdale/goBank2/token"
	"github.com/izaakdale/goBank2/util"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func (server *Server) login(ctx *gin.Context) {

	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get user
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// validate password encryption
	err = util.VerifyPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// placeholder key for dev
	key := util.RandomString(32)
	// return paseto with claims
	maker, err := token.NewPasetoMaker(key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// placeholer duration for dev
	tokenTimeout := time.Minute
	paseto, err := maker.CreateToken(req.Username, tokenTimeout)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginResponse{
		paseto,
	}
	ctx.JSON(http.StatusOK, rsp)
}
