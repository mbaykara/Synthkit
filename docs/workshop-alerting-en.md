---
title: Grafana Alerting in Action (English)
description: Grafana Alerting in Action – three labs on one shared stack in 60 minutes.
---

# Grafana Alerting in Action

[Deutsche Fassung](workshop-alerting-de.md) · Other sessions: [Foundation](workshop-foundation-en.md) ·
[Assistant](workshop-assistant-en.md)

**One shared stack, 60 minutes, three labs:** create an alert rule from a query you can explain,
route it with labels to exactly one team, and reduce notification noise. The labs belong to deck
slides **8, 12 and 16**. The timings below are a facilitation suggestion; the slide numbers come
from the lab cards.

The facilitator runs Synthkit on Kubernetes. All attendees use the same Grafana Cloud stack with
their own Grafana login. Kubernetes access, ingest tokens and the Synthkit control password stay
with the facilitator. This session stands alone: attendees who missed Foundation need nothing
from it.

The lab card queries use Micrometer/Spring metric names. On this stack the synthetic JVM service
`shop-catalog` emits them. The values are modelled, not measured from a real application.

## What you take away

- Create a Grafana-managed alert rule with a threshold and pending period you can justify.
- Explain the rule lifecycle, including the No data and Error states.
- Route with a few labels so a notification reaches exactly one contact point.
- Reduce noise with grouping, silences and mute timings.

## Facilitator: preparation before the 60 minutes

The `grafana-cloud-workshop` blueprint supplies the data. Synthkit and Helm create **no** folders,
alert rules, contact points, notification policies, silences, mute timings, users or teams. These
stack artefacts must be created and checked in Grafana before the session.

1. **At least seven days of lead time.** Synthkit cannot backfill history. Lab 1 asks attendees to
   look at the last seven days, so the generator must run through the
   [Kubernetes deployment](kubernetes.md) at least seven days before the session. Run the Lab 1
   query yourself over `Last 7 days` and check the data has no gaps.
2. Through the facilitator-only [control plane](control-plane.md),
   [reset the exercise state](kubernetes.md#reset-the-exercise-not-the-telemetry): workshop
   blueprint only, volume multiplier `1`, no active scenarios. This session needs no Synthkit
   scenario.
3. Create `{folder}` and give attendees editor rights on it.
4. Create the contact point `workshop-webhook` pointing at `{webhook-url}` and test it. The target
   must be visible on the shared screen so everyone can see a notification arrive. Also create
   `workshop-email` (catch-all) and `workshop-servicenow`.
5. Build the prepared policy tree from Lab 2 and leave it in place.
6. Create the Lab 3 noisy rule and leave it **paused**. It is enabled on cue during the lab.
7. Send the attendee lab card pack with the invitation.

### Hand out the completed session sheet

| Lab card field | Value for this session |
|---|---|
| `{namespace}` | `workshop-shop` |
| `{service}` | `shop-catalog` (JVM service with 12 pods, synthetic cluster `workshop-prod`) |
| `{team}` | One routing value chosen by the facilitator, for example `workshop-a` |
| `{folder}` | Enter the actual rules folder attendees can write to |
| `{webhook-url}` | Enter the actual webhook target that is visible on screen |
| `{runbook-link}` | Any URL; a placeholder runbook is fine |

**Start check:** without seven days of data for the Lab 1 query, a tested `workshop-webhook` and
the paused Lab 3 rule, the session is not ready. Attendees do not change other attendees' policy
branches, contact points or the control state.

## Schedule: 60 minutes

| Minutes | Section | Outcome |
|---|---|---|
| 0–8 | An alert is a validated query with a decision; threshold, error budget, adaptive; rule anatomy and lifecycle | Terms for Lab 1 clear |
| 8–22 | Lab 1, slide 8: Create an alert rule | Rule evaluates, threshold justified |
| 22–27 | Labelling strategy, from rule to notification | Four to five labels, policy tree understood |
| 27–39 | Lab 2, slide 12: Label and route it | Exactly one notification at the right target |
| 39–44 | Noise control, escalation, operating at scale | Grouping, silences, alerts as code placed |
| 44–56 | Lab 3, slide 16: Quiet the noise | One notification instead of twelve |
| 56–60 | Close and next steps | What pages at 3am is named |

## Lab 1 · Deck slide 8 · Create an alert rule

### Scenario

Alert when a service in `workshop-shop` runs hot on CPU. Look at the last seven days first, then
choose a threshold you could defend to a colleague.

### Starting query

Your own query, any panel you can explain, or this one:

```promql
avg by (service_name) (process_cpu_usage{namespace="workshop-shop"})
```

Kubernetes alternative (cAdvisor, all shop pods):

```promql
sum by (pod) (rate(container_cpu_usage_seconds_total{namespace="workshop-shop"}[5m]))
```

### Steps

1. Choose your starting query.
2. Set a condition from what the history shows, and a pending period that outlasts normal noise.
3. Add a summary annotation with a runbook link, and `severity` and `team` labels.
4. Stretch: ask Assistant to draft the same rule, then compare its threshold with yours.

### Suggested settings – adjust from the history

- Condition: above `0.8` to start, then justify it or move it.
- Evaluation every minute (`1m`); pending period `5m`.
- No data and error handling: choose deliberately, do not leave defaults unread.

### Labels and annotations

```text
labels:      team={team} severity=warning
             service_name=shop-catalog namespace=workshop-shop
annotations: summary     = CPU high on shop-catalog - check recent deploys and load
             runbook_url = {runbook-link}
```

### Expected result and completion check

A rule in **Normal** state that evaluates on schedule, and a sentence explaining why your
threshold and pending period are what they are. If you did the stretch step, expect Assistant's
threshold to be reasonable but generic: it has not seen your incident history. That contrast is
the point. You can **explain exactly when the rule fires and why**.

### For the facilitator

`process_cpu_usage` is steady between 0.30 and 0.72 per pod; the average across the twelve pods
sits near 0.5. The `0.8` start therefore never fires, and the Normal state is by design. Attendees
must argue the threshold from the history. The values are not tied to any scenario.

### If you get stuck

Use the query, threshold and labels exactly as printed above; they work on the shared stack.

## Lab 2 · Deck slide 12 · Label and route it

### Prepared policy tree – already on the stack

```text
Default policy -> workshop-email (catch-all)
|- team = {team}
|    -> workshop-webhook
|    group_by: [service_name]
|
|- team = {team}, severity = critical
     -> workshop-servicenow
```

### Steps

1. Confirm your rule carries the agreed labels: `team` and `severity` at minimum.
2. Build the policy branch: match `team` for the route, then `severity` within it.
   **Alerting > Notification policies > new child policy**.
3. Send a test notification to the shared contact point and confirm it arrives once: press
   **Test** on `workshop-webhook`, then watch `{webhook-url}`.

### Expected result and completion check

Exactly one notification arrives at `{webhook-url}`, carrying your labels, and nothing arrives at
the catch-all email. If it arrives twice, two branches match: check `continue matching` on the
parent. If it only reaches the catch-all, your labels and your matcher do not agree.
The alert **reaches the intended contact point, and no other**.

### A note on contact points

The critical branch points at a ServiceNow contact point, because that is where many teams send
incidents. Grafana Alerting can notify anything that accepts a webhook, so the same tree works
with Teams, email, Grafana IRM or your own tooling. The routing decision is the part that matters;
the destination is yours.

### If you get stuck

Use the prepared policy tree and the shared webhook contact point rather than untangling your own.

## Lab 3 · Deck slide 16 · Quiet the noise

### Scenario

The facilitator enables a prepared noisy rule. It evaluates per pod, so one memory problem
produces twelve alert instances. Ungrouped, that is twelve notifications for one story.

### The noisy rule – already created, enabled on cue

```promql
(jvm_memory_used_bytes{namespace="workshop-shop", area="heap"}
  / jvm_memory_max_bytes{namespace="workshop-shop", area="heap"})
  > 0.7
```

### Steps

1. Set grouping on your policy branch so related instances arrive as one notification:
   `group_by [team, service_name]`, `group_wait 30s`, `group_interval 5m`, `repeat_interval 4h`.
2. Add a silence or mute timing for a planned change: matcher `service_name = shop-catalog`,
   duration `2h`, comment `planned change`.
3. Follow the alert end to end: rule, policy branch, contact point, and on to whatever handles
   escalation.

### Expected result and completion check

One notification listing every affected pod, instead of one per pod. Write down your before and
after counts; you will be asked for them in the debrief. After the silence is applied the rule
still fires and is still visible in the UI, but no notification is sent. That distinction matters.
**One notification tells the whole story, and you can say what pages at 3am.**

### For the facilitator

The "cue" is enabling the paused rule in Grafana, not a Synthkit scenario. `shop-catalog` reports
74 to 88 percent of a 512 MiB heap on all twelve pods, permanently, so the rule fires twelve
instances immediately. All twelve share `service_name="shop-catalog"`, so
`group_by [team, service_name]` collapses them into one notification. Expected: 12 before, 1 after.

### If you get stuck

Use the prepared noisy rule once the facilitator has enabled it, and fix it with grouping alone. Grouping is the change that
produces the visible win.

## Close and next run

Which rule would you keep tomorrow? Does it pass the four tests: actionable, symptom-based, owned,
justified? Which labels does your team need for routing, triage and analytics?

The facilitator pauses the Lab 3 rule again, removes attendee silences and, if needed, deletes
attendee rules from `{folder}`. The generator can keep running; the history stays available for
the next class.

## Reference: alerting data on this stack

Blueprint `grafana-cloud-workshop`, synthetic cluster `workshop-prod`, namespace `workshop-shop`:

| Metric | Labels | Model |
|---|---|---|
| `process_cpu_usage` | `service_name`, `namespace`, `pod` (12) | 0.30–0.72 per pod |
| `jvm_memory_used_bytes` | `area="heap"`, `service_name`, `namespace`, `pod` (12) | 74–88 % of 512 MiB |
| `jvm_memory_max_bytes` | `area="heap"`, `service_name`, `namespace`, `pod` (12) | constant 512 MiB |
| `container_cpu_usage_seconds_total` | `namespace`, `pod` (all shop pods) | cAdvisor |

| Observation | Check first |
|---|---|
| Lab 1 query returns less than seven days | Generator runtime; there is no backfill |
| Lab 1 rule does not fire at `0.8` | Expected; the average sits near 0.5 |
| Lab 3 rule does not fire | Rule still paused, wrong namespace or `area="heap"` |
| Duplicate notification | Two matching branches, `continue matching` on the parent |
| Only the catch-all receives it | Rule labels and branch matchers do not agree |
| Everything empty | Dry run, blueprint selection, paused pod and readiness |

Further diagnosis: [Control Plane](control-plane.md) and [Troubleshooting](troubleshooting.md).
Attendees report problems to the facilitator instead of changing tokens, Helm values or shared
resources.
