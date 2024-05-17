package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type City struct {
	Name      string
	Latitude  float64
	Longitude float64
}

type Suggestion struct {
	Name      string  `json:"name"`
	Latitude  string  `json:"latitude"`
	Longitude string  `json:"longitude"`
	Score     float64 `json:"score"`
}

var cities []City

func main() {
	loadCities() // Load cities data from the TSV file on startup
	http.HandleFunc("/suggestions", suggestionsHandler)
	// Start the HTTP server
	fmt.Println("Server listening on port 9090...")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func loadCities() {
	file, err := os.Open("cities_canada-usa.tsv")
	if err != nil {
		log.Fatalf("Error opening TSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'         // Set the delimiter to Tab for TSV files
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record

	// Read and discard the header line
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Error reading TSV header: %v", err)
	}

	for {
		record, err := reader.Read()
		if record == nil {
			fmt.Println("EOF")
			break
		}
		if err != nil {
			if err == csv.ErrFieldCount || err == csv.ErrBareQuote {
				continue // skip bad records
			}
			if err == csv.ErrTrailingComma {
				break // stop at EOF
			}
			log.Fatalf("Error reading TSV records: %v", err)
		}
		lat, _ := strconv.ParseFloat(record[4], 64)
		long, _ := strconv.ParseFloat(record[5], 64)
		cities = append(cities, City{
			Name:      record[1], // Name
			Latitude:  lat,       // Latitude
			Longitude: long,      // Longitude
		})
	}
}

func suggestionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	searchTerm := query.Get("q")
	latitude := query.Get("latitude")
	longitude := query.Get("longitude")

	suggestions := searchForSuggestions(searchTerm, latitude, longitude)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]Suggestion{"suggestions": suggestions})
}

func searchForSuggestions(searchTerm, latitude, longitude string) []Suggestion {
	var suggestions []Suggestion

	// This is a simplified example where we only check if the city name contains the search term.
	for _, city := range cities {
		if strings.Contains(strings.ToLower(city.Name), strings.ToLower(searchTerm)) {
			score := calculateScore(city, searchTerm, latitude, longitude)
			suggestions = append(suggestions, Suggestion{
				Name:      city.Name,
				Latitude:  strconv.FormatFloat(city.Latitude, 'f', 5, 64),
				Longitude: strconv.FormatFloat(city.Longitude, 'f', 5, 64),
				Score:     score,
			})
		}
	}

	// Sort suggestions by score in descending order
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	return suggestions
}

// calculateScore is a placeholder function to calculate the score of each city suggestion.
// Implement the logic to calculate the score based on the search term similarity and geographic proximity.
func calculateScore(city City, searchTerm, latitude, longitude string) float64 {
	// Placeholder score calculation: score is higher if the search term matches exactly.
	matchName := strings.ToLower(city.Name) == strings.ToLower(searchTerm)

	// Convert latitude and longitude strings to float64 before comparison
	lat := strconv.FormatFloat(city.Latitude, 'f', 5, 64)
	long := strconv.FormatFloat(city.Longitude, 'f', 5, 64)

	// Adjusted return statement with simplified conditional logic
	if matchName && lat == latitude && long == longitude {
		return 1.0
	} else if (matchName && lat == latitude) || (matchName && long == longitude) {
		return 0.9
	} else if matchName {
		return 0.8
	} else if lat == latitude && long == longitude {
		return 0.7
	} else if lat == latitude || long == longitude {
		return 0.6
	} else {
		return 0.5
	}
}
