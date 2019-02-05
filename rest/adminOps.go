package rest

import (
	"encoding/json"
	"github.com/toravir/glm/db"
	"net/http"
)

func writeAddPurchaseResp(w http.ResponseWriter, httpStatus int, respStatus int, respStr string) {
	var resp billingAddPurchaseResp

	resp.Status = respStatus
	resp.StatusString = respStr
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Debug().Caller().AnErr("Error", err).Msg("When Encoding resp")
	}
}

func addPurchaseHdlr(w http.ResponseWriter, r *http.Request) {
	body, ok := fetchPayload(w, r)
	if !ok {
		return
	}

	var req billingAddPurchaseReq
	if err := json.Unmarshal(body, &req); err != nil {
		handleUnmarshallErr("billingAddPurchaseReq", err, w)
		return
	}

	err := db.AddCustomerPurchase(req.CustomerName, req.FeatureName,
		req.LicenseCount, req.UsagePeriod)
	httpStatus := 200
	respStatus := 200
	respStr := "Sucessful"

	if err != nil {
		logger.Error().Caller().AnErr("Error", err).Msg("When adding Purchase")
		httpStatus = 400
		respStatus = 400
		respStr = "Add Purchase Failed"
	}

	writeAddPurchaseResp(w, httpStatus, respStatus, respStr)
	return
}
