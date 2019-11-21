package actions

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/vuongabc92/octocv"
	"github.com/vuongabc92/octocv/config"
	"github.com/vuongabc92/octocv/helpers"
	"github.com/vuongabc92/octocv/http/mail"
	"github.com/vuongabc92/octocv/http/router/auth"
	"github.com/vuongabc92/octocv/http/session"
	"github.com/vuongabc92/octocv/models"
	"github.com/vuongabc92/octocv/worker"
	"golang.org/x/crypto/bcrypt"
)

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
		verifyToken []byte
	)

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

	// Generate verify token, the token will be changed only and only the param change.
	// Token will be expired after an hour (base on updated_at).
	sha1 := sha1.New()
	sha1.Write([]byte(req.Email + helpers.Now().String()))
	verifyToken = sha1.Sum(nil)

	user = models.User{
		Email:    req.Email,
		Password: string(hash),
		Status:   config.UserStatusNew,
		VerifyToken: fmt.Sprintf("%x", verifyToken),
	}

	if err = userRepo.Insert(ctx.Context, &user); err != nil {
		ctx.Logger.Error(err)
		return err
	}

	mailData := mail.RegisterConfirmationMail{
		MailTo:    req.Email,
		VerifyURL: ctx.Url("get_email_verify", "email", req.Email, "signature", fmt.Sprintf("%x", verifyToken)),
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

	// Encrypt password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		msgBag := session.NewMessageBag()
		msgBag.Add("email.fails", helpers.Trans("login_fails"))
		ctx.Flash.AddError(msgBag)
		ctx.Redirect(ctx.Url("get_login"))
		return nil
	}

	// TODO: redirect user.

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
		token         []byte
		resetPassword models.ResetPassword
	)

	msgBag.Add("email.sent", helpers.Trans("forgot_pass_sent"))
	if _, err = userRepo.FindByEmail(ctx.Context, req.Email); err != nil {
		ctx.Flash.AddSuccess(msgBag)
		ctx.Redirect(ctx.Url("get_forgot_password"))
		return err
	}

	// Generate reset password token
	token, err = bcrypt.GenerateFromPassword([]byte(req.Email+helpers.Now().String()), bcrypt.MinCost)
	if err != nil {
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
			Token: fmt.Sprintf("%s", token),
		}

		if err = resetPassRepo.Insert(ctx.Context, &resetPassword); err != nil {
			ctx.Logger.Error(err)
			return err
		}
	} else {
		// Update user reset password token
		resetPassword.Token = fmt.Sprintf("%s", token)
		if err = resetPassRepo.Update(ctx.Context, resetPassword.ID, &resetPassword); err != nil {
			ctx.Logger.Error(err)
			return err
		}
	}

	mailData := mail.ForgotPasswordMail{
		MailTo:    req.Email,
		VerifyURL: "http://octocv.co/password/reset/e0791177f94a0dbe6d04a45e4153507b2c559ca6d12880c2e4064070e3d52bad/master@yopmail.com",
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
