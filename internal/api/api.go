// Package api is responsible for creating and initializing endpoints for link database and users.
//
package api

import (
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	sdk "github.com/polyse/database-sdk"
	"github.com/rs/zerolog/log"
)

// API structure containing the necessary server settings and responsible for starting and stopping it.
type API struct {
	e            *echo.Echo
	addr         string
	dbCollection string
	dbClient     *sdk.DBClient
}

// AppConfig structure containing the server settings necessary for its operation.
type AppConfig struct {
	NetInterface string
	Timeout      time.Duration
	DBClient     *sdk.DBClient
	DBCollection string
}

func (ac *AppConfig) checkConfig() {
	log.Debug().Msg("checking api application config")

	if ac.NetInterface == "" {
		ac.NetInterface = "localhost:9900"
	}
	if ac.Timeout <= 0 {
		ac.Timeout = 10 * time.Millisecond
	}
	if ac.DBCollection == "" {
		ac.DBCollection = "default"
	}
}

// SearchRequest is strust for storage and validate query param.
type SearchRequest struct {
	Query string `validate:"required" query:"q"`
	// Page  int    `validate:"gte=1" query:"page"`
}

// TemplateData is strust for send data in "search.html" template.
type TemplateData struct {
	Data  []sdk.ResponseData
	Query string
}

// Validator - to add custom validator in echo.
type Validator struct {
	validator *validator.Validate
}

// Validate add go-playground/validator in echo.
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// TemplateRenderer is a custom html/template renderer for Echo framework.
type TemplateRenderer struct {
	*template.Template
}

// Render renders a template document.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Template.ExecuteTemplate(w, name, data)
}

// NewApp returns a new ready-to-launch API object with adjusted settings.
func NewApp(appCfg AppConfig) (*API, error) {
	appCfg.checkConfig()

	log.Debug().Interface("api app config", appCfg).Msg("starting initialize api application")

	e := echo.New()

	a := &API{
		e:            e,
		addr:         appCfg.NetInterface,
		dbClient:     appCfg.DBClient,
		dbCollection: appCfg.DBCollection,
	}

	e.Use(logMiddleware)
	e.Renderer = &TemplateRenderer{
		Template: template.Must(template.ParseGlob("./internal/web/*.html")),
	}
	e.Validator = &Validator{validator: validator.New()}

	e.GET("/healthcheck", a.handleHealthcheck)
	e.GET("/", a.handleIndex)
	e.GET("/search", a.handleSearch)
	e.GET("/about", a.handleAbout)
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
	var err error
	req := &SearchRequest{}

	if err = c.Bind(req); err != nil {
		log.Debug().Err(err).Msg("handleSearch Bind err")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err = c.Validate(req); err != nil {
		if req.Query == "" {
			return a.handleIndex(c)
		}

		log.Debug().Err(err).Msg("handleSearch Validate err")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	data := TemplateData{
		Query: req.Query,
	}
	data.Data, err = a.dbClient.GetData(a.dbCollection, req.Query, 0, 0)
	if err != nil {
		log.Debug().Err(err).Msg("handleSearch GetData err")
	}

	return c.Render(http.StatusOK, "search.html", data)
}

func (a *API) handleAbout(c echo.Context) error {
	return c.File("./internal/web/about.html")
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
