// đây là nơi "lắp ráp" tất cả các bộ phận lại với nhau
package app

import (
	"github.com/Hidas2004/TaskFlow/internal/handlers/v1handler"
	"github.com/Hidas2004/TaskFlow/internal/repositories"
	"github.com/Hidas2004/TaskFlow/internal/routes/v1routes"
	"github.com/Hidas2004/TaskFlow/internal/services/v1services"
)

type AuthModule struct {
	routes v1routes.AuthRoutes
}

func NewAuthModule(ctx *ModuleContext) *AuthModule {
	// 1. Khởi tạo Repo
	authRepo := repositories.NewUserRepository(ctx.DB)

	// 2. Khởi tạo Service (QUAN TRỌNG: Truyền ctx.Config vào thay vì jwtSecret)
	authService := v1services.NewAuthService(authRepo, ctx.Config)

	// 3. Khởi tạo Handler
	authHandler := v1handler.NewAuthHandler(authService)

	// 4. Khởi tạo Routes
	authRoutes := v1routes.NewAuthRoutes(authHandler)

	return &AuthModule{
		routes: *authRoutes,
	}
}
func (am *AuthModule) GetRoutes() *v1routes.AuthRoutes {
	return &am.routes
}
