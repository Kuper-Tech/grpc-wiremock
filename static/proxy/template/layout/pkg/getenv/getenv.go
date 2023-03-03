package getenv

import (
	"os"
	"strconv"
)

func GetPort() string {
	rawPort := os.Getenv("GRPC_TO_HTTP_PROXY_PORT")

	_, err := strconv.ParseInt(rawPort, 10, 0)
	if err != nil {
		return defaultPort
	}

	return rawPort
}
