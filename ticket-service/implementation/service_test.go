package implementation_test

import (
	"context"
	"errors"
	"example.com/ticket-service/implementation"
	"example.com/ticket-service/repository"
	"example.com/ticket-service/ticket"
	"github.com/go-kit/log"
	"testing"
)

func TestTicketLifecycle(t *testing.T) {
	ctx := context.Background()
	s := implementation.NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	id, err := s.Create(ctx, ticket.Ticket{EventName: "DevOps", Seat: "A1", Price: 1500, Sold: true})
	if err != nil || id == "" {
		t.Fatalf("create: %q %v", id, err)
	}
	got, err := s.GetByID(ctx, id)
	if err != nil || got.Sold || got.EventName != "DevOps" || got.Seat != "A1" || got.Price != 1500 {
		t.Fatalf("get: %+v %v", got, err)
	}
	if err := s.Sell(ctx, id); err != nil {
		t.Fatal(err)
	}
	got, err = s.GetByID(ctx, id)
	if err != nil || !got.Sold {
		t.Fatalf("sold: %+v %v", got, err)
	}
	if err := s.Sell(ctx, id); !errors.Is(err, ticket.ErrTicketAlreadySold) {
		t.Fatalf("resell: %v", err)
	}
}
func TestMissingTicket(t *testing.T) {
	s := implementation.NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	if _, err := s.GetByID(context.Background(), "missing"); !errors.Is(err, ticket.ErrTicketNotFound) {
		t.Fatal(err)
	}
	if err := s.Sell(context.Background(), "missing"); !errors.Is(err, ticket.ErrTicketNotFound) {
		t.Fatal(err)
	}
}
