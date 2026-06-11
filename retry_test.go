package buildkite

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestDo_POST_IsRetriedOn429 reproduces Bug #1: the GET-only guard in the
// current Do() means POST requests that hit 429 fail immediately.
func TestDo_POST_IsRetriedOn429(t *testing.T) {
	callCount := 0

	ms, client, teardown := newMockServerAndClient(t)
	defer teardown()

	ms.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req, err := client.NewRequest(context.Background(), http.MethodPost, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	_, err = client.Do(req, nil)
	if err != nil {
		t.Errorf("expected POST to be retried after 429, got error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 server calls (initial + 1 retry), got %d", callCount)
	}
}

// TestWithMaxRetries_Zero verifies that WithMaxRetries(0) disables retries entirely.
func TestWithMaxRetries_Zero(t *testing.T) {
	callCount := 0

	ms, client, teardown := newMockServerAndClient(t)
	defer teardown()

	ms.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("RateLimit-Reset", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	if err := WithMaxRetries(0)(client); err != nil {
		t.Fatalf("WithMaxRetries(0): %v", err)
	}

	req, err := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	_, _ = client.Do(req, nil)

	if callCount != 1 {
		t.Errorf("expected exactly 1 call (no retries when maxRetries=0), got %d", callCount)
	}
}

// TestWithMaxRetries_LimitsRetries verifies the client stops retrying after the configured limit.
func TestWithMaxRetries_LimitsRetries(t *testing.T) {
	callCount := 0

	ms, client, teardown := newMockServerAndClient(t)
	defer teardown()

	ms.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("RateLimit-Reset", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	if err := WithMaxRetries(2)(client); err != nil {
		t.Fatalf("WithMaxRetries(2): %v", err)
	}

	req, err := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	_, err = client.Do(req, nil)
	if err == nil {
		t.Fatal("expected error after exhausting retries, got nil")
	}
	// 1 initial + 2 retries = 3 total
	if callCount != 3 {
		t.Errorf("expected 3 server calls (1 initial + 2 retries), got %d", callCount)
	}
}

// TestDo_POST_BodyIntactOnRetry verifies the request body is fully replayed on each retry.
func TestDo_POST_BodyIntactOnRetry(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	callCount := 0
	var receivedBodies []string

	ms, client, teardown := newMockServerAndClient(t)
	defer teardown()

	ms.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		body, _ := io.ReadAll(r.Body)
		receivedBodies = append(receivedBodies, string(body))
		if callCount == 1 {
			w.Header().Set("RateLimit-Reset", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	req, err := client.NewRequest(context.Background(), http.MethodPost, "/test", payload{Name: "my-pipeline"})
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	_, err = client.Do(req, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 server calls, got %d", callCount)
	}
	for i, body := range receivedBodies {
		if body == "" {
			t.Errorf("call %d received empty body — request body was not replayed on retry", i+1)
		}
		if receivedBodies[0] != body {
			t.Errorf("call %d body %q differs from call 1 body %q", i+1, body, receivedBodies[0])
		}
	}
}

// TestDo_ContextCancelledDuringRetryWait verifies a cancelled context interrupts
// the retry sleep promptly.
func TestDo_ContextCancelledDuringRetryWait(t *testing.T) {
	ms, client, teardown := newMockServerAndClient(t)
	defer teardown()

	ms.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("RateLimit-Reset", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	start := time.Now()
	_, err = client.Do(req, nil)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error after context cancellation, got nil")
	}
	if elapsed > 2*time.Second {
		t.Errorf("expected prompt return on context cancel, took %v", elapsed)
	}
}

// TestDo_ExhaustedRetries_ReturnsFinal429AsErrorResponse verifies the final 429
// flows through as a proper *ErrorResponse, not a generic string error.
func TestDo_ExhaustedRetries_ReturnsFinal429AsErrorResponse(t *testing.T) {
	ms, client, teardown := newMockServerAndClient(t)
	defer teardown()

	ms.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("RateLimit-Reset", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = fmt.Fprint(w, `{"message":"rate limit exceeded"}`)
	})

	if err := WithMaxRetries(1)(client); err != nil {
		t.Fatalf("WithMaxRetries(1): %v", err)
	}

	req, err := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	resp, err := client.Do(req, nil)
	if err == nil {
		t.Fatal("expected error after retry exhaustion, got nil")
	}
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
	if errResp.Response.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429 status in ErrorResponse, got %d", errResp.Response.StatusCode)
	}
	if resp == nil {
		t.Error("expected non-nil Response returned alongside error")
	}
}

// TestWithRateLimitNotify_AttemptNumbers verifies the callback fires on each retry
// with 1-based sequential attempt numbers.
func TestWithRateLimitNotify_AttemptNumbers(t *testing.T) {
	callCount := 0

	ms, client, teardown := newMockServerAndClient(t)
	defer teardown()

	ms.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 3 {
			w.Header().Set("RateLimit-Reset", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	var notifyAttempts []int
	if err := WithRateLimitNotify(func(attempt int, _ time.Duration) {
		notifyAttempts = append(notifyAttempts, attempt)
	})(client); err != nil {
		t.Fatalf("WithRateLimitNotify: %v", err)
	}

	req, err := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	_, err = client.Do(req, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(notifyAttempts) != 2 {
		t.Errorf("expected notify called twice (2 retries), got %d times", len(notifyAttempts))
	}
	for i, attempt := range notifyAttempts {
		if want := i + 1; attempt != want {
			t.Errorf("notifyAttempts[%d] = %d, want %d", i, attempt, want)
		}
	}
}

// TestRetryDelay covers all branches of the retryDelay helper.
func TestRetryDelay(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		attempt int
		wantMin time.Duration
		wantMax time.Duration
	}{
		{
			name:    "RateLimit-Reset:0 gives 500ms buffer + jitter",
			headers: map[string]string{"RateLimit-Reset": "0"},
			attempt: 0,
			wantMin: 500 * time.Millisecond,
			wantMax: 1500 * time.Millisecond,
		},
		{
			name:    "RateLimit-Reset:5 gives 5.5s + jitter",
			headers: map[string]string{"RateLimit-Reset": "5"},
			attempt: 0,
			wantMin: 5500 * time.Millisecond,
			wantMax: 6500 * time.Millisecond,
		},
		{
			name:    "RateLimit-Reset > 120 is clamped to 120s",
			headers: map[string]string{"RateLimit-Reset": "999"},
			attempt: 0,
			wantMin: 120*time.Second + 500*time.Millisecond,
			wantMax: 120*time.Second + 1500*time.Millisecond,
		},
		{
			name:    "negative RateLimit-Reset falls back to exponential",
			headers: map[string]string{"RateLimit-Reset": "-1"},
			attempt: 0,
			wantMin: 1 * time.Second,
			wantMax: 2 * time.Second,
		},
		{
			name:    "Retry-After used when RateLimit-Reset absent",
			headers: map[string]string{"Retry-After": "0"},
			attempt: 0,
			wantMin: 500 * time.Millisecond,
			wantMax: 1500 * time.Millisecond,
		},
		{
			name:    "RateLimit-Reset takes precedence over Retry-After",
			headers: map[string]string{"RateLimit-Reset": "0", "Retry-After": "60"},
			attempt: 0,
			wantMin: 500 * time.Millisecond,
			wantMax: 1500 * time.Millisecond,
		},
		{
			name:    "exponential backoff at attempt 0 (no header)",
			headers: map[string]string{},
			attempt: 0,
			wantMin: 1 * time.Second,
			wantMax: 2 * time.Second,
		},
		{
			name:    "exponential backoff capped at 30s (attempt 10)",
			headers: map[string]string{},
			attempt: 10,
			wantMin: 30 * time.Second,
			wantMax: 31 * time.Second,
		},
		{
			name:    "no overflow at high attempt count (attempt 100)",
			headers: map[string]string{},
			attempt: 100,
			wantMin: 30 * time.Second,
			wantMax: 31 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{Header: make(http.Header)}
			for k, v := range tt.headers {
				resp.Header.Set(k, v)
			}
			got := retryDelay(resp, tt.attempt)
			if got < tt.wantMin || got >= tt.wantMax {
				t.Errorf("retryDelay() = %v, want in [%v, %v)", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}
