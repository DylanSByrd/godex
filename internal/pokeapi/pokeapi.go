package pokeapi

const (
	baseUrl = "https://pokeapi.co/api/v2"
	locationAreaPath = "/location-area"
)

type ResourceList struct {
	Count int
	Next *string
	Previous *string
	Results []struct {
		Name string
		Url string
	}
}

