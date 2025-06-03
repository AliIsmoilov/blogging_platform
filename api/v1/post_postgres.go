package v1

import (
	"net/http"
	"time"

	"blogging_platform/api/models"
	"blogging_platform/storage/repo"

	"github.com/gin-gonic/gin"
)

func (h *handlerV1) CreatePost(ctx *gin.Context) {
	var req models.CreatePostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "error while binding data",
		})
		return
	}

	data, err := h.strg.Postgres().Post().Create(ctx, &repo.CreatePostReq{
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

	ctx.JSON(http.StatusCreated, parsePostRepoToApi(data))
}

func (h *handlerV1) GetAllPosts(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "")
	offset := ctx.DefaultQuery("offset", "")
	query := ctx.DefaultQuery("query", "")

	data, err := h.strg.Postgres().Post().GetAll(ctx, &repo.GetAllUserReq{
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

func parsePostRepoToApi(user *repo.PostModelResp) models.PostModelResp {
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
