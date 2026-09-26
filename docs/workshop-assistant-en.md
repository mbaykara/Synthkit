---
title: Grafana Assistant in Action (English)
description: Grafana Assistant in Action – three labs on one shared stack in 60 minutes.
---

# Grafana Assistant in Action

[Deutsche Version](workshop-assistant-de.md) · Other sessions:
[Foundation](workshop-foundation-en.md) ([Deutsch](workshop.md)) ·
[Alerting](workshop-alerting-en.md) ([Deutsch](workshop-alerting-de.md))

**One shared stack, 60 minutes, three labs:** build a dashboard with Assistant, investigate a
symptom to a hypothesis you can defend, and make your method repeatable as a Rule, a Skill and an
Automation. The labs belong to deck slides **7, 11 and 17**. The timings below are a facilitation
suggestion; the slide numbers come from the lab cards.

This session stands alone. Attendees who missed Foundation need nothing from it: the cards supply
the service and every prompt. The facilitator runs Synthkit on Kubernetes; everyone works in the
same Grafana Cloud stack with their own login. The shop is simulated. Assistant drafts and
explores, the engineer validates: every lab includes a step where you check the work yourself.

## What you take away

- Build and iterate a dashboard from prompts, and validate the queries behind it.
- Run an evidence-backed investigation and validate a finding yourself in Explore.
- State the impact only after validating it, and name its limits.
- Create an Assistant Rule, a Skill and a scheduled Automation for your method.
- Treat shared AI artefacts responsibly: personal first, shared only with an owner.

## Facilitator: preparation before the 60 minutes

The `grafana-cloud-workshop` blueprint supplies the data. Helm does not enable Assistant and does
not create dashboards, Rules, Skills, Automations, users or annotations. Check the following no
later than the day before.

1. Keep the [Kubernetes deployment](kubernetes.md) running with all four ingest signals. The
   namespace needs **24 hours** of live metrics, logs and traces: start the generator at least a
   day before the session. Synthkit does not backfill history.
2. Assistant is enabled and entitled on the stack. Run one prompt end to end yourself. With an
   attendee account, confirm that **Settings → Rules, Skills and Automations** are reachable;
   Lab 3 needs all three.
3. **Important:** delete any leftover tenant-wide Rules (scope **Everybody**) from previous
   cohorts. They silently change how Assistant behaves for the next group.
4. Check whether **Investigations** is available. If not, brief the conversational fallback on the
   Lab 2 card.
5. Through the facilitator-only [control plane](control-plane.md),
   [reset the exercise state](kubernetes.md#reset-the-exercise-not-the-telemetry): workshop
   blueprint only, volume multiplier `1`, no active scenarios. Historical telemetry is kept.
6. Seed the Lab 2 fault: after healthy history, activate `grafana-cloud-workshop/payment-regression`
   in the Scenarios view. Record the activation time with **date and timezone** as `{fault-time}`,
   let it run for about **ten minutes**, then deactivate it explicitly. The scenario does not
   expire on its own. It raises failed payment requests and `shop-payment` response time; it
   changes synthetic data only, not real systems. Do not activate it during Foundation or Alerting
   on the same stack unless you intend to.
7. Prepare the `{change}`. Synthkit emits **no** deployment or change event. Either state the
   change verbally (for example "a payment release at `{fault-time}`") or create a Grafana
   annotation at that time beforehand. Otherwise Assistant cannot find the change in telemetry;
   say so openly.
8. After ingest delay, verify the change yourself: more errors in the `shop-payment` logs and a
   higher mean duration in the same window.
9. Send attendees the filled lab card pack with the invitation.

**Naming warning:** in Lab 3, "Rule" means an **Assistant behaviour rule**, not an alert rule.
Say this explicitly; the mix-up is common.

### Hand out the filled session sheet

| Field on the lab card | Value for this session |
|---|---|
| `{service.namespace}` | `workshop-shop` |
| `{service}` | `shop-payment` |
| `{fault-time}` | Actual activation time with date and timezone |
| `{change}` | The change stated verbally or the prepared annotation, e.g. a payment release |
| `{dashboard name - lab 1}` | `Shop - Application Health` plus your initials |
| Assistant | Enabled and tested; Investigations available or conversational fallback |

`workshop-shop` is the **synthetic** service namespace, not the namespace of the actual Synthkit
pod. Signals available for `shop-payment`: metrics (`http_server_request_duration_seconds` with
the labels `service`, `service_name`, `namespace`), app logs with `status` and `outcome` plus
`trace_id` as structured metadata, traces along `shop-storefront → shop-checkout → shop-payment`,
and CPU profiles.

## Schedule: 60 minutes

| Minutes | Section | Outcome |
|---|---|---|
| 0–5 | Assistant accelerates a loop you already use; what Assistant is | A context-aware agent inside Grafana, under existing RBAC |
| 5–20 | Lab 1, slide 7: Build with Assistant | A saved dashboard where every query can be explained |
| 20–27 | Investigations and prompt craft | Naming symptom, scope, timeframe and what changed |
| 27–42 | Lab 2, slide 11: Investigate a symptom | A defended hypothesis and one impact statement |
| 42–47 | Rules, Skills, Automations, MCP | Which one fits when; governance |
| 47–57 | Lab 3, slide 17: Make it repeatable | Rule, Skill and Automation, as far as time allows |
| 57–60 | Close | Next use named |

## Lab 1 · Deck slide 7 · Build with Assistant

### Scenario

Use a service you know on this stack, or the prepared applications in `workshop-shop`. Treat them
as a service you have just inherited.

### Steps

1. Ask Assistant about the health of your service: errors, latency and saturation over the last hour.
2. Ask it to show the queries it used, then to create a dashboard from them.
3. Ask it to explain one query, and confirm the service scope and time range are right.

Copy these in order:

> Tell me about the performance of the applications in workshop-shop and the cluster they run on.
> Show me the queries you used.

> Create a dashboard from those queries, with rows for errors, latency, and saturation.

> Explain the latency query. What exactly does it measure, and over what time range?

**Iterate, do not restart:** stay in the same conversation and refine: add a panel, change a
visualisation, filter to one pod. Starting over loses the context. Save the dashboard as
**Shop - Application Health** with your initials, so nobody overwrites someone else's dashboard.

For orientation: saturation on this stack comes from Kubernetes metrics (cAdvisor and
kube-state-metrics) for every pod in `workshop-shop` on cluster `workshop-prod`. The
`process_cpu_usage` metric exists only for `shop-catalog`. Which panels Assistant proposes is not
fixed; check each one yourself.

### Expected result and completion check

A dashboard with panels for errors, latency and saturation, where you can say what every panel
measures and over what window. Expect one or two panels to need a follow-up prompt; naming the
data source with `@` reduces empty panels.
**A dashboard exists and you can explain every query on it.**

### Finished early?

> Create an alert rule that fires when CPU usage for any service in workshop-shop stays above 90%
> for five minutes.

Review the draft before you save it: metric, filters, threshold and pending period. Saving it is
not required for this session.

### If you get stuck

Use the prompts above verbatim against `workshop-shop`. If a panel returns no data, name the data
source explicitly with `@` and ask again.

## Lab 2 · Deck slide 11 · Investigate a symptom

### Scenario

Since `{fault-time}`, `shop-payment` in `workshop-shop` has shown an elevated error rate. A change
went out at that time: `{change}`. Your job is not to fix it; it is to reach a hypothesis you can
defend, and then say what it cost.

### Steps

1. Describe one real symptom with the service, the timeframe, and what changed.
2. Run the investigation, then read the findings and the evidence behind each one.
3. Validate one finding in Explore before you accept it.
4. Ask Assistant to quantify the impact: how many requests were affected, and for how long.

Prompts:

> Investigate the elevated error rate on shop-payment in workshop-shop over the last 3 hours.
> It started after {fault-time}.

> Show me the query behind this finding.

Run that query yourself in Explore and check scope, data source and time range.

> How many users or requests were affected, and for how long?

The facilitator replaces `{fault-time}` with the actual time, including date and timezone.

### If Investigations is not available, run it conversationally

> Show the error rate for shop-payment over the last hour. When exactly did it change?

> Show the error logs for shop-payment in that window. What patterns do you see?

> Use traces to find which call is failing.

If `{fault-time}` is more than an hour ago, replace "the last hour" with the actual window.

### On impact: validate first, then quantify

The data is synthetic: there are no real users and no revenue. Counted error log events or error
spans are simulated requests in this exercise. The histogram count is modelled observations,
**not** request throughput; do not derive a user count from it. An unvalidated impact number gets
repeated in management reports long after the detail is forgotten. Hence the order: check the
finding first, then state the impact and write down what it rests on.

### Expected result and completion check

A root-cause hypothesis you can defend or reject, with a query you have personally run, plus a
one-line impact statement. Rejecting the hypothesis is a valid and valuable outcome. A dependency
and a coincidence in time along the trace path are not yet proof of cause.
**You can defend the hypothesis with evidence, and state the impact.**

### If you get stuck

Investigate the prepared scenario and its known fault in the window from `{fault-time}`. If the
change does not appear in the findings, that is because it exists only verbally or as an
annotation, not because of your prompt.

## Lab 3 · Deck slide 17 · Make it repeatable

### Scenario

Three short steps that turn today's method into something your team keeps: a Rule for always, a
Skill for a specific job, an Automation for when.

### Step 1: Create a Rule

**Assistant → Settings → Create rule.** Keep it concrete, and keep the scope to **Just me** today.

> When investigating issues, always check traces before logs. Start with Tempo, then drill into Loki.

> When a critical issue is found, always name the owning team and suggest creating an incident via
> the IRM integration.

This is an Assistant behaviour rule, not an alert rule.

### Step 2: Create a Skill

Write your own or create one from a template. Example:

> Name: payment triage
>
> Use when: alerts or questions about shop-payment health
>
> Instructions:
>
> 1. Open the Shop - Application Health dashboard (with my initials) from Lab 1.
> 2. Compare the error rate for the last hour with the previous hour.
> 3. Check for deployments, config or error rate changes in the window.
> 4. If errors correlate with a change or error increase, summarise: symptom, timing, suspected
>    change, evidence.
> 5. If severity is critical, name the owning team to page.

### Step 3: Schedule it as an Automation

Enable the skill's slash command and run it once. Then **Assistant → Settings → Automations →
New**: add a schedule and a prompt, or call the skill you just built. Run it once manually to
check it works.

### Expected result and completion check

A Rule that visibly changes how a new conversation behaves; a Skill that produces the first step
you would actually take, not a generic checklist; and one scheduled Automation with a successful
manual run. You will probably not finish all three; that is expected, and you can complete the
rest on your own stack afterwards.
**A Rule shapes every conversation, and a scheduled Skill produces your first step.**

### If you get stuck

Create the Rule and adapt the Skill above to the prepared scenario dashboard and `shop-payment`.
If you cannot reach Settings, that is a fallback, **not** a saved artefact; do not improvise
permissions during the lab.

The facilitator checks enablement and access beforehand with the
[official Assistant guide](https://grafana.com/docs/grafana-cloud/platform/grafana-assistant/introduction/)
and [enablement and access](https://grafana.com/docs/grafana-cloud/platform/grafana-assistant/get-started/grafana-cloud/);
UI availability and permissions can differ per stack.

## Close and next run

Which query on your dashboard could you explain, and which not? Which hypothesis did you validate
or reject, and what does your impact statement rest on? Which Rule or Skill would you share with
the team, and who would own it?

The facilitator confirms that `payment-regression` is deactivated and
[resets the exercise state](kubernetes.md#reset-the-exercise-not-the-telemetry). Before the next
cohort, delete every Rule scoped **Everybody** again and prepare a fresh `{fault-time}`. Old
telemetry does not need to be deleted.

## Facilitator troubleshooting

| Observation | Check first |
|---|---|
| Assistant answers differently than expected | Leftover Everybody Rules, the person's own Rules |
| Panels without data | Name the data source with `@`, check time range and filters |
| No elevated errors on `shop-payment` | Scenario window, ingest delay, `{fault-time}` and its timezone |
| The change is missing from the findings | Expected without an annotation; state it verbally |
| Investigations missing | Conversational fallback on the Lab 2 card |
| Settings not reachable | Role and Assistant access of the attendee account |
| Everything empty | Dry run, blueprint selection, paused pod and readiness |

Further diagnosis: [Control Plane](control-plane.md) and [Troubleshooting](troubleshooting.md).
Attendees report problems to the facilitator instead of changing tokens, Helm values or shared
resources.
