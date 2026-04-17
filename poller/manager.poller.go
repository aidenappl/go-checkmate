package poller

import (
	"context"
	"crypto/tls"
	"errors"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

type Spec struct {
	ID       string // unique key (e.g., req.Service + ":" + req.Name)
	Name     string
	Service  string
	Endpoint string
	Interval time.Duration
	Timeout  time.Duration
}

type Result struct {
	ID         string
	Endpoint   string
	StartedAt  time.Time
	Duration   time.Duration
	StatusCode int
	Err        error
}

type Sink interface {
	HandleResult(Result)
}

type Manager struct {
	mu     sync.Mutex
	jobs   map[string]*job
	client *http.Client
	sink   Sink
	ctx    context.Context
	cancel context.CancelFunc
}

type job struct {
	spec   Spec
	cancel context.CancelFunc
}

func New(sink Sink) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		jobs: make(map[string]*job),
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				ForceAttemptHTTP2:   true,
				TLSHandshakeTimeout: 5 * time.Second,
				TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
			},
		},
		sink:   sink,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *Manager) Close() { m.cancel() }

func (m *Manager) AddOrUpdate(spec Spec) error {
	if spec.Interval < time.Second {
		return errors.New("interval must be >= 1s")
	}
	if spec.Timeout == 0 || spec.Timeout > spec.Interval {
		spec.Timeout = minDuration(5*time.Second, spec.Interval/2)
	}
	if spec.ID == "" {
		return errors.New("missing ID")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// If exists, replace
	if existing, ok := m.jobs[spec.ID]; ok {
		existing.cancel()
		delete(m.jobs, spec.ID)
	}

	ctx, cancel := context.WithCancel(m.ctx)
	j := &job{spec: spec, cancel: cancel}
	m.jobs[spec.ID] = j

	go m.run(ctx, j)
	return nil
}

func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if j, ok := m.jobs[id]; ok {
		j.cancel()
		delete(m.jobs, id)
	}
}

func (m *Manager) run(ctx context.Context, j *job) {
	// small jitter so all jobs do not fire at once
	jitter := time.Duration(rand.Int64N(int64(j.spec.Interval / 5)))
	select {
	case <-time.After(jitter):
	case <-ctx.Done():
		return
	}

	ticker := time.NewTicker(j.spec.Interval)
	defer ticker.Stop()

	// do an immediate check once
	m.doCheck(ctx, j.spec)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.doCheck(ctx, j.spec)
		}
	}
}

func (m *Manager) doCheck(parent context.Context, spec Spec) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(parent, spec.Timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, spec.Endpoint, nil)
	req.Header.Set("User-Agent", "checkmate/1.0")

	res, err := m.client.Do(req)
	if err != nil {
		m.sink.HandleResult(Result{
			ID: spec.ID, Endpoint: spec.Endpoint,
			StartedAt: start, Duration: time.Since(start), Err: err,
		})
		return
	}
	_ = res.Body.Close()

	m.sink.HandleResult(Result{
		ID:         spec.ID,
		Endpoint:   spec.Endpoint,
		StartedAt:  start,
		Duration:   time.Since(start),
		StatusCode: res.StatusCode,
		Err:        nil,
	})
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
