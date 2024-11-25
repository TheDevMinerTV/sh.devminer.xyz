package web

import (
	"fmt"
)

var F = fmt.Sprintf

func getCurlScript(baseUrl, name string, execute bool) string {
	if !execute {
		return fmt.Sprintf("curl -sSL %s/%s -o %s", baseUrl, name, name)
	}

	return fmt.Sprintf("curl -sSL %s/%s | sh", baseUrl, name)
}

func getWgetScript(baseUrl, name string, execute bool) string {
	if !execute {
		return fmt.Sprintf("wget -qO %s %s/%s", name, baseUrl, name)
	}

	return fmt.Sprintf("wget -qO- %s/%s | sh", baseUrl, name)
}
