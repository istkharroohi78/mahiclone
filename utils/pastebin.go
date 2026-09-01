package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const BaseBatbin = "https://batbin.me/"

// Post handles HTTP POST requests
func Post(url string, data string) (map[string]interface{}, error) {
	reqBody := []byte(data)
	resp, err := http.Post(url, "text/plain", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// Fallback to text parsing if API changes
		return map[string]interface{}{"message": string(body), "success": true}, nil
	}

	return result, nil
}

// AnjaliBin pastes text to batbin and returns the URL
func AnjaliBin(text string) (string, error) {
	resp, err := Post(BaseBatbin+"api/v2/paste", text)
	if err != nil {
		return "", err
	}

	if success, ok := resp["success"].(bool); ok && success {
		if msg, ok := resp["message"].(string); ok {
			return BaseBatbin + msg, nil
		}
	}
	return "", nil
}
