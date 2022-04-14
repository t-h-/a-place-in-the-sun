package infra

import (
	"backend/sunnyness"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"context"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/log"

	"golang.org/x/time/rate"
)

const ApiKey = "591b7934afcf484fa3191051223101"

var (
	ApiErr = errors.New("Unable to handle Api Request")
)

type api struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
	logger      log.Logger
}

func NewApi(logger log.Logger) *api {
	return &api{
		client:      http.DefaultClient,
		Ratelimiter: rate.NewLimiter(rate.Every(2*time.Second), 200),
		logger:      log.With(logger, "cache", "apiTODO"),
	}
}

func (api *api) QueryPoints(points []*sunnyness.Point) {
	for _, p := range points {
		err := api.QueryPoint(p)
		if err != nil {
			// TODO error handling
		}
	}
}

func (api *api) QueryPoint(p *sunnyness.Point) error {
	reqURL := "http://api.weatherapi.com/v1/current.json"
	req, _ := http.NewRequest("GET", reqURL, nil)
	coords := fmt.Sprintf("%v,%v", p.Lat, p.Lng)
	q := req.URL.Query()
	q.Add("key", ApiKey)
	q.Add("q", coords)
	q.Add("aqi", "no")
	req.URL.RawQuery = q.Encode()
	resp, err := api.do(req)
	if err != nil {
		// TODO error handling
		return ApiErr
	}

	res, err := unmarshal(resp)
	if err != nil {
		// TODO error handling
		return err
	}
	p.Val = float32(100 - res.Current.Cloud)
	return nil
}

func (c *api) do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func unmarshal(resp *http.Response) (Response, error) {
	defer resp.Body.Close()
	fdsa, _ := ioutil.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(fdsa, &result); err != nil {
		fmt.Println("Can not unmarshal JSON") // TODO proper logging and error handling
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
