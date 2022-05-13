package weatherapi

import (
	s "backend/shared"
	"errors"
	"sync"
)

// TODO longterm: translate/unify all lower level errors to here defined errors
var (
	ApiErr = errors.New("Unable to handle Api Request")
)

//go:generate mockgen -destination=../mocks/mock_api.go -package=mocks . WeatherService
type WeatherService interface {
	QueryPoint(p *s.Point, wg *sync.WaitGroup, cc chan struct{})
}
