package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/api/response"
    "net/http"
)

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, response.Success(data))
}

func Error(ctx *gin.Context, code response.ErrorCode) {
	ctx.JSON(http.StatusOK, response.Error(code))
}

func InternalError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, nil)
}
