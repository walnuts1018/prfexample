package router

import (
	"net/http"

	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	"github.com/walnuts1018/PRFExample/server/config"
	"github.com/walnuts1018/PRFExample/server/router/handler"
)

func NewRouter(cfg config.Server, handler handler.Handler) http.Handler {
	e := echo.New()
	e.Use(echootel.NewMiddleware(cfg.Origin.Hostname()))

	e.GET("/readyz", handler.Readiness)
	e.GET("/livez", handler.Liveness)

	apiv1 := e.Group("/api").Group("/v1")
	apiv1.Use(handler.SessionMiddleware)
	{
		webauthn := apiv1.Group("/webauthn")
		// webauthn.GET("", handler.GetWebAuthnCredentials)
		{
			registration := webauthn.Group("/registration")
			registration.GET("/creation", handler.GetWebAuthnCredentialCreation)
			registration.POST("/create", handler.CreateWebAuthnCredential)
		}

		{
			verification := webauthn.Group("/verification")
			verification.GET("/assertion", handler.GetWebAuthnCredentialAssertion)
			verification.POST("/verify", handler.VerifyWebAuthnCredentialAssertion)
		}
	}
	{
		data := apiv1.Group("/data")
		data.GET("", handler.ListEncryptedData)
		data.GET(":id", handler.GetEncryptedData)
		data.POST("", handler.SaveEncryptedData)
	}
	return e
}
