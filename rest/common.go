package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/toravir/glm/db"
	"io"
	"io/ioutil"
	"net/http"
)

func validateUrl(w http.ResponseWriter, r *http.Request) (string, bool) {
	params := mux.Vars(r)
	custName := params["customerName"]
	if _, ok := db.IsValidCustomer(custName); !ok {
		logger.Debug().Caller().Str("CustomerName", custName).Str("Status", "Invalid CustomerName").
			Msg("failing request")
		w.WriteHeader(404)
		return custName, false
	}
	return custName, true
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
