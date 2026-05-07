package rajaongkir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RajaOngkirClient struct {
	apiKey   string
	baseURL  string
	client   *http.Client
}

type Province struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

type City struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

type CostRequest struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Weight      int    `json:"weight"`
	Courier     string `json:"courier"`
}

type CostResponse struct {
	RajaOngkir struct {
		Query struct {
			Origin      string `json:"origin"`
			Destination string `json:"destination"`
			Weight      int    `json:"weight"`
			Courier     string `json:"courier"`
		} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
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

func NewRajaOngkirClient(apiKey, baseURL string) *RajaOngkirClient {
	return &RajaOngkirClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *RajaOngkirClient) GetProvinces() ([]Province, error) {
	url := fmt.Sprintf("%s/province", r.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("key", r.apiKey)
	
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
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

func (r *RajaOngkirClient) GetCities(provinceID string) ([]City, error) {
	url := fmt.Sprintf("%s/city?province=%s", r.baseURL, provinceID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("key", r.apiKey)
	
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
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

func (r *RajaOngkirClient) CalculateCost(req *CostRequest) (*CostResponse, error) {
	url := fmt.Sprintf("%s/cost", r.baseURL)
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	httpReq.Header.Set("key", r.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	
	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, err
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

func (r *RajaOngkirClient) GetCityByID(cityID string) (*City, error) {
	url := fmt.Sprintf("%s/city?id=%s", r.baseURL, cityID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("key", r.apiKey)
	
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		RajaOngkir struct {
			Results City `json:"results"`
		} `json:"rajaongkir"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	return &result.RajaOngkir.Results, nil
}