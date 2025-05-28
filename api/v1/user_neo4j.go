package v1

import (
	"blogging_platform/api/models"
	"blogging_platform/storage/repo"
	"fmt"
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

	data, err := h.strg.Neo4j().User().Create(ctx, &repo.UserModelRespMongo{
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

func (h *handlerV1) GetAllUsersNeo4j(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "")
	offset := ctx.DefaultQuery("offset", "")
	query := ctx.DefaultQuery("query", "")

	fmt.Println("HEYYYYYY")
	data, err := h.strg.Neo4j().User().GetAll(ctx, &repo.GetAllUserReq{
		Limit:  limit,
		Offset: offset,
		Query:  query,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error while getting all users Mongo",
		})
		return
	}

	resp := models.GetAllUsersResp{
		Count: int64(len(data)),
	}

	for _, elem := range data {
		user := parseUserMongoToApi(elem)
		resp.Users = append(resp.Users, &user)
	}

	ctx.JSON(http.StatusOK, resp)
}
