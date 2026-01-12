package envi

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func InitEnvi() error {
	filepath := "config/envi/.env"
	if err := godotenv.Load(filepath); err != nil {
		return err
	}
	fmt.Println("✅ Load Part:", filepath)
	return nil
}

func GetEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
