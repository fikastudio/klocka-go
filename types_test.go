package klocka_test

import (
	"net/http"
	"testing"

	"github.com/fikastudio/klocka-go"
)

func TestSignature(t *testing.T) {
	payload := []byte("Hello")
	secret := "TheiS8ee"
	headers := klocka.ConstructHeaders(payload, secret)

	// passes
	if err := klocka.VerifyRequest(headers, payload, secret); err != nil {
		t.Fail()
	}

	// malformed headers
	if err := klocka.VerifyRequest(http.Header{}, payload, secret); err == nil {
		t.Fail()
	}

	// malformed body
	if err := klocka.VerifyRequest(headers, []byte(""), secret); err == nil {
		t.Fail()
	}

	// malformed secret
	if err := klocka.VerifyRequest(headers, payload, "another"); err == nil {
		t.Fail()
	}
}
