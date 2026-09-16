---
title: Grafana Cloud Workshop
description: A 90-minute shared-stack introduction to metrics, logs, traces, profiles, and investigating an incident.
---

# Grafana Cloud workshop: one stack, 90 minutes

One instructor runs Synthkit on Kubernetes; everyone explores the **same Grafana Cloud stack**.
Learners need a browser and their own Grafana login, not Kubernetes access, ingest tokens, or the
Synthkit control password. This is a synthetic shop, not a deployed shop application: clicking a
website does not generate requests. Synthkit continuously generates the teaching data.

By the end, learners can choose an appropriate signal, scope a query, follow a request across
services, read a CPU profile, and distinguish missing instrumentation from broken ingestion.
No DOMM installation, qualification fixtures, Terraform, dashboards, alerts, users, teams, or
incident-management resources are created by this blueprint.

## Instructor preparation (before the 90 minutes)

1. Follow [Kubernetes deployment](kubernetes.md) with the `grafana-cloud-workshop` blueprint and
   all four ingest lanes configured. Start at least 15 minutes early so rate and comparison
   windows contain data. Check fresh metrics, logs, traces, and profiles in Grafana, not just a
   green Pod. Run the queries below yourself.
2. Give learners permission to query the four data sources and use Explore. Test the exact
   participant role beforehand; a dashboard-only account may lack the required access. Share
   the stack URL and the exact Prometheus, Loki, Tempo, and Pyroscope data source names. No
   participant needs a write token or permission to change data sources.
3. Keep the instructor-only [operator UI](control-plane.md) available over port-forward. Reset
   control state before the session, confirm only this blueprint is enabled, the volume
   multiplier is `1`, and no scenarios or ad-hoc failures are active. Historical telemetry stays.
4. Use one instructor for mutations. Learners work in Explore or personal scratch dashboards;
   do not overwrite shared dashboards or alter the data sources during class.
5. Rehearse the incident and recovery. Record a healthy absolute time interval and a degraded
   interval, so a transient ingest outage does not prevent discussing previously captured data.

### Optional click-through setup

Shared IDs make correlation possible; they do **not** configure Grafana links automatically.
An administrator should set these up and test them before class, or use the manual fallback:

- Loki to Tempo: a derived field of type **Label** matching `^trace_id$`, with an internal link
  to the workshop Tempo data source and query `${__value.raw}`. Here `trace_id` is structured
  metadata, not part of the JSON body or an indexed stream label.
- Tempo to Loki: map `service.name` to `service_name` and `service.namespace` to `namespace`.
  For a custom query, use `{${__tags}} | trace_id="${__trace.traceId}"` and allow a small time
  margin around the trace. Do not require matching `span_id`: these application logs carry the
  request's root span ID, not necessarily the selected downstream span's ID.
- Tempo to Pyroscope: select the workshop Pyroscope data source, map `service.name` to
  `service_name` and `service.namespace` to `service_namespace`, and select CPU profiles.
  Test an actual Go server span; Go span profiles are CPU-only.

These are instructor configuration tasks, not Helm side effects. Provisioned Cloud data sources
may be read-only; use the supported clone/provisioning workflow rather than changing production
configuration. See Grafana's [trace/log correlation instructions](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/configure-trace-to-logs/)
and [trace/profile configuration](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/configure-trace-to-profiles/).

## Teaching estate

The blueprint is `grafana-cloud-workshop`, the synthetic cluster is `workshop-prod`, and the
application namespace is `workshop-shop`. These labels describe simulated infrastructure; they
are not the namespace of the real Synthkit Pod.

| Service | Metrics | Application logs | Traces | Profiles |
|---|---|---|---|---|
| `shop-storefront` | Yes | Yes | Yes | Yes |
| `shop-checkout` | Yes | Yes | Yes | Yes |
| `shop-payment` | Yes | Yes | Yes | Yes |
| `shop-inventory` | Yes | No | No | No |
| `shop-shipping` | Yes | Yes | No | No |

The traced request graph is `shop-storefront → shop-checkout → shop-payment`. Inventory and
shipping are separate workloads with intentionally uneven adoption. Their absent signals are
not an exercise failure. Kubernetes infrastructure signals are separate from application
instrumentation; seeing infrastructure data does not prove that an application has tracing.

## Agenda

| Minutes | Activity | Learner outcome |
|---|---|---|
| 0–10 | Orientation | Choose a data source and time range; explain four signals |
| 10–25 | Metrics | Filter/group measurements and compare services |
| 25–40 | Logs | Inspect events and extract a request's trace ID |
| 40–55 | Traces | Read a request graph and correlate a log to its trace |
| 55–65 | Profiles | Find CPU-consuming functions and explain flame graph width |
| 65–85 | Investigation | Combine evidence, then verify recovery |
| 85–90 | Recap | Identify the next instrumentation improvement |

## 0–10: orientation

Open Explore, select the instructor's Prometheus data source, and set **Last 15 minutes**.
The instructor briefly introduces metrics as aggregate measurements, logs as individual events,
traces as request paths, and profiles as where execution time is spent. Grafana queries and
visualizes the data stored in the corresponding backends; changing the Explore data source
changes the query language.

In pairs, predict what you could investigate with inventory's metrics alone, and what you would
need logs, traces, or profiles to establish. Keep this question for the recap.

## 10–25: metrics

Select Prometheus in Explore, use code mode, and run:

```promql
go_goroutines{blueprint="grafana-cloud-workshop"}
```

Inspect `service`, `namespace`, and `cluster` labels. Narrow the selector with
`service="shop-checkout"`, then remove it to compare services. A gauge describes a sampled
state; it is not a count of requests.

Compare each service's recent mean observed HTTP duration, in seconds:

```promql
sum by (service) (
  rate(http_server_request_duration_seconds_sum{blueprint="grafana-cloud-workshop"}[5m])
)
/
sum by (service) (
  rate(http_server_request_duration_seconds_count{blueprint="grafana-cloud-workshop"}[5m])
)
```

Expected: nonempty measurements for the five services, with healthy latency comparatively stable.
Explain why counters need a time window and why dividing duration sum by observation count gives
a mean. These blueprint DSL histograms model observations; their count is **not** literal
end-user request throughput. Do not use it to assert real request volume or an error ratio.
No dashboard, recording rule, or Tempo metrics-generator configuration is required for these
queries. Immediately after startup/restart, allow enough samples for `rate()`.

Checkpoint: show a service-filtered graph and explain how missing data differs from zero.

## 25–40: logs

Switch to Loki and run:

```logql
{blueprint="grafana-cloud-workshop", source="app", service_name="shop-checkout"} | json
```

Expand a line. Inspect JSON fields such as `msg`, `route`, `status`, and `outcome`, alongside
the `level` stream label. Find `trace_id` and `span_id` in structured metadata. Keep one
`trace_id` for the next exercise. Trace IDs are deliberately not indexed stream labels.

Compare adoption:

```logql
{blueprint="grafana-cloud-workshop", source="app", service_name="shop-shipping"} | json
```

Shipping has logs, but no corresponding emitted traces. Inventory has no application log stream;
its absence is intentional. Never assume an identifier alone proves that a trace was ingested.

Checkpoint: explain what a request log adds to the metrics graph and why an empty inventory
log query does not prove a Loki outage.

## 40–55: traces and correlation

Select Tempo in Explore and run this TraceQL search:

```traceql
{ resource.service.name = "shop-checkout" && resource.service.namespace = "workshop-shop" }
```

Open a result. Locate the storefront, checkout, and payment spans; compare their durations and
parent/child relationships. A service dependency is not automatically the cause of a problem.
TraceQL's resource selectors scope the search to this teaching estate; see the official
[TraceQL query examples](https://grafana.com/docs/grafana/latest/datasources/tempo/query-editor/traceql-query-examples/).

Use the log's trace link if configured. Otherwise switch to Tempo's **Trace ID** query mode and
paste the `trace_id` collected in the logs exercise. To move back manually, select Loki and use
the same trace's time range, replacing the placeholder below:

```logql
{blueprint="grafana-cloud-workshop", source="app"} | trace_id="REPLACE_WITH_TRACE_ID"
```

Expected: events sharing that request ID. A trace-ID pipeline filter uses structured metadata;
`|= "TRACE_ID"` searches the body and is not a suitable fallback here. If links fail but the
manual query works, fix correlation configuration, not ingestion.

Checkpoint: identify one request across two signals and state which services lack this option.

## 55–65: profiles

Open Profiles Drilldown or select Pyroscope in Explore. Select profile type
`process_cpu:cpu:nanoseconds:cpu:nanoseconds` and filter:

```text
{service_name="shop-checkout", service_namespace="workshop-shop"}
```

Use the same time range as your other signals. Inspect the widest frames and their callers.
Width represents CPU contribution, not a request timeline; compare total and self contribution.
See Grafana's [profile query editor](https://grafana.com/docs/grafana/latest/datasources/pyroscope/query-editor/).

If the instructor configured trace-to-profile links, open a checkout server span and inspect
its CPU profile. Otherwise query the service and time range manually: this is a service/time
comparison, not proof of a specific span's profile. Memory or goroutine profiles are not a
substitute for the CPU span-profile exercise.

Checkpoint: explain what a CPU-heavy function tells you that a slow trace alone cannot.

## 65–85: a bounded investigation

Only the instructor changes Synthkit state. This drill changes synthetic telemetry, not real
shop services or the Kubernetes host's CPU consumption. It creates no Grafana incident object.

1. **Minutes 65–68:** capture a healthy interval. In the operator UI's Scenarios view, activate
   `grafana-cloud-workshop/checkout-regression` and note the activation time. It combines
   checkout-targeted latency, error, and CPU-hotspot modes. There is no random activation.
2. **Minutes 68–76:** participants investigate in pairs. Compare checkout's mean-duration
   graph with payment. Find checkout error logs using the query below. Search for checkout
   error spans, inspect the request path, and compare checkout CPU profiles before/during
   the change. Write down observations separately from hypotheses.
3. **At minute 76:** the instructor deactivates `checkout-regression` in the UI, even if some
   groups are still investigating. This bounds the live fault to about eight minutes; the
   scenario does not automatically expire. Groups can continue using the recorded interval.
4. **Minutes 76–85:** refresh recent data and verify recovery. The five-minute metrics window
   mixes healthy and degraded samples until it advances. Logs and traces retain the past
   errors; their continued presence in a historical window is not an active failure. Compare
   the recorded healthy/degraded/recovered intervals and present one evidence-backed finding.

```logql
{blueprint="grafana-cloud-workshop", source="app", service_name="shop-checkout", level="error"} | json
```

```traceql
{ resource.service.name = "shop-checkout" && resource.service.namespace = "workshop-shop" && status = error }
```

Facilitator notes: expected checkout evidence is higher modeled latency, request errors, and
CPU-hotspot profile changes. Some request-level effects can appear elsewhere in the same graph;
do not teach that every correlated error proves a second independent fault. The scenario is
deliberately constructed evidence, not a guarantee that synthetic measurements precisely match
one another numerically or reproduce a real application's causal performance model.

## 85–90: recap and repeat

Ask each pair: Which signal first showed the problem? Which added an explanation? What could
you have concluded for metrics-only inventory? What instrumentation would you add next?

The instructor confirms the scenario is inactive and restores defaults using the control reset
described in [Kubernetes operations](kubernetes.md#reset-the-exercise-not-the-telemetry). Record
the session's time window if it will be reused. A new class uses fresh windows; it does not need
telemetry deletion. Participants may retain query links without saving shared dashboards.

## When data is missing

| Observation | Check first |
|---|---|
| Inventory has no logs/traces/profiles | Expected adoption boundary in the table above |
| Shipping has no traces/profiles | Expected; logs alone do not imply traced requests |
| All services lack one signal | Correct data source/time range, then instructor sink status and credentials |
| All signals absent | Instructor checks dry-run mode, selected blueprint, paused Pod, and readiness |
| Logs and traces exist but links fail | Manual trace-ID lookup, then data-source mappings |
| Profiles absent only | Profile credentials, expected service/type, then delivery status; never assume a green metrics lane proves profiles |
| Old errors remain after recovery | Query a fresh interval or compare explicit before/after windows |

For implementation provenance, the exact declarations are in
[`grafana-cloud-workshop.yaml`](https://github.com/mbaykara/Synthkit/blob/main/blueprints/grafana-cloud-workshop.yaml).
The [control plane](control-plane.md) describes diagnostics; [troubleshooting](troubleshooting.md)
describes delivery failures. Learners report problems to the instructor rather than modifying
tokens, Helm values, or shared control state.
