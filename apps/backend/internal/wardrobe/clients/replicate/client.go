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
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ReplicateClient struct {
	apiToken string
	baseURL  string
	client   *http.Client
	modelVer string
	logger   *logger.Logger
}

func New(cfg *config.Config, logger *logger.Logger) (*ReplicateClient, error) {
	return &ReplicateClient{
		apiToken: cfg.Replicate.Token,
		baseURL:  cfg.Replicate.BaseURL,
		modelVer: cfg.Replicate.ModelVersion,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
	}, nil
}

func (c *ReplicateClient) PostTryOn(ctx context.Context, personPath, garmentPath string) (string, error) {

	c.logger.Info("Starting try-on prediction")
	c.logger.Debug("Person image path=%s", personPath)
	c.logger.Debug("Garment image path=%s", garmentPath)

	personURL, err := c.uploadFile(ctx, personPath)
	if err != nil {
		c.logger.Error("Failed uploading person image: %v", err)
		return "", err
	}

	c.logger.Debug("Person uploaded url=%s", personURL)

	garmentURL, err := c.uploadFile(ctx, garmentPath)
	if err != nil {
		c.logger.Error("Failed uploading garment image: %v", err)
		return "", err
	}

	c.logger.Debug("Garment uploaded url=%s", garmentURL)

	predID, err := c.createPrediction(ctx, personURL, garmentURL)
	if err != nil {
		c.logger.Error("Failed creating prediction: %v", err)
		return "", err
	}

	c.logger.Info("Prediction created id=%s", predID)

	result, err := c.waitPrediction(ctx, predID)
	if err != nil {
		c.logger.Error("Prediction failed: %v", err)
		return "", err
	}

	c.logger.Success("Try-on prediction completed result=%s", result)

	return result, nil
}

func (c *ReplicateClient) uploadFile(ctx context.Context, path string) (string, error) {

	c.logger.Info("Uploading file to Replicate")
	c.logger.Debug("File path=%s", path)

	file, err := os.Open(path)
	if err != nil {
		c.logger.Error("Failed opening file: %v", err)
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		c.logger.Error("Failed creating form file: %v", err)
		return "", err
	}

	size, err := io.Copy(part, file)
	if err != nil {
		c.logger.Error("Failed copying file data: %v", err)
		return "", err
	}

	c.logger.Debug("File copied size=%d bytes", size)

	writer.Close()

	url := c.baseURL + "/files"

	c.logger.Trace("HTTP POST %s", url)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		body,
	)
	if err != nil {
		c.logger.Error("Failed creating request: %v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Token "+c.apiToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Upload request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	c.logger.Debug("Upload response status=%d", resp.StatusCode)

	var r domain.ReplicateUploadResp

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		c.logger.Error("Failed decoding upload response: %v", err)
		return "", err
	}

	c.logger.Success("File uploaded url=%s", r.URL)

	return r.URL, nil
}

func (c *ReplicateClient) createPrediction(ctx context.Context, personURL, garmentURL string) (string, error) {

	c.logger.Info("Creating prediction")

	c.logger.Debug("Person URL=%s", personURL)
	c.logger.Debug("Garment URL=%s", garmentURL)
	c.logger.Debug("Model version=%s", c.modelVer)

	payload := map[string]interface{}{
		"version": c.modelVer,
		"input": map[string]interface{}{
			"person":  personURL,
			"garment": garmentURL,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		c.logger.Error("Failed marshaling payload: %v", err)
		return "", err
	}

	c.logger.Trace("Prediction payload=%s", string(data))

	url := c.baseURL + "/predictions"

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		bytes.NewBuffer(data),
	)
	if err != nil {
		c.logger.Error("Failed creating prediction request: %v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Token "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Prediction request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	c.logger.Debug("Prediction response status=%d", resp.StatusCode)

	var r domain.ReplicatePredictionResp

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		c.logger.Error("Failed decoding prediction response: %v", err)
		return "", err
	}

	c.logger.Success("Prediction created id=%s", r.ID)

	return r.ID, nil
}

func (c *ReplicateClient) waitPrediction(ctx context.Context, id string) (string, error) {

	c.logger.Info("Waiting prediction result id=%s", id)

	for {

		status, url, err := c.getPrediction(ctx, id)
		if err != nil {
			c.logger.Error("Failed getting prediction status: %v", err)
			return "", err
		}

		c.logger.Debug("Prediction status=%s", status)

		switch status {

		case "succeeded":

			c.logger.Success("Prediction succeeded result=%s", url)

			return url, nil

		case "failed":

			c.logger.Error("Prediction failed id=%s", id)

			return "", errors.New("prediction failed")

		case "processing", "starting":

			c.logger.Trace("Prediction still running... waiting 2s")
			time.Sleep(2 * time.Second)

		default:

			c.logger.Warn("Unknown prediction status=%s", status)
			time.Sleep(2 * time.Second)

		}
	}
}

func (c *ReplicateClient) getPrediction(ctx context.Context, id string) (string, string, error) {

	url := c.baseURL + "/predictions/" + id

	c.logger.Trace("HTTP GET %s", url)

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url,
		nil,
	)
	if err != nil {
		c.logger.Error("Failed creating request: %v", err)
		return "", "", err
	}

	req.Header.Set("Authorization", "Token "+c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Prediction status request failed: %v", err)
		return "", "", err
	}
	defer resp.Body.Close()

	c.logger.Debug("Prediction poll response status=%d", resp.StatusCode)

	var r domain.ReplicatePredictionResp

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		c.logger.Error("Failed decoding prediction status: %v", err)
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
	default:
		status = ""
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

	c.logger.Trace("Prediction parsed status=%s output=%s", status, output)

	return status, output, nil
}
