package db

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/BurntSushi/toml"
)

// writeConfigWithBackup writes the configuration back to the file, with backup handling
func writeConfigWithBackup(filename string, config *DBConfig) error {
	backupFilename := filename + ".bak"

	// Check if the file exists
	if _, err := os.Stat(filename); err == nil {
		// Create a backup of the existing file
		if err := os.Rename(filename, backupFilename); err != nil {
			return fmt.Errorf("error creating backup file: %v", err)
		}
	}

	// Write the new config file
	if err := writeConfig(filename, config); err != nil {
		// If writing fails, restore the backup
		if _, err := os.Stat(backupFilename); err == nil {
			if err := os.Rename(backupFilename, filename); err != nil {
				return fmt.Errorf("error restoring backup file: %v", err)
			}
		}
		return fmt.Errorf("error writing config file: %v", err)
	}

	// If write succeeds, remove the backup file
	if _, err := os.Stat(backupFilename); err == nil {
		if err := os.Remove(backupFilename); err != nil {
			return fmt.Errorf("error removing backup file: %v", err)
		}
	}

	return nil
}

// writeConfig writes the configuration back to the file
func writeConfig(filename string, config *DBConfig) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating config file: %v", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error encoding config file: %v", err)
	}

	return nil
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefTUVWXYZ01ghituvwxyzABCDEFGHI!JKLMNOPQRS234567jklmnopqrs89"
	result := make([]byte, length)
	for i := range result {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomInt.Int64()]
	}
	return string(result), nil
}
