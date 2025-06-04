package v1

import (
	"blogging_platform/api/models"
	"blogging_platform/storage/repo"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *handlerV1) CreatePostMongo(ctx *gin.Context) {
	var req models.CreatePostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "error while binding data",
		})
		return
	}

	data, err := h.strg.MongoStorage().Post().Create(ctx, &repo.PostModelRespMongo{
		UserId:  req.UserId,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error while creating post",
		})
		return
	}

	ctx.JSON(http.StatusCreated, parseMongoPostRepoToApi(data))
}

func (h *handlerV1) GetAllPostsMongo(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "")
	offset := ctx.DefaultQuery("offset", "")
	query := ctx.DefaultQuery("query", "")

	data, err := h.strg.MongoStorage().Post().GetAll(ctx, &repo.GetAllUserReq{
		Limit:  limit,
		Offset: offset,
		Query:  query,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error while getting all posts PostgreSQL",
		})
		return
	}

	resp := models.GetAllPostsResp{
		Count: data.Count,
	}

	for _, elem := range data.Posts {
		post := parsePostRepoToApi(elem)
		resp.Posts = append(resp.Posts, &post)
	}

	ctx.JSON(http.StatusOK, resp)
}

func parseMongoPostRepoToApi(user *repo.PostModelRespMongo) models.PostModelResp {
	resp := models.PostModelResp{
		Id:        user.Id,
		UserId:    user.UserId,
		Title:     user.Title,
		Content:   user.Content,
		CreatedAt: user.CreatedAt.Format(time.RFC1123Z),
	}

	if !user.UpdatedAt.IsZero() {
		tm := user.UpdatedAt.Format(time.RFC1123Z)
		resp.UpdatedAt = &tm
	}
	return resp
}
