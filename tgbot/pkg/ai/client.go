package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
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
	Result          string   `json:"result"`
	Recommendations []string `json:"recommendations"`
}

func (c *Client) GetDiagnosises(
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
		Post("/api/v1/ai_backend/diagnose/")
	if err != nil {
		return GetRecommendationsResponse{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return GetRecommendationsResponse{}, fmt.Errorf("get diagnosis: unexpected status code: %s", resp.Body())
	}

	return respObj, nil
}

type GetAnalysisRequest struct {
	Image string `json:"image"`
}

type GetAnalysisResponse struct {
	Result    string `json:"result"`
	Analytics string `json:"analytics"`
}

func (c *Client) SendAnalysis(
	ctx context.Context,
	photo io.ReadCloser,
) (GetAnalysisResponse, error) {
	var respObj GetAnalysisResponse
	buf := bytes.Buffer{}
	_, err := io.Copy(&buf, photo)
	if err != nil {
		return GetAnalysisResponse{}, err
	}

	imageStr := base64.StdEncoding.EncodeToString(buf.Bytes())

	resp, err := c.client.R().
		SetHeader("Content-Type", "image").
		SetBody(GetAnalysisRequest{
			Image: imageStr,
		}).
		SetResult(&respObj).
		Post("/api/v1/ai_backend/analyze/")
	if err != nil {
		return GetAnalysisResponse{}, err
	}

	if resp.StatusCode() != http.StatusOK {
		return GetAnalysisResponse{}, fmt.Errorf("send analysis unexpected status code: %s", resp.Body())
	}

	return respObj, nil
}
