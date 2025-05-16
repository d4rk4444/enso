package ai

import (
	"encoding/json"
	"enso/src"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type VerifyResponse struct {
	Body    string
	Cookies []*http.Cookie
	Token   string
}

func GetNonce(proxy string) string {
	getResp, err := src.MakeRequest(src.RequestOptions{
		Method:   "GET",
		URL:      "https://enso.brianknows.org/api/auth/nonce",
		ProxyURL: proxy,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		log.Fatalf("Error GET request: %v", err)
	}

	return string(getResp.Body)
}

func Verify(message src.StructMessage, signature string, proxy string) (*VerifyResponse, error) {
	requestBody := map[string]interface{}{
		"message":   message,
		"signature": signature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error decode JSON: %v", err)
	}

	postResp, err := src.MakeRequest(src.RequestOptions{
		Method: "POST",
		URL:    "https://enso.brianknows.org/api/auth/verify",
		Headers: map[string]string{
			"Origin":       "https://enso.brianknows.org",
			"Content-Type": "application/json",
			"Priority":     "u=1, i",
			"Referer":      "https://enso.brianknows.org/builds",
		},
		Body:     jsonData,
		ProxyURL: proxy,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("error executing POST request: %v", err)
	}

	// Create result
	result := &VerifyResponse{
		Body:    string(postResp.Body),
		Cookies: postResp.Cookies,
	}

	// Look for token in cookies
	for _, cookie := range postResp.Cookies {
		// Look for possible cookie names containing token
		if cookie.Name == "brian_token" ||
			cookie.Name == "auth_token" ||
			cookie.Name == "token" ||
			cookie.Name == "session" ||
			cookie.Name == "connect.sid" ||
			cookie.Name == "jwt" ||
			strings.Contains(cookie.Name, "token") ||
			strings.Contains(cookie.Name, "auth") {
			result.Token = cookie.Value
			break
		}
	}

	return result, nil
}

func Search(query string, kbId string, userId string, proxy string) string {
	requestBody := map[string]interface{}{
		"query": query,
		"kbId":  kbId,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Sprintf("Error code JSON: %v", err)
	}

	postResp, err := src.MakeRequest(src.RequestOptions{
		Method: "POST",
		URL:    "https://enso.brianknows.org/api/search",
		Headers: map[string]string{
			"accept":          "*/*",
			"accept-language": "en-US,en;q=0.9,ru;q=0.8",
			"content-type":    "text/plain;charset=UTF-8",
			"priority":        "u=1, i",
			"x-enso-user-id":  userId,
			"referer":         "https://enso.brianknows.org/search",
		},
		Body:     jsonData,
		ProxyURL: proxy,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		return fmt.Sprintf("Error POST request: %v", err)
	}

	return string(postResp.Body)
}

func Build(query string, chain int, proxy string, token string) string {
	requestBody := map[string]interface{}{
		"query": query,
		"chain": chain,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Sprintf("Error code JSON: %v", err)
	}

	postResp, err := src.MakeRequest(src.RequestOptions{
		Method: "POST",
		URL:    "https://enso.brianknows.org/api/builds",
		Headers: map[string]string{
			"accept":       "*/*",
			"content-type": "application/json",
			"cookie":       "brian-token=" + token,
		},
		Body:     jsonData,
		ProxyURL: proxy,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		return fmt.Sprintf("Error POST request: %v", err)
	}

	return string(postResp.Body)
}

func Points(txHash string, action string, chainId int, userId string, proxy string, token string) string {
	requestBody := map[string]interface{}{
		"txHash":  txHash,
		"action":  action,
		"chainId": chainId,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Sprintf("Error code JSON: %v", err)
	}

	postResp, err := src.MakeRequest(src.RequestOptions{
		Method: "POST",
		URL:    "https://enso.brianknows.org/api/points",
		Headers: map[string]string{
			"accept":         "*/*",
			"content-type":   "text/plain;charset=UTF-8",
			"x-enso-user-id": userId,
			"cookie":         "brian-token=" + token,
		},
		Body:     jsonData,
		ProxyURL: proxy,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		return fmt.Sprintf("Error POST request: %v", err)
	}

	return string(postResp.Body)
}

func TrackProject(projectSlug string, projectType string, userID string, zealyUserID string, proxy string) string {
	requestBody := map[string]interface{}{
		"projectSlug": projectSlug,
		"projectType": projectType,
		"userId":      userID,
		"zealyUserId": zealyUserID,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Sprintf("Error code JSON: %v", err)
	}

	postResp, err := src.MakeRequest(src.RequestOptions{
		Method: "POST",
		URL:    "https://speedrun.enso.build/api/track-project-creation",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:     jsonData,
		ProxyURL: proxy,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		return fmt.Sprintf("Error POST request: %v", err)
	}

	if postResp.StatusCode != 200 {
		return fmt.Sprintf("Error code status: %d", postResp.StatusCode)
	}

	var response struct {
		Message string `json:"message"`
	}

	if err := json.Unmarshal(postResp.Body, &response); err != nil {
		return fmt.Sprintf("Error parse response: %v", err)
	}

	return response.Message
}
