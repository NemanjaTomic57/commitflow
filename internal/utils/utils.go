package utils

import (
	"io"
	"log"
	"net/http"
)

func ExtractBodyFromHTTPResponse(resp *http.Response) []byte {
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("extractBodyFromHTTPResponse() -> error reading response body:", err)
	}

	return body
}
