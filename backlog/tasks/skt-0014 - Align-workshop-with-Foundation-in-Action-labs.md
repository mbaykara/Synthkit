---
id: SKT-0014
title: Align workshop with Foundation in Action labs
status: Done
assignee:
  - '@codex'
created_date: '2026-09-20 22:15'
updated_date: '2026-09-20 22:28'
labels: []
dependencies: []
type: docs
ordinal: 37000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The shared-stack workshop must follow the supplied German Foundation session cards rather than separate signal tutorials. Preserve 90 minutes and teach evidence validation before action.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 German guide maps all three lab cards to slides 9, 14 and 27 with outputs and fallback paths.
- [x] #2 Instructor preparation fixes a shared service and absolute window with a dashboard and Explore-link checklist; mandatory cross-signal validation and personal Quickstart are explicit.
- [x] #3 Synthetic-data limits, Assistant prerequisites and safe repeat/reset are documented; documentation checks pass.
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Rewrite docs/workshop.md around the supplied cards, retain exact known queries and operational safeguards, verify Assistant instructions against official docs, review and run docs-check and hygiene.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Documentation-only scope: no blueprint/config fields or emission behavior changed, so schema generation and renderer inventory diff are not applicable. Read-only review confirmed query namespace labels against implementation and all three supplied lab contracts. Official Grafana Assistant documentation checked for personal Quickstart and context workflow. make docs-check hygiene, git diff --check and full make gate passed. No live stack validation or provisioning performed; actual dashboard, Explore links, absolute window and participant Assistant access remain instructor prerequisites.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Replaced the signal-by-signal agenda with German Foundation labs for slides 9, 14 and 27, a 90-minute facilitator schedule, pre-captured bounded checkout scenario, required two-signal validation and personal Quickstart examples. Preserved synthetic-data and safe-reset limits. Full repository gate and documentation checks passed.
<!-- SECTION:FINAL_SUMMARY:END -->
