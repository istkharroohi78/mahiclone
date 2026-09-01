package platforms

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"ANJALI/utils"

	"github.com/PuerkitoBio/goquery"
)

type RessoAPI struct {
	BaseURL string
	Regex   *regexp.Regexp
}

func NewRessoAPI() *RessoAPI {
	return &RessoAPI{
		BaseURL: "https://m.resso.com/",
		Regex:   regexp.MustCompile(`^(https:\/\/m\.resso\.com\/)(.*)$`),
	}
}

func (r *RessoAPI) Valid(link string) bool {
	return r.Regex.MatchString(link)
}

func (r *RessoAPI) Track(url string, playID bool) (map[string]string, string, error) {
	if playID {
		url = r.BaseURL + url
	}

	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil, "", fmt.Errorf("failed to fetch resso")
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", err
	}

	var title, desc string
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if prop, _ := s.Attr("property"); prop == "og:title" {
			title, _ = s.Attr("content")
		}
		if prop, _ := s.Attr("property"); prop == "og:description" {
			desc, _ = s.Attr("content")
			desc = strings.Split(desc, "·")[0]
		}
	})

	if desc == "" {
		return nil, "", fmt.Errorf("description empty")
	}

	ytResults := utils.SearchYouTube(title, 1)
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
