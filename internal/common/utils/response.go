package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
}

func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		}

		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}
func ResponseSuccess(ctx *gin.Context, statusCode int, message string, data ...any) {
	resp := APIResponse{
		Status:  "success",
		Message: message,
	}
	if len(data) > 0 && data[0] != nil {
		if items, ok := data[0].(map[string]any); ok {
			if pagination, exist := items["pagination"]; exist {
				resp.Pagination = pagination
			}

			if data, exist := items["data"]; exist {
				resp.Data = data
			} else {
				resp.Data = items
			}
		} else {
			resp.Data = data[0]
		}
	}
	ctx.JSON(statusCode, resp)
}
func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.JSON(status, gin.H{})
}
func ResponseWValidator(ctx *gin.Context, handleValidation any) {
	ctx.JSON(http.StatusBadRequest, handleValidation)
}
