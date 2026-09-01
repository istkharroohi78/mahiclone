package platforms

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"ANJALI/utils" // Assuming SearchYouTube is implemented here

	"github.com/PuerkitoBio/goquery"
)

type AppleAPI struct {
	BaseURL string
	Regex   *regexp.Regexp
}

func NewAppleAPI() *AppleAPI {
	return &AppleAPI{
		BaseURL: "https://music.apple.com/in/playlist/",
		Regex:   regexp.MustCompile(`^(https:\/\/music\.apple\.com\/)(.*)$`),
	}
}

func (a *AppleAPI) Valid(link string) bool {
	return a.Regex.MatchString(link)
}

func (a *AppleAPI) Track(url string, playID bool) (map[string]string, string, error) {
	if playID {
		url = a.BaseURL + url
	}

	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil, "", fmt.Errorf("failed to fetch apple music")
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", err
	}

	var search string
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if prop, _ := s.Attr("property"); prop == "og:title" {
			search, _ = s.Attr("content")
		}
	})

	if search == "" {
		return nil, "", fmt.Errorf("track not found")
	}

	// Using the SearchYouTube helper from utils
	ytResults := utils.SearchYouTube(search, 1)
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

func (a *AppleAPI) Playlist(url string, playID bool) ([]string, string, error) {
	if playID {
		url = a.BaseURL + url
	}

	parts := strings.Split(url, "playlist/")
	playlistID := ""
	if len(parts) > 1 {
		playlistID = parts[1]
	}

	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil, "", fmt.Errorf("failed to fetch playlist")
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", err
	}

	var results []string
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if prop, _ := s.Attr("property"); prop == "music:song" {
			content, _ := s.Attr("content")

			// Extract song name from Apple URL
			parts := strings.Split(content, "album/")
			if len(parts) > 1 {
				name := strings.Split(parts[1], "/")[0]
				name = strings.ReplaceAll(name, "-", " ")
				results = append(results, name)
			}
		}
	})

	return results, playlistID, nil
}
