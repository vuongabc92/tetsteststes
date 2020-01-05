package locales

// America english
var EnUsTexts = map[string]string{
	"login":                    "Login",
	"log_in":                   "Log in",
	"register":                 "Register",
	"create_an_acc":            "Create an account",
	"login_gg":                 "Login with Google",
	"login_fb":                 "Login with Facebook",
	"fname":                    "Full Name",
	"email":                    "Email",
	"pass":                     "Password",
	"pass_char":                "6 + characters",
	"by_register":              "Creating an account means you're okay with our",
	"terms_service":            "Terms of Service",
	"and":                      "and",
	"or":                       "Or",
	"privacy_policy":           "Privacy Policy",
	"login_to":                 "Log in to Octocv",
	"forgot_pass":              "Forgot password?",
	"not_mem":                  "Not a member?",
	"already_amem":             "Already a member?",
	"register_now":             "Register now",
	"forgot_pass_instruction1": "Enter the email address you used when you joined and we’ll send you instructions to reset your password.",
	"forgot_pass_instruction2": "For security reasons, we do NOT store your password. So rest assured that we will never send your password via email.",
	"forgot_pass_sent":         "- Reset password instruction was sent.",
	"reset_pass":               "Reset your password",
	"new_pass":                 "New Password",
	"reset_pass_instruction1":  "For security, your password should be at least 8 characters long. Contain numbers, letters, uppercase, lowercase...",
	"reset_pass_instruction2":  "This is good: Str0#gP@ssword!23, T0M##$ABC123. And these are some very bad: 123456, 111111.",

	// Mail subjects
	"register_confirmation_subject": "Octocv - Verify Email Address",
	"forgot_password_subject":       "Octocv - Reset Password",

	// Mail content
	"mail_confirmation_hi":           "Hi, Hello and Welcome :)",
	"mail_confirmation_we_happy":     "We're very excited to have you on board! Please verify your Email.",
	"mail_confirmation_verify_btn":   "Verify Your Email",
	"mail_confirmation_btn_fails":    "If the button above does not work, please copy this link and paste into your browser:",
	"mail_confirmation_support_team": "Octocv Team",
	"mail_forgot_hi":                 "Hi there",
	"mail_forgot_we_hear":            "We heard you need a password reset. Click on the button below to reset your password",
	"mail_forgot_submit_btn":         "Reset My Password",
	"mail_forgot_not_you":            "If you did not make this request then you can safely ignore this email :)",
	"mail_forgot_link_fail":          "If the button above does not work, please copy this link and paste into your browser",
	"mail_forgot_support_team":       "Octocv Team",

	//Validation messages
	"field_required":             "- This field is required.",
	"unallowed_char":             "- Contains unallowed character.",
	"email_invalid":              "- Email is invalid.",
	"email_taken":                "- This email is already taken.",
	"pass_minmax":                "- Must be between 6 and 64.",
	"fname_min":                  "- Must be greater than 3 characters.",
	"login_fails":                " / Password is invalid.",
	"verify_email_expired":       "Sorry! Verify email token was expired.",
	"verify_email_token_invalid": "Sorry! Verify email token is invalid.",

	// Page title
	"home_page_title": "Octocv - We Help People Find Their Dream Job",
}
