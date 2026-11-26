package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"example.com/ticket-service/transport"
)

func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/tickets").Handler(kithttp.NewServer(
		endpoints.Create, decodeCreateRequest, encodeResponse, options...,
	))
	r.Methods("GET").Path("/tickets/{id}").Handler(kithttp.NewServer(
		endpoints.GetByID, decodeGetByIDRequest, encodeResponse, options...,
	))
	r.Methods("PATCH").Path("/tickets/{id}/sell").Handler(kithttp.NewServer(
		endpoints.Sell, decodeSellRequest, encodeResponse, options...,
	))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Ticket); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return transport.GetByIDRequest{ID: vars["id"]}, nil
}

func decodeSellRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return transport.SellRequest{ID: vars["id"]}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
