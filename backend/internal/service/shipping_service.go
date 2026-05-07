package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"erp-cosmetics-backend/internal/config"

	"go.uber.org/zap"
)

type ShippingService interface {
	GetProvinces(ctx context.Context) ([]Province, error)
	GetCities(ctx context.Context, provinceID string) ([]City, error)
	CalculateCost(ctx context.Context, req *CostRequest) (*CostResponse, error)
}

type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProvinceID  string `json:"province_id"`
	ProvinceName string `json:"province_name"`
	PostalCode  string `json:"postal_code"`
}

type CostRequest struct {
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	Weight      int    `json:"weight" binding:"required,min=1"`
	Courier     string `json:"courier" binding:"required"`
}

type CostResponse struct {
	RajaOngkir struct {
		Results []CourierResult `json:"results"`
	} `json:"rajaongkir"`
}

type CourierResult struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Costs []Cost `json:"costs"`
}

type Cost struct {
	Service     string       `json:"service"`
	Description string       `json:"description"`
	Cost        []CostDetail `json:"cost"`
}

type CostDetail struct {
	Value int    `json:"value"`
	Etd   string `json:"etd"`
	Note  string `json:"note"`
}

type shippingService struct {
	cfg    *config.Config
	logger *zap.Logger
	client *http.Client
}

func NewShippingService(cfg *config.Config, logger *zap.Logger) ShippingService {
	return &shippingService{
		cfg:    cfg,
		logger: logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *shippingService) GetProvinces(ctx context.Context) ([]Province, error) {
	url := fmt.Sprintf("%s/province", s.cfg.RajaOngkirBaseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("key", s.cfg.RajaOngkirAPIKey)
	
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to get provinces", zap.Error(err))
		return nil, errors.New("failed to fetch provinces")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		RajaOngkir struct {
			Results []Province `json:"results"`
		} `json:"rajaongkir"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	return result.RajaOngkir.Results, nil
}

func (s *shippingService) GetCities(ctx context.Context, provinceID string) ([]City, error) {
	url := fmt.Sprintf("%s/city?province=%s", s.cfg.RajaOngkirBaseURL, provinceID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("key", s.cfg.RajaOngkirAPIKey)
	
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to get cities", zap.Error(err))
		return nil, errors.New("failed to fetch cities")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		RajaOngkir struct {
			Results []City `json:"results"`
		} `json:"rajaongkir"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	return result.RajaOngkir.Results, nil
}

func (s *shippingService) CalculateCost(ctx context.Context, req *CostRequest) (*CostResponse, error) {
	url := fmt.Sprintf("%s/cost", s.cfg.RajaOngkirBaseURL)
	
	payload := map[string]interface{}{
		"origin":      req.Origin,
		"destination": req.Destination,
		"weight":      req.Weight,
		"courier":     req.Courier,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	
	httpReq.Header.Set("key", s.cfg.RajaOngkirAPIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := s.client.Do(httpReq)
	if err != nil {
		s.logger.Error("Failed to calculate cost", zap.Error(err))
		return nil, errors.New("failed to calculate shipping cost")
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result CostResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}