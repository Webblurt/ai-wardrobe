package replicate

import (
	"ai-wardrobe/internal/app/config"
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
}

func New(cfg *config.Config) (*ReplicateClient, error) {
	return &ReplicateClient{
		apiToken: cfg.Replicate.Token,
		baseURL:  cfg.Replicate.BaseURL,
		modelVer: cfg.Replicate.ModelVersion,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

func (c *ReplicateClient) PostTryOn(ctx context.Context, personPath, garmentPath string) (string, error) {

	personURL, err := c.uploadFile(ctx, personPath)
	if err != nil {
		return "", err
	}

	garmentURL, err := c.uploadFile(ctx, garmentPath)
	if err != nil {
		return "", err
	}

	predID, err := c.createPrediction(ctx, personURL, garmentURL)
	if err != nil {
		return "", err
	}

	return c.waitPrediction(ctx, predID)
}

func (c *ReplicateClient) uploadFile(ctx context.Context, path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+"/files",
		body,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+c.apiToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r domain.ReplicateUploadResp

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	return r.URL, nil
}

func (c *ReplicateClient) createPrediction(ctx context.Context, personURL, garmentURL string) (string, error) {

	payload := map[string]interface{}{
		"version": c.modelVer,
		"input": map[string]interface{}{
			"person":  personURL,
			"garment": garmentURL,
		},
	}

	data, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL+"/predictions",
		bytes.NewBuffer(data),
	)
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

	var r domain.ReplicatePredictionResp

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	return r.ID, nil
}

func (c *ReplicateClient) waitPrediction(ctx context.Context, id string) (string, error) {

	for {

		status, url, err := c.getPrediction(ctx, id)
		if err != nil {
			return "", err
		}

		switch status {

		case "succeeded":
			return url, nil

		case "failed":
			return "", errors.New("prediction failed")

		case "processing", "starting":
			time.Sleep(2 * time.Second)
		}
	}
}

func (c *ReplicateClient) getPrediction(ctx context.Context, id string) (string, string, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.baseURL+"/predictions/"+id,
		nil,
	)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Token "+c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var r domain.ReplicatePredictionResp

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", "", err
	}

	// status
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

	// output
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

	return status, output, nil
}
