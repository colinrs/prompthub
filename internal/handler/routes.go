// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	category "github.com/colinrs/prompthub/internal/handler/category"
	prompt "github.com/colinrs/prompthub/internal/handler/prompt"
	user "github.com/colinrs/prompthub/internal/handler/user"
	"github.com/colinrs/prompthub/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: category.CreateCategoryHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/list",
				Handler: category.GetCategoryListHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1/category"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/create",
				Handler: prompt.CreatePromptHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/delete",
				Handler: prompt.DeletePromptHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/get",
				Handler: prompt.GetPromptHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/like",
				Handler: prompt.LikePromptHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/list",
				Handler: prompt.ListPromptsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/save",
				Handler: prompt.SavePromptHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/search",
				Handler: prompt.SearchPromptsHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/update",
				Handler: prompt.UpdatePromptHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1/prompt"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPut,
				Path:    "/change_password",
				Handler: user.ChangePasswordHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/logout",
				Handler: user.LogoutHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/prompt/list",
				Handler: user.GetUserPromptsListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/prompts/like",
				Handler: user.GetUserLikePromptsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/prompts/save",
				Handler: user.GetUserSavePromptsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/refresh_token",
				Handler: user.RefreshTokenHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterUserHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/update",
				Handler: user.UpdateUserHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1/user"),
	)
}
