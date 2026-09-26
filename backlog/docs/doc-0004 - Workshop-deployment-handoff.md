---
id: doc-0004
title: Workshop deployment handoff
type: guide
created_date: '2026-09-17 05:28'
updated_date: '2026-09-26 19:07'
---
# Workshop deployment handoff

Updated: 2026-09-26. Status: live. One continuous emitter runs the `grafana-cloud-workshop` blueprint through the Helm chart and backs all three enablement sessions (Foundation, Alerting, Assistant; guides in German and English under docs/workshop*.md).

## Current state

- Image: fork publish of revision 2ad290d093660ed58335787fc59ecc27395bfd7b, index sha256:976de9dd9018c28edc61f732ecdad937f9dc767711047f0b20236647e7d9bf47 (linux/amd64 and linux/arm64, public pull). SKT-0015 Done.
- Deployed 2026-09-26 with the documented values/Secret flow to the instructor-owned cluster, release and namespace `synthkit`. Pod ready, zero restarts, all four sinks delivering; workshop metrics, logs, traces and profiles verified in the target stack. The operator keeps the private record (context, stack, values file, credentials file) outside this repository.
- SKT-0016 Done: `shop-catalog` 12-pod JVM gauges for Alerting labs and the `payment-regression` scenario for the Assistant investigation lab.
- Control state: defaults, no active scenarios.

## Constraints

- No backfill: 7 days of history (Alerting Lab 1, Foundation facilitator sheet) exist from 2026-10-03.
- Grafana-side artefacts (dashboards, Explore links, alert folder, contact points, policy tree, paused noisy rule, Assistant rules/skills cleanup, change annotation) are instructor preparation; the chart creates none.
- Scenarios do not expire; activate before a session, deactivate explicitly, reset control state between cohorts. Never run a second emitter for these identities.
- Foundation guide keeps its 90-minute agenda while the course catalogue lists 60 minutes; Alerting and Assistant guides use 60.

## Next

- Before each session follow the matching guide's facilitator checklist and verify the listed queries against fresh data.
- Upgrades: publish from main, verify index digest and per-platform revision, update the values file, `helm upgrade --install --wait`.
- Undeploy with `helm uninstall`; the PVC and external Secret are retained by design.
