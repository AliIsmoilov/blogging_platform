package v1

import (
	"blogging_platform/api/models"
	"blogging_platform/storage/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handlerV1) CreateUserNeo4j(ctx *gin.Context) {
	var req models.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "error while binding data",
		})
		return
	}

	data, err := h.strg.Neo4j().User().CreateUserNeo4j(ctx, &repo.UserModelRespMongo{
		FullName:    *req.FullName,
		Email:       req.Email,
		PhoneNumber: *req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error while creating user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, parseUserMongoToApi(data))
}
