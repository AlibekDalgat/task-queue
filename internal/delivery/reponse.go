package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

type createResponse struct {
	TaskID uint32 `json:"task_id"`
}

type getResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
}
