package handler

import (
	"net/http"
	"strconv"

	"github.com/dilaragorum/ticket-api/internal/ticket/service"
	"github.com/labstack/echo/v4"
)

var (
	WarnMessageWhenNameIsEmpty              = "Name cannot be empty."
	WarnMessageWhenNameIsDuplicated         = "This name is already used"
	WarnMessageWhenDescriptionIsEmpty       = "Description cannot be empty."
	WarnMessageWhenAllocationIsBelowThanOne = "Allocation cannot be below than one."

	WarnMessageWhenInvalidID         = "Id need to be valid"
	WarnMessageWhenTicketWasNotFound = "Ticket was not found"

	WarnMessageWhenPurchaseTicketMoreThanAvailable = "Quantity of ticket wanted to be purchased is " +
		"higher than available ones"
	WarnMessageWhenQuantityLowerThanOne = "Quantity cannot be lower than one"
	WarnInternalServerError             = "an error occurred please try again later"
)

type DefaultHandler struct {
	service service.Service
}

func NewDefaultTicketHandler(e *echo.Echo, service service.Service) *DefaultHandler {
	t := DefaultHandler{service: service}

	e.GET("/ticket/:id", t.GetTicket)
	e.POST("/ticket_options", t.CreateTicketOption)
	e.POST("/ticket_options/:id/purchases", t.PurchaseFromTicketOption)

	return &t
}

// CreateTicketOption
// @Tags ticket
// @Summary      Create Ticket Option
// @Description  Create a ticket_option with an allocation of tickets available to purchase
// @Param requestBody body CreateTicketOptionRequestBody true "Create Ticket Option Request Body"
// @Accept       json
// @Produce      json
// @Success      201  {object}  ticket.Ticket
// @Failure      400              {string}  string
// @Failure      500              {string}  string
// @Router       /ticket_options [post]
func (t *DefaultHandler) CreateTicketOption(c echo.Context) error {
	options := new(CreateTicketOptionRequestBody)

	if err := c.Bind(&options); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	ticketOptions, err := t.service.CreateTicketOption(c.Request().Context(), options.Name, options.Desc, options.Allocation)
	if err != nil {
		switch err {
		case service.ErrNameIsEmpty:
			return c.String(http.StatusBadRequest, WarnMessageWhenNameIsEmpty)
		case service.ErrDescriptionIsEmpty:
			return c.String(http.StatusBadRequest, WarnMessageWhenDescriptionIsEmpty)
		case service.ErrAllocationIsLowerThanOne:
			return c.String(http.StatusBadRequest, WarnMessageWhenAllocationIsBelowThanOne)
		case service.ErrNameIsDuplicate:
			return c.String(http.StatusBadRequest, WarnMessageWhenNameIsDuplicated)
		default:
			return c.String(http.StatusInternalServerError, WarnInternalServerError)
		}
	}

	return c.JSON(http.StatusCreated, *ticketOptions)
}

// GetTicket
// @Tags ticket
// @Summary      Get ticket by ticket id
// @Description  Get specified ticket with ID from available tickets
// @Produce      json
// @Param        id   path      int  true  "Ticket ID"
// @Success      200  {object}  ticket.Ticket
// @Failure      400              {string}  string
// @Failure      404              {string}  string
// @Failure      500              {string}  string
// @Router       /ticket/{id} [get]
func (t *DefaultHandler) GetTicket(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidID)
	}

	ticket, err := t.service.GetTicket(c.Request().Context(), id)
	if err != nil {
		switch err {
		case service.ErrTicketWasNotFound:
			return c.String(http.StatusNotFound, WarnMessageWhenTicketWasNotFound)
		case service.ErrIDLowerThanOne:
			return c.String(http.StatusBadRequest, WarnMessageWhenInvalidID)
		default:
			return c.String(http.StatusInternalServerError, WarnInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, ticket)
}

// PurchaseFromTicketOption
// @Tags ticket
// @Summary      Purchase from Ticket Option
// @Description  Purchase a quantity of tickets from the allocation of the given ticket_option
// @Accept       json
// @Param requestBody body CreatePurchaseTicketOptionRequestBody true "Purchase Ticket Option Request Body"
// @Param        id   path      int  true  "Ticket ID"
// @Success      200
// @Failure      400              {string}  string
// @Failure      500              {string}  string
// @Router       /ticket_options/{id}/purchases [post]
func (t *DefaultHandler) PurchaseFromTicketOption(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, WarnMessageWhenInvalidID)
	}

	purchasedTicketOption := new(CreatePurchaseTicketOptionRequestBody)
	if err = c.Bind(&purchasedTicketOption); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = t.service.PurchaseFromTicketOption(c.Request().Context(), id, purchasedTicketOption.Quantity, purchasedTicketOption.UserID)
	if err != nil {
		switch err {
		case service.ErrPurchaseTicketMoreThanAvailable:
			return c.String(http.StatusBadRequest, WarnMessageWhenPurchaseTicketMoreThanAvailable)
		case service.ErrQuantityLowerThanOne:
			return c.String(http.StatusBadRequest, WarnMessageWhenQuantityLowerThanOne)
		case service.ErrIDLowerThanOne:
			return c.String(http.StatusBadRequest, WarnMessageWhenInvalidID)
		default:
			return c.String(http.StatusInternalServerError, WarnInternalServerError)
		}
	}

	return c.NoContent(http.StatusOK)
}
