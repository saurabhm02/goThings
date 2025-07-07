package handlers

import (
	"go-auth/internals/middlewares"
	"go-auth/internals/models"
	"go-auth/internals/services"

	"github.com/gorilla/mux"
)

type Handler struct {
	UserHandler *UserHandler
}

func NewHandler(userService *services.UserService) *Handler {
	return &Handler{
		UserHandler: NewUserHandler(userService),
	}
}

func (h *Handler) RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", h.UserHandler.RegisterUser).Methods("POST")
	authRouter.HandleFunc("/login", h.UserHandler.LoginUser).Methods("POST")

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.RoleAuthorizationMiddleware(models.RoleAdmin))
	adminRouter.HandleFunc("/update/role", h.UserHandler.DemoteAdminToUser).Methods("PUT")
	adminRouter.HandleFunc("/delete/user", h.UserHandler.RemoveUser).Methods("DELETE")

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(middlewares.RoleAuthorizationMiddleware(models.RoleUser))
	userRouter.HandleFunc("/delete/user", h.UserHandler.RemoveUser).Methods("DELETE")

	return router
}
