package config

import "time"

const (
	// Environments
	Development string = "environment"
	Staging     string = "staging"
	Production  string = "production"

	// Context
	ContextKeyRequestID ContextKey = "requestID"

	// Resource path
	ResourceFolder   string = "resources"
	AssetPath        string = ResourceFolder + "/assets"
	ViewPath         string = ResourceFolder + "/views"
	ViewFrontendPath string = ViewPath + "/frontend"
	AssetVersion     string = "1.0"

	// User
	UserStatusNew       UserStatus = 1
	UserStatusConfirmed UserStatus = 2

	// Reset password
	ResetPasswordStatusNew  ResetPasswordStatus = 1
	ResetPasswordStatusDone ResetPasswordStatus = 2

	// User's verify email token expired at: number of minutes
	ExpiredAt = time.Minute * 60 // Number of seconds (1 hour)
)
