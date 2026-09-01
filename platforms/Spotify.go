package platforms

import (
	"fmt"
	"regexp"

	"ANJALI/utils"
)

type SpotifyAPI struct {
	Regex *regexp.Regexp
}

func NewSpotifyAPI() *SpotifyAPI {
	return &SpotifyAPI{
		Regex: regexp.MustCompile(`^(https:\/\/open\.spotify\.com\/)(.*)$`),
	}
}

func (s *SpotifyAPI) Valid(link string) bool {
	return s.Regex.MatchString(link)
}

// In Go, typically you'd use a package like github.com/zmb3/spotify/v2.
// For direct conversion of your flow (fetching track name and searching YT):
func (s *SpotifyAPI) TrackFallbackAPI(url string) (map[string]string, string, error) {
	// Custom implementation using a public free API or scrape to avoid forcing OAuth implementation
	// Since you use SpotifyDown fallback in youtube.py, we can utilize Spotify API if OAuth token exists.
	// We'll simulate the search process assuming we fetched the track title "Track Name - Artist".

	// SIMULATED FETCH (Replace with actual Spotify Go Client)
	trackName := "Extracted Spotify Track"

	ytResults := utils.SearchYouTube(trackName, 1)
	if len(ytResults) == 0 {
		return nil, "", fmt.Errorf("youtube search failed")
	}

	result := ytResults[0]
	trackDetails := map[string]string{
		"title":        result.Title,
		"link":         result.Link,
		"vidid":        result.ID,
		"duration_min": result.Duration,
		"thumb":        result.Thumbnail,
	}

	return trackDetails, result.ID, nil
}
