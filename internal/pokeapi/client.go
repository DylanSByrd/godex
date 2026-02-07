package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client {
		httpClient: http.Client {
			Timeout: timeout,
		},
	}
}

func (c* Client) RequestLocationArea(pageUrl *string) (ResourceList, error) {
	url := baseUrl + locationAreaPath
	if pageUrl != nil {
		url = *pageUrl
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ResourceList{}, fmt.Errorf("Error making location area request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return ResourceList{}, fmt.Errorf("Error requesting location area: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ResourceList{}, fmt.Errorf("Error reading location area response body: %w", err)
	}

	var resourceList ResourceList
	err = json.Unmarshal(body, &resourceList)
	if err != nil {
		return ResourceList{}, fmt.Errorf("Error unmarshalling location area response: %w", err)
	}

	return resourceList, nil
}
