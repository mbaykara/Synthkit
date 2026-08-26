---
id: SKT-0012
title: Rename DOMM qualification scopes to Northstar business namespaces
status: Done
assignee:
  - '@codex'
created_date: '2026-08-26 10:02'
updated_date: '2026-08-26 10:14'
labels: []
dependencies: []
modified_files:
  - blueprints/domm-qualification.yaml
  - internal/blueprint/domm_qualification_test.go
  - docs/domm-qualification.md
  - >-
    backlog/tasks/skt-0012 -
    Rename-DOMM-qualification-scopes-to-Northstar-business-namespaces.md
priority: high
type: enhancement
ordinal: 35000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Replace maturity-revealing ownership scopes in the bundled DOMM qualification blueprint with one coherent fictional retailer namespace set, while preserving telemetry shape and maturity evidence.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 All production and non-production qualification service_namespace values use the approved northstar-* mapping and active blueprint content contains no old domm maturity scope names
- [x] #2 Telemetry adoption, scenarios, service counts, alert-boundary semantics, and absence of DOMM-native scoring facts remain unchanged
- [x] #3 The bundled blueprint loads by name and its focused contract test and dry-run inventory pass
- [x] #4 The full Synthkit gate passes and the resulting commit is suitable for DOMM to pin and mirror byte-for-byte
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [x] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [x] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Update the bundled blueprint identity fields and focused contract expectations using the approved Northstar mapping. 2. Verify no active old maturity scope names remain and compare the dry-run inventory for identity-only changes. 3. Run the complete repository gate, commit on main, and publish the commit for DOMM pinning.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Updated the bundled blueprint, focused contract test, and demo guide to the approved Northstar mapping. The focused contract test passes, active blueprint/docs contain no old maturity scope names, and the dry-run loads 6 constructs and 9 workloads with Northstar ownership labels.

Validation complete: focused bundled-blueprint contract passed; DRY_RUN loaded 6 constructs and 9 workloads with Northstar labels; make gate passed. No blueprint schema or config struct changed, so blueprint-schema regeneration was not required.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Renamed DOMM qualification ownership scopes to fictional Northstar Commerce namespaces without changing telemetry semantics; verified with the focused contract test, dry-run inventory, and full make gate.
<!-- SECTION:FINAL_SUMMARY:END -->
