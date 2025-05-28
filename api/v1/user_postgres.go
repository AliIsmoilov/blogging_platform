package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"blogging_platform/api/models"
	"blogging_platform/storage/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *handlerV1) CreateUser(ctx *gin.Context) {
	var req models.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "error while binding data",
		})
		return
	}

	data, err := h.strg.Postgres().User().Create(ctx, &repo.CreateUserReq{
		Email:       req.Email,
		Password:    req.Password,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error while creating user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, parseUserRepoToApi(data))
}

func (h *handlerV1) UpdateUser(ctx *gin.Context) {
	var req models.UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "error while binding data",
		})
		return
	}

	data, err := h.strg.Postgres().User().Update(ctx, &repo.UpdateUserReq{
		Id:          req.Id,
		PhoneNumber: req.PhoneNumber,
		FullName:    req.FullName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error while updating user",
		})
		return
	}

	ctx.JSON(http.StatusOK, parseUserRepoToApi(data))
}

func (h *handlerV1) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "strconv atoi not worked",
		})
		return
	}
	data, err := h.strg.Postgres().User().GetById(ctx, int64(idInt))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "user not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "user get by id not worked",
		})
		return
	}
	ctx.JSON(http.StatusOK, parseUserRepoToApi(data))
}

func (h *handlerV1) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "strconv atoi not worked",
		})
		return
	}
	err = h.strg.Postgres().User().Delete(ctx, idInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "user not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "user get by id not worked",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func (h *handlerV1) GetAllUsers(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "")
	offset := ctx.DefaultQuery("offset", "")
	query := ctx.DefaultQuery("query", "")

	data, err := h.strg.Postgres().User().GetAll(ctx, &repo.GetAllUserReq{
		Limit:  limit,
		Offset: offset,
		Query:  query,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error while getting all users PostgreSQL",
		})
		return
	}

	resp := models.GetAllUsersResp{
		Count: data.Count,
	}

	for _, elem := range data.Users {
		user := parseUserRepoToApi(elem)
		resp.Users = append(resp.Users, &user)
	}

	ctx.JSON(http.StatusOK, resp)
}

func parseUserRepoToApi(user *repo.UserModelResp) models.UserModelResp {
	resp := models.UserModelResp{
		Id:        user.Id,
		Email:     user.Email,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt.Format(time.RFC1123Z),
	}
	if user.FullName.Valid {
		resp.FullName = &user.FullName.String
	}
	if user.PhoneNumber.Valid {
		resp.PhoneNumber = &user.PhoneNumber.String
	}
	if user.UpdatedAt.Valid {
		tm := user.UpdatedAt.Time.Format(time.RFC1123Z)
		resp.UpdatedAt = &tm
	}
	return resp
}

func parseUserMongoToApi(user *repo.UserModelRespMongo) models.UserModelResp {
	updatedAt := user.UpdatedAt.Format(time.RFC1123Z)

	resp := models.UserModelResp{
		Id:          user.Id,
		Email:       user.Email,
		Balance:     user.Balance,
		FullName:    &user.FullName,
		PhoneNumber: &user.PhoneNumber,
		CreatedAt:   user.CreatedAt.Format(time.RFC1123Z),
		UpdatedAt:   &updatedAt,
	}

	return resp
}
