package weatherapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	s "backend/shared"

	// "github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	// "github.com/go-kit/kit/ratelimit"
	kithttp "github.com/go-kit/kit/transport/http"
	// "github.com/sony/gobreaker"
	// "golang.org/x/time/rate"
)

type proxyService struct {
	context.Context
	WeatherEndpoint endpoint.Endpoint
	WeatherService
	logger log.Logger
}

func (srv proxyService) QueryPoint(p *s.Point, wg *sync.WaitGroup, cc chan struct{}) {
	defer wg.Done()
	defer func() { <-cc }()
	response, err := srv.WeatherEndpoint(srv.Context, GetWeatherRequest{
		Lat: p.Lat,
		Lng: p.Lng,
	})
	if err != nil {
		level.Debug(srv.logger).Log("msg", "error querying weather api endpoint", "err", err)
		p.Val = -1
		return
	}

	resp := response.(GetWeatherResponse)

	sun := 100 - float32(resp.Current.Cloud)
	p.Val = sun
	level.Debug(srv.logger).Log("msg", "queried point", "lat", p.Lat, "lng", p.Lng, "val", sun)
}

type ServiceMiddleware func(WeatherService) WeatherService

func NewProxyingMiddleware(ctx context.Context, logger log.Logger) ServiceMiddleware {
	return func(srv WeatherService) WeatherService {
		var e endpoint.Endpoint
		e = makeWeatherEndpoint(ctx)
		// e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		// e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Duration(int64(900/float32(s.Config.WeatherApiMaxRequestsPerSecond)))*time.Millisecond), s.Config.WeatherApiMaxRequestsPerSecond))(e)
		return proxyService{ctx, e, srv, logger}
	}
}

func makeWeatherEndpoint(ctx context.Context) endpoint.Endpoint {
	u, err := url.Parse(s.Config.WeatherApiUrl)
	// u, err := url.Parse("https://reqbin.com/echo/get/json")
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/v1/current.json" // TODO hmmm...
	}
	c := http.Client{
		Timeout: time.Duration(s.Config.WeatherApiClientTimeoutSec) * time.Second,
	}
	return kithttp.NewClient(
		"GET", u,
		encodeGetWeatherRequest,
		decodeGetWeatherResponse,
		kithttp.SetClient(&c), // XXX SetClient does actually not expect the http.Client, weird this compiles
	).Endpoint()
}

func decodeGetWeatherResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	defer resp.Body.Close()
	var response GetWeatherResponse
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Request not successfull. Status Code: %v", resp.StatusCode))
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Current == nil {
		return nil, errors.New("mandatory field 'current' not in response")
	}
	return response, nil
}

func encodeGetWeatherRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(GetWeatherRequest)
	k := s.Config.WeatherApiKey
	coords := fmt.Sprintf("%v,%v", req.Lat, req.Lng)
	q := r.URL.Query()
	q.Add("key", k)
	q.Add("q", coords)

	r.URL.RawQuery = q.Encode()

	return nil
}

type GetWeatherRequest struct {
	Lat float32
	Lng float32
}

type GetWeatherResponse struct {
	Location *Location
	Current  *Current
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
