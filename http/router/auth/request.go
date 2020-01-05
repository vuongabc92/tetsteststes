package auth

type RegisterRequest struct {
	FullName string `form:"full_name" validate:"required,alpha_space,min=3"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6,max=64"`
}

type LoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type ForgotPasswordRequest struct {
	Email string `form:"email"`
}

type VerifyEmailRequest struct {
	Token string `form:"token"`
}

type ResetPasswordRequest struct {
	Token string `form:"token"`
}

type UpdatePasswordRequest struct {
	Token    string `form:"token" validate:"required"`
	Password string `form:"password" validate:"required,min=6,max=64"`
}
