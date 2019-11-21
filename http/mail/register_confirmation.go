package mail

type RegisterConfirmationMail struct {
	MailTo    string
	VerifyURL string
}

type ForgotPasswordMail struct {
	MailTo    string
	VerifyURL string
}
