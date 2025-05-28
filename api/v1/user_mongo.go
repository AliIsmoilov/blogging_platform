package v1

import (
	"blogging_platform/api/models"
	"blogging_platform/storage/repo"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handlerV1) GetAllUsersMongo(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "")
	offset := ctx.DefaultQuery("offset", "")
	query := ctx.DefaultQuery("query", "")

	dataMongo, err := h.strg.MongoStorage().User().GetAll(ctx, &repo.GetAllUserReq{
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
		Count: dataMongo.Count,
	}

	for _, elem := range dataMongo.UsersMongo {
		fmt.Println("Mongo User: ", elem)
		user := parseUserMongoToApi(elem)
		resp.Users = append(resp.Users, &user)
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *handlerV1) CreateUserMongo(ctx *gin.Context) {
	var req models.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "error while binding data",
		})
		return
	}

	data, err := h.strg.MongoStorage().User().Create(ctx, &repo.UserModelRespMongo{
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
