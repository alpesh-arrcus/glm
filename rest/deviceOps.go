package rest

import (
	"encoding/json"
	"github.com/toravir/glm/db"
	"net/http"
	"time"
)

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
	custName, ok := validateUrl(w, r)
	if !ok {
		return
	}
	body, ok := fetchPayload(w, r)
	if !ok {
		return
	}

	var req deviceHBPunchIn
	if err := json.Unmarshal(body, &req); err != nil {
		handleUnmarshallErr("deviceHBHdlr", err, w)
		return
	}
	if custName != req.CustomerName {
		w.WriteHeader(503)
		return
	}

	if !validateCustomerSecret(w, custName, req.CustomerSecret) {
		return
	}

	httpStatus, respStatus := 200, 200
	respStr := "Successful"

	expiredLics, err := db.DeviceHeartBeat(req.CustomerName, req.Fingerprint, req.AutoRealloc)
	if err != nil {
		httpStatus, respStatus = 400, 400
		respStr = "Error Processing Heartbeat."
	}
	tsNow := time.Now().UTC().Format(time.RFC3339)

	var resp deviceHBPunchOut
	resp.Status = respStatus
	resp.StatusString = respStr
	resp.CurTime = tsNow
	resp.ExpiredLics = expiredLics

	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding resp")
	}
}
