package httputils

import (
	"github.com/vuongabc92/octocv"
	"github.com/vuongabc92/octocv/config"
	"github.com/vuongabc92/octocv/errdefs"
	"github.com/vuongabc92/octocv/http/log"
	"github.com/vuongabc92/octocv/render"
	"github.com/vuongabc92/octocv/template"
	gotpl "html/template"
	"net/http"
)

var htmlRender = render.NewMultipleHTML()
var templateEngine = template.NewTemplateEngine()

func init() {
	htmlRender.AddFromFilesFuncs("error.404", gotpl.FuncMap{}, viewPath("error/404.html"))
	htmlRender.AddFromFilesFuncs("error.5xx", gotpl.FuncMap{}, viewPath("error/5xx.html"))

	templateEngine.HTMLRender = htmlRender
}

// APIFunc is an adapter to allow the use of ordinary functions as server endpoints.
// Any function that has the appropriate signature can be registered as an server endpoint (e.g. getVersion).
type HandlerFunc func(ctx *octocv.Context, vars map[string]string) error

// MakeErrorHandler makes an HTTP handler that decodes error and
// returns it in the response.
func MakeErrorHandler(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Error logger
		logFiles := make(map[string]string)
		logFiles[log.ErrorLog] = *config.ErrorLogFile
		logger := log.NewLog(logFiles)

		statusCode := errdefs.GetHTTPErrorStatusCode(err)
		logger.Errorf("Error status code: %d, error message: %s", statusCode, err.Error())

		instance := templateEngine.HTMLRender.Instance("error.5xx", nil)
		if statusCode == http.StatusNotFound {
			instance = templateEngine.HTMLRender.Instance("error.404", nil)
		}

		w.WriteHeader(statusCode)
		if err := instance.Render(w); err != nil {
			panic("Can not load error page. Error: " + err.Error())
		}
	}
}

func viewPath(name string) string {
	return config.ViewFrontendPath + "/" + name
}
