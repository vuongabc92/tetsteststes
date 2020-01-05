package actions

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/vuongabc92/octocv"
	"github.com/vuongabc92/octocv/config"
	"github.com/vuongabc92/octocv/helpers"
	"github.com/vuongabc92/octocv/http/mail"
	"github.com/vuongabc92/octocv/http/router/auth"
	"github.com/vuongabc92/octocv/http/session"
	"github.com/vuongabc92/octocv/models"
	"github.com/vuongabc92/octocv/worker"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const VerifyEmailJWTIssuer string = "VerifyEmail"
const ResetPasswordJWTIssuer string = "ResetPassword"

type VerifyEmailJWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type ResetPasswordJWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type Auth struct {
}

// Add new user to system.
// This method implement auth backend.
// User will be redirected back to register page if there is any validation errors or
// user will be redirected to a configured page if register success.
func (a Auth) Register(ctx *octocv.Context, req auth.RegisterRequest) error {
	var (
		repoFactory = ctx.GetRepositoryFactory()
		userRepo    = repoFactory.User()
		err         error
		user        models.User
		hash        []byte
	)

	flag.Parse()

	// Is email used?
	if user, err = userRepo.FindByEmail(ctx.Context, req.Email); err == nil {
		// Found an account with the requested email, that means the email was already taken.
		// We not gonna allow this.
		msgBag := session.NewMessageBag()
		msgBag.Add("email.existed", helpers.Trans("email_taken"))
		ctx.Flash.AddError(msgBag)
		ctx.Redirect(ctx.Url("get_register"))
	}

	// Something went wrong, maybe it's a fatal error.
	if !helpers.IsMongoNoDocumentError(err) {
		ctx.Logger.Error(err)
		return err
	}

	// Encrypt password
	if hash, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	// JWT token
	verifyEmailJWTClaims := VerifyEmailJWTClaims{
		req.Email,
		jwt.StandardClaims{
			ExpiresAt: helpers.Now().Add(config.ExpiredAt).Unix(),
			Issuer:    VerifyEmailJWTIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, verifyEmailJWTClaims)

	// Sign and get the complete encoded token as a string using the secret
	var tokenString string
	if tokenString, err = token.SignedString([]byte(*config.JWTSecretKey)); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	user = models.User{
		Email:       req.Email,
		Password:    string(hash),
		Status:      config.UserStatusNew,
		VerifyToken: tokenString,
	}

	if err = userRepo.Insert(ctx.Context, &user); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	mailData := mail.RegisterConfirmationMail{
		MailTo:    req.Email,
		VerifyURL: ctx.FullUrl("get_verify_email", "token", user.VerifyToken),
	}

	var msg []byte
	if msg, err = json.Marshal(mailData); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	// Publish a message to mail worker to send register confirmation email.
	// The message contains user's register information.
	if err = ctx.Redis.Publish(worker.RegisterConfirmationChannel, msg).Err(); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	// Todo: Redirect to welcomepage
	ctx.Redirect(ctx.Url("get_home"))

	return nil
}

// Log user into system
// This method implement auth backend
// User will be redirected back to login page if there is any validation errors or user
// will be redirected to a configured page if success.
func (a Auth) Login(ctx *octocv.Context, req auth.LoginRequest) error {
	var (
		repoFactory = ctx.GetRepositoryFactory()
		userRepo    = repoFactory.User()
		err         error
		user        models.User
	)

	// Find user by email address
	if user, err = userRepo.FindByEmail(ctx.Context, req.Email); err != nil {
		msgBag := session.NewMessageBag()
		msgBag.Add("email.fails", helpers.Trans("login_fails"))
		ctx.Flash.AddError(msgBag)
		ctx.Redirect(ctx.Url("get_login"))
		return nil
	}

	// Compare password user input and the one in DB
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		msgBag := session.NewMessageBag()
		msgBag.Add("email.fails", helpers.Trans("login_fails"))
		ctx.Flash.AddError(msgBag)
		ctx.Redirect(ctx.Url("get_login"))
		return nil
	}

	// TODO: create session or cookie and redirect user.
	ctx.Redirect(ctx.Url("get_home"))

	return nil
}

// Check and send user instruction to get back their password.
// This method implement auth backend.
// User will be redirected back to forgot password form and will see a message
// that "the instruction email has sent" ever the email is invalid, it because security reason,
// user must not know the other user's mail.
func (a Auth) ForgotPassword(ctx *octocv.Context, req auth.ForgotPasswordRequest) error {
	var (
		repoFactory   = ctx.GetRepositoryFactory()
		userRepo      = repoFactory.User()
		resetPassRepo = repoFactory.ResetPassword()
		msgBag        = session.NewMessageBag()
		err           error
		resetPassword models.ResetPassword
	)

	msgBag.Add("email.sent", helpers.Trans("forgot_pass_sent"))
	if _, err = userRepo.FindByEmail(ctx.Context, req.Email); err != nil {
		ctx.Flash.AddSuccess(msgBag)
		ctx.Redirect(ctx.Url("get_forgot_password"))
		return err
	}

	// JWT token
	resetPasswordJWTClaims := ResetPasswordJWTClaims{
		req.Email,
		jwt.StandardClaims{
			ExpiresAt: helpers.Now().Add(config.ExpiredAt).Unix(),
			Issuer:    ResetPasswordJWTIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, resetPasswordJWTClaims)

	// Sign and get the complete encoded token as a string using the secret
	var tokenString string
	if tokenString, err = token.SignedString([]byte(*config.JWTSecretKey)); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	if resetPassword, err = resetPassRepo.FindByEmail(ctx.Context, req.Email); err != nil {
		if !helpers.IsMongoNoDocumentError(err) {
			ctx.Logger.Errorf("Error unknown. Error: %s", err.Error())
			return err
		}

		// Insert new reset password
		resetPassword = models.ResetPassword{
			Email: req.Email,
			Token: tokenString,
			Status: config.ResetPasswordStatusNew,
		}

		if err = resetPassRepo.Insert(ctx.Context, &resetPassword); err != nil {
			ctx.Logger.Error(err)
			return err
		}
	} else {
		// Update user reset password token
		resetPassword.Token = tokenString
		resetPassword.Status = config.ResetPasswordStatusNew
		if err = resetPassRepo.Update(ctx.Context, &resetPassword); err != nil {
			ctx.Logger.Error(err)
			return err
		}
	}

	mailData := mail.ForgotPasswordMail{
		MailTo:    req.Email,
		VerifyURL: ctx.FullUrl("get_reset_password", "token", tokenString),
	}

	var msg []byte
	if msg, err = json.Marshal(mailData); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	if err = ctx.Redis.Publish(worker.ForgotPasswordChannel, msg).Err(); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	ctx.Flash.AddSuccess(msgBag)
	ctx.Redirect(ctx.Url("get_forgot_password"))

	return nil
}

// Verify user register email.
// This method implements auth backend.
// When user register there is a confirmation email is sent to the registered email,
// the email contains a token that will be expired for x minutes.
// If the verify fails, user will be redirected to an error page that show why it fails or redirected to
// welcome page if valid.
func (a Auth) VerifyEmail(ctx *octocv.Context, req auth.VerifyEmailRequest) error {
	var (
		repoFactory = ctx.GetRepositoryFactory()
		userRepo    = repoFactory.User()
		err         error
		user        models.User
	)

	flag.Parse()

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(*config.JWTSecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if user, err = userRepo.FindByVerifyToken(ctx.Context, token.Raw); err == nil {
			if user.Email == claims["email"] {
				// Update user status to verified
				user.Status = config.UserStatusConfirmed
				if err = userRepo.Update(ctx.Context, &user); err != nil {
					ctx.Logger.Error(err)
					return err
				}

				// TODO: Create welcome page
				ctx.Redirect("Welcome page")
			}
		}
	}

	err = errors.New("resent password token is invalid or expired, please try again")
	ctx.Logger.Error(err)
	return err
}

// Verify reset password email.
// This method implements auth backend.
// Verify reset password token, if the token is valid then display reset password form.
func (a Auth) ResetPassword(ctx *octocv.Context, req auth.ResetPasswordRequest) error {
	var (
		repoFactory       = ctx.GetRepositoryFactory()
		resetPasswordRepo = repoFactory.ResetPassword()
		err               error
		resetPassword     models.ResetPassword
	)

	flag.Parse()

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(*config.JWTSecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if resetPassword, err = resetPasswordRepo.FindByToken(ctx.Context, token.Raw); err == nil {
			if resetPassword.Email == claims["email"] && resetPassword.Status == config.ResetPasswordStatusNew {
				resp := auth.ResetPasswordResponse{Token: token.Raw}
				ctx.HTML(http.StatusOK, "auth.reset-password", resp)
				return nil
			}
		}
	}

	err = errors.New("verify token is invalid or expired, please try again")
	ctx.Logger.Error(err)
	return err
}

// Update user password.
// This method implements auth backend.
// Update user's password when they click on forgot password link and get a reset password instruction.
func (a Auth) UpdatePassword(ctx *octocv.Context, req auth.UpdatePasswordRequest) error {
	var (
		repoFactory       = ctx.GetRepositoryFactory()
		resetPasswordRepo = repoFactory.ResetPassword()
		err               error
		resetPassword     models.ResetPassword
	)

	flag.Parse()

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(*config.JWTSecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if resetPassword, err = resetPasswordRepo.FindByToken(ctx.Context, token.Raw); err == nil {
			if resetPassword.Email == claims["email"] {
				var (
					userRepo = repoFactory.User()
					user     models.User
					hash     []byte
				)

				if user, err = userRepo.FindByEmail(ctx.Context, resetPassword.Email); err != nil {
					return err
				}

				// Encrypt password
				if hash, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
					ctx.Logger.Error(err)
					return err
				}

				user.Password = string(hash)
				if err = userRepo.Update(ctx.Context, &user); err != nil {
					return err
				}

				resetPassword.Status = config.ResetPasswordStatusDone
				if err = resetPasswordRepo.Update(ctx.Context, &resetPassword); err != nil {
					return err
				}

				ctx.Redirect(ctx.Url("get_home"))
				return nil
			}
		}
	}

	err = errors.New("verify token is invalid or expired, please try again")
	ctx.Logger.Error(err)
	return err
}
