---
id: SKT-0018
title: Instructor Control Plane dashboard layout
status: Done
assignee:
  - '@claude'
created_date: '2026-09-26 19:55'
updated_date: '2026-09-26 19:55'
labels: []
dependencies: []
type: enhancement
ordinal: 41000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The server-side control dashboard worked but was a plain table of small buttons with no view of scenario impact. Instructors asked for a friendlier console with clear buttons and business charts, named Control Plane, defaulting to three hours.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 synthkit-control-dash -layout control-plane renders status tiles, one card per scenario with colored Start/Stop buttons on the per-item activate/deactivate endpoints, Business Charts impact charts per affected service and collapsed fleet controls, with a 3h default window and a -folder option
- [x] #2 The generated dashboard renders with live data on the target stack and is verified from a snapshot image
- [x] #3 Generator tests cover the layout and validation; tools docs describe the new flags; make gate and docs-check pass
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [x] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [x] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add EChartsPanel, StatusTile, WithActionColor, collapsible rows. 2. control-plane layout and flags. 3. Tests. 4. Generate, push, snapshot and review. 5. Docs, gate, commit.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Palette: impact lines use validated categorical slots 1 and 2 (theme-stepped light/dark); buttons use status serious/good with text labels. Impact queries pin the stack's Prometheus and Loki datasources because the default datasource can be label-restricted. Stress preset omitted from the instructor layout. Snapshot reviewed: status None/1/Healthy, aligned scenario cards, four ECharts impact charts showing the earlier payment test. Blueprint change is scenario summary text only (no fields), so blueprint-schema and dump are not applicable. make gate and docs-check pass.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added a control-plane layout to the control dashboard generator: status tiles, scenario cards with colored server-side Start/Stop buttons, Business Charts impact charts and collapsed fleet controls with a 3h default. Verified with generator tests, make gate, docs-check and a reviewed live snapshot.
<!-- SECTION:FINAL_SUMMARY:END -->
