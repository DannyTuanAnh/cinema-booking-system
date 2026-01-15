package key

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	repo "cinema.com/demo/bff/repository"
	"cinema.com/demo/bff/utils"
	pkgFile "cinema.com/demo/pkg/file"
)

func GenerateAPIKey(path string, clientType string, maxReq int, winSec int, database *sql.DB) error {
	apiRepo := repo.NewApiKeyRepo(database)

	exits, _ := apiRepo.ExistsByClient(clientType)
	if !exits {

		plaintext := utils.GenerateApiKey(clientType)
		hash := utils.HashApiKey(plaintext)

		if err := apiRepo.Insert(clientType, hash, maxReq, winSec); err != nil {
			return err
		}
		if clientType == "web" {

			// append to .env
			err := pkgFile.WriteEnv(path, "VITE_API_KEY", plaintext)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("API key created successfully for client: %s", clientType)

			return nil

		} else {

			// Append to .env
			keyEnv := fmt.Sprintf("API_KEY_%s", strings.ToUpper(clientType))

			err := pkgFile.WriteEnv(path, keyEnv, plaintext)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("API key created successfully for client: %s", clientType)

			return nil

		}
	}

	return fmt.Errorf("API key already exists for client: %s", clientType)
}
