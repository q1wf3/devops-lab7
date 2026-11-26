package grpc

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	"example.com/ticket-service/pb"
	"example.com/ticket-service/ticket"
	"example.com/ticket-service/transport"
)

type grpcServer struct {
	create  grpctransport.Handler
	getByID grpctransport.Handler
	sell    grpctransport.Handler
	pb.UnimplementedTicketServiceServer
}

func NewGRPCServer(endpoints transport.Endpoints) pb.TicketServiceServer {
	options := []grpctransport.ServerOption{}

	return &grpcServer{
		create: grpctransport.NewServer(
			endpoints.Create, decodeGRPCCreateRequest, encodeGRPCCreateResponse, options...,
		),
		getByID: grpctransport.NewServer(
			endpoints.GetByID, decodeGRPCGetByIDRequest, encodeGRPCGetByIDResponse, options...,
		),
		sell: grpctransport.NewServer(
			endpoints.Sell, decodeGRPCSellRequest, encodeGRPCSellResponse, options...,
		),
	}
}

func (s *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, resp, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateResponse), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	_, resp, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetByIDResponse), nil
}

func (s *grpcServer) Sell(ctx context.Context, req *pb.SellRequest) (*pb.SellResponse, error) {
	_, resp, err := s.sell.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SellResponse), nil
}

// --- декодеры/энкодеры ---

func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateRequest)
	return transport.CreateRequest{
		Ticket: ticket.Ticket{
			EventName: req.EventName,
			Seat:      req.Seat,
			Price:     req.Price,
		},
	}, nil
}

func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.CreateResponse)
	var errStr string
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}
	return &pb.CreateResponse{Id: resp.ID, Error: errStr}, nil
}

func decodeGRPCGetByIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetByIDRequest)
	return transport.GetByIDRequest{ID: req.Id}, nil
}

func encodeGRPCGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.GetByIDResponse)
	var errStr string
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}
	return &pb.GetByIDResponse{
		Ticket: &pb.Ticket{
			Id: resp.Ticket.ID, EventName: resp.Ticket.EventName,
			Seat: resp.Ticket.Seat, Price: resp.Ticket.Price,
			Sold: resp.Ticket.Sold,
		},
		Error: errStr,
	}, nil
}

func decodeGRPCSellRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SellRequest)
	return transport.SellRequest{ID: req.Id}, nil
}

func encodeGRPCSellResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.SellResponse)
	var errStr string
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}
	return &pb.SellResponse{Error: errStr}, nil
}
