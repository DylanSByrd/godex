package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dylansbyrd/godex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache pokecache.Cache
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client {
		httpClient: http.Client {
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
}

func (client* Client) RequestLocationArea(pageUrl *string) (ResourceList, error) {
	url := baseUrl + locationAreaPath
	if pageUrl != nil {
		url = *pageUrl
	}

	if cachedBytes, exists := client.cache.Get(url); exists {
		var resourceList ResourceList
		err := json.Unmarshal(cachedBytes, &resourceList)
		if err == nil {
			fmt.Println("Using cached result")
			return resourceList, nil
		} else {
			fmt.Printf("Failed to decode cached result. Making new request for %v. Error: %v\n", url, err)
		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ResourceList{}, fmt.Errorf("Error making location area request: %w", err)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return ResourceList{}, fmt.Errorf("Error requesting location area: %w", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return ResourceList{}, fmt.Errorf("Error reading location area response body: %w", err)
	}

	var resourceList ResourceList
	err = json.Unmarshal(data, &resourceList)
	if err != nil {
		return ResourceList{}, fmt.Errorf("Error unmarshalling location area response: %w", err)
	}

	client.cache.Add(url, data)
	return resourceList, nil
}
