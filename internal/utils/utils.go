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

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("extractBodyFromResponse() -> request status code error: %s\n%s", resp.Status, string(body))
	}

	return body
}
