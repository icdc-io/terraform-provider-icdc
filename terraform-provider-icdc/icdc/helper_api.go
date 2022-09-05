package icdc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func requestApi(method, url string, body io.Reader) (*json.Decoder, error) {

	client := &http.Client{Timeout: 100 * time.Second}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", os.Getenv("API_GATEWAY"), url), body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AUTH_TOKEN")))
	req.Header.Set("X_MIQ_GROUP", fmt.Sprintf("%s.%s", os.Getenv("ACCOUNT"), os.Getenv("ROLE")))

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	decodedBody := json.NewDecoder(r.Body)

	/* ahrechushkin: so, we need to close the body, but we can't do it here, because we need to return the body to the caller
	defer r.Body.Close()
	*/
	return decodedBody, nil
}
