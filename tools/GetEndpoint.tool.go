package tools

import (
	"fmt"
	"net/http"
)

func GetEndpoint(endpoint string, expectedStatus string) error {
	httpReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return nil
}
