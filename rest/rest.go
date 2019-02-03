package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/toravir/glm/config"
	. "github.com/toravir/glm/context"
	"github.com/toravir/glm/db"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var ctxCache Context

func ListenAndServe(ctx Context) error {
	ctxCache = ctx
	logger := config.GetLogger(ctx)
	listenAddr := config.GetGLMListenAddress(ctx)
	logger.Debug().Caller().Str("listenAddr", listenAddr).Msg("In Rest...")
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

func deviceInitHdlr(w http.ResponseWriter, r *http.Request) {
	logger := config.GetLogger(ctxCache)
	params := mux.Vars(r)
	custName := params["customerName"]

	if !db.IsValidCustomer(custName) {
		logger.Debug().Caller().Str("CustomerName", custName).Str("Status", "Invalid CustomerName").
			Msg("failing request")
		w.WriteHeader(404)
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240)) //MAX Payload is 10K
	if err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Reading req payload")
		return
	}
	if err := r.Body.Close(); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When closing req payload")
		return
	}
	var req deviceInitReq
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Unmarshalling request")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding Error resp")
		}
		return
	}
	if custName != req.CustomerName {
		//Something fishy - payload has different name from the url ??
		//Don't return any payload - may be an attacker ??
		w.WriteHeader(503)
		return
	}

	var resp deviceInitResp

	if !db.IsValidCustomerSecret(custName, req.CustomerSecret) {
		w.WriteHeader(401)
		resp.Status = 401
		resp.StatusString = "Unauthorized"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding Error resp")
		}
		return
	}

	status, isnew := db.AddDevice(custName, req.Fingerprint)
	if !status {
		w.WriteHeader(200)
		resp.Status = 400
		resp.StatusString = "Backend Device registration Failed!"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding resp")
		}
		return
	}
	resp.Status = 200
	if isnew {
		resp.Status = 201 //indicate we have registered this device for first time
	}
	resp.StatusString = "Successful"
	resp.CurTime = time.Now().UTC().Format(time.RFC3339)

	w.WriteHeader(resp.Status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding resp")
	}
	return
}

func deviceHBHdlr(w http.ResponseWriter, r *http.Request) {
	//logger := config.GetLogger(ctxCache)
}

func licenseAllocHdlr(w http.ResponseWriter, r *http.Request) {
	//logger := config.GetLogger(ctxCache)
}

func licenseFreeHdlr(w http.ResponseWriter, r *http.Request) {
	//logger := config.GetLogger(ctxCache)
}

func urlRouter() *mux.Router {
	//logger := config.GetLogger(ctxCache)
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{customerName}/deviceInit", deviceInitHdlr)
	rtr.HandleFunc("/{customerName}/licenseAlloc", licenseAllocHdlr)
	rtr.HandleFunc("/{customerName}/licenseFree", licenseFreeHdlr)
	rtr.HandleFunc("/{customerName}/deviceHB", deviceHBHdlr)
	return rtr
}
