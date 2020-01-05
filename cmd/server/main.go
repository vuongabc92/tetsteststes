package main

import (
	"flag"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/vuongabc92/octocv/config"
	database "github.com/vuongabc92/octocv/database/mongodb"
	"github.com/vuongabc92/octocv/http"
	"github.com/vuongabc92/octocv/http/actions"
	"github.com/vuongabc92/octocv/http/log"
	"github.com/vuongabc92/octocv/http/mail"
	"github.com/vuongabc92/octocv/http/middleware"
	"github.com/vuongabc92/octocv/http/router"
	authrouter "github.com/vuongabc92/octocv/http/router/auth"
	homerouter "github.com/vuongabc92/octocv/http/router/home"
	"github.com/vuongabc92/octocv/lang"
	"github.com/vuongabc92/octocv/render"
	tpl "github.com/vuongabc92/octocv/template"
	"github.com/vuongabc92/octocv/worker"
	"html/template"
)

func main() {
	flag.Parse()

	// Init session
	store := sessions.NewCookieStore([]byte(*config.SessionKey))

	// Connect to DB
	dbConfig := database.MongoDBConfig{
		ConnectionString: *config.MongoDBConnectionStr,
		DBName:           *config.MongoDBName,
	}
	dbConnection := database.NewMongoDBConnection(dbConfig)
	dbConnection.Connect()

	// Logging
	logFiles := make(map[string]string)
	logFiles[log.ErrorLog] = *config.ErrorLogFile
	logFiles[log.InfoLog] = *config.InfoLogFile
	logFiles[log.DebugLog] = *config.DebugLogFile
	logger := log.NewLog(logFiles)

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     *config.RedisAddress,
		Password: *config.RedisPassword,
		DB:       0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		logger.Errorf("Can not connect to redis. Error: %s", err.Error())
		panic("Can not connect to to redis. Error: " + err.Error())
	}

	// New server
	opts := http.Options{
		SessionStore: store,
		SessionName:  *config.SessionName,
		DB:           dbConnection,
		Log:          logger,
		Redis:        redisClient,
	}
	s := http.NewServer(opts)

	// Init routers
	routers := []router.Router{
		authrouter.NewRouter(actions.Auth{}),
		homerouter.NewRouter(actions.Auth{}),
	}
	s.InitRouter(routers...)

	// Init middlewares
	if err := initMiddlewares(s); err != nil {
		s.Log.Errorf("Can not init middleware. Error: %s", err)
		panic("Can not init middleware. Error: " + err.Error())
	}

	// Load view
	templateEngine := tpl.NewTemplateEngine()
	templateEngine.HTMLRender = loadTemplate(s.Router())
	s.SetAddr(*config.HttpAddr).SetTemplateEngine(templateEngine)

	// Load multiple language
	langMsg := lang.NewMessage()
	if err := langMsg.Load(); err != nil {
		s.Log.Errorf("Can not load multiple language. Error: %s", err.Error())
		panic("Can not load multiple language. Error: " + err.Error())
	}

	workerOpts := &worker.Option{
		Redis: redisClient,
		SMTPMail: &mail.SMTPMail{
			Host:     *config.SmtpHost,
			Port:     *config.SmtpPort,
			Username: *config.SmtpUsername,
			Password: *config.SmtpPassword,
		},
		TemplateEngine: templateEngine,
		Logger:         logger,
	}
	worker := worker.NewWorker(workerOpts)
	worker.Run()

	// Run server
	s.Serve()
}

func loadTemplate(mux *mux.Router) render.MultipleHTML {
	r := render.NewMultipleHTML()
	templateFuncs := getTemplateFuncs(mux)

	// Register view template
	r.AddFromFilesFuncs("auth.login", templateFuncs, viewPath("layouts/__auth.html"), viewPath("auth/login.html"))
	r.AddFromFilesFuncs("auth.register", templateFuncs, viewPath("layouts/__auth.html"), viewPath("auth/register.html"))
	r.AddFromFilesFuncs("auth.forgot-password", templateFuncs, viewPath("layouts/__auth.html"), viewPath("auth/forgot-password.html"))
	r.AddFromFilesFuncs("auth.reset-password", templateFuncs, viewPath("layouts/__auth.html"), viewPath("auth/reset-password.html"))
	r.AddFromFilesFuncs("mail.register-confirmation", templateFuncs, viewPath("mail/register-confirmation.html"))
	r.AddFromFilesFuncs("mail.forgot-password", templateFuncs, viewPath("mail/forgot-password.html"))
	r.AddFromFilesFuncs("page.home", templateFuncs, viewPath("layouts/__master.html"), viewPath("pages/home.html"))

	return r
}

func getTemplateFuncs(r *mux.Router) template.FuncMap {
	return template.FuncMap{
		"url": func(name string, queryPairs ...string) string {
			url, _ := r.Get(name).URL(queryPairs...)
			return url.EscapedPath()
		},
		"full_url": func(name string, queryPairs ...string) string {
			url, _ := r.Get(name).URL(queryPairs...)
			flag.Parse()
			return *config.BaseUrl + url.EscapedPath()
		},
		"form_msg":            tpl.MessageBagGet,
		"form_has_msg":        tpl.MessageBagHas,
		"form_text":           tpl.FormDataText,
		"trans":               tpl.Trans,
		"asset_url":           tpl.FrontendAsset,
		"frontend_full_asset": tpl.FrontendFullAsset,
		"support_email":       tpl.SupportEmailAddress,
	}
}

func viewPath(name string) string {
	return config.ViewFrontendPath + "/" + name
}

func initMiddlewares(s *http.Server) error {
	s.UseMiddleware(middleware.NewRequestIDMiddleware())
	s.UseMiddleware(middleware.NewLangMiddleware())

	return nil
}
