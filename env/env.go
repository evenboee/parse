package env

import (
	"os"

	"github.com/evenboee/parse"
)

func GetString(key string, def ...string) string {
	return getEnv(key, def)
}

func Get[T any](key string, def ...string) T {
	return parse.Must[T](getEnv(key, def))
}

func ShouldGet[T any](key string, def ...string) (T, error) {
	return parse.Try[T](getEnv(key, def))
}

func getEnv(key string, def []string) string {
	val := os.Getenv(key)
	if val == "" && len(def) != 0 {
		val = def[0]
	}

	return val
}
