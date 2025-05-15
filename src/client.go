package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TrackProjectRequest struct {
	ProjectSlug string `json:"projectSlug"`
	ProjectType string `json:"projectType"`
	UserID      string `json:"userId"`
	ZealyUserID string `json:"zealyUserId"`
}

type TrackProjectResponse struct {
	Message string `json:"message"`
}

func TrackProjectWithProxy(proxy string, projectSlug, projectType, userID, zealyUserID string) (string, error) {
	proxyParts := strings.Split(proxy, "@")
	if len(proxyParts) != 2 {
		return "", fmt.Errorf("Error format proxy")
	}

	auth := proxyParts[0]
	hostPort := proxyParts[1]

	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", hostPort))
	if err != nil {
		return "", fmt.Errorf("Error with link proxy: %v", err)
	}

	if auth != "" {
		proxyURL.User = url.UserPassword(strings.Split(auth, ":")[0], strings.Split(auth, ":")[1])
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	payload := TrackProjectRequest{
		ProjectSlug: projectSlug,
		ProjectType: projectType,
		UserID:      userID,
		ZealyUserID: zealyUserID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Error JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://speedrun.enso.build/api/track-project-creation", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("Error request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error process request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error code status: %d", resp.StatusCode)
	}

	var response TrackProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("Error parse response: %v", err)
	}

	return response.Message, nil
}
