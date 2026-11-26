package ticket

import "context"

// Ticket представляет билет на мероприятие
type Ticket struct {
	ID        string  `json:"id,omitempty"`
	EventName string  `json:"event_name"`
	Seat      string  `json:"seat"`
	Price     float64 `json:"price"`
	Sold      bool    `json:"sold"`
}

// Service описывает бизнес-логику сервиса билетов
type Service interface {
	Create(ctx context.Context, ticket Ticket) (string, error)
	GetByID(ctx context.Context, id string) (Ticket, error)
	Sell(ctx context.Context, id string) error
}

// Repository описывает хранилище билетов
type Repository interface {
	CreateTicket(ctx context.Context, ticket Ticket) error
	GetTicketByID(ctx context.Context, id string) (Ticket, error)
	SellTicket(ctx context.Context, id string) error
}
