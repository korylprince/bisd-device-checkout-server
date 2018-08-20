package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-device-checkout-server/api"
)

//POST /devices/:bagTag/checkout
func handleCheckoutDevice(w http.ResponseWriter, r *http.Request) *handlerResponse {
	type request struct {
		Note string `json:"note,omitempty"`
	}

	otherID := mux.Vars(r)["otherID"]
	bagTag := mux.Vars(r)["bagTag"]

	var req *request
	d := json.NewDecoder(r.Body)

	err := d.Decode(&req)
	if err != nil || req == nil {
		return handleError(http.StatusBadRequest, fmt.Errorf("Could not decode json: %v", err))
	}

	err = api.CheckoutDevice(r.Context(), otherID, bagTag, req.Note)
	if resp := checkAPIError(err); resp != nil {
		return resp
	}

	return &handlerResponse{Code: http.StatusOK, Body: map[string]string{"Status": "OK"}}
}
