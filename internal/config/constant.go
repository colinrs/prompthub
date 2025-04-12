package config

const (
	UserStatusNormal  = 1
	UserStatusLocked  = 2
	UserStatusPending = 3
	UserStatusDelete  = 4

	CategoryStatusNormal    = 1
	CategoryStatusUnDisplay = 2

	PromptsStatusNormal = 1
	PromptsStatusLocked = 2
)

const (
	LikeAction   = "like"
	UnlikeAction = "unlike"
	SaveAction   = "save"
	UnSaveAction = "unSave"

	Newest     = "newest"
	Popularity = "popular"

	EmailVerificationEvent = "emailVerification"
	ForgotPasswordEvent    = "forgotPassword"
)

const (
	VerificationCode             = "verification_code"
	VerificationEmailLimitKey    = "verification_limit_key"
	SingleVerificationEmailLimit = "single_verification_limit"
)
