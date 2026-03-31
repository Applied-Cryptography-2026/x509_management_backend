package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/your-org/x509-clean-architecture/config"
	"github.com/your-org/x509-clean-architecture/infrastructure/datastore"
	"github.com/your-org/x509-clean-architecture/infrastructure/router"
	"github.com/your-org/x509-clean-architecture/registry"
)

func main() {
	// 1. Load YAML configuration
	config.ReadConfig()

	// 2. Open database connection
	db, err := datastore.NewDB()
	if err != nil {
		log.Fatalln("app: failed to open database:", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 3. Build the DI composition root
	r := registry.NewRegistry(db)

	// 4. Bootstrap Echo HTTP framework
	e := echo.New()
	e = router.NewRouter(e, r.NewAppController())

	// 5. Start the HTTP server
	addr := ":" + config.C.Server.Address
	fmt.Println("x509-clean-architecture server listening at http://localhost" + addr)
	if err := e.Start(addr); err != nil {
		log.Fatalln("app: server failed:", err)
	}
}
