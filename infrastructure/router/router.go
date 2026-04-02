package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/your-org/x509-clean-architecture/interface/controller"
)

// NewRouter wires all HTTP routes and middleware to the Echo instance.
func NewRouter(e *echo.Echo, ac controller.AppController) *echo.Echo {
	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// --- Certificate routes ---
	cert := e.Group("/certificates")
	cert.GET("", func(c echo.Context) error { 
		return ac.Certificate.GetCertificates(controller.NewEchoContext(c)) 
	})

	cert.GET("/:id", func(c echo.Context) error { 
		return ac.Certificate.GetCertificate(controller.NewEchoContext(c)) 
	})

	cert.POST("", func(c echo.Context) error { 
		return ac.Certificate.ImportCertificate(controller.NewEchoContext(c)) 
	})

	cert.DELETE("/:id", func(c echo.Context) error { 
		return ac.Certificate.DeleteCertificate(controller.NewEchoContext(c)) 
	})

	cert.POST("/:id/revoke", func(c echo.Context) error { 
		return ac.Certificate.RevokeCertificate(controller.NewEchoContext(c)) 
	})

	cert.GET("/expiring", func(c echo.Context) error { 
		return ac.Certificate.GetExpiringCertificates(controller.NewEchoContext(c)) 
	})

	cert.POST("/validate", func(c echo.Context) error { 
		return ac.Certificate.ValidateCertificate(controller.NewEchoContext(c)) 
	})

	// --- CSR routes ---
	csr := e.Group("/csrs")
	csr.GET("", func(c echo.Context) error { 
		return ac.CSR.GetCSRs(controller.NewEchoContext(c)) 
	})

	csr.GET("/:id", func(c echo.Context) error { 
		return ac.CSR.GetCSR(controller.NewEchoContext(c)) 
	})

	csr.POST("", func(c echo.Context) error { 
		return ac.CSR.SubmitCSR(controller.NewEchoContext(c)) 
	})

	csr.POST("/:id/approve", func(c echo.Context) error { 
		return ac.CSR.ApproveCSR(controller.NewEchoContext(c)) 
	})

	csr.POST("/:id/reject", func(c echo.Context) error { 
		return ac.CSR.RejectCSR(controller.NewEchoContext(c)) 
	})

	// --- Chain routes ---
	chain := e.Group("/chains")
	chain.GET("", func(c echo.Context) error { 
		return ac.Chain.GetChains(controller.NewEchoContext(c)) 
	})

	chain.GET("/:id", func(c echo.Context) error { 
		return ac.Chain.GetChain(controller.NewEchoContext(c)) 
	})

	chain.POST("", func(c echo.Context) error { 
		return ac.Chain.CreateChain(controller.NewEchoContext(c)) 
	})

	chain.DELETE("/:id", func(c echo.Context) error { 
		return ac.Chain.DeleteChain(controller.NewEchoContext(c)) 
	})
	
	chain.POST("/:id/validate", func(c echo.Context) error { 
		return ac.Chain.ValidateChain(controller.NewEchoContext(c)) 
	})

	return e
}