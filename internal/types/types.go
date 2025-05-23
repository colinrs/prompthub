// Code generated by goctl. DO NOT EDIT.
package types

type Category struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CategoryResponse struct {
	List     []Category `json:"list"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Total    int        `json:"total"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,max=30"`
	NewPassword string `json:"newPassword" validate:"required,max=30"`
}

type CreateCategoryRequest struct {
	Name  string `json:"name" validate:"max=50"`
	Color string `json:"color,optional" validate:"omitempty,max=50"`
}

type CreatePromptRequest struct {
	CategoryId uint   `json:"categoryId" validate:"min=1"`
	Title      string `json:"title" validate:"max=50"`
	Content    string `json:"content" validate:"max=10240"`
}

type CreatedBy struct {
	UserID   uint   `json:"userId"`
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
}

type DeletePromptRequest struct {
	PromptID uint `form:"promptId" validate:"omitempty,min=1"`
}

type ForgotPasswordnRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
}

type ForgotPasswordnResponse struct {
}

type GetCategoryListRequest struct {
	Page     int `form:"page,optional,default=1"`
	PageSize int `form:"pageSize,optional,default=10"`
}

type GetPromptRequest struct {
	PromptID uint `form:"promptId" validate:"omitempty,min=1"`
}

type LikePromptRequest struct {
	PromptID uint   `json:"promptId" validate:"omitempty,min=1"`
	Action   string `json:"action"` // like or unlike
}

type LikePromptResponse struct {
}

type ListPromptResponse struct {
	List     []Prompt `json:"list"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Total    int      `json:"total"`
}

type ListPromptsRequest struct {
	Page     int `form:"page,optional,default=1"`
	PageSize int `form:"pageSize,optional,default=10"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=50,email"`
	Password string `json:"password" validate:"required,max=30,min=6"`
}

type LoginResponse struct {
	UserId    uint   `json:"userId"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expireAt"`
}

type Prompt struct {
	Id            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CategoryID    uint      `json:"categoryId"`
	CategoryColor string    `json:"categoryColor"`
	Category      string    `json:"category"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
	CreatedBy     CreatedBy `json:"createdBy"`
	Likes         int       `json:"likes"`
	Liked         bool      `json:"liked"`
	Saved         bool      `json:"saved"`
}

type RefreshTokenRequeste struct {
	Token string `json:"token"`
}

type RefreshTokenResponse struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expireAt"`
}

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"max=50"`
	Password string `json:"password" validate:"min=6"`
	Email    string `json:"email" validate:"email,max=80"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email,max=100"`
	Code        string `json:"code" validate:"required,max=6"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type SavePromptRequest struct {
	PromptID uint   `json:"promptId" validate:"omitempty,min=1"`
	Action   string `json:"action"` // save or unsave
}

type SavePromptResponse struct {
}

type SearchPromptsRequest struct {
	Title      string `form:"title,optional"`
	Content    string `form:"content,optional"`
	CategoryID uint   `form:"categoryId,optional"`
	Sort       string `form:"sort,optional"` // popular,newest
	Page       int    `form:"page,optional,default=1"`
	PageSize   int    `form:"pageSize,optional,default=10"`
}

type SearchPromptsResponse struct {
	List     []Prompt `json:"list"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Total    int      `json:"total"`
}

type SendVerificationCodeRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
	Event string `json:"event" validate:"required,oneof=forgotPassword emailVerification"`
}

type SendVerificationCodeResponse struct {
}

type UpdatePromptRequest struct {
	ID         uint   `json:"id" validate:"omitempty,min=1"`
	CategoryId uint   `json:"categoryId" validate:"min=1"`
	Title      string `json:"title" validate:"max=50"`
	Content    string `json:"content" validate:"max=10240"`
}

type UpdateUserRequest struct {
	Name   string `json:"name,optional" validate:"omitempty,max=50"`
	Avatar string `json:"avatar" validate:"omitempty,max=300"`
}

type UserPromptsResponse struct {
	List     []Prompt `json:"list"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Total    int      `json:"total"`
}

type VerificationCodeRequest struct {
	Code  string `form:"code" validate:"required,min=6"`
	Email string `form:"email" validate:"required,email,max=100"`
}

type VerificationCodeResponse struct {
}
