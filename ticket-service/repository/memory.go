package repository

import (
	"context"
	"sync"

	"example.com/ticket-service/ticket"
)

type memoryRepo struct {
	mu      sync.RWMutex
	tickets map[string]ticket.Ticket
}

func NewMemoryRepository() ticket.Repository {
	return &memoryRepo{
		tickets: make(map[string]ticket.Ticket),
	}
}

func (r *memoryRepo) CreateTicket(ctx context.Context, t ticket.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tickets[t.ID] = t
	return nil
}

func (r *memoryRepo) GetTicketByID(ctx context.Context, id string) (ticket.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tickets[id]
	if !ok {
		return ticket.Ticket{}, ticket.ErrTicketNotFound
	}
	return t, nil
}

func (r *memoryRepo) SellTicket(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tickets[id]
	if !ok {
		return ticket.ErrTicketNotFound
	}
	t.Sold = true
	r.tickets[id] = t
	return nil
}
