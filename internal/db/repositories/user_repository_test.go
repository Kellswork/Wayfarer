package repositories

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kellswork/wayfarer/internal/db"
	"github.com/kellswork/wayfarer/internal/db/migration"
	"github.com/kellswork/wayfarer/internal/db/models"
	"github.com/kellswork/wayfarer/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateUser(t *testing.T) {
	// 1. setup pahse
	// 1.1 setup db connectiom
	// err := godotenv.Load("../../../.env")
	// require.NoError(t, err)
	// testDBUrl := os.Getenv("TEST_DB_URL")
	ctx := context.Background()

	container, err := utils.CreatePostgresContainer(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf(" failed to terminate container: %s", err)
		}
	})

	dbStore, err := db.ConnectDatabase(container.ConnectionString)
	require.NoError(t, err)

	migrator := migration.New(dbStore)
	err = migrator.Up()
	assert.NoError(t, err)

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

	userRepo := newUserRepository(dbStore)
	err = userRepo.Create(ctx, &sampleUser)
	assert.NoError(t, err)

	resp, err := userRepo.GetByEmail(ctx, sampleUser.Email)
	assert.NoError(t, err)

	assert.Equal(t, sampleUser.ID, resp.ID)
	assert.Equal(t, sampleUser.Email, resp.Email)
	assert.Equal(t, sampleUser.FirstName, resp.FirstName)
	assert.Equal(t, sampleUser.LastName, resp.LastName)
	assert.Equal(t, sampleUser.Password, resp.Password)
	assert.Equal(t, sampleUser.IsAdmin, resp.IsAdmin)

	err = deleteUser(dbStore, sampleUser.ID)
	require.NoError(t, err)
}

func Test_EmailExist(t *testing.T) {
	// 1. setup pahse
	// 1.1 setup db connectiom
	err := godotenv.Load("../../../.env")
	require.NoError(t, err)
	testDBUrl := os.Getenv("TEST_DB_URL")
	ctx := context.Background()
	dbStore, err := db.ConnectDatabase(testDBUrl)
	require.NoError(t, err)

	defer db.CloseDatabase(dbStore)

	migrator := migration.New(dbStore)
	err = migrator.Up()
	assert.NoError(t, err)

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

	userRepo := newUserRepository(dbStore)
	err = userRepo.Create(ctx, &sampleUser)
	assert.NoError(t, err)

	resp, err := userRepo.EmailExists(ctx, sampleUser.Email)
	assert.NoError(t, err)

	assert.Equal(t, resp, true)

	err = deleteUser(dbStore, sampleUser.ID)
	require.NoError(t, err)
}

func deleteUser(db *sql.DB, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}
