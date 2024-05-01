package repositories

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kellswork/wayfarer/internal/db"
	"github.com/kellswork/wayfarer/internal/db/migration"
	"github.com/kellswork/wayfarer/internal/db/models"
	"github.com/kellswork/wayfarer/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
	suite.Suite
	pgContiner *utils.PostgresContainer
	userRepo   *userRespository
	ctx        context.Context
	dbStore    *sql.DB
}

func (suite *UserRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	container, err := utils.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContiner = container

	dbStore, err := db.ConnectDatabase(suite.pgContiner.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	suite.dbStore = dbStore

	migrator := migration.New(dbStore)
	err = migrator.Up()
	if err != nil {
		log.Fatal(err)
	}

	userReporistory := newUserRepository(dbStore)
	if err != nil {
		log.Fatal(err)
	}
	suite.userRepo = userReporistory
}

func (suite *UserRepoTestSuite) TearDownSuite() {
	if err := suite.pgContiner.Terminate(suite.ctx); err != nil {
		log.Fatalf(" failed to terminate postgres container: %s", err)
	}
}

func (suite *UserRepoTestSuite) Test_CreateUser() {
	t := suite.T()
	// 1. setup pahse
	sampleUser := models.User{
		ID:        uuid.NewString(),
		Email:     "kellasw@gmail.com",
		FirstName: "Kelechi",
		LastName:  "Ogbonna",
		Password:  "1234",
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.userRepo.Create(suite.ctx, &sampleUser)
	assert.NoError(t, err)

	resp, err := suite.userRepo.GetByEmail(suite.ctx, sampleUser.Email)
	assert.NoError(t, err)

	assert.Equal(t, sampleUser.ID, resp.ID)
	assert.Equal(t, sampleUser.Email, resp.Email)
	assert.Equal(t, sampleUser.FirstName, resp.FirstName)
	assert.Equal(t, sampleUser.LastName, resp.LastName)
	assert.Equal(t, sampleUser.Password, resp.Password)
	assert.Equal(t, sampleUser.IsAdmin, resp.IsAdmin)

	err = deleteUser(suite.dbStore, sampleUser.ID)
	require.NoError(t, err)
}

func (suite *UserRepoTestSuite) Test_EmailExist() {
	t := suite.T()
	// 1. setup pahse
	// 1.1 setup db connectiom
	// err := godotenv.Load("../../../.env")
	// require.NoError(t, err)
	// testDBUrl := os.Getenv("TEST_DB_URL")
	// ctx := context.Background()
	// dbStore, err := db.ConnectDatabase(testDBUrl)
	// require.NoError(t, err)

	// defer db.CloseDatabase(dbStore)

	// migrator := migration.New(dbStore)
	// err = migrator.Up()
	// assert.NoError(t, err)

	sampleUser := models.User{
		ID:        uuid.NewString(),
		Email:     "kellasw@gmail.com",
		FirstName: "Kelechi",
		LastName:  "Ogbonna",
		Password:  "1234",
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// userRepo := newUserRepository(dbStore)
	err := suite.userRepo.Create(suite.ctx, &sampleUser)
	assert.NoError(t, err)

	resp, err := suite.userRepo.EmailExists(suite.ctx, sampleUser.Email)
	assert.NoError(t, err)

	assert.Equal(t, resp, true)

	err = deleteUser(suite.dbStore, sampleUser.ID)
	require.NoError(t, err)
}

func deleteUser(db *sql.DB, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
