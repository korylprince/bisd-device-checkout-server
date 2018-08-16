package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/korylprince/bisd-device-checkout-server/api"
)

type handlerResponse struct {
	Code int
	Body interface{}
	User *api.User
	Err  error
}

type returnHandler func(http.ResponseWriter, *http.Request) *handlerResponse

const logTemplate = "{{.Date}} {{.Method}} {{.Path}}{{if .Query}}?{{.Query}}{{end}} {{.Code}} ({{.Status}}){{if .User}}, User: {{.User.Username}}{{end}}{{if .Err}}, Error: {{.Err}}{{end}}\n"

type logData struct {
	Date   string
	User   *api.User
	Status string
	Code   int
	Method string
	Path   string
	Query  string
	Err    error
}

func logMiddleware(next returnHandler, writer io.Writer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := next(w, r)

		err := template.Must(template.New("log").Parse(logTemplate)).Execute(writer, &logData{
			Date:   time.Now().Format("2006-01-02:15:04:05 -0700"),
			User:   resp.User,
			Status: http.StatusText(resp.Code),
			Code:   resp.Code,
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Err:    resp.Err,
		})

		if err != nil {
			panic(err)
		}
	})
}

func jsonMiddleware(next returnHandler) returnHandler {
	return func(w http.ResponseWriter, r *http.Request) *handlerResponse {
		var resp *handlerResponse

		if r.Method != "GET" {
			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				resp = handleError(http.StatusBadRequest, errors.New("Could not parse Content-Type"))
				goto serve
			}
			if mediaType != "application/json" {
				resp = handleError(http.StatusBadRequest, errors.New("Content-Type not application/json"))
				goto serve
			}
		}

		w.Header().Set("Content-Type", "application/json")
		resp = next(w, r)

	serve:
		w.WriteHeader(resp.Code)
		e := json.NewEncoder(w)
		err := e.Encode(resp.Body)
		if err != nil {
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could encode json: %v", err))
		}
		return resp
	}
}

func authMiddleware(next returnHandler, s SessionStore) returnHandler {
	return func(w http.ResponseWriter, r *http.Request) *handlerResponse {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			return handleError(http.StatusUnauthorized, errors.New("No Authorization header"))
		}

		if !strings.HasPrefix(auth, `Session id="`) || len(auth) < 13 {
			return handleError(http.StatusBadRequest, errors.New("Invalid Authorization header"))
		}

		id := auth[12 : len(auth)-1]

		sess, err := s.Check(id)
		if err != nil {
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could not check session key: %v", err))
		}
		if sess == nil {
			return handleError(http.StatusUnauthorized, errors.New("Could not find session"))
		}

		ctx := context.WithValue(r.Context(), api.UserKey, sess.User)
		resp := next(w, r.WithContext(ctx))
		resp.User = sess.User

		return resp
	}
}

func txMiddleware(next returnHandler, inventoryDB, skywardDB *sql.DB) returnHandler {
	return func(w http.ResponseWriter, r *http.Request) *handlerResponse {
		//create inventory tx
		itx, err := inventoryDB.Begin()
		if err != nil {
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could not begin Inventory transaction: %v", err))
		}
		ctx := context.WithValue(r.Context(), api.InventoryTransactionKey, itx)

		//create skyward tx
		stx, err := skywardDB.Begin()
		if err != nil {
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could not begin Skyward transaction: %v", err))
		}
		ctx = context.WithValue(ctx, api.SkywardTransactionKey, stx)

		resp := next(w, r.WithContext(ctx))

		//commit skyward tx
		if err = stx.Commit(); err != nil {
			//rollback skyward tx
			if rErr := stx.Rollback(); rErr != nil && rErr != sql.ErrTxDone {
				return handleError(http.StatusInternalServerError, fmt.Errorf("Could not rollback Skyward transaction: %v", rErr))
			}
			//rollback inventory tx
			if rErr := itx.Rollback(); rErr != nil && rErr != sql.ErrTxDone {
				return handleError(http.StatusInternalServerError, fmt.Errorf("Could not rollback Inventory transaction: %v", rErr))
			}
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could not commit Skyward transaction: %v", err))
		}

		//commit inventory tx
		if err = itx.Commit(); err != nil {
			if rErr := itx.Rollback(); rErr != nil && rErr != sql.ErrTxDone {
				return handleError(http.StatusInternalServerError, fmt.Errorf("Could not rollback Inventory transaction: %v", rErr))
			}
			return handleError(http.StatusInternalServerError, fmt.Errorf("Could not commit Inventory transaction: %v", err))
		}

		return resp
	}
}
