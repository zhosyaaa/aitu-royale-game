package forms

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ForgotPasswordForm struct {
	Email string `json:"email" binding:"required,email"`
}
type ResetPasswordForm struct {
	Email           string `json:"email,omitempty"`
	Password        string `json:"password" binding:"required" json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required" json:"password_confirm,omitempty"`
}

type CheckCode struct {
	Email string `json:"email,omitempty"`
	Code  string `json:"code"`
}
