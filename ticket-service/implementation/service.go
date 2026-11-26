package implementation

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"

	"example.com/ticket-service/ticket"
)

type service struct {
	repository ticket.Repository
	logger     log.Logger
}

// NewService создаёт новый экземпляр сервиса билетов
func NewService(repo ticket.Repository, logger log.Logger) ticket.Service {
	return &service{
		repository: repo,
		logger:     logger,
	}
}

func (s *service) Create(ctx context.Context, t ticket.Ticket) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	id := uuid.New().String()
	t.ID = id
	t.Sold = false // новый билет всегда непродан

	if err := s.repository.CreateTicket(ctx, t); err != nil {
		_ = level.Error(logger).Log("err", err)
		return "", ticket.ErrCmdRepository
	}
	return id, nil
}

func (s *service) GetByID(ctx context.Context, id string) (ticket.Ticket, error) {
	logger := log.With(s.logger, "method", "GetByID")

	t, err := s.repository.GetTicketByID(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return t, ticket.ErrTicketNotFound
	}
	return t, nil
}

func (s *service) Sell(ctx context.Context, id string) error {
	logger := log.With(s.logger, "method", "Sell")

	// Бизнес-правило: нельзя продать уже проданный билет
	t, err := s.repository.GetTicketByID(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return ticket.ErrTicketNotFound
	}
	if t.Sold {
		return ticket.ErrTicketAlreadySold
	}

	if err := s.repository.SellTicket(ctx, id); err != nil {
		_ = level.Error(logger).Log("err", err)
		return ticket.ErrCmdRepository
	}
	return nil
}
