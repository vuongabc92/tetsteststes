package locales

// America english
var EnUsTexts = map[string]string{
	"login":                    "Login",
	"register":                 "Register",
	"create_an_acc":            "Create an account",
	"login_gg":                 "Login with Google",
	"login_fb":                 "Login with Facebook",
	"fname":                    "Full Name",
	"email":                    "Email",
	"pass":                     "Password",
	"pass_char":                "6 + characters",
	"already_register":         "Already have an account?",
	"by_register":              "By registering, you agree to our",
	"terms_service":            "Terms of Service",
	"and":                      "and",
	"or":                       "Or",
	"privacy_policy":           "Privacy Policy",
	"login_to":                 "Log in to Octocv",
	"forgot_pass":              "Forgot password?",
	"new_acc":                  "Need a new account?",
	"forgot_pass_instruction1": "Enter the email address you used when you joined and we’ll send you instructions to reset your password.",
	"forgot_pass_instruction2": "For security reasons, we do NOT store your password. So rest assured that we will never send your password via email.",
	"forgot_pass_sent":         "- Reset password instruction was sent.",

	// Mail subjects
	"register_confirmation_subject": "Octocv - Verify Email Address",
	"forgot_password_subject":       "Octocv - Reset Password",

	//Validation messages
	"field_required": "- This field is required.",
	"unallowed_char": "- Contains unallowed character.",
	"email_invalid":  "- Email is invalid.",
	"email_taken":    "- This email is already taken.",
	"pass_minmax":    "- Must be between 6 and 64.",
	"fname_min":      "- Must be greater than 3 characters.",
	"login_fails":    " / Password is invalid.",
}
