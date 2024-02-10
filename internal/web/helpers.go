package web

import (
	"fmt"
)

var F = fmt.Sprintf

func getCurlScript(baseUrl, name string) string {
	return fmt.Sprintf("curl -sSL %s/%s | sh", baseUrl, name)
}

func getWgetScript(baseUrl, name string) string {
	return fmt.Sprintf("wget -qO- %s/%s | sh", baseUrl, name)
}
