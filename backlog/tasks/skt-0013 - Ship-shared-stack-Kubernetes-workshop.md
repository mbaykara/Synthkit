---
id: SKT-0013
title: Ship shared-stack Kubernetes workshop
status: Done
assignee: []
created_date: '2026-09-16 19:23'
updated_date: '2026-09-16 19:25'
labels: []
dependencies: []
type: feature
ordinal: 36000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide a repeatable Helm deployment and generic ninety-minute Grafana Cloud workshop without qualification fixture dependencies.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Helm chart validates pinned images, external secrets, single-writer persistence and separate liveness/readiness with render regression coverage.
- [x] #2 Generic workshop blueprint models uneven telemetry adoption and instructor-triggered incidents with offline emission tests.
- [x] #3 Ninety-minute shared-stack guide covers deployment, exercises, recovery, pause/resume and retained telemetry.
- [x] #4 Relevant local gates pass and new checks are wired into CI; live verification limitations are explicit.
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [x] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [x] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Implement chart and blueprint with separate file ownership, then integrate guide and CI. Verify offline emission and intentional signal absence, incident recovery, readiness and persistence, Helm rendering, docs and full repository gate. Retain historical telemetry. No live cluster mutation or upstream PR.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented in 1563b19. Final make gate passed: build, vet, tests, race, RW2 compatibility, 629 Go-file SPDX checks and public-surface scan. Ten Helm regression tests and docs-check passed. Selected workshop dry-run reports 317 metric names. Integration tests prove five-service adoption gaps, correlation, scenario effects and recovery. Production image built locally; restricted-filesystem dry-run succeeded, scenario persisted across container restart, reset succeeded, missing profile credentials failed readiness. No blueprint fields/config structs changed, so schema regeneration is not applicable. Local test container is stopped. No live Kubernetes/Grafana mutations; cluster storage, registry pulls, actual ingestion and participant permissions remain deployment acceptance checks. Existing inherited image publisher fails for mixed-case fork repository paths; documented explicit lowercase build/push alternative. No upstream PR.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Shipped Helm chart, generic uneven-adoption blueprint, ninety-minute shared-stack guide and fail-closed profile readiness in 1563b19. Local gate, Helm/docs checks and container persistence/reset smoke passed. Historical telemetry remains untouched. Live deployment and exact-SHA hosted CI are not claimed.
<!-- SECTION:FINAL_SUMMARY:END -->
