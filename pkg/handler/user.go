package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/digitalhouse-dev/dh-kit/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/ncostamagna/axul_user/internal/user"
)

//NewHTTPServer is a server handler
func NewHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {

	r := mux.NewRouter()
	r.Use(commonMiddleware)

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllHandler,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetHandler,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/users/{id}/login", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Login),
		decodeLoginHandler,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Store),
		decodeStoreHandler,
		encodeResponse,
		opts...,
	)).Methods("POST")

	return r

}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func decodeStoreHandler(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.StoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetHandler(_ context.Context, r *http.Request) (interface{}, error) {
	p := mux.Vars(r)
	req := user.GetReq{
		ID: p["id"],
	}

	return req, nil
}

func decodeGetAllHandler(_ context.Context, r *http.Request) (interface{}, error) {
	req := user.GetAllReq{}

	return req, nil
}

func decodeLoginHandler(_ context.Context, r *http.Request) (interface{}, error) {
	req := user.LoginReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	p := mux.Vars(r)
	req.ID = p["id"]

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var resp response.Response
	switch err {
	case user.NotFound:
		resp = response.NotFound(err.Error())
		break
	case user.FieldIsRequired, user.InvalidAuthentication:
		resp = response.BadRequest(err.Error())
		break
	default:
		resp = response.InternalServerError(err.Error())
		break
	}

	w.WriteHeader(resp.StatusCode())

	_ = json.NewEncoder(w).Encode(resp)
}
