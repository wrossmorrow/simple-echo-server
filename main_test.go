package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetEchoSetting(t *testing.T) {
	r := Request{Headers: http.Header{}}
	if v, ok := getEchoSetting(r, "Delay"); ok {
		t.Errorf("expected no X-Echo-Delay header, got %q", v)
	}

	r = Request{Headers: http.Header{"X-Echo-Delay": {"2"}}}
	v, ok := getEchoSetting(r, "Delay")
	if !ok {
		t.Fatal("expected X-Echo-Delay header to be found")
	}
	if v != "2" {
		t.Errorf("expected delay value %q, got %q", "2", v)
	}

	r = Request{Headers: http.Header{"X-Echo-Status": {"404"}}}
	v, ok = getEchoSetting(r, "Status")
	if !ok || v != "404" {
		t.Errorf("expected status value %q, got %q (ok=%v)", "404", v, ok)
	}
}

func TestGetStatus(t *testing.T) {
	r := Request{Headers: http.Header{}}
	if got := getStatus(r); got != 200 {
		t.Errorf("expected default status 200, got %v", got)
	}

	r = Request{Headers: http.Header{"X-Echo-Status": {"404"}}}
	if got := getStatus(r); got != 404 {
		t.Errorf("expected status 404, got %v", got)
	}

	r = Request{Headers: http.Header{"X-Echo-Status": {"abc"}}}
	if got := getStatus(r); got != 200 {
		t.Errorf("expected default status 200 for invalid value, got %v", got)
	}

	r = Request{Headers: http.Header{"X-Echo-Status": {"99"}}}
	if got := getStatus(r); got != 200 {
		t.Errorf("expected default status 200 for out-of-range value, got %v", got)
	}

	r = Request{Headers: http.Header{"X-Echo-Status": {"600"}}}
	if got := getStatus(r); got != 200 {
		t.Errorf("expected default status 200 for out-of-range value, got %v", got)
	}
}

func TestDoDelay(t *testing.T) {
	r := Request{Headers: http.Header{}}
	start := time.Now()
	doDelay(r)
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("unexpected delay without header: %v", elapsed)
	}

	r = Request{Headers: http.Header{"X-Echo-Delay": {"abc"}}}
	start = time.Now()
	doDelay(r)
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("unexpected delay for invalid value: %v", elapsed)
	}

	r = Request{Headers: http.Header{"X-Echo-Delay": {"1ms"}}}
	start = time.Now()
	doDelay(r)
	if elapsed := time.Since(start); elapsed < time.Millisecond {
		t.Errorf("expected delay, elapsed %v", elapsed)
	}
}

func TestEchoDefaultResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	echo(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected status 200, got %v", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal JSON body: %v", err)
	}
	if got["Status"].(float64) != 200 {
		t.Errorf("expected JSON Status 200, got %v", got["Status"])
	}
}

func TestEchoCustomStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Echo-Status", "503")
	rec := httptest.NewRecorder()
	echo(rec, req)

	if rec.Code != 503 {
		t.Errorf("expected status 503, got %v", rec.Code)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal JSON body: %v", err)
	}
	if got["Status"].(float64) != 503 {
		t.Errorf("expected JSON Status 503, got %v", got["Status"])
	}
}

func TestEchoNoBodyForStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Echo-Status", "204")
	rec := httptest.NewRecorder()
	echo(rec, req)

	if rec.Code != 204 {
		t.Errorf("expected status 204, got %v", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body for 204, got %q", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Echo-Status", "304")
	rec = httptest.NewRecorder()
	echo(rec, req)

	if rec.Code != 304 {
		t.Errorf("expected status 304, got %v", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("expected empty body for 304, got %q", rec.Body.String())
	}
}
