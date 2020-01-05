package httputils

import (
	"fmt"
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
	tplFuncs := gotpl.FuncMap{
		"asset_url": template.FrontendAsset,
	}
	htmlRender.AddFromFilesFuncs("errors.error", tplFuncs, viewPath("error/404.html"))

	templateEngine.HTMLRender = htmlRender
}

// APIFunc is an adapter to allow the use of ordinary functions as server endpoints.
// Any function that has the appropriate signature can be registered as an server endpoint (e.g. getVersion).
type HandlerFunc func(ctx *octocv.Context, vars map[string]string) error

type ErrorResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorTitle   string `json:"error_title"`
	ErrorContent string `json:"error_content"`
}

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

		errResp := ErrorResponse{
			ErrorCode:    statusCode,
			ErrorTitle:   "Shit happens",
			ErrorContent: err.Error(),
		}
		if statusCode == http.StatusNotFound {
			errResp = ErrorResponse{
				ErrorCode:    http.StatusNotFound,
				ErrorTitle:   "Page not found",
				ErrorContent: fmt.Sprintf("The requested URL %s was not found on this server", r.RequestURI),
			}
		}

		instance := templateEngine.HTMLRender.Instance("errors.error", errResp)

		w.WriteHeader(statusCode)
		if err := instance.Render(w); err != nil {
			panic("Can not load error page. Error: " + err.Error())
		}
	}
}

func viewPath(name string) string {
	return config.ViewFrontendPath + "/" + name
}
