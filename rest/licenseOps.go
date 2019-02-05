package rest

import (
	"encoding/json"
	"github.com/toravir/glm/db"
	"net/http"
	"time"
)

func writeLicenseAllocResp(w http.ResponseWriter, httpStatus int, respStatus int, respStr string, tm string) {
	var resp licenseAllocResp

	resp.Status = respStatus
	resp.StatusString = respStr
	resp.CurTime = tm
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding resp")
	}
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

	var req licenseAllocReq
	if err := json.Unmarshal(body, &req); err != nil {
		handleUnmarshallErr("licenseAllocReq", err, w)
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
	tm := time.Now().UTC().Format(time.RFC3339)
	if !db.AllocateLicense(custName, req.Fingerprint, req.FeatureName) {
		httpStatus = 401
		respStatus = 401
		respStr = "No License Available"
	}
	writeLicenseAllocResp(w, httpStatus, respStatus, respStr, tm)
	return
}

func licenseFreeHdlr(w http.ResponseWriter, r *http.Request) {
}
