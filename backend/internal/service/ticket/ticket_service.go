package ticker_services

import (
	"context"

	"cinema.com/demo/internal/model"
	"cinema.com/demo/internal/repository"
)

type TicketService struct {
	ticketRepo repository.TicketRepository
}

func NewTicketService(ticketRepo repository.TicketRepository) *TicketService {
	return &TicketService{ticketRepo: ticketRepo}
}

func (s *TicketService) GetTicketByUserID(ctx context.Context, user_id int) ([]model.GetTicketByUserIdResponse, error) {
	return s.ticketRepo.GetTicketByUserID(ctx, user_id)
}
