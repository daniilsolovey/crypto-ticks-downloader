package operator

import (
	"testing"

	"github.com/daniilsolovey/crypto-ticks-downloader/internal/config"
	"github.com/daniilsolovey/crypto-ticks-downloader/internal/database"
	"github.com/reconquest/pkg/log"
)

const (
	CONFIG_PATH = "../../test/testdata/config.yaml"
)

func createDatabase() *database.Database {
	config, err := config.Load(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	database := database.NewDatabase(
		config.Database.Name, config.Database.User, config.Database.Password,
	)

	return database
}

func TestOperator_StartDatabase(
	t *testing.T,
) {
	config, err := config.Load(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	testDatabase := createDatabase()
	// err = testDatabase.CreateDatabase(config.Database.Name)
	// assert.NoError(t, err)
	log.Warning("testDatabase ", testDatabase)
	log.Warning("config ", config)
	defer testDatabase.Close()
	// defer testDatabase.Drop(config.Database.Name)
}
