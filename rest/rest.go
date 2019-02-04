package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
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
	err = http.Serve(l, router)
	if err != nil {
		panic(err)
	}
	return nil
}

func validateUrl(w http.ResponseWriter, r *http.Request) (string,bool) {
	params := mux.Vars(r)
	custName := params["customerName"]

	if !db.IsValidCustomer(custName) {
		logger.Debug().Caller().Str("CustomerName", custName).Str("Status", "Invalid CustomerName").
			Msg("failing request")
		w.WriteHeader(404)
		return custName,false
	}
	return custName,true
}

func fetchPayload(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10240)) //MAX Payload is 10K
	if err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Reading req payload")
		w.WriteHeader(500)
		return body, false
	}
	if err := r.Body.Close(); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When closing req payload")
		w.WriteHeader(500)
		return body, false
	}
	return body, true
}

func handleUnmarshallErr(msgType string, err error, w http.ResponseWriter) {
	logger.Debug().Caller().AnErr("Error", err).Msgf("When Unmarshalling %s", msgType)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(422) // unprocessable entity
	if err := json.NewEncoder(w).Encode(err); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding Error resp")
	}
}

func validateCustomerSecret(w http.ResponseWriter, custName string, secret string) bool {
	if !db.IsValidCustomerSecret(custName, secret) {
		writeDeviceInitResp(w, 401, 401, "Unauthorized", "")
		return false
	}
	return true
}

func writeDeviceInitResp(w http.ResponseWriter, httpStatus int, respStatus int, respStr string, tm string) {
	var resp deviceInitResp

	resp.Status = respStatus
	resp.StatusString = respStr
	resp.CurTime = tm
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding resp")
	}
}

func deviceInitHdlr(w http.ResponseWriter, r *http.Request) {
	custName, ok := validateUrl(w, r)
	if !ok {
		return
	}
	body, ok := fetchPayload(w, r)
	if !ok {
		return
	}

	var req deviceInitReq
	if err := json.Unmarshal(body, &req); err != nil {
		handleUnmarshallErr("deviceInitReq", err, w)
		return
	}
	if custName != req.CustomerName {
		//Something fishy - payload has different name from the url ??
		//Don't return any payload - may be an attacker ??
		w.WriteHeader(503)
		return
	}

	if !validateCustomerSecret(w, custName, req.CustomerSecret) {
		return
	}

	status, isnew := db.AddDevice(custName, req.Fingerprint)
	httpStatus, respStatus := 200, 200
	respStr := "Successful"
	tm := time.Now().UTC().Format(time.RFC3339)
	if !status {
		respStatus = 400
		respStr = "Backend Device registration Failed!"
	} else {
		if isnew {
			httpStatus = 201
			respStatus = 201 //indicate we have registered this device for first time
		}
	}
	writeDeviceInitResp(w, httpStatus, respStatus, respStr, tm)
	return
}

func deviceHBHdlr(w http.ResponseWriter, r *http.Request) {
}

func licenseAllocHdlr(w http.ResponseWriter, r *http.Request) {
	custName, ok := validateUrl(w, r)
	if !ok {
		return
	}
	body, ok := fetchPayload(w, r)
	if !ok {
		return
	}

	var req deviceInitReq
	if err := json.Unmarshal(body, &req); err != nil {
		handleUnmarshallErr("deviceInitReq", err, w)
		return
	}
	if custName != req.CustomerName {
		//Something fishy - payload has different name from the url ??
		//Don't return any payload - may be an attacker ??
		w.WriteHeader(503)
		return
	}

	if !validateCustomerSecret(w, custName, req.CustomerSecret) {
		return
	}

	status, isnew := db.AddDevice(custName, req.Fingerprint)
	httpStatus, respStatus := 200, 200
	respStr := "Successful"
	tm := time.Now().UTC().Format(time.RFC3339)
	if !status {
		respStatus = 400
		respStr = "Backend Device registration Failed!"
	} else {
		if isnew {
			httpStatus = 201
			respStatus = 201 //indicate we have registered this device for first time
		}
	}
	writeDeviceInitResp(w, httpStatus, respStatus, respStr, tm)
	return
}

func licenseFreeHdlr(w http.ResponseWriter, r *http.Request) {
}

func urlRouter() *mux.Router {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/{customerName}/deviceInit", deviceInitHdlr)
	rtr.HandleFunc("/{customerName}/licenseAlloc", licenseAllocHdlr)
	rtr.HandleFunc("/{customerName}/licenseFree", licenseFreeHdlr)
	rtr.HandleFunc("/{customerName}/deviceHB", deviceHBHdlr)
	return rtr
}
