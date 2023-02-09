package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dilaragorum/ticket-api/internal/ticket"
	"github.com/dilaragorum/ticket-api/internal/ticket/mocks"
	"github.com/dilaragorum/ticket-api/internal/ticket/repository"
	"github.com/dilaragorum/ticket-api/internal/ticket/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Create Ticket Option Unit Tests
func Test_Should_Return_Successfully_Create_When_TicketOption_Is_Valid(t *testing.T) {
	// Given
	ticketOption := ticket.Ticket{ID: 1, Name: "example", Desc: "sample description", Allocation: 100}
	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockRepository.
		EXPECT().CreateTicketOption(gomock.Any(), "example", "sample description", 100).
		Return(&ticketOption, nil).Times(1)

	ticketOptService := service.NewDefaultService(mockRepository)

	// When
	actualTicketOption, err := ticketOptService.CreateTicketOption(context.TODO(), "example", "sample description", 100)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, ticketOption, *actualTicketOption)
}

func Test_Should_Return_Error_When_Creating_TicketOption_Is_Not_Valid(t *testing.T) {
	type testCase struct {
		testName                  string
		ticketName                string
		ticketDescription         string
		ticketAllocation          int
		mockRepositoryTimes       int
		mockRepositoryErr         error
		expectedCreatingStatusErr error
	}

	dbWrapperErr := errors.New("test DBWrapper Error")
	testCases := []testCase{
		{
			testName:                  "Test_Should_Return_Error_When_TicketOption_Name_Is_Empty",
			ticketName:                "",
			ticketDescription:         "sample description",
			ticketAllocation:          100,
			mockRepositoryTimes:       0,
			mockRepositoryErr:         nil,
			expectedCreatingStatusErr: service.ErrNameIsEmpty,
		},
		{
			testName:                  "Test_Should_Return_Error_When_TicketOption_Description_Is_Empty",
			ticketName:                "example",
			ticketDescription:         "",
			ticketAllocation:          100,
			mockRepositoryTimes:       0,
			mockRepositoryErr:         nil,
			expectedCreatingStatusErr: service.ErrDescriptionIsEmpty,
		},
		{
			testName:                  "Test_Should_Return_Error_When_Allocation_Is_Lower_Than_One",
			ticketName:                "example",
			ticketDescription:         "sample description",
			ticketAllocation:          0,
			mockRepositoryTimes:       0,
			mockRepositoryErr:         nil,
			expectedCreatingStatusErr: service.ErrAllocationIsLowerThanOne,
		},
		{
			testName:                  "Test_Should_Return_Database_Error",
			ticketName:                "example",
			ticketDescription:         "sample description",
			ticketAllocation:          100,
			mockRepositoryTimes:       1,
			mockRepositoryErr:         dbWrapperErr,
			expectedCreatingStatusErr: dbWrapperErr,
		},
		{
			testName:                  "Test_Should_Return_Err_Duplicated_Ticket_Name",
			ticketName:                "example",
			ticketDescription:         "sample description",
			ticketAllocation:          100,
			mockRepositoryTimes:       1,
			mockRepositoryErr:         repository.ErrDBDuplicatedTicketName,
			expectedCreatingStatusErr: service.ErrNameIsDuplicate,
		},
	}
	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			// Given
			mockRepository := mocks.NewMockRepository(gomock.NewController(t))
			mockRepository.EXPECT().
				CreateTicketOption(gomock.Any(), test.ticketName, test.ticketDescription, test.ticketAllocation).
				Return(nil, test.mockRepositoryErr).Times(test.mockRepositoryTimes)

			svc := service.NewDefaultService(mockRepository)

			// When
			option, err := svc.CreateTicketOption(context.TODO(), test.ticketName, test.ticketDescription, test.ticketAllocation)

			// Then
			assert.Equal(t, test.expectedCreatingStatusErr, err)
			assert.Nil(t, option)
		})
	}
}

// Get Ticket Unit Tests
func Test_Should_Return_Success_When_Get_TicketOption(t *testing.T) {
	// Given
	expectedTicket := ticket.Ticket{
		ID:         1,
		Name:       "Example",
		Desc:       "Sample Description",
		Allocation: 100,
	}

	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockRepository.EXPECT().GetTicket(gomock.Any(), 1).Return(&expectedTicket, nil).Times(1)

	ticketService := service.NewDefaultService(mockRepository)

	// When
	actualTicket, err := ticketService.GetTicket(context.TODO(), 1)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expectedTicket, *actualTicket)
}

func Test_Should_Return_Error_Cases_When_Get_TicketOption(t *testing.T) {
	t.Run("Test_Should_Return_Error_When_Id_Is_Lower_Than_One", func(t *testing.T) {
		// Given
		ticketService := service.NewDefaultService(nil)

		// When
		getTicket, err := ticketService.GetTicket(context.TODO(), 0)

		// Then
		assert.Equal(t, err, service.ErrIDLowerThanOne)
		assert.Nil(t, getTicket)
	})
	t.Run("Test_Should_Return_Error_When_Ticket_Was_Not_Found", func(t *testing.T) {
		// Given
		mockRepository := mocks.NewMockRepository(gomock.NewController(t))
		mockRepository.EXPECT().GetTicket(gomock.Any(), 864).Return(nil, repository.ErrDBTicketNotFound).Times(1)

		ticketService := service.NewDefaultService(mockRepository)

		// When
		getTicket, err := ticketService.GetTicket(context.TODO(), 864)

		// Then
		assert.Equal(t, service.ErrTicketWasNotFound, err)
		assert.Nil(t, getTicket)
	})
	t.Run("Test_Should_Return_Error_When_Other_DB_Problems", func(t *testing.T) {
		// Given
		mockRepository := mocks.NewMockRepository(gomock.NewController(t))
		mockRepository.
			EXPECT().GetTicket(gomock.Any(), 1).Return(nil, errors.New("test")).Times(1)

		ticketService := service.NewDefaultService(mockRepository)

		// When
		getTicket, err := ticketService.GetTicket(context.TODO(), 1)

		// Then
		assert.Error(t, err)
		assert.Nil(t, getTicket)
	})
}

// Purchase Ticket Unit Tests
func Test_Should_Return_Success_When_User_Can_Purchase_Specified_Ticket(t *testing.T) {
	// Given
	expectedTicket := ticket.Ticket{
		ID:         1,
		Name:       "Example",
		Desc:       "Sample Description",
		Allocation: 100,
	}

	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockRepository.EXPECT().GetTicket(gomock.Any(), 1).Return(&expectedTicket, nil).Times(1)
	mockRepository.
		EXPECT().PurchaseFromTicketOption(gomock.Any(), expectedTicket.ID, 20, "406c1d05-bbb2-4e94-b183-7d208c2692e1").
		Return(nil).Times(1)

	ticketService := service.NewDefaultService(mockRepository)

	// When
	err := ticketService.PurchaseFromTicketOption(context.TODO(), 1, 20, "406c1d05-bbb2-4e94-b183-7d208c2692e1")

	// Then
	assert.Nil(t, err)
}

func Test_Should_Return_Error_When_User_Want_To_Purchase_Specified_Ticket(t *testing.T) {
	t.Run("Test_Should_Return_Err_Quantity_Lower_Than_One_When_Quantity_Lower_Than_One", func(t *testing.T) {
		// Given
		ticketService := service.NewDefaultService(nil)

		// When
		err := ticketService.PurchaseFromTicketOption(context.TODO(), 1, 0, "406c1d05-bbb2-4e94-b183-7d208c2692e1")

		// Then
		assert.Equal(t, service.ErrQuantityLowerThanOne, err)
	})

	t.Run("Test_Should_Return_Err_When_Get_Ticket_Option", func(t *testing.T) {
		// Given
		mockRepository := mocks.NewMockRepository(gomock.NewController(t))
		mockRepository.EXPECT().GetTicket(gomock.Any(), 1999).Return(nil, errors.New("test")).Times(1)
		ticketService := service.NewDefaultService(mockRepository)

		// When
		err := ticketService.PurchaseFromTicketOption(context.TODO(), 1999, 100, "userId")

		// Then
		assert.Error(t, err)
	})

	t.Run("Test_Should_Return_Err_Purchase_Ticket_More_Than_Available", func(t *testing.T) {
		// Given
		ticketTest := ticket.Ticket{
			ID:         1,
			Name:       "sample",
			Desc:       "example desc",
			Allocation: 50,
		}

		mockRepository := mocks.NewMockRepository(gomock.NewController(t))
		mockRepository.EXPECT().GetTicket(gomock.Any(), 1).Return(&ticketTest, nil).Times(1)

		ticketService := service.NewDefaultService(mockRepository)

		// When
		err := ticketService.PurchaseFromTicketOption(context.TODO(), 1, 100, "test")

		// Then
		assert.Equal(t, service.ErrPurchaseTicketMoreThanAvailable, err)
	})

	t.Run("Test_Should_Return_Error_When_User_Cannot_Purchase_Specified_Ticket", func(t *testing.T) {
		getTicketResponse := ticket.Ticket{
			ID:         1,
			Name:       "sample getTicketResponse",
			Desc:       "example getTicketResponse description",
			Allocation: 100,
		}

		mockRepository := mocks.NewMockRepository(gomock.NewController(t))
		mockRepository.EXPECT().GetTicket(gomock.Any(), 1).Return(&getTicketResponse, nil).Times(1)
		mockRepository.EXPECT().PurchaseFromTicketOption(gomock.Any(), 1, 50, "test").Return(errors.New("test")).Times(1)

		defaultService := service.NewDefaultService(mockRepository)
		err := defaultService.PurchaseFromTicketOption(context.TODO(), 1, 50, "test")

		assert.Error(t, err)
	})

	t.Run("Test_Should_Return_Error_When_Get_Ticket_Option_With_Specified_ID", func(t *testing.T) {
		mockRepository := mocks.NewMockRepository(gomock.NewController(t))
		mockRepository.EXPECT().GetTicket(gomock.Any(), 1).Return(nil, errors.New("test")).Times(1)

		defaultService := service.NewDefaultService(mockRepository)
		err := defaultService.PurchaseFromTicketOption(context.TODO(), 1, 20, "testUserId")

		assert.Error(t, err)
	})
}
