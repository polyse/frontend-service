// Package api is responsible for creating and initializing endpoints for link database and users.
//
package api

import (
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// API structure containing the necessary server settings and responsible for starting and stopping it.
type API struct {
	e    *echo.Echo
	addr string
}

// AppConfig structure containing the server settings necessary for its operation.
type AppConfig struct {
	NetInterface string
	Timeout      time.Duration
	DB           string
}

func (ac *AppConfig) checkConfig() {
	log.Debug().Msg("checking api application config")

	if ac.NetInterface == "" {
		ac.NetInterface = "localhost:9900"
	}
	if ac.Timeout <= 0 {
		ac.Timeout = 10 * time.Millisecond
	}
}

// SearchRequest is strust for storage and validate query param.
type SearchRequest struct {
	Query  string `validate:"required" query:"q"`
	Limit  int    `validate:"gte=0" query:"limit"`
	Offset int    `validate:"gte=0" query:"offset"`
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// NewApp returns a new ready-to-launch API object with adjusted settings.
func NewApp(appCfg AppConfig) (*API, error) {
	appCfg.checkConfig()

	log.Debug().Interface("api app config", appCfg).Msg("starting initialize api application")

	e := echo.New()

	a := &API{
		e:    e,
		addr: appCfg.NetInterface,
	}

	e.Use(logMiddleware)
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./internal/web/*.html")),
	}

	e.GET("/healthcheck", a.handleHealthcheck)
	e.GET("/", a.handleIndex)
	e.GET("/search", a.handleSearch)
	e.Static("/", "./internal/web")

	log.Debug().Msg("endpoints registered")

	return a, nil
}

func (a *API) handleHealthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}

func (a *API) handleIndex(c echo.Context) error {
	return c.File("./internal/web/index.html")
}

func (a *API) handleSearch(c echo.Context) error {
	// request := &SearchRequest{}
	// if err := c.Bind(request); err != nil {
	// 	log.Debug().Err(err).Msg("handleSearch Bind err")
	// 	return echo.NewHTTPError(http.StatusBadRequest)
	// }

	// return c.String(http.StatusOK, request.Query)

	return c.File("./internal/web/search.html")
}

// Run start the server.
func (a *API) Run() error {
	return a.e.Start(a.addr)
}

// Close stop the server.
func (a *API) Close() error {
	log.Debug().Msg("shutting down server")
	return a.e.Close()
}

func logMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		start := time.Now()

		err := next(c)

		stop := time.Now()

		log.Debug().
			Str("remote", req.RemoteAddr).
			Str("user_agent", req.UserAgent()).
			Str("method", req.Method).
			Str("path", c.Path()).
			Int("status", res.Status).
			Dur("duration", stop.Sub(start)).
			Str("duration_human", stop.Sub(start).String()).
			Msgf("called url %s", req.URL)

		return err
	}
}
