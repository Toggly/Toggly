package api

import (
	"github.com/Toggly/core/internal/pkg/storage"
)

// Engine type
type Engine struct {
	Project ProjectAPI
}

// NewEngine creates new API engine
func NewEngine(storage storage.DataStorage) *Engine {
	return &Engine{
		Project: ProjectAPI{Storage: storage},
	}
}
