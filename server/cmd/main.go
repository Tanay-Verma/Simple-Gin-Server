package main

import (
	"log"
	"server/db"
	"server/internal/user"
	"server/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not establish connection to the database: %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHadler := user.NewHandler(userSvc)

	router.InitRouter(userHadler)

	err = router.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Could not start the router")
	}
}
