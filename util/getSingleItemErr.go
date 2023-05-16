package util

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckError(ctx *gin.Context, err error) bool {
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return false
	}

	return true
}
