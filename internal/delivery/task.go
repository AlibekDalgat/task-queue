package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"task-queue/internal/models"
)

func (h *Handler) createTask(c *gin.Context) {
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Неверное содержание json")
		return
	}
	id, err := h.services.Create(input.InputData)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, createResponse{TaskID: id})
}

func (h *Handler) getTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Неверный парамметр id")
		return
	}
	task, err := h.services.Get(uint32(id))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getResponse{
		Status: task.Status,
		Result: task.Result,
	})
}
