package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/labstack/gommon/log"

	"github.com/dilaragorum/ticket-api/internal/ticket"
)

var (
	ErrDBTicketNotFound       = errors.New("ticket not found")
	ErrDBDuplicatedTicketName = errors.New(`pq: duplicate key value violates unique constraint "tickets_name_uindex"`)
)

type Repository interface {
	CreateTicketOption(ctx context.Context, name, description string, allocation int) (*ticket.Ticket, error)
	GetTicket(ctx context.Context, id int) (*ticket.Ticket, error)
	PurchaseFromTicketOption(ctx context.Context, id, quantity int, userID string) error
}

type DefaultRepository struct {
	database *gorm.DB
}

func NewDefaultRepository(database *gorm.DB) *DefaultRepository {
	return &DefaultRepository{
		database: database,
	}
}

func (df *DefaultRepository) CreateTicketOption(ctx context.Context, name, description string, allocation int) (*ticket.Ticket, error) {
	ticket := ticket.Ticket{
		Name:       name,
		Desc:       description,
		Allocation: allocation,
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond) //nolint:gomnd
	defer cancel()

	err := df.database.WithContext(timeoutCtx).Model(&ticket).Create(&ticket).Error

	if err != nil {
		if err.Error() == ErrDBDuplicatedTicketName.Error() {
			return nil, ErrDBDuplicatedTicketName
		}
		log.Error(err)
		return nil, err
	}

	return &ticket, nil
}

func (df *DefaultRepository) GetTicket(ctx context.Context, id int) (*ticket.Ticket, error) {
	ticket := ticket.Ticket{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond) //nolint:gomnd
	defer cancel()

	if err := df.database.WithContext(timeoutCtx).Model(&ticket).First(&ticket, "id = ?", id).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDBTicketNotFound
		}

		log.Error(err)
		return nil, err
	}

	return &ticket, nil
}

func (df *DefaultRepository) PurchaseFromTicketOption(ctx context.Context, id, quantity int, userID string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond) //nolint:gomnd
	defer cancel()

	tx := df.database.WithContext(timeoutCtx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := tx.Model(ticket.Ticket{}).Where("id = ?", id).
		Update("allocation", gorm.Expr("allocation - ?", quantity)).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	purchase := ticket.Purchase{
		UserID:   userID,
		TicketID: id,
		Quantity: quantity,
		Model:    gorm.Model{},
	}

	if err = tx.Model(&ticket.Purchase{}).Create(&purchase).Error; err != nil {
		log.Error(err)
		return err
	}

	if err = tx.Commit().Error; err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
