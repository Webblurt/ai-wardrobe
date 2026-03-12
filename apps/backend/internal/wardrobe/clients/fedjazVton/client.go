package fedjazvton

import (
	"ai-wardrobe/internal/app/config"
	"ai-wardrobe/internal/platform/logger"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FedjazVtonClient struct {
	apiToken string
	baseURL  string
	client   *http.Client
	logger   *logger.Logger
}

func New(cfg *config.Config, logger *logger.Logger) (*FedjazVtonClient, error) {

	logger.Info("Initializing FedjazVton client")

	return &FedjazVtonClient{
		apiToken: cfg.FedjazVton.Token,
		baseURL:  cfg.FedjazVton.BaseURL,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		logger: logger,
	}, nil
}

func (c *FedjazVtonClient) PostTryOn(
	ctx context.Context,
	personPath string,
	garmentPath string,
) ([]byte, error) {

	id, err := c.submit(ctx, personPath, garmentPath)
	if err != nil {
		return nil, err
	}

	c.logger.Info("TryOn submitted correlationId=%s", id)

	return c.waitResult(ctx, id)
}

func (c *FedjazVtonClient) submit(
	ctx context.Context,
	personPath string,
	garmentPath string,
) (string, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	personFile, err := os.Open(personPath)
	if err != nil {
		return "", err
	}
	defer personFile.Close()

	garmentFile, err := os.Open(garmentPath)
	if err != nil {
		return "", err
	}
	defer garmentFile.Close()

	personPart, err := writer.CreateFormFile("person", filepath.Base(personPath))
	if err != nil {
		return "", err
	}

	io.Copy(personPart, personFile)

	garmentPart, err := writer.CreateFormFile("garment", filepath.Base(garmentPath))
	if err != nil {
		return "", err
	}

	io.Copy(garmentPart, garmentFile)

	// optional params
	writer.WriteField("category", "dresses")
	writer.WriteField("steps", "30")

	writer.Close()

	url := c.baseURL + "/tryon"

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("x-api-key", c.apiToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("submit failed: %s", string(data))
	}

	var r struct {
		CorrelationID string `json:"correlationId"`
	}

	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", err
	}

	if r.CorrelationID == "" {
		return "", errors.New("empty correlationId")
	}

	return r.CorrelationID, nil
}

func (c *FedjazVtonClient) waitResult(
	ctx context.Context,
	id string,
) ([]byte, error) {

	c.logger.Info("Waiting result id=%s", id)

	for {

		img, done, err := c.getResult(ctx, id)
		if err != nil {
			return nil, err
		}

		if done {
			return img, nil
		}

		time.Sleep(2 * time.Second)
	}
}

func (c *FedjazVtonClient) getResult(
	ctx context.Context,
	id string,
) ([]byte, bool, error) {

	url := c.baseURL + "/tryon/" + id

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, false, err
	}

	req.Header.Set("x-api-key", c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case 200:

		img, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, false, err
		}

		return img, true, nil

	case 202:

		c.logger.Debug("Inference still processing")
		return nil, false, nil

	case 404:

		return nil, false, errors.New("task not found")

	default:

		body, _ := io.ReadAll(resp.Body)
		return nil, false, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}
}
