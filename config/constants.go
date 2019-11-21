package config

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

	UserStatusNew       UserStatus = 1
	UserStatusConfirmed UserStatus = 2
)
