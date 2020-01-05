package http

import (
	"context"
	"flag"
	"github.com/go-redis/redis"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/vuongabc92/octocv"
	"github.com/vuongabc92/octocv/config"
	db "github.com/vuongabc92/octocv/database/mongodb"
	"github.com/vuongabc92/octocv/errdefs"
	"github.com/vuongabc92/octocv/http/httputils"
	"github.com/vuongabc92/octocv/http/log"
	"github.com/vuongabc92/octocv/http/middleware"
	"github.com/vuongabc92/octocv/http/router"
	"github.com/vuongabc92/octocv/http/session"
	"github.com/vuongabc92/octocv/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Config provides the configuration for the server
type Options struct {
	SessionStore sessions.Store
	SessionName  string

	// Database
	DB *db.MongoDBConnection

	// Logger
	Log *log.Log

	// Redis
	Redis *redis.Client
}

// Server contains instance details for the server
type Server struct {
	Options

	// Http server
	srv *http.Server

	// Mux router
	router *mux.Router

	// Routers hold a list of router, each one contains one or many routes
	routers []router.Router

	//List middle wares that system use
	middlewares []middleware.Middleware

	//Template engine
	templateEngine *template.Engine
}

// New server instant
func NewServer(opts Options) *Server {
	flag.Parse()
	return &Server{
		Options: opts,
		srv: &http.Server{
			WriteTimeout: *config.WriteTimeout,
			ReadTimeout:  *config.ReadTimeout,
			IdleTimeout:  *config.IdleTimeout,
		},
	}
}

// UseMiddleware appends a new middleware to the request chain.
// This needs to be called before the API routes are configured.
func (s *Server) UseMiddleware(m middleware.Middleware) {
	s.middlewares = append(s.middlewares, m)
}

//Set server address
func (s *Server) SetAddr(addr string) *Server {
	s.srv.Addr = addr
	return s
}

//Set template engine
func (s *Server) SetTemplateEngine(templateEngine *template.Engine) *Server {
	s.templateEngine = templateEngine
	return s
}

// InitRouter initializes the list of routers for the server.
func (s *Server) InitRouter(routers ...router.Router) *mux.Router {
	flag.Parse()

	s.routers = append(s.routers, routers...)

	m := s.createMux()
	m.Schemes("https")

	// CSRF protection
	csrfErrorHandler := httputils.MakeErrorHandler(errdefs.ForbiddenError{})
	if *config.ENV == config.Development {
		csrf := csrf.Protect([]byte(*config.CsrfKey), csrf.FieldName(*config.CsrfFieldName),
			csrf.Secure(false),
			csrf.ErrorHandler(csrfErrorHandler),
		)
		m.Use(csrf)
	} else {
		csrf := csrf.Protect([]byte(*config.CsrfKey), csrf.FieldName(*config.CsrfFieldName),
			csrf.ErrorHandler(csrfErrorHandler),
		)
		m.Use(csrf)
	}

	s.srv.Handler = m
	s.router = m

	return m
}

// Run server with graceful shutdown
func (s *Server) Serve() {
	flag.Parse()

	// Run our server.yml in a goroutine so that it doesn't block.
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			panic("Can not run server. Error: " + err.Error())
		}
	}()

	s.Log.Info("Server has started at address: " + s.srv.Addr)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C),
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/).
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), *config.GracefullTimeout)
	defer cancel()
	// Doesn't block if no connections, but
	// will otherwise wait until the timeout deadline.
	s.srv.Shutdown(ctx)
	// Optionally, we could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if our application should wait for other services
	// to finalize based on context cancellation.
	s.Log.Info("shutting down")
	os.Exit(0)
}

func (s *Server) createMux() *mux.Router {
	m := mux.NewRouter()

	//Register assets path: css, javascript,...
	assetsPath := http.Dir(config.AssetPath)
	m.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assetsPath)))

	for _, router := range s.routers {
		for _, r := range router.Routes() {
			f := s.makeHTTPHandler(r.Handler())
			m.Path(r.Path()).Methods(r.Method()).Handler(f).Name(r.Name())
			s.Log.Infof("Register route: %s, method: %s and name: %s", r.Path(), r.Method(), r.Name())
		}
	}

	notFoundHandler := httputils.MakeErrorHandler(errdefs.PageNotFoundError{})
	m.NotFoundHandler = notFoundHandler
	m.MethodNotAllowedHandler = notFoundHandler

	return m
}

// Handle HTTP request
func (s *Server) makeHTTPHandler(handler httputils.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Handle middle ware
		handlerFunc := s.handlerWithGlobalMiddlewares(handler)

		ctx := s.NewContext(w, r)

		// Temporary save form post data as flash message to read
		// After submit form and get error
		if strings.ToLower(r.Method) == "post" {
			_ = r.ParseForm()

			formData := session.NewFormData()
			for k, v := range r.Form {
				if k == "password" || k == "Password" || k == "pass" {
					continue
				}

				formData.AddData(k, v)
			}

			ctx.Flash.AddFormData(formData)
		}

		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}

		if err := handlerFunc(ctx, vars); err != nil {
			statusCode := errdefs.GetHTTPErrorStatusCode(err)
			if statusCode >= 500 {
				ctx.Logger.Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
			}
			httputils.MakeErrorHandler(err)(w, r)
		}
	}
}

func (s *Server) Router() *mux.Router {
	return s.router
}

// Get session using a request and response.
func (s *Server) getSession(w http.ResponseWriter, r *http.Request) *session.Session {
	ss, _ := s.SessionStore.Get(r, s.SessionName)
	return &session.Session{
		Session: ss,
		Request: r,
		Writer:  w,
	}
}

func (s *Server) NewContext(w http.ResponseWriter, r *http.Request) *octocv.Context {
	// Get session
	ss := s.getSession(w, r)

	// New context for every request
	ctx := &octocv.Context{
		Context:        r.Context(),
		Writer:         w,
		Request:        r,
		TemplateEngine: s.templateEngine,
		Router:         s.router,
		Session:        ss,
		Flash:          session.NewFlash(ss),
		DB:             s.DB,
		Logger:         s.Log,
		Redis:          s.Redis,
	}

	return ctx
}
