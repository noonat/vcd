package vcd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/pkg/errors"
)

const (
	captainForeverURL = "http://www.captainforever.com/captainforever.php?cfe="
)

type handlerErrorFunc func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) error

// handleError turns an handlerErrorFunc into an httprouter.Handle function.
// If the wrapped function returns an error, it logs the error and writes a
// 500 Internal Server Error response.
func handleError(fn handlerErrorFunc) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if err := fn(w, req, params); err != nil {
			log.Printf("%+v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

// handleNew handles GET by returning the new vessel form, and handles the
// form POST by inserting a new vessel into the database.
func handleNew(db *sql.DB) httprouter.Handle {
	return handleError(func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) error {
		if req.Method == http.MethodPost {
			host, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				host = "169.254.1.1"
			}
			ip := net.ParseIP(host)
			if err := req.ParseForm(); err != nil {
				return errors.Wrap(err, "error parsing form")
			}
			formVessel, err := NewVesselFromString(req.Form.Get("vessel_data"))
			if err != nil {
				return errors.Wrap(err, "error parsing vessel from form")
			}
			vessel, err := QueryVesselByCFE(context.Background(), db, formVessel.CFE)
			if err != nil {
				return err
			} else if vessel == nil {
				vessel, err = CreateVessel(context.Background(), db, ip, formVessel.CFE, formVessel.Data)
				if err != nil {
					return err
				}
			}
			http.Redirect(w, req, fmt.Sprintf("/vessels/%d", vessel.ID), http.StatusTemporaryRedirect)
			return nil
		}
		if err := newTmpl.ExecuteTemplate(w, "layout", newTmplArgs{}); err != nil {
			return errors.Wrap(err, "error rendering new template")
		}
		return nil
	})
}

// handleList returns the list of vessels.
func handleList(db *sql.DB) httprouter.Handle {
	return handleError(func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) error {
		vessels, err := QueryVessels(context.Background(), db)
		if err != nil {
			return errors.Wrap(err, "error querying vessels")
		}
		if err := listTmpl.ExecuteTemplate(w, "layout", vessels); err != nil {
			return errors.Wrap(err, "error rendering list template")
		}
		return nil
	})
}

// handleVessel renders the vessel.
func handleVessel(db *sql.DB) httprouter.Handle {
	return handleError(func(w http.ResponseWriter, req *http.Request, params httprouter.Params) error {
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		vessel, err := QueryVesselByID(context.Background(), db, id)
		if err != nil {
			return errors.Wrap(err, "error querying vessel")
		} else if vessel == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		if err := showTmpl.ExecuteTemplate(w, "layout", vessel); err != nil {
			return errors.Wrap(err, "error rendering show template")
		}
		return nil
	})
}

// handleVesselPilot logs a click and redirects to the Captain Forever page.
func handleVesselPilot(db *sql.DB) httprouter.Handle {
	return handleError(func(w http.ResponseWriter, req *http.Request, params httprouter.Params) error {
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		vessel, err := QueryVesselByID(context.Background(), db, id)
		if err != nil {
			return errors.Wrap(err, "error querying vessel")
		} else if vessel == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		http.Redirect(w, req, captainForeverURL+url.QueryEscape(vessel.CFE), http.StatusTemporaryRedirect)
		return nil
	})
}
