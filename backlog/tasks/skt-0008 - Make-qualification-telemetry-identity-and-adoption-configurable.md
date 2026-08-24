---
id: SKT-0008
title: Make qualification telemetry identity and adoption configurable
status: In Progress
assignee:
  - '@codex'
created_date: '2026-08-24 13:15'
updated_date: '2026-08-24 13:30'
labels: []
dependencies: []
modified_files:
  - BLUEPRINT-SCHEMA.md
  - internal/blueprint/schema.go
  - internal/blueprintschema/fielddocs.json
  - internal/construct/k8scluster/cadvisor.go
  - internal/construct/k8scluster/cadvisor_test.go
  - internal/workload/app/app.go
  - internal/workload/app/app_test.go
  - internal/workload/app/interp.go
  - internal/workload/app/profiles.go
  - internal/workload/app/profiles_test.go
  - internal/workload/app/project.go
  - internal/workload/app/servicegraph_test.go
  - signals/apm.md
  - signals/k8s.md
  - signals/profiles.md
  - signals/traces.md
priority: high
type: enhancement
ordinal: 27000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Provide generic, opt-in application and collector behavior needed to model honest observability maturity without emitting DOMM facts or adding blueprint-specific logic.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Opt-in cAdvisor enrichment adds service_namespace, service_name, and deployment_environment_name while default output remains unchanged.
- [ ] #2 An app can explicitly disable traces while omitted configuration preserves the enabled default.
- [ ] #3 An omitted service version remains absent from emitted resource attributes and span-derived identity.
- [ ] #4 SDK-pushed application profiles carry service_namespace alongside service_name.
- [ ] #5 Focused tests, dry-run inventory, and make gate pass.
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [ ] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [ ] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Review the existing local implementation against architecture and signal contracts. 2. Complete focused test and documentation coverage without blueprint-specific logic. 3. Run the full repository gate and dry-run inventory. 4. Commit directly to main and push to the mbaykara/Synthkit fork.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation existed locally before task capture; publication was requested for the DOMM qualification harness.

Focused package tests passed. Blueprint schema regenerated. Full make gate passed with localhost socket permission. DOMM qualification dry-run loaded one blueprint with 6 constructs and 9 workloads, emitted canonical cAdvisor identity, omitted undeclared service.version, and emitted service_namespace on SDK profiles.
<!-- SECTION:NOTES:END -->
