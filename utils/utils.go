package utils

import (
	"os"
	"strconv"
)

func TryGetEnvString(key string, defVal string) string {
	strVal := os.Getenv(key)
	if strVal == "" {
		strVal = defVal
	}

	return strVal
}

func TryGetEnvInt64(key string, defVal int64) int64 {
	strVal := os.Getenv(key)
	if strVal == "" {
		return defVal
	}

	intVal, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil || intVal <= 0 {
		return defVal
	}

	return intVal
}
