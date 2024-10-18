package ai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewClient(baseURL string) *Client {
	client := resty.New()
	client.SetBaseURL(baseURL)

	return &Client{
		client: client,
	}
}

type GetRecommendationsRequest struct {
	Text string `json:"text"`
}

type Recommendation string

type GetRecommendationsResponse struct {
	Result          string           `json:"result"`
	Recommendations []Recommendation `json:"recommendations"`
}

func (c *Client) GetRecommendations(
	ctx context.Context,
	userInput string,
) (GetRecommendationsResponse, error) {
	var respObj GetRecommendationsResponse

	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(GetRecommendationsRequest{
			Text: userInput,
		}).
		SetResult(&respObj).
		Post("/api/v1/ai_backend/recommendations")
	if err != nil {
		return GetRecommendationsResponse{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return GetRecommendationsResponse{}, fmt.Errorf("unexpected status code: %s", resp.Body())
	}

	return respObj, nil
}
