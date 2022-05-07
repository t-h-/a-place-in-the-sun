package weatherapi

import (
	s "backend/shared"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"context"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"

	"golang.org/x/time/rate"
)

var (
	ApiErr = errors.New("Unable to handle Api Request")
)

//go:generate mockgen -destination=../mocks/mock_api.go -package=mocks . WeatherApi
type WeatherApi interface {
	QueryPoints(points []*s.Point)
}

type api struct {
	client               *http.Client
	Ratelimiter          *rate.Limiter
	logger               log.Logger
	ApiKey               string
	MaxRequestsPerSecond int
	MaxRequestBurst      int
}

func NewApi(apiKey string, maxRequestsPerSecond int, maxRequestBurst int, logger log.Logger) *api {
	c := http.Client{
		Timeout: 7 * time.Second,
	}
	return &api{
		client:               &c,
		Ratelimiter:          rate.NewLimiter(rate.Every(time.Duration(1000/maxRequestsPerSecond)*time.Millisecond), maxRequestBurst),
		logger:               log.With(logger, "api", "weatherapi.com"),
		MaxRequestsPerSecond: maxRequestsPerSecond,
		MaxRequestBurst:      maxRequestBurst,
		ApiKey:               apiKey,
	}
}

func (api *api) QueryPoints(points []*s.Point) {
	var wg sync.WaitGroup
	var flag bool = true
	for _, p := range points {
		wg.Add(1)
		go api.QueryPoint(p, &wg)
	}
	go func() {
		wg.Wait()
		flag = false
		level.Debug(api.logger).Log("msg", "done querying")
	}()

	for flag {
		level.Debug(api.logger).Log("msg", "open", "routines", runtime.NumGoroutine())
		time.Sleep(500 * time.Millisecond)
	}
}

func (api *api) QueryPoint(p *s.Point, wg *sync.WaitGroup) error {
	defer wg.Done()
	reqURL := "http://api.weatherapi.com/v1/current.json"
	req, _ := http.NewRequest("GET", reqURL, nil)
	coords := fmt.Sprintf("%v,%v", p.Lat, p.Lng)
	q := req.URL.Query()
	q.Add("key", api.ApiKey)
	q.Add("q", coords)
	req.URL.RawQuery = q.Encode()
	resp, err := api.do(req)
	if err != nil {
		level.Error(api.logger).Log("msg", "error while executing request", "error", err, "lat", p.Lat, "lng", p.Lng)
		return ApiErr
	}

	if resp.StatusCode != 200 {
		level.Error(api.logger).Log("msg", "Not able to handle request", "status_code", resp.StatusCode, "lat", p.Lat, "lng", p.Lng)
		return ApiErr
	}

	res, err := api.unmarshal(resp)
	if err != nil {
		level.Error(api.logger).Log("msg", "Can not unmarshal JSON", "error", err, "lat", p.Lat, "lng", p.Lng)
		return ApiErr
	}
	p.Val = float32(100 - res.Current.Cloud)
	return nil
}

func (api *api) do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := api.Ratelimiter.Wait(ctx)
	if err != nil {
		level.Error(api.logger).Log("msg", "error when waiting for ratelimiter", "error", err)
		return nil, err
	}
	level.Debug(api.logger).Log("msg", "requesting point", "req", req.URL)
	resp, err := api.client.Do(req)
	if err != nil {
		level.Error(api.logger).Log("msg", "request failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (api *api) unmarshal(resp *http.Response) (Response, error) {
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return Response{}, err
	}
	return result, nil
}

type Response struct {
	Location Location
	Current  Current
}

type Location struct {
	Name            string  `json:"name"`
	Region          string  `json:"region"`
	Country         string  `json:"country"`
	Lat             float32 `json:"lat"`
	Lon             float32 `json:"lon"`
	Tz_id           string  `json:"tz_id"`
	Localtime_epoch int64   `json:"localtime_epoch"`
	Localtime       string  `json:"localtime"`
}

type Current struct {
	Last_updated_epoch int64     `json:"last_updated_epoch"`
	Last_updated       string    `json:"last_updated"`
	Temp_c             float32   `json:"temp_c"`
	Temp_f             float32   `json:"temp_f"`
	Is_day             int       `json:"is_day"`
	Condition          Condition `json:"condition"`
	Wind_mph           float32   `json:"wind_mph"`
	Wind_kph           float32   `json:"wind_kph"`
	Wind_degree        int       `json:"wind_degree"`
	Wind_dir           string    `json:"wind_dir"`
	Pressure_mb        float32   `json:"pressure_mb"`
	Pressure_in        float32   `json:"pressure_in"`
	Precip_mm          float32   `json:"precip_mm"`
	Precip_in          float32   `json:"precip_in"`
	Humidity           int       `json:"humidity"`
	Cloud              int       `json:"cloud"`
	Feelslike_c        float32   `json:"feelslike_c"`
	Feelslike_f        float32   `json:"feelslike_f"`
	Vis_km             float32   `json:"vis_km"`
	Vis_miles          float32   `json:"vis_miles"`
	Uv                 float32   `json:"uv"`
	Gust_mph           float32   `json:"gust_mph"`
	Gust_kph           float32   `json:"gust_kph"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}
