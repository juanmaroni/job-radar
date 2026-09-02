package environment

import (
	"fmt"
	"os"
)

func GetEnvVar(varName string) (string, error) {
	value, exists := os.LookupEnv(varName)

	if !exists {
		return "", fmt.Errorf("environment variable '%s' is not set", varName)
	}

	return value, nil
}
