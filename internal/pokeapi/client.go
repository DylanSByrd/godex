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

func fetchFromCache[T any](cache* pokecache.Cache, key string) (T, bool) {
	var result T
	if cachedBytes, exists := cache.Get(key); exists {
		err := json.Unmarshal(cachedBytes, &result)
		if err != nil {
			fmt.Printf("Failed to decode cached %T result for %s. Error: %v", result, key, err)
			return result, false
		}

		fmt.Printf("Found cached response\n")
		return result, true
	}

	return result, false
}

func getResourceAs[T any](client *Client, url string) (T, error) {
	var resource T

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resource, fmt.Errorf("Error creating request: %w", err)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return resource, fmt.Errorf("Error received as response: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		return resource, fmt.Errorf("Request failed with status code %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return resource, fmt.Errorf("Error reading response body: %w", err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return resource, fmt.Errorf("Error unmarshalling response: %w", err)
	}

	client.cache.Add(url, data)
	return resource, nil
}

func (client* Client) RequestLocationArea(pageUrl *string) (ResourceList, error) {
	url := baseUrl + locationAreaPath
	if pageUrl != nil {
		url = *pageUrl
	}

	if cachedEntry, exists := fetchFromCache[ResourceList](&client.cache, url); exists {
		return cachedEntry, nil
	}

	resourceList, err := getResourceAs[ResourceList](client, url)
	return resourceList, err
}

func (client* Client) RequestLocationAreaDetails(areaName string) (LocationArea, error) {
	url := baseUrl + locationAreaPath + "/" + areaName

	if cachedEntry, exists := fetchFromCache[LocationArea](&client.cache, url); exists {
		return cachedEntry, nil
	}

	locationArea, err := getResourceAs[LocationArea](client, url)
	return locationArea, err
}

func (client* Client) RequestPokemonDetails(pokemonName string) (PokemonDetails, error) {
	url := baseUrl + pokemonPath + "/" + pokemonName

	if cachedEntry, exists := fetchFromCache[PokemonDetails](&client.cache, url); exists {
		return cachedEntry, nil
	}

	pokemonDetails, err := getResourceAs[PokemonDetails](client, url)
	return pokemonDetails, err
}
