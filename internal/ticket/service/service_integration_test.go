package service_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	ticket2 "github.com/dilaragorum/ticket-api/internal/ticket"
	"github.com/dilaragorum/ticket-api/internal/ticket/database"
	"github.com/dilaragorum/ticket-api/internal/ticket/repository"
	"github.com/dilaragorum/ticket-api/internal/ticket/service"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	svc            *service.DefaultService
	container      *dockertest.Resource
	connectionPool *gorm.DB
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.container, suite.connectionPool = createContainer()
	defaultRepository := repository.NewDefaultRepository(suite.connectionPool)
	suite.svc = service.NewDefaultService(defaultRepository)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	_ = suite.container.Close()
}

func (suite *IntegrationTestSuite) Test_Should_Insert_New_Ticket() {
	// When
	option, err := suite.svc.CreateTicketOption(context.TODO(), "ticket", "description", 100)

	// Then
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, option.ID)
	assert.Equal(suite.T(), "ticket", option.Name)
	assert.Equal(suite.T(), "description", option.Desc)
	assert.Equal(suite.T(), 100, option.Allocation)
}

func (suite *IntegrationTestSuite) Test_Should_Get_Ticket_With_ID() {
	// Given
	ticket := ticket2.Ticket{
		Name:       "example2",
		Desc:       "sample description2",
		Allocation: 100,
	}

	err := suite.connectionPool.Model(&ticket).Create(&ticket).Error
	if err != nil {
		suite.T().Error(err)
	}

	// When
	option, err := suite.svc.GetTicket(context.TODO(), ticket.ID)

	// Then
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, option.ID)
	assert.Equal(suite.T(), "example2", option.Name)
	assert.Equal(suite.T(), "sample description2", option.Desc)
	assert.Equal(suite.T(), 100, option.Allocation)
}

func (suite *IntegrationTestSuite) Test_Should_Purchase_From_Ticket() {
	// Given
	ticket := ticket2.Ticket{
		Name:       "example3",
		Desc:       "sample description3",
		Allocation: 100,
	}

	err := suite.connectionPool.Model(&ticket).Create(&ticket).Error
	if err != nil {
		suite.T().Error(err)
	}

	// When
	if err := suite.svc.PurchaseFromTicketOption(context.TODO(), ticket.ID, 50, "406c1d05-bbb2-4e94-b183-7d208c2692e1"); err != nil {
		suite.T().Error(err)
	}

	// Then
	assert.Nil(suite.T(), err)
}

func createContainer() (*dockertest.Resource, *gorm.DB) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	pool.MaxWait = 1 * time.Minute

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.1",
		Env: []string{
			"POSTGRES_DB=ticket_app",
			"POSTGRES_USER=ticket_user",
			"POSTGRES_PASSWORD=postgres",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	port := container.GetPort("5432/tcp")
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", port)
	os.Setenv("POSTGRES_USER", "ticket_user")
	os.Setenv("POSTGRES_PASSWORD", "postgres")
	os.Setenv("POSTGRES_DB", "ticket_app")

	var connectionPool *gorm.DB

	pool.Retry(func() error { //nolint:errcheck
		connectionPool, err = database.Setup()
		if err != nil {
			return err
		}

		conn, err := connectionPool.DB()
		if err != nil {
			return err
		}

		err = conn.Ping()
		if err != nil {
			return err
		}

		database.Migrate()

		return nil
	})

	return container, connectionPool
}
