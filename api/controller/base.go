package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rekaime/r-mio/api/response"
	"io"
	"net/http"
	"time"
)

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, response.Success(data))
}

func Error(ctx *gin.Context, code response.ErrorCode) {
	ctx.JSON(http.StatusOK, response.Error(code))
}

func InternalError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, nil)
}

func Data(ctx *gin.Context, content any) {
	switch v := content.(type) {
	case []byte:
		if len(v) == 0 {
			ctx.Status(http.StatusNoContent)
			return
		}
		contentType := http.DetectContentType(v)
		ctx.Data(http.StatusOK, contentType, v)

	default:
		Stream(ctx, v)
	}
}

func Stream(ctx *gin.Context, content any) {
	switch v := content.(type) {
	case []byte:
		reader := bytes.NewReader(v)
		http.ServeContent(ctx.Writer, ctx.Request, "", time.Now(), reader)

	case io.ReadSeeker:
		http.ServeContent(ctx.Writer, ctx.Request, "", time.Now(), v)

	case io.Reader:
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Status(http.StatusOK)
		_, _ = io.Copy(ctx.Writer, v)

	default:
		ctx.String(http.StatusInternalServerError, "unsupported content type")
	}
}
