package weatherapi

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
	"github.com/jsperandio/goexptCloudRun/pkg/text"
)

const fixedCountry = "Brazil"

type Request struct {
	City    string
	State   string
	Country string
}

func NewRequest(lc entity.Location) Request {
	return Request{
		City:    text.StripAccents(lc.City),
		State:   text.StripAccents(lc.StateName),
		Country: fixedCountry,
	}
}

func (rq Request) Query(apiKey string) url.Values {
	query := url.Values{}
	query.Set("key", apiKey)
	query.Set("q", rq.location())
	query.Set("aqi", "no")

	return query
}

func (rq Request) validateResolved(loc Location) error {
	if !text.EqualIgnoringAccents(loc.Name, rq.City) {
		return rq.mismatchError(loc)
	}

	if rq.State != "" && !text.EqualIgnoringAccents(loc.Region, rq.State) {
		return rq.mismatchError(loc)
	}

	return nil
}

func (rq Request) mismatchError(loc Location) error {
	return fmt.Errorf(
		"%w: asked for %s, %s and got %s, %s",
		ErrLocationMismatch,
		rq.City,
		rq.State,
		loc.Name,
		loc.Region,
	)
}

func (rq Request) location() string {
	parts := make([]string, 0, 3)

	parts = append(parts, rq.City)

	if rq.State != "" {
		parts = append(parts, rq.State)
	}

	parts = append(parts, rq.Country)

	return strings.Join(parts, ",")
}
