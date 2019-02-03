package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/toravir/glm/config"
	. "github.com/toravir/glm/context"
	"github.com/toravir/glm/db"
	"log"
	"net"
	"net/http"
)

var ctxCache Context

func ListenAndServe(ctx Context) error {
	ctxCache = ctx
	logger := config.GetLogger(ctx)
	listenAddr := config.GetGLMListenAddress(ctx)
	logger.Debug().Str("listenAddr", listenAddr).Msg("In Rest...")
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("Listen failed: %s\n", err)
	}

	router := urlRouter()
	err = http.Serve(l, router)
	if err != nil {
		panic(err)
	}
	return nil
}

func urlRouter() *mux.Router {
	cus := db.GetCustomerNames()
	fmt.Println("Customers:", cus)
	return mux.NewRouter()
}
