package replicate

import (
	"ai-wardrobe/internal/app/config"
	"ai-wardrobe/internal/platform/logger"
	"ai-wardrobe/internal/wardrobe/domain"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ReplicateClient struct {
	apiToken string
	baseURL  string
	modelVer string
	client   *http.Client
	logger   *logger.Logger
}

func New(cfg *config.Config, logger *logger.Logger) (*ReplicateClient, error) {

	logger.Info("Initializing Replicate client")

	return &ReplicateClient{
		apiToken: cfg.Replicate.Token,
		baseURL:  cfg.Replicate.BaseURL,
		modelVer: cfg.Replicate.ModelVersion,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		logger: logger,
	}, nil
}

func (c *ReplicateClient) PostTryOn(ctx context.Context, personURL, garmentURL string) (string, error) {

	c.logger.Info("Starting try-on prediction")
	c.logger.Debug("Person URL=%s", personURL)
	c.logger.Debug("Garment URL=%s", garmentURL)

	predID, err := c.createPrediction(ctx, personURL, garmentURL)
	if err != nil {
		return "", err
	}

	c.logger.Info("Prediction created id=%s", predID)

	return c.waitPrediction(ctx, predID)
}

func (c *ReplicateClient) createPrediction(ctx context.Context, personURL, garmentURL string) (string, error) {

	payload := map[string]interface{}{
		"version": c.modelVer,
		"input": map[string]interface{}{
			"human_img": personURL,
			"garm_img":  garmentURL,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	url := c.baseURL + "/predictions"

	c.logger.Trace("POST %s", url)
	c.logger.Trace("Payload=%s", string(data))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	c.logger.Debug("Prediction response status=%d", resp.StatusCode)
	c.logger.Trace("Prediction response body=%s", string(body))

	if resp.StatusCode >= 300 {
		return "", errors.New("replicate prediction creation failed")
	}

	var r domain.ReplicatePredictionResp

	if err := json.Unmarshal(body, &r); err != nil {
		return "", err
	}

	if r.ID == "" {
		return "", errors.New("replicate returned empty prediction id")
	}

	c.logger.Success("Prediction created id=%s", r.ID)

	return r.ID, nil
}

func (c *ReplicateClient) waitPrediction(ctx context.Context, id string) (string, error) {

	c.logger.Info("Waiting prediction result id=%s", id)

	for {

		status, output, err := c.getPrediction(ctx, id)
		if err != nil {
			return "", err
		}

		c.logger.Debug("Prediction status=%s", status)

		switch status {

		case "succeeded":

			c.logger.Success("Prediction succeeded result=%s", output)

			return output, nil

		case "failed", "canceled":

			return "", errors.New("prediction failed")

		case "processing", "starting":

			time.Sleep(2 * time.Second)

		default:

			c.logger.Warn("Unknown status=%s", status)
			time.Sleep(2 * time.Second)
		}
	}
}

func (c *ReplicateClient) getPrediction(ctx context.Context, id string) (string, string, error) {

	url := c.baseURL + "/predictions/" + id

	c.logger.Trace("GET %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Token "+c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	c.logger.Debug("Prediction poll status=%d", resp.StatusCode)

	if resp.StatusCode >= 300 {
		c.logger.Error("Prediction poll failed body=%s", string(body))
		return "", "", errors.New("replicate poll failed")
	}

	var r domain.ReplicatePredictionResp

	if err := json.Unmarshal(body, &r); err != nil {
		return "", "", err
	}

	var status string

	switch v := r.Status.(type) {
	case string:
		status = v
	case float64:
		status = strconv.Itoa(int(v))
	case int:
		status = strconv.Itoa(v)
	}

	var output string

	switch v := r.Output.(type) {

	case string:
		output = v

	case []interface{}:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				output = s
			}
		}

	case []string:
		if len(v) > 0 {
			output = v[0]
		}
	}

	c.logger.Trace("Parsed status=%s output=%s", status, output)

	return status, output, nil
}
