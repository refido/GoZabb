package helpers

import (
	"io"
	"os"
)

// EnsureEnvFile copies .env.example → .env if .env does not exist.
func EnsureEnvFile() error {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		src, err := os.Open(".env.example")
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(".env")
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		return err
	}
	return nil
}
