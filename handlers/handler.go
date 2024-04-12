package handlers

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"ssshekhu53/folder-lock/services"
)

type Handler interface {
	Init(ctx *gofr.Context) (interface{}, error)
	Lock(ctx *gofr.Context) (interface{}, error)
	Unlock(ctx *gofr.Context) (interface{}, error)
	Help(_ *gofr.Context) (interface{}, error)
}

type handler struct {
	service services.FolderLock
}

func New(service services.FolderLock) Handler {
	return &handler{service: service}
}

func (h *handler) Init(ctx *gofr.Context) (interface{}, error) {
	password := ctx.Param("password")
	if password == "" {
		return nil, errors.New("password is required")
	}

	err := h.service.Init(password)
	if err != nil {
		return nil, err
	}

	return "folder lock initialized", nil
}

func (h *handler) Lock(ctx *gofr.Context) (interface{}, error) {
	err := h.service.Lock()
	if err != nil {
		return nil, err
	}

	return "folder locked", nil
}

func (h *handler) Unlock(ctx *gofr.Context) (interface{}, error) {
	password := ctx.Param("password")
	if password == "" {
		return nil, errors.New("password is required")
	}

	err := h.service.Unlock(password)
	if err != nil {
		return nil, err
	}

	return "folder unlocked", nil
}

func (h *handler) Help(_ *gofr.Context) (interface{}, error) {
	return `Folder Lock CLI Tool

Usage:
  folder-lock [command]

Available Commands:
  init      Create a directory named private and initialize the folder lock
  lock      Hide the private directory
  unlock    Unhide the private directory`, nil
}
