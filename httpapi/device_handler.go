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
	bagTag := mux.Vars(r)["bagTag"]

	//read Other ID
	var req *CheckoutRequest
	d := json.NewDecoder(r.Body)

	err := d.Decode(&req)
	if err != nil {
		return handleError(http.StatusBadRequest, fmt.Errorf("Could not decode JSON: %v", err))
	}

	err = api.CheckoutDevice(r.Context(), bagTag, req.UserID)
	if resp := checkAPIError(err); resp != nil {
		return resp
	}

	return &handlerResponse{Code: http.StatusOK, Body: map[string]string{"Status": "OK"}}
}
