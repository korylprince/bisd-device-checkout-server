package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-device-checkout-server/api"
)

//POST /devices/:bagTag/checkout
func handleCheckoutDevice(w http.ResponseWriter, r *http.Request) *handlerResponse {
	otherID := mux.Vars(r)["otherID"]
	bagTag := mux.Vars(r)["bagTag"]

	err := api.CheckoutDevice(r.Context(), otherID, bagTag)
	if resp := checkAPIError(err); resp != nil {
		return resp
	}

	return &handlerResponse{Code: http.StatusOK, Body: map[string]string{"Status": "OK"}}
}
