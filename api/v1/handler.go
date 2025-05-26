package v1

import (
	"github.com/AliIsmoilov/blogging_platform/config"
	"github.com/AliIsmoilov/blogging_platform/storage"
)

type handlerV1 struct {
	cfg  *config.Config
	strg storage.StorageI
}

type HandleV1 struct {
	Cfg  *config.Config
	Strg storage.StorageI
}

func New(h *HandleV1) *handlerV1 {
	return &handlerV1{
		cfg:  h.Cfg,
		strg: h.Strg,
	}
}
