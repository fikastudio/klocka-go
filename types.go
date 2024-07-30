package klocka

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	RegionID      string            `json:"regionId"`
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
	RegionID      string            `json:"regionId"`
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

type APIError struct {
	err    error
	body   []byte
	status int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("got status code %d: %v", e.status, e.err.Error())
}

func (e *APIError) Unwrap() error { return e.err }
func (e *APIError) Status() int   { return e.status }
func (e *APIError) Body() []byte  { return e.body }

func ConstructHeaders(payload []byte, secret string) http.Header {
	timestamp := time.Now()
	sig := ComputeSignature(timestamp, payload, secret)
	sigBase64 := hex.EncodeToString(sig)
	headers := http.Header{}

	headers.Set("x-klocka-timestamp", fmt.Sprintf("%d", timestamp.Unix()))
	headers.Set("x-klocka-signature", sigBase64)

	return headers
}

func VerifyRequest(headers http.Header, payload []byte, secret string) error {
	timestampStr, err := strconv.ParseInt(headers.Get("x-klocka-timestamp"), 10, 64)
	if err != nil {
		return fmt.Errorf("x-klocka-timestamp is invalid: %v", err)
	}
	timestamp := time.Unix(timestampStr, 0)

	sig := ComputeSignature(timestamp, payload, secret)
	sigBase64 := hex.EncodeToString(sig)
	if sigBase64 != headers.Get("x-klocka-signature") {
		return errors.New("x-klocka-signature does not match")
	}

	return nil
}

func ComputeSignature(t time.Time, payload []byte, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d", t.Unix())))
	mac.Write([]byte("."))
	mac.Write(payload)
	return mac.Sum(nil)
}
