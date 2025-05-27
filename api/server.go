package api

import (
	v1 "github.com/AliIsmoilov/blogging_platform/api/v1"
	"github.com/AliIsmoilov/blogging_platform/config"
	"github.com/AliIsmoilov/blogging_platform/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Cfg  *config.Config
	Strg storage.StorageI
}

func New(h *Handler) *gin.Engine {
	engine := gin.Default()

	handlerV1 := v1.New(&v1.HandleV1{
		Cfg:  h.Cfg,
		Strg: h.Strg,
	})

	apiV1 := engine.Group("/v1")
	apiV1.GET("/user/:id", handlerV1.GetUserById)
	apiV1.POST("/user", handlerV1.CreateUser)
	apiV1.PUT("/user", handlerV1.UpdateUser)
	apiV1.DELETE("/user/:id", handlerV1.DeleteUser)
	apiV1.GET("/users", handlerV1.GetAllUsers)

	apiV1.POST("/user/mongo", handlerV1.CreateUserMongo)
	apiV1.GET("/users/mongo", handlerV1.GetAllUsersMongo)

	return engine
}
