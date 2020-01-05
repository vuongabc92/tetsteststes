package config

import (
	"flag"
	"time"
)

var (
	ENV              = flag.String("ENV", Development, "Application environment. There are 3 env: development, staging, production, so depend on each env the the system may has some differences.")
	BaseUrl          = flag.String("base-url", "http://localhost:8080", "Application base url.")
	HttpAddr         = flag.String("server-address", ":8080", "HTTP service address.")
	WriteTimeout     = flag.Duration("write-timeout", time.Second*15, "The duration for server write timeout.")
	ReadTimeout      = flag.Duration("read-timeout", time.Second*15, "The duration for server read timeout.")
	IdleTimeout      = flag.Duration("idle-timeout", time.Second*30, "The duration for server idle timeout.")
	GracefullTimeout = flag.Duration("graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m.")

	// CSRF protection
	CsrfKey       = flag.String("crsf-key", "b9ea44e85b6e0fd4fc649f345f7c389d0f26692e86c9cc46b01117a6c3b9d694", "CSRF key for crsf protection.")
	CsrfFieldName = flag.String("crsf-field-name", "_csrf_token", "CSRF field name in the html form template.")

	// Session cookie
	SessionKey  = flag.String("session-key", "2m6Mc#T$7<ok-SBBR10mn^o7E$),U8", "Session secret key.")
	SessionName = flag.String("session-name", "_octocv_session", "Session name.")

	// Mongo DB
	MongoDBConnectionStr = flag.String("mongodb_connection_str", "mongodb://octocv:root@172.17.12.2:27017/octocv", "Connection string to connect to MongoDB.")
	MongoDBName          = flag.String("mongodb_name", "octocv", "The database name.")

	//Log
	ErrorLogFile = flag.String("error_log_file", "resources/storage/log/error.log", "File holds all error log.")
	DebugLogFile = flag.String("debug_log_file", "resources/storage/log/debug.log", "File holds all debug log.")
	InfoLogFile  = flag.String("info_log_file", "resources/storage/log/info.log", "File holds all info log.")

	// Mail
	SmtpHost            = flag.String("smtp_host", "smtp.mailtrap.io", "SMTP host address.")
	SmtpPort            = flag.Int("smtp_port", 25, "SMTP port.")
	SmtpUsername        = flag.String("smtp_username", "57cea03441dddf", "SMTP authenticate username.")
	SmtpPassword        = flag.String("smtp_password", "ed8b8702e74bcb", "SMTP authenticated password.")
	NoReplyEmailAddress = flag.String("noreply_email", "no-reply@octocv.co", "No-reply email address.")
	SupportEmailAddress = flag.String("support_email", "support@octocv.co", "Support email address.")

	// Redis
	RedisAddress  = flag.String("redis_address", "172.17.12.2:6379", "Redis host string.")
	RedisPassword = flag.String("redis_pass", "", "Redis connection password.")

	// JWT
	JWTSecretKey = flag.String("jwt-secret-key", "P5dUPJ72AfMaut2qzt5YENqM4YXQfcPbGhQ5NDEULxbqfLFCYAEwEexY8VyNaXaY", "JWT secret token for generate secure auth token")
)
