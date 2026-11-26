package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/ticket-service/ticket"
)

// DTO для запросов и ответов
type CreateRequest struct {
	Ticket ticket.Ticket
}

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

type GetByIDResponse struct {
	Ticket ticket.Ticket `json:"ticket"`
	Err    error         `json:"error,omitempty"`
}

type SellRequest struct {
	ID string `json:"id"`
}

type SellResponse struct {
	Err error `json:"error,omitempty"`
}

// Endpoints — набор эндпоинтов сервиса билетов
type Endpoints struct {
	Create  endpoint.Endpoint
	GetByID endpoint.Endpoint
	Sell    endpoint.Endpoint
}

func MakeEndpoints(s ticket.Service) Endpoints {
	return Endpoints{
		Create:  makeCreateEndpoint(s),
		GetByID: makeGetByIDEndpoint(s),
		Sell:    makeSellEndpoint(s),
	}
}

func makeCreateEndpoint(s ticket.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Ticket)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s ticket.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		t, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Ticket: t, Err: err}, nil
	}
}

func makeSellEndpoint(s ticket.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SellRequest)
		err := s.Sell(ctx, req.ID)
		return SellResponse{Err: err}, nil
	}
}
