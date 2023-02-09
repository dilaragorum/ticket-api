package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dilaragorum/ticket-api/internal/ticket"
	"github.com/dilaragorum/ticket-api/internal/ticket/handler"
	"github.com/dilaragorum/ticket-api/internal/ticket/mocks"
	"github.com/dilaragorum/ticket-api/internal/ticket/service"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TicketOption Unit Tests

func Test_Should_Return_Status_Created_When_TicketOption_Is_Valid(t *testing.T) {
	// Given
	requestBody := `{"name":"example","desc":"sample description","allocation":100}`
	req := httptest.NewRequest(http.MethodPost, "/ticket_options", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	expectedCreatedTicketOption := ticket.Ticket{ID: 1, Name: "example", Desc: "sample description", Allocation: 100}
	mockService := mocks.NewMockService(gomock.NewController(t))
	mockService.EXPECT().
		CreateTicketOption(gomock.Any(), "example", "sample description", 100).
		Return(&expectedCreatedTicketOption, nil).Times(1)

	ticketOptHandler := handler.NewDefaultTicketHandler(e, mockService)

	// When
	err := ticketOptHandler.CreateTicketOption(c)

	// Then
	assert.Nil(t, err)

	var actualCreatedTicketOption ticket.Ticket
	_ = json.NewDecoder(rec.Body).Decode(&actualCreatedTicketOption)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, expectedCreatedTicketOption, actualCreatedTicketOption)
}

func Test_Should_Return_Status_BadRequest_When_TicketOption_Is_Not_Valid(t *testing.T) {
	type testCase struct {
		name                string
		ticketRequest       handler.CreateTicketOptionRequestBody
		ticketStatusErr     error
		expectedWarnMessage string
	}

	testCases := []testCase{
		{
			name:                "Test_Should_Return_BadRequest_When_TicketOptions_Name_Is_Empty",
			ticketRequest:       handler.CreateTicketOptionRequestBody{Name: "", Desc: "Sample Description", Allocation: 100},
			ticketStatusErr:     service.ErrNameIsEmpty,
			expectedWarnMessage: handler.WarnMessageWhenNameIsEmpty,
		},
		{
			name:                "Test_Should_Return_BadRequest_When_TicketOptions_Description_Is_Empty",
			ticketRequest:       handler.CreateTicketOptionRequestBody{Name: "Ticket", Desc: "", Allocation: 100},
			ticketStatusErr:     service.ErrDescriptionIsEmpty,
			expectedWarnMessage: handler.WarnMessageWhenDescriptionIsEmpty,
		},
		{
			name:                "Test_Should_Return_BadRequest_When_TicketOptions_Allocation_Below_Than_One",
			ticketRequest:       handler.CreateTicketOptionRequestBody{Name: "Ticket", Desc: "Ticket Description", Allocation: 0},
			ticketStatusErr:     service.ErrAllocationIsLowerThanOne,
			expectedWarnMessage: handler.WarnMessageWhenAllocationIsBelowThanOne,
		},
		{
			name:                "Test_Should_Return_BadRequest_When_TicketOptions_Name_Is_Duplicate_One",
			ticketRequest:       handler.CreateTicketOptionRequestBody{Name: "Ticket", Desc: "Ticket Description", Allocation: 0},
			ticketStatusErr:     service.ErrNameIsDuplicate,
			expectedWarnMessage: handler.WarnMessageWhenNameIsDuplicated,
		},
		{
			name:                "Test_Should_Return_BadRequest_When_Client_Request_Is_Not_Json",
			ticketRequest:       handler.CreateTicketOptionRequestBody{Name: "Ticket", Desc: "Ticket Description", Allocation: 0},
			ticketStatusErr:     service.ErrAllocationIsLowerThanOne,
			expectedWarnMessage: handler.WarnMessageWhenAllocationIsBelowThanOne,
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Given
			ticketRequestBody, _ := json.Marshal(test.ticketRequest)
			req := httptest.NewRequest(http.MethodPost, "/ticket_options", bytes.NewBuffer(ticketRequestBody))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, res)

			mockService := mocks.NewMockService(gomock.NewController(t))
			mockService.
				EXPECT().
				CreateTicketOption(gomock.Any(), test.ticketRequest.Name, test.ticketRequest.Desc, test.ticketRequest.Allocation).
				Return(nil, test.ticketStatusErr).
				Times(1)

			ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

			// When
			err := ticketHandler.CreateTicketOption(c)

			assert.Nil(t, err)
			assert.Equal(t, test.expectedWarnMessage, res.Body.String())
			assert.Equal(t, http.StatusBadRequest, res.Code)
		})
	}
}

func Test_Should_Return_Status_Internal_Server_Error_When_Ticket_Option_Create(t *testing.T) {
	// Given
	requestBody := `{"name":"Ticket", "desc": "Ticket Description", "allocation": 0}`
	req := httptest.NewRequest(http.MethodPost, "/ticket_options", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	mockService := mocks.NewMockService(gomock.NewController(t))
	mockService.EXPECT().
		CreateTicketOption(gomock.Any(), "Ticket", "Ticket Description", 0).
		Return(nil, errors.New("test Error")).Times(1)

	ticketOptHandler := handler.NewDefaultTicketHandler(e, mockService)

	// When
	err := ticketOptHandler.CreateTicketOption(c)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, handler.WarnInternalServerError, rec.Body.String())
}

// GetTicket Unit Tests

func Test_Should_Return_Status_OK_When_GetTicket(t *testing.T) {
	// Given
	req := httptest.NewRequest(http.MethodGet, "/ticket/1", http.NoBody)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/ticket/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	expectedTicket := ticket.Ticket{
		ID:         1,
		Name:       "example",
		Desc:       "sample description",
		Allocation: 100,
	}

	mockService := mocks.NewMockService(gomock.NewController(t))
	mockService.EXPECT().
		GetTicket(gomock.Any(), 1).
		Return(&expectedTicket, nil).Times(1)

	ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

	// When
	err := ticketHandler.GetTicket(c)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var actualTicket ticket.Ticket
	_ = json.NewDecoder(rec.Body).Decode(&actualTicket)
	assert.Equal(t, expectedTicket, actualTicket)
}

func Test_Should_Return_Status_BadRequest_When_Get_Ticket(t *testing.T) {
	t.Run("Test_Should_Return_BadRequest_When_Invalid_id - Cannot be converted to int", func(t *testing.T) {
		// Given
		req := httptest.NewRequest(http.MethodGet, "/ticket/test", http.NoBody)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/ticket/:id")
		c.SetParamNames("id")
		c.SetParamValues("test")

		ticketHandler := handler.NewDefaultTicketHandler(e, nil)

		// When
		err := ticketHandler.GetTicket(c)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("Test_Should_Return_BadRequest_When_Invalid_id - id is lower than one", func(t *testing.T) {
		// Given
		req := httptest.NewRequest(http.MethodGet, "/ticket/0", http.NoBody)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/ticket/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")

		mockService := mocks.NewMockService(gomock.NewController(t))
		mockService.EXPECT().
			GetTicket(gomock.Any(), 0).Return(nil, service.ErrIDLowerThanOne).Times(1)

		ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

		// When
		err := ticketHandler.GetTicket(c)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, handler.WarnMessageWhenInvalidID, rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func Test_Should_Return_Status_NotFound_When_Get_Ticket(t *testing.T) {
	// Given
	req := httptest.NewRequest(http.MethodGet, "/ticket/12323", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/ticket/:id")
	c.SetParamNames("id")
	c.SetParamValues("12323")

	mockService := mocks.NewMockService(gomock.NewController(t))

	mockService.EXPECT().GetTicket(gomock.Any(), 12323).
		Return(nil, service.ErrTicketWasNotFound).Times(1)

	ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

	// When
	err := ticketHandler.GetTicket(c)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, handler.WarnMessageWhenTicketWasNotFound, rec.Body.String())
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func Test_Should_Return_Status_Internal_Server_Error_When_Ticket_Get(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ticket/1", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/ticket/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService := mocks.NewMockService(gomock.NewController(t))
	mockService.EXPECT().GetTicket(gomock.Any(), 1).
		Return(nil, errors.New("test Error")).Times(1)

	ticketHandler := handler.NewDefaultTicketHandler(e, mockService)
	err := ticketHandler.GetTicket(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, handler.WarnInternalServerError, rec.Body.String())
}

// Purchase Ticket Option Unit Tests

func Test_Should_Return_Status_Created_When_Purchase_From_Ticket_Option(t *testing.T) {
	// Given
	requestBody := `{"quantity":2,"user_id":"406c1d05-bbb2-4e94-b183-7d208c2692e1"}`
	req := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchases", bytes.NewBufferString(requestBody))
	rec := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/ticket_options/:id/purchases")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService := mocks.NewMockService(gomock.NewController(t))
	mockService.
		EXPECT().PurchaseFromTicketOption(gomock.Any(), 1, 2, "406c1d05-bbb2-4e94-b183-7d208c2692e1").
		Return(nil).Times(1)

	ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

	// When
	err := ticketHandler.PurchaseFromTicketOption(c)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func Test_Should_Return_Bad_Request_When_Purchase_From_Ticket_Option(t *testing.T) {
	// TODO: testcase e refactor edelim

	t.Run("Test_Should_Return_BadRequest_When_Invalid_id - Cannot be converted to int", func(t *testing.T) {
		// Given
		requestBody := `{"quantity":2,"user_id":"406c1d05-bbb2-4e94-b183-7d208c2692e1"}`
		req := httptest.NewRequest(http.MethodPost, "/ticket_options/test/purchases", bytes.NewBufferString(requestBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/ticket_options/:id/purchases")
		c.SetParamNames("id")
		c.SetParamValues("abc")

		ticketHandler := handler.NewDefaultTicketHandler(e, nil)

		// When
		err := ticketHandler.PurchaseFromTicketOption(c)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, handler.WarnMessageWhenInvalidID, rec.Body.String())
	})
	t.Run("Test_Should_Return_BadRequest_When_Request_Is_Not_JSON", func(t *testing.T) {
		// Given
		requestBody := "test"
		req := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchases", bytes.NewBufferString(requestBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/ticket_options/:id/purchases")
		c.SetParamNames("id")
		c.SetParamValues("1")

		ticketHandler := handler.NewDefaultTicketHandler(e, nil)

		// When
		err := ticketHandler.PurchaseFromTicketOption(c)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("Test_Should_Return_BadRequest_When_Purchase_Ticket_More_Than_Available", func(t *testing.T) {
		// Given
		requestBody := `{"quantity":1000,"user_id":"406c1d05-bbb2-4e94-b183-7d208c2692e1"}`
		req := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchases", bytes.NewBufferString(requestBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/ticket_options/:id/purchases")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockService := mocks.NewMockService(gomock.NewController(t))
		mockService.
			EXPECT().PurchaseFromTicketOption(gomock.Any(), 1, 1000, "406c1d05-bbb2-4e94-b183-7d208c2692e1").
			Return(service.ErrPurchaseTicketMoreThanAvailable).Times(1)

		ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

		// When
		err := ticketHandler.PurchaseFromTicketOption(c)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, handler.WarnMessageWhenPurchaseTicketMoreThanAvailable, rec.Body.String())
	})
	t.Run("Test_Should_Return_BadRequest_When_Quantity_Lower_Than_One", func(t *testing.T) {
		// Given
		requestBody := `{"quantity":0,"user_id":"406c1d05-bbb2-4e94-b183-7d208c2692e1"}`
		req := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchases", bytes.NewBufferString(requestBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/ticket_options/:id/purchases")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockService := mocks.NewMockService(gomock.NewController(t))
		mockService.
			EXPECT().PurchaseFromTicketOption(gomock.Any(), 1, 0, "406c1d05-bbb2-4e94-b183-7d208c2692e1").
			Return(service.ErrQuantityLowerThanOne).Times(1)

		ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

		// When
		err := ticketHandler.PurchaseFromTicketOption(c)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, handler.WarnMessageWhenQuantityLowerThanOne, rec.Body.String())
	})
	t.Run("Test_Should_Return_BadRequest_When_Id_Lower_Than_One", func(t *testing.T) {
		// Given
		requestBody := `{"quantity":1,"user_id":"406c1d05-bbb2-4e94-b183-7d208c2692e1"}`
		req := httptest.NewRequest(http.MethodPost, "/ticket_options/0/purchases", bytes.NewBufferString(requestBody))
		rec := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/ticket_options/:id/purchases")
		c.SetParamNames("id")
		c.SetParamValues("0")

		mockService := mocks.NewMockService(gomock.NewController(t))
		mockService.
			EXPECT().PurchaseFromTicketOption(gomock.Any(), 0, 1, "406c1d05-bbb2-4e94-b183-7d208c2692e1").
			Return(service.ErrIDLowerThanOne).Times(1)

		ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

		// When
		err := ticketHandler.PurchaseFromTicketOption(c)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, handler.WarnMessageWhenInvalidID, rec.Body.String())
	})
}

func Test_Should_Return_Internal_Server_Error_When_Purchase_From_Ticket_Option(t *testing.T) {
	// Given
	requestBody := `{"quantity":2,"user_id":"406c1d05-bbb2-4e94-b183-7d208c2692e1"}`
	req := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchases", bytes.NewBufferString(requestBody))
	rec := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/ticket_options/:id/purchases")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService := mocks.NewMockService(gomock.NewController(t))
	mockService.
		EXPECT().PurchaseFromTicketOption(gomock.Any(), 1, 2, "406c1d05-bbb2-4e94-b183-7d208c2692e1").
		Return(errors.New("test")).Times(1)

	ticketHandler := handler.NewDefaultTicketHandler(e, mockService)

	// When
	err := ticketHandler.PurchaseFromTicketOption(c)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, handler.WarnInternalServerError, rec.Body.String())
}
