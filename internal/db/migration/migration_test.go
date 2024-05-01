package migration

import (
	"context"
	"database/sql"
	"testing"

	_ "embed"

	"github.com/kellswork/wayfarer/internal/utils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_migration(t *testing.T) {
	// err := godotenv.Load("../../../.env")
	// require.NoError(t, err)
	// dbUrl := os.Getenv("TEST_DB_URL")

	ctx := context.Background()
	container, err := utils.CreatePostgresContainer(ctx)
	require.NoError(t, err)
	db, err := sql.Open("postgres", container.ConnectionString)
	require.NoError(t, err)

	m := New(db)
	err = m.Up()
	assert.NoError(t, err)

	// err = m.Down()
	// assert.NoError(t, err)
}
