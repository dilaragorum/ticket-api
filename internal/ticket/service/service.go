package service

import (
	"context"
	"errors"

	"github.com/dilaragorum/ticket-api/internal/ticket"
	"github.com/dilaragorum/ticket-api/internal/ticket/repository"
)

var (
	ErrNameIsEmpty              = errors.New("name should not be empty")
	ErrNameIsDuplicate          = errors.New("ticket name exists already")
	ErrDescriptionIsEmpty       = errors.New("description should not be empty")
	ErrAllocationIsLowerThanOne = errors.New("allocation should be higher than zero ")

	ErrTicketWasNotFound = errors.New("ticket does not exist")
	ErrIDLowerThanOne    = errors.New("id must not be lower than one")

	ErrPurchaseTicketMoreThanAvailable = errors.New("quantity of ticket wanted to be purchased must " +
		"not be more than available ones")
	ErrQuantityLowerThanOne = errors.New("quantity must not be lower than one")
)

type Service interface {
	CreateTicketOption(ctx context.Context, name, description string, allocation int) (*ticket.Ticket, error)
	GetTicket(ctx context.Context, id int) (*ticket.Ticket, error)
	PurchaseFromTicketOption(ctx context.Context, id, quantity int, userID string) error
}

type DefaultService struct {
	repository repository.Repository
}

func NewDefaultService(repository repository.Repository) *DefaultService {
	return &DefaultService{repository: repository}
}

func (s *DefaultService) CreateTicketOption(ctx context.Context, name, description string, allocation int) (*ticket.Ticket, error) {
	if name == "" {
		return nil, ErrNameIsEmpty
	}

	if description == "" {
		return nil, ErrDescriptionIsEmpty
	}

	if allocation < 1 {
		return nil, ErrAllocationIsLowerThanOne
	}

	option, err := s.repository.CreateTicketOption(ctx, name, description, allocation)
	if err != nil {
		if errors.Is(err, repository.ErrDBDuplicatedTicketName) {
			return nil, ErrNameIsDuplicate
		}
		return nil, err
	}

	return option, nil
}

func (s *DefaultService) GetTicket(ctx context.Context, id int) (*ticket.Ticket, error) {
	if id < 1 {
		return nil, ErrIDLowerThanOne
	}

	t, err := s.repository.GetTicket(ctx, id)
	if err != nil {
		switch err {
		case repository.ErrDBTicketNotFound:
			return nil, ErrTicketWasNotFound
		default:
			return nil, err
		}
	}

	return t, nil
}

func (s *DefaultService) PurchaseFromTicketOption(ctx context.Context, id, quantity int, userID string) error {
	if quantity < 1 {
		return ErrQuantityLowerThanOne
	}

	ticketOption, err := s.GetTicket(ctx, id)
	if err != nil {
		return err
	}

	if ticketOption.Allocation < quantity {
		return ErrPurchaseTicketMoreThanAvailable
	}

	if err := s.repository.PurchaseFromTicketOption(ctx, id, quantity, userID); err != nil {
		return err
	}

	return nil
}
