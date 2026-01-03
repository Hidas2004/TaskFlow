package app

import (
	"github.com/Hidas2004/TaskFlow/internal/config"
	"gorm.io/gorm"
)

type ModuleContext struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewModuleContext(db *gorm.DB, cfg *config.Config) *ModuleContext {
	return &ModuleContext{
		DB:     db,
		Config: cfg,
	}
}
