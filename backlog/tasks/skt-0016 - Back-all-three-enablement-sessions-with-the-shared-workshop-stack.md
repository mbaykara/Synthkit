---
id: SKT-0016
title: Back all three enablement sessions with the shared workshop stack
status: Done
assignee:
  - '@claude'
created_date: '2026-09-26 18:40'
updated_date: '2026-09-26 18:47'
labels: []
dependencies: []
type: feature
ordinal: 39000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The enablement programme has three 60-minute hands-on sessions (Foundation, Alerting, Assistant) delivered in German and English on one shared Grafana Cloud stack. Only Foundation was mapped to the workshop blueprint, and only in German. The Alerting lab cards assume Micrometer/Spring JVM metrics that the Go-only shop did not emit, and the Assistant facilitator sheet pointed at an external OpenTelemetry Demo deployment.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Workshop blueprint emits process_cpu_usage and per-pod jvm_memory_used_bytes/jvm_memory_max_bytes (area=heap, 12 pods, ratio above 0.7) for a shop-catalog JVM service in workshop-shop, verified by dry-run dump
- [x] #2 A payment-regression instructor scenario raises shop-payment errors and latency for the Assistant investigation lab and recovers on deactivation
- [x] #3 Foundation, Alerting and Assistant each have a German and an English guide with a filled facilitator value sheet, lab mapping and fallbacks
- [x] #4 New metric names carry PENDING provenance entries in cantfind.md; make gate and docs-check pass
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [x] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [x] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add shop-catalog JVM workload (inline gauges, pod enum x12, service_name const label) and payment-regression scenario to blueprints/grafana-cloud-workshop.yaml. 2. Add cantfind.md PENDING entries for process_cpu_usage and jvm_memory_max_bytes. 3. Update workshop integration tests for the sixth service and new scenario. 4. Dry-run dump check. 5. Write EN Foundation plus DE/EN Alerting and Assistant guides, register in docs navigation. 6. make gate, docs-check.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Blueprint-only change: shop-catalog (runtime jvm, 12 replicas) declares process_cpu_usage, jvm_memory_used_bytes and jvm_memory_max_bytes inline with a pod enum equal to the deterministic k8s pod names; service_name/namespace are auto-stamped (a declared service_name label is rejected as reserved). runtime_jvm profile deliberately not attached to avoid a pod-less jvm_memory_used_bytes. New tests TestWorkshopAlertingFleet (12 pods joined to kube_pod_info, heap ratio in (0.7,1), CPU below 0.8) and TestWorkshopPaymentScenario (payment latency and error logs up, inventory unchanged, recovery) pass; workshop inventory 317 -> 320 names. cantfind SK-84 process_cpu_usage, SK-85 jvm_memory_max_bytes. No Go struct change, so blueprint-schema not applicable. Docs: Foundation EN plus Alerting and Assistant DE/EN, Enablement Workshops nav. make gate, docs-check, forbidden-words and git diff --check pass. Not verified live: no deployment exists; Synthkit has no backfill, so Alerting Lab 1 needs the generator running 7 days before the session. Foundation guide remains 90 minutes while the course catalogue says 60.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Extended the shared workshop blueprint with a 12-pod JVM shop-catalog fleet for the Alerting labs and a payment-regression scenario for the Assistant investigation lab, and wrote German and English guides for all three sessions. Verified with new integration tests, make gate, docs-check and forbidden-words; live stack verification remains a deployment prerequisite.
<!-- SECTION:FINAL_SUMMARY:END -->
