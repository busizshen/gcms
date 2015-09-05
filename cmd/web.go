package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
	"github.com/robvdl/pongo2gin"

	"github.com/robvdl/gcms/admin"
	"github.com/robvdl/gcms/auth"
	"github.com/robvdl/gcms/config"
)

// CmdWeb starts the web server
var CmdWeb = cli.Command{
	Name:        "web",
	Usage:       "Start the web server",
	Description: "Run gcms web to start the server",
	Action:      runWeb,
	Flags:       []cli.Flag{},
}

// setupMiddleware is an internal method where we setup GIN middleware
func setupMiddleware(r *gin.Engine) {
	// TODO: CACHE_URL should come from an environment variable but this requires
	// validating and parsing of the connection url into it's base components.
	store, err := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte(config.Config.Session_Secret))
	if err != nil {
		log.Fatalln("Failed to connect to Redis.", err)
	}

	r.Use(
		secure.Secure(secure.Options{ // TODO: we should get these from config
			AllowedHosts:          []string{},
			SSLRedirect:           false,
			SSLHost:               "",
			SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
			STSSeconds:            315360000,
			STSIncludeSubdomains:  true,
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
		}),
		sessions.Sessions("session", store),
		auth.UserMiddleware(),
	)
}

// setupRoutes is an internal method where we setup application routes
func setupRoutes(r *gin.Engine) {
	// serve static folder
	r.Static("/static", "./static")

	// FIXME: this is just the temporary location of the login page
	r.GET("/", auth.LoginPage)
	r.POST("/", auth.LoginPage)

	// sessionResource is a special api resource with POST and DELETE endpoints
	// doing a POST to this API with a username and password logs a user in,
	// doing a DELETE to this API logs a user out.
	sessionResource := r.Group("/api/session")
	{
		sessionResource.POST("", auth.LoginAPI)
		sessionResource.DELETE("", auth.LogoutAPI)
	}

	adminRoutes := r.Group("/admin")
	{
		adminRoutes.GET("", admin.AdminPage)
	}
}

// csrfFailed is called by nosurf when the csrf token check fails
func csrfFailed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	fmt.Fprintln(w, nosurf.Reason(r)) // reason of the failure
}

// runWeb is an starts the GIN application
func runWeb(ctx *cli.Context) {
	r := gin.Default()
	r.HTMLRender = pongo2gin.Default() // Use Pongo2 for templates

	setupMiddleware(r)
	setupRoutes(r)

	// Initialise nosurf for csrf token support.
	csrfHandler := nosurf.New(r)
	csrfHandler.SetFailureHandler(http.HandlerFunc(csrfFailed))
	csrfHandler.ExemptRegexp("/api/(.*)") // ignore API urls for the time being

	// Start the Gin application with nosurf (for csrf protection).
	// This is an alternative way to start up the Gin application.
	http.ListenAndServe(":"+config.Config.Port, csrfHandler)
}
