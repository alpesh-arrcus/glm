package main

import (
	. "github.com/toravir/glm/config"
	. "github.com/toravir/glm/db"
	. "github.com/toravir/glm/rest"
	"log"
)

func main() {
    ctx := ParseCmdLineArgs()

	ctx = InitLicenseDb(ctx)

	err := ListenAndServe(ctx)
	if err != nil {
		log.Fatal("serveRest() unexpectedly returned.. Exiting..", err)
	}
}
