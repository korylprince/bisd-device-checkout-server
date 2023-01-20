package httpapi

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-device-checkout-server/api"
)

// NewRouter returns an HTTP router for the HTTP API
func NewRouter(w io.Writer, config *api.AuthConfig, apikey string, s SessionStore, inventoryDB, skywardDB *sql.DB) http.Handler {

	//construct middleware
	var m = func(h returnHandler) http.Handler {
		return logMiddleware(jsonMiddleware(txMiddleware(authMiddleware(h, s), inventoryDB, skywardDB)), w)
	}
	var mk = func(h returnHandler) http.Handler {
		return logMiddleware(jsonMiddleware(txMiddleware(authKeyMiddleware(h, apikey), inventoryDB, skywardDB)), w)
	}

	r := mux.NewRouter()

	r.Path("/students").Queries("status", "true").Methods("GET").Handler(m(handleReadStudentStatuses))
	r.Path("/students").Methods("GET").Handler(m(handleReadStudentList))
	r.Path("/students/{otherID:[0-9]{6}}/status").Methods("GET").Handler(m(handleReadStudentStatus))
	r.Path("/students/{otherID:[0-9]{6}}/devices/{bagTag:[0-9]{4}}").Methods("POST").Handler(m(handleCheckoutDevice))

	r.Path("/auth").Methods("POST").Handler(logMiddleware(jsonMiddleware(txMiddleware(handleAuthenticate(config, s), inventoryDB, skywardDB)), w))

	r.Path("/nosession/students").Queries("status", "true").Methods("GET").Handler(mk(handleReadStudentStatuses))
	r.Path("/nosession/students").Methods("GET").Handler(mk(handleReadStudentList))
	r.Path("/nosession/students/{otherID:[0-9]{6}}/status").Methods("GET").Handler(mk(handleReadStudentStatus))

	r.NotFoundHandler = m(notFoundHandler)

	return http.StripPrefix("/api/1.4", r)
}
