---
id: SKT-0010
title: Model cAdvisor ownership-label enrichment
status: Done
assignee:
  - '@codex'
created_date: '2026-08-23 22:35'
updated_date: '2026-08-23 22:43'
labels: []
dependencies: []
modified_files:
  - BLUEPRINT-SCHEMA.md
  - internal/blueprint/schema.go
  - internal/blueprintschema/fielddocs.json
  - internal/construct/k8scluster/cadvisor.go
  - internal/construct/k8scluster/cadvisor_test.go
  - signals/k8s.md
priority: high
type: bug
ordinal: 26000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Allow a qualification blueprint to model collector-side promotion of the Kubernetes namespace into service_namespace on cAdvisor container series without changing the raw default signal shape.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Raw cAdvisor output remains unchanged unless the blueprint explicitly enables ownership-label enrichment.
- [x] #2 When enabled, cAdvisor container and pod-network series with namespace also carry the same service_namespace value.
- [x] #3 The DOMM qualification blueprint enables the enrichment for its fresh-stack telemetry.
- [x] #4 Focused tests and the signal contract document the optional label.
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [x] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [x] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add an opt-in k8s_monitoring feature that promotes namespace to service_namespace on cAdvisor series.
2. Cover default-off and enabled container/network behavior with focused tests and update the signal contract.
3. Enable the feature in the DOMM qualification blueprint, validate both repositories, and resume live ingestion.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented opt-in namespace ownership-label promotion on cAdvisor container and pod-network series. Default-off and enabled tests pass. DOMM dry-run inventory and live Grafana query confirmed namespace equals service_namespace for all six qualification scopes.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added opt-in collector-enrichment modeling for cAdvisor ownership labels, regenerated blueprint docs, and verified with make gate, dry-run inventory, and live DOMM stack queries.
<!-- SECTION:FINAL_SUMMARY:END -->
