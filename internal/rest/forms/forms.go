package forms

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ForgotPasswordForm struct {
	Email string `json:"email" binding:"required,email"`
}
type ResetPasswordForm struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}
