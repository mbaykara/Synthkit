---
title: Grafana Cloud Foundation in Action (English)
description: Grafana Cloud Foundation in Action – three labs on one shared stack in 90 minutes.
---

# Grafana Cloud Foundation in Action

[Deutsche Fassung](workshop.md)

Other sessions: [Alerting](workshop-alerting-en.md), [Assistant](workshop-assistant-en.md)

**One shared stack, 90 minutes, three labs:** find a signal, test a question against two signal
types, and keep the useful path as a personal quickstart. The labs belong to deck slides
**9, 14 and 27**. The timings below are a facilitation suggestion; the slide numbers come from the
lab cards.

The facilitator runs Synthkit on Kubernetes. All participants use the same Grafana Cloud stack with
their own Grafana login and browser. Kubernetes access, the ingest token and the Synthkit control
password stay with the facilitator. The shop is simulated: there is no shop website whose visits
generate requests. Synthkit delivers the telemetry continuously.

## What you take away

- Find the service, dashboard and time range that matter.
- Investigate a question across metrics and a second signal type; place traces and profiles deliberately.
- Validate queries and findings before you act.
- Reuse what works: first as a personal Assistant prompt, later with reviewed team standards.

## Facilitator: preparation before the 90 minutes

The `grafana-cloud-workshop` blueprint supplies the data. Helm does **not** create Grafana dashboards,
Explore links, Assistant quickstarts, users, teams, alert rules or DOMM fixtures, and it does not
enable Assistant. The following stack artefacts must be prepared and checked before the session.

1. Start the [Kubernetes deployment](kubernetes.md) with all four ingest signals. Plan at least
   **30 minutes of lead time**, longer on first setup. Run the queries further down yourself and
   check for fresh data in metrics, logs, traces and profiles; a green pod is not enough. For the
   optional comparison with yesterday, yesterday's window must also contain data.
2. Test the dashboard, Explore and all four data sources with the actual participant role.
   Check Assistant access and creating a personal quickstart separately. If they are not
   available, announce the card fallback for Lab 3; it does not replace evidence of a saved prompt.
3. Through the facilitator-only [control plane access](control-plane.md),
   [reset the exercise state](kubernetes.md#reset-the-exercise-not-the-telemetry): workshop blueprint only,
   volume multiplier `1`, no active scenarios or ad-hoc failures. Historical telemetry is retained.
4. Create a scenario dashboard, for example **Foundation – Observe checkout**. At least one
   time series panel with the metric query from Lab 2, unit seconds, service `shop-checkout`.
   Name the panel "Checkout: mean modelled HTTP duration" and highlight it as the entry point.
   Show the signal, but do not give away a cause. An existing dashboard is also suitable if it
   queries exactly this synthetic data. No dashboard JSON is shipped here.
5. Collect healthy history, then activate `grafana-cloud-workshop/checkout-regression` in the
   Scenarios view. Note the activation time, let it run for about **eight minutes**, then
   deactivate it explicitly. The scenario does not expire on its own. It combines checkout
   latency, errors and a CPU hotspot; it changes synthetic data, not real shop systems or the
   CPU load of the Kubernetes host, and it creates no Grafana incident object.
6. After the ingest delay, confirm the change is actually visible. Choose an **absolute window with
   date and time zone** that contains both healthy history and the change. Example sequence:
   ten minutes healthy, eight minutes of scenario, then at least five minutes of recovery.
   The five-minute `rate()` range smooths transitions. Save or share the window in the dashboard
   and in Explore; the signal must already be visible before Lab 1. Do not activate it during Lab 2.
7. Copy the dashboard link and two prepared Explore links from the Grafana UI: metrics and logs
   (or traces), each for the same service and absolute time range. Open them with a participant
   account and check that data source, filters and time window are preserved.
   Do not hand out invented data source UIDs or example URLs to participants.

### Hand out the completed session sheet

The facilitator replaces the open fields before the start. `workshop-shop` is the **synthetic**
service namespace, not the namespace of the actual Synthkit pod.

| Field on the lab card | Value for this session |
|---|---|
| `{service}` | `shop-checkout` |
| `{namespace}` | `workshop-shop` |
| `{window}` | Enter the actually verified start and end time with date and time zone |
| `{dashboard}` | Enter the link to the verified checkout dashboard |
| Data sources | Enter the actual names of Prometheus, Loki, Tempo and Pyroscope |
| Explore fallback | Enter the verified metric and log/trace links for the same window |
| Assistant | Available and personal quickstarts tested, or card fallback |

**Start check:** without a visible signal and verified dashboard/Explore links, Lab 1/2 is not
ready. If ingest is disrupted, use a previously verified, still saved window and state its date
openly. Only one person changes the generator; participants change neither shared dashboards,
data sources nor control state.

## Agenda: 90 minutes

| Minutes | Section | Outcome |
|---|---|---|
| 0–10 | Orientation, stack and session sheet | Service, dashboard, data sources and window found |
| 10–25 | Lab 1, slide 9: Find your signal | One question worth investigating, no diagnosis yet |
| 25–30 | Debrief | Impact in business language, uncertainty named |
| 30–55 | Lab 2, slide 14: Investigate and validate | Two executed queries across two signal types |
| 55–65 | Discuss evidence, trace/profile demo | Explain the claim and limits of additional signals |
| 65–85 | Lab 3, slide 27: Create a personal quickstart | Save personally, run, check and refine |
| 85–90 | Close and reuse | Next use and limits named |

## Lab 1 · Deck slide 9 · Find your signal

### Scenario

You have taken over `shop-checkout`. Something changed during the `{window}` given on the
session sheet. You do not fix it yet and you do not diagnose it yet: you choose a question
worth following.

### Steps

1. Open a service dashboard you know, if it shows the workshop data; otherwise `{dashboard}`.
2. Set `{window}`. For the shared investigation, use the absolute time range from the session
   sheet, not a moving "Last 15 minutes".
3. Name one signal and one question. Keep service and time range for Lab 2.

### Expected result and completion check

One sentence with service, time range and signal, for example:
"For `shop-checkout`, the mean modelled HTTP duration rises in the marked window;
do errors also occur in the same window?" Add the actual times from your view.
You can name the **service, time range and question**. No cause has been established yet.

### In the debrief

Phrase the possible business meaning: "Checkout was slower in the observed window; whether
purchases were affected still needs to be checked."
The statement "a fifth of customers were affected for twenty minutes" would be inadmissible
without further evidence. This blueprint provides neither real customer numbers nor lost revenue.
Observed duration, possible impact and unknown extent stay separate.

### If you get stuck

Open the prepared scenario dashboard with the highlighted latency panel.
First describe only **what** changed **when** – not why.

## Lab 2 · Deck slide 14 · Investigate and validate

### Scenario and steps

Take your question from Lab 1 into Explore. Keep `shop-checkout` and `{window}` fixed throughout.
Do not change service and time range at the same time.

1. Open the relevant dashboard panel in Explore or the prepared metric link.
2. Check or refine the scoped query. Explain what it actually measures.
3. **You must switch to a second signal type:** from metrics to logs or traces.
   Query builder and Assistant may help with drafting; only the executed query and its result
   are evidence. Compare the same service in the same time window.

| Investigation checklist | Record for this session |
|---|---|
| Service | `shop-checkout`, namespace `workshop-shop` |
| Time range | `{window}` including date and time zone |
| Query | Data source, filters, measure and unit |
| History | Keep the two useful queries and their results or links |

### First query: metric

Run in the Prometheus data source:

```promql
sum by (service) (
  rate(http_server_request_duration_seconds_sum{blueprint="grafana-cloud-workshop",service="shop-checkout",namespace="workshop-shop"}[5m])
)
/
sum by (service) (
  rate(http_server_request_duration_seconds_count{blueprint="grafana-cloud-workshop",service="shop-checkout",namespace="workshop-shop"}[5m])
)
```

The ratio of duration to observation count gives the **mean modelled HTTP duration in seconds**,
not p95, error rate or the duration of each individual request. `[5m]` is the calculation window
at each graph point, not the whole Explore time range. After a restart, `rate()` needs enough
samples. The DSL histogram count is **not real user or request throughput**; do not derive a
customer count or error ratio from it. Missing data is not zero.

### Second query: logs

Switch to the Loki data source, leave `{window}` unchanged and run:

```logql
{blueprint="grafana-cloud-workshop",source="app",service_name="shop-checkout",namespace="workshop-shop"} | json
```

Check `msg`, `route`, `status` and `outcome` in the JSON, and the stream label `level`.
For the question about errors, then narrow down:

```logql
{blueprint="grafana-cloud-workshop",source="app",service_name="shop-checkout",namespace="workshop-shop",level="error"} | json
```

The metric shows the aggregated trend; logs show individual events and their details.
Errors in the same window can support a hypothesis, but they do not prove a cause of the latency.
An empty error result only meaningfully refutes the hypothesis if the general log query returns
data and data source, service and time window are correct.

### Alternative or deeper dive: traces

In Tempo, with `{window}` unchanged:

```traceql
{ resource.service.name = "shop-checkout" && resource.service.namespace = "workshop-shop" }
```

Open a result and look at the checkout spans in the request path
`shop-storefront → shop-checkout → shop-payment`. For error spans, also add
`&& status = error` inside the braces. Dependency and coincidence in time are not proof of cause.
Synthetic metrics and spans need not match numerically.

Logs carry `trace_id` and `span_id` as **structured metadata**, not in the JSON body or as
indexed stream labels. Without preconfigured links: copy the ID from the log and search for it
in Tempo's trace ID mode. Back in Loki, in the same time range:

```logql
{blueprint="grafana-cloud-workshop",source="app",service_name="shop-checkout",namespace="workshop-shop"} | trace_id="REPLACE_WITH_TRACE_ID"
```

`|= "TRACE_ID"` searches the body and is not a substitute here. A trace ID alone does not
guarantee that the trace was actually ingested.

### Expected result and completion check

Show **two executed queries across two signal types**, both for the same service and the same
window. For each query, note what it says and its limit, for example:
"The metric shows increased mean duration; the logs show error events in the same window.
That supports a simultaneous degradation, but does not yet explain its cause."
A contradiction is a good result: you tested your idea instead of just confirming it.

### If you get stuck

Use the prepared Explore links and compare the two named signals.
Let Assistant explain or draft the query, but check filters, time range, unit and result yourself.
Switching to the second signal is not an optional extra.

## 55–65: Additional evidence with traces and profiles

The facilitator briefly shows the trace path and opens Profiles Drilldown or Pyroscope in Explore.
Profile type `process_cpu:cpu:nanoseconds:cpu:nanoseconds`, same `{window}`, selector:

```text
{service_name="shop-checkout",service_namespace="workshop-shop"}
```

Which functions contribute to CPU time? The width of a flame graph bar stands for CPU share,
not for a request timeline. A CPU hotspot alone does not explain every latency. Without a tested
trace-to-profile link, we compare service and time range, not the profile of a specific span.
The Go span profiles in this exercise are CPU-only.

The demo complements Lab 2; it does not replace the two queries you ran yourself.

## Lab 3 · Deck slide 27 · Create a personal quickstart

### Scenario and steps

You will repeat this investigation. Save a prompt so that the next run already has context.
Keep it personal for now.

1. Choose a recurring task for `shop-checkout`, such as a morning health check. For the first
   validation, use the known `{window}` from Lab 2.
2. Write **action + service scope + time range + @context**. After typing `@`, pick the data source
   or dashboard that actually exists from the selection. `@workshop-shop-logs` is only valid if
   that object really exists; the namespace itself does not create a data source. Replace the
   context placeholders below before saving.
3. Under **Grafana Assistant → Settings → Quickstart prompts → Create Quickstart Prompt**, enter
   title and prompt, set **Scope: Just me**, leave Enabled switched on and save.
4. Run the saved quickstart. Check at least one generated query and its result: service, window,
   data source, measure and actual evidence. Compare with Lab 2. Refine the prompt if needed and
   run it again. Do not switch it to **Everybody**.

Saving personally and choosing context follow the
[official Assistant guide](https://grafana.com/docs/grafana-cloud/platform/grafana-assistant/introduction/).
The facilitator checks [enablement and access](https://grafana.com/docs/grafana-cloud/platform/grafana-assistant/get-started/grafana-cloud/)
in advance; UI availability and permissions can differ per stack.

### Two worked prompts

**Repeat the investigation – first with the absolute lab window:**

> Check shop-checkout in namespace workshop-shop during {window} with @PROMETHEUS and @LOKI.
> Show the mean modelled HTTP duration and matching error logs. Use the same service and time
> filters. Show both queries and results; separate observations, hypotheses and open questions.
> Do not make any changes.

`@PROMETHEUS` and `@LOKI` are placeholders for real selected context objects.
After a successful check, a personal health-check variant can use "last 30 minutes".
It may show healthy data after the scenario is deactivated – that is not a malfunction.
A quickstart is not an automatically running monitor.

**Compare – only with existing history:**

> Compare the mean modelled HTTP duration of shop-checkout in namespace workshop-shop
> in the last completed hour with the same hour yesterday using @PROMETHEUS.
> Show the exact time windows, queries and differences. Point out missing data instead of
> treating it as zero. Do not claim a cause without further evidence.

Only use this comparison if both windows actually contain data. Otherwise compare two verified
healthy/degraded time windows from the same session. Do not derive an "error rate" from the DSL
histogram count; for a start, investigate error events. A ratio from logs would at most be a ratio
of generated log events, not of customers.

### Expected result and completion check

A **personally saved and executed prompt**, at least one checked result and one named
refinement or a reasoned decision to keep it as it is.
You can explain when you will use it again and what data it needs.
Share it with the team only when it works reliably and someone owns it.

### If you get stuck

Draft the prompt on the lab card and talk through the context selection, the expected query and
the result check out loud. Without Assistant access or save permission, this is a deliberate
fallback, **not** a successfully saved quickstart. Do not improvise access during the lab by
creating new shared permissions.

## Close and next run

Which question did you test? What could the second signal add or refute?
Which statement remained unsupported? Which personal quickstart would you reuse?

The facilitator confirms that the scenario is deactivated and resets the
[exercise state](kubernetes.md#reset-the-exercise-not-the-telemetry). The generator may keep
running. Old telemetry does not need to be deleted and is no obstacle for the next run.
For a new class, prepare fresh windows and re-check all links;
a pod restart alone does not replace a reset of the persistent control state.

## Reference: intentionally uneven instrumentation

Blueprint `grafana-cloud-workshop`, synthetic cluster `workshop-prod`, namespace `workshop-shop`:

| Service | Metrics | App logs | Traces | Profiles |
|---|---|---|---|---|
| `shop-storefront` | Yes | Yes | Yes | Yes |
| `shop-checkout` | Yes | Yes | Yes | Yes |
| `shop-payment` | Yes | Yes | Yes | Yes |
| `shop-inventory` | Yes | No | No | No |
| `shop-shipping` | Yes | Yes | No | No |
| `shop-catalog` | Yes (JVM gauges only, 12 pods) | No | No | No |

`shop-catalog` backs the Alerting session and is not used in the Foundation labs.

Optional transfer after the labs: what could you say about inventory with metrics alone?
Missing application signals are intentional here. Infrastructure telemetry does not prove
application instrumentation; a log ID does not prove an available trace.

### Optional links and troubleshooting for the facilitator

Shared IDs do not configure Grafana links. Test them before the session or use the manual
trace ID search from Lab 2:

- Loki → Tempo: derived field of type label, matcher `^trace_id$`, internal Tempo link,
  query `${__value.raw}`. The ID comes from structured metadata.
- Tempo → Loki: map `service.name` to `service_name`, `service.namespace` to `namespace`;
  custom query `{${__tags}} | trace_id="${__trace.traceId}"`. Allow a small time buffer.
  No mandatory `span_id` match: app logs carry the root span ID of the request.
- Tempo → Pyroscope: map `service.name` to `service_name`, `service.namespace` to
  `service_namespace`, select the CPU profile and test on a Go server span.

See [Trace to logs](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/configure-trace-to-logs/)
and [Trace to profiles](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/configure-trace-to-profiles/).
Provisioned data sources can be read-only; use the supported provisioning or clone path and do not
change production configuration during the session.

| Observation | Check first |
|---|---|
| Inventory without logs or shipping without traces | Expected limits per the table |
| All services missing one signal type | Data source and window, then ingest status and credentials |
| Everything empty | Dry run, blueprint selection, paused pod and readiness |
| Logs/traces present, links broken | Manual ID search, then data source mappings |
| Only profiles missing | Profile credentials, service/type and delivery status |
| Old errors visible after recovery | A historical window still shows historical errors |

The source of the data shapes is the
[`grafana-cloud-workshop` blueprint](https://github.com/mbaykara/Synthkit/blob/main/blueprints/grafana-cloud-workshop.yaml).
Further diagnosis: [Control Plane](control-plane.md) and [Troubleshooting](troubleshooting.md).
Participants report problems to the facilitator instead of changing tokens, Helm values or shared resources.
