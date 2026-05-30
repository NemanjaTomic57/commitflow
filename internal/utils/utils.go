package utils

import (
	"io"
	"log"
	"net/http"
)

func ExtractBodyFromResponse(resp *http.Response) []byte {
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("extractBodyFromResponse() -> error reading response body:", err)
	}

	return body
}
