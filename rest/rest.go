package rest

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/toravir/glm/config"
	. "github.com/toravir/glm/context"
	"log"
	"net"
	"net/http"
)

var ctxCache Context
var logger *zerolog.Logger

func ListenAndServe(ctx Context) error {
	ctxCache = ctx
	logger = config.GetLogger(ctx)
	listenAddr := config.GetGLMListenAddress(ctx)
	logger.Debug().Caller().Str("listenAddr", listenAddr).Msg("In Rest...")
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("Listen failed: %s\n", err)
	}

	router := urlRouter()
	isHttps, key, crt := config.GetHttpConfig(ctx)
	if isHttps {
		err = http.ServeTLS(l, router, crt, key)
	} else {
		err = http.Serve(l, router)
	}
	if err != nil {
		panic(err)
	}
	return nil
}

func urlRouter() *mux.Router {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{customerName}/deviceInit", deviceInitHdlr)
	rtr.HandleFunc("/{customerName}/licenseAlloc", licenseAllocHdlr)
	rtr.HandleFunc("/{customerName}/licenseFree", licenseFreeHdlr)
	rtr.HandleFunc("/{customerName}/deviceHB", deviceHBHdlr)
	rtr.HandleFunc("/addPurchase", addPurchaseHdlr)
	return rtr
}
