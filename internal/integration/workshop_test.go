// SPDX-License-Identifier: AGPL-3.0-only

package integration

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rknightion/synthkit/internal/blueprint"
	"github.com/rknightion/synthkit/internal/control"
	"github.com/rknightion/synthkit/internal/core/coretest"
	"github.com/rknightion/synthkit/internal/runner"
	"github.com/rknightion/synthkit/internal/sink/otlp"
	psink "github.com/rknightion/synthkit/internal/sink/pyroscope"
)

type workshopProfiles struct {
	mu     sync.Mutex
	series []psink.Series
}

func (p *workshopProfiles) Write(_ context.Context, batch []psink.Series) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.series = append(p.series, batch...)
	return nil
}

type workshopCapture struct {
	m *coretest.MetricCapture
	l *coretest.LogCapture
	t *coretest.TraceCapture
	p *workshopProfiles
	r *runner.Runner
}

func runWorkshop(t *testing.T, now time.Time, incident bool) workshopCapture {
	t.Helper()
	data, err := os.ReadFile("../../blueprints/grafana-cloud-workshop.yaml")
	if err != nil {
		t.Fatal(err)
	}
	res, err := blueprint.Load(data, runner.Catalog())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Incidents) != 0 {
		t.Fatal("workshop must not activate incidents automatically")
	}
	c := workshopCapture{m: &coretest.MetricCapture{}, l: &coretest.LogCapture{}, t: &coretest.TraceCapture{}, p: &workshopProfiles{}}
	r := runner.New(runner.Sinks{Metrics: c.m, Logs: c.l, Traces: c.t, Profiles: c.p}, runner.Catalog(), runner.Options{})
	c.r = r
	if err := r.AddBlueprint(res); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { r.DrainQueues(context.Background()) })
	state := control.DefaultState()
	if incident {
		state.ActiveScenarios = []string{"grafana-cloud-workshop/checkout-regression"}
	}
	r.ApplyControl(state)
	for i := range 12 {
		if err := r.MasterTick(context.Background(), now.Add(time.Duration(i*5)*time.Second)); err != nil {
			t.Fatal(err)
		}
	}
	if err := r.RunOnce(context.Background(), now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestWorkshopAdoptionAndCorrelation(t *testing.T) {
	// A weekend midnight is deliberate: the workshop must not depend on business hours.
	c := runWorkshop(t, time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC), false)
	metrics, logs, traces, profiles := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
	traceIDs, spanIDs := map[string]bool{}, map[string]map[string]bool{}
	runtimeMetrics := map[string]bool{}
	for _, s := range c.m.All() {
		if strings.HasPrefix(s.Name, "aws_") {
			t.Fatalf("unrelated cloud telemetry: %s", s.Name)
		}
		if s.Name == "go_goroutines" {
			runtimeMetrics[s.Labels["service"]] = true
		}
	}
	for _, s := range c.m.Find("http_server_request_duration_seconds_count") {
		metrics[s.Labels["service"]] = true
	}
	for _, r := range c.t.Resources {
		service, _ := r.Attrs["service.name"].(string)
		traces[service] = true
		if spanIDs[service] == nil {
			spanIDs[service] = map[string]bool{}
		}
		for _, s := range r.Spans {
			traceIDs[s.TraceID] = true
			spanIDs[service][s.SpanID] = true
		}
	}
	correlatedLogs, correlatedProfiles := 0, 0
	for _, stream := range c.l.Streams {
		service := stream.Labels["service_name"]
		logs[service] = true
		for _, line := range stream.Lines {
			if id := line.Meta["trace_id"]; id != "" {
				if !traceIDs[id] {
					t.Fatalf("log links to absent trace %s", id)
				}
				correlatedLogs++
			}
		}
	}
	for _, s := range c.p.series {
		labels := map[string]string{}
		for _, l := range s.Labels {
			labels[l.Name] = l.Value
		}
		service := labels["service_name"]
		profiles[service] = true
		if labels["__name__"] != "process_cpu" || labels["source"] != "" || labels["pyroscope_spy"] != "gospy" {
			t.Fatalf("unexpected CPU SDK profile labels: %v", labels)
		}
		for _, sample := range s.Profile.Sample {
			for _, label := range sample.Label {
				if s.Profile.StringTable[label.Key] == "span_id" {
					id := s.Profile.StringTable[label.Str]
					if !spanIDs[service][id] {
						t.Fatalf("%s profile links to absent local span %s", service, id)
					}
					correlatedProfiles++
				}
			}
		}
	}
	for _, service := range []string{"shop-storefront", "shop-checkout", "shop-payment", "shop-inventory", "shop-shipping"} {
		if !metrics[service] || !runtimeMetrics[service] {
			t.Errorf("%s has no application metrics", service)
		}
		wantLogs := service != "shop-inventory"
		wantTracing := service != "shop-inventory" && service != "shop-shipping"
		if logs[service] != wantLogs || traces[service] != wantTracing || profiles[service] != wantTracing {
			t.Errorf("%s adoption logs=%v traces=%v profiles=%v", service, logs[service], traces[service], profiles[service])
		}
	}
	if correlatedLogs == 0 || correlatedProfiles == 0 {
		t.Fatalf("missing correlation: logs=%d profiles=%d", correlatedLogs, correlatedProfiles)
	}
	t.Logf("workshop capture: %d metric series, %d names, %d trace resources, %d profile series", len(c.m.All()), len(c.m.Names()), len(c.t.Resources), len(c.p.series))
}

func TestWorkshopScenario(t *testing.T) {
	now := time.Date(2026, 9, 16, 11, 0, 0, 0, time.UTC)
	base, hot := runWorkshop(t, now, false), runWorkshop(t, now, true)
	latency := func(c workshopCapture, service string) float64 {
		var sum, count float64
		for _, s := range c.m.Find("http_server_request_duration_seconds_sum") {
			if s.Labels["service"] == service {
				sum = s.Value
			}
		}
		for _, s := range c.m.Find("http_server_request_duration_seconds_count") {
			if s.Labels["service"] == service {
				count = s.Value
			}
		}
		if count == 0 {
			t.Fatalf("%s has no histogram observations", service)
		}
		return sum / count
	}
	if latency(hot, "shop-checkout") < 2*latency(base, "shop-checkout") {
		t.Fatal("scenario does not increase checkout metric latency")
	}
	if latency(hot, "shop-inventory") != latency(base, "shop-inventory") {
		t.Fatal("scenario changes unrelated inventory metrics")
	}
	errors := func(c workshopCapture) int {
		n := 0
		for _, s := range c.l.Streams {
			if s.Labels["service_name"] == "shop-checkout" {
				for _, l := range s.Lines {
					var body map[string]any
					if err := json.Unmarshal([]byte(l.Body), &body); err != nil {
						t.Fatal(err)
					}
					if body["status"] == "500" {
						n++
					}
				}
			}
		}
		return n
	}
	if errors(hot) <= errors(base) {
		t.Fatal("scenario does not increase checkout errors")
	}
	traceStats := func(c workshopCapture) (time.Duration, int) {
		var duration time.Duration
		count, failures := 0, 0
		for _, r := range c.t.Resources {
			if r.Attrs["service.name"] == "shop-checkout" {
				for _, s := range r.Spans {
					if s.Kind == otlp.KindServer {
						duration += s.End.Sub(s.Start)
						count++
						if s.Status == otlp.StatusError {
							failures++
						}
					}
				}
			}
		}
		if count == 0 {
			t.Fatal("missing checkout server spans")
		}
		return duration / time.Duration(count), failures
	}
	baseDuration, baseErrors := traceStats(base)
	hotDuration, hotErrors := traceStats(hot)
	if hotDuration <= baseDuration || hotErrors <= baseErrors {
		t.Fatal("scenario does not degrade checkout trace latency and errors")
	}
	cpu := func(c workshopCapture) int64 {
		var total int64
		for _, s := range c.p.series {
			checkout := false
			for _, l := range s.Labels {
				if l.Name == "service_name" && l.Value == "shop-checkout" {
					checkout = true
				}
			}
			if checkout {
				for _, sample := range s.Profile.Sample {
					total += sample.Value[0]
				}
			}
		}
		return total
	}
	if cpu(hot) <= cpu(base) {
		t.Fatal("scenario does not increase checkout CPU profiles")
	}
	// Deactivation on the same runner must restore subsequent observations without
	// deleting historical telemetry or resetting the cumulative histogram.
	beforeSum, beforeCount := workshopHistogram(hot, "shop-checkout")
	hot.r.ApplyControl(control.DefaultState())
	if err := hot.r.RunOnce(context.Background(), now.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	afterSum, afterCount := workshopHistogram(hot, "shop-checkout")
	if afterCount <= beforeCount {
		t.Fatal("no observations after scenario deactivation")
	}
	recovered := (afterSum - beforeSum) / (afterCount - beforeCount)
	if recovered >= latency(base, "shop-checkout")*1.5 {
		t.Fatalf("checkout did not recover: %f", recovered)
	}
}

func workshopHistogram(c workshopCapture, service string) (sum, count float64) {
	for _, s := range c.m.Find("http_server_request_duration_seconds_sum") {
		if s.Labels["service"] == service {
			sum = s.Value
		}
	}
	for _, s := range c.m.Find("http_server_request_duration_seconds_count") {
		if s.Labels["service"] == service {
			count = s.Value
		}
	}
	return sum, count
}
