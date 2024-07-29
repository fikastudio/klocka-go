package klocka

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Task struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	MaxDuration   Duration          `json:"maxDuration"`
	URL           string            `json:"url"`
	AllowOverlap  bool              `json:"allowOverlap"`
	Interval      Duration          `json:"interval"`
	Cron          *string           `json:"cron"`
	HttpMethod    string            `json:"httpMethod"`
	HttpHeaders   http.Header       `json:"httpHeaders"`
	Region        string            `json:"region"`
	Meta          map[string]string `json:"meta"`
	OkStatusCodes []int             `json:"okStatusCodes"`
}

type TaskInput struct {
	ID            *string           `json:"id,omitempty"`
	Name          string            `json:"name"`
	MaxDuration   Duration          `json:"maxDuration"`
	URL           string            `json:"url"`
	AllowOverlap  bool              `json:"allowOverlap"`
	Interval      Duration          `json:"interval"`
	Cron          *string           `json:"cron"`
	HttpMethod    string            `json:"httpMethod"`
	HttpHeaders   http.Header       `json:"httpHeaders"`
	Region        string            `json:"region"`
	Meta          map[string]string `json:"meta"`
	OkStatusCodes []int             `json:"okStatusCodes"`
}

type PaginatedResponse struct {
	Data     []Task
	PageInfo PageInfo `json:"pageInfo"`
}

type PageInfo struct{}

type ListOpts struct {
	PerPage *uint8
	Page    *int
}

func (li *ListOpts) Encode() string {
	vals := url.Values{}
	if li.PerPage != nil {
		vals.Set("perPage", fmt.Sprintf("%d", *li.PerPage))
	}
	if li.Page != nil {
		vals.Set("page", fmt.Sprintf("%d", *li.Page))
	}

	return vals.Encode()
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return
	}

	var id int64
	id, err = json.Number(string(b)).Int64()
	d.Duration = time.Duration(id)

	return
}

func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}
