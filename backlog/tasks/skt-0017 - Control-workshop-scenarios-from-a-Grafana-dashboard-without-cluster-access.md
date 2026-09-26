---
id: SKT-0017
title: Control workshop scenarios from a Grafana dashboard without cluster access
status: Done
assignee:
  - '@claude'
created_date: '2026-09-26 19:45'
updated_date: '2026-09-26 19:46'
labels: []
dependencies: []
type: feature
ordinal: 40000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Instructors needed kubectl port-forward to switch scenarios. The stock control dashboard uses browser-direct fetch actions, which require an internet-facing control endpoint and cross-origin Basic auth. Grafana's server-side infinity actions keep the control Service private when combined with Private Data Source Connect.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 synthkit-control-dash supports -action-mode infinity with -ds-uid, emitting server-side infinity actions; default fetch output is unchanged; tests cover both modes and flag validation
- [x] #2 A dashboard click activates and clears a workshop scenario on a live deployment through an Infinity datasource over PDC with the control Service still ClusterIP
- [x] #3 docs describe the flags, the vizActionsAuth prerequisite, the PDC region-format requirement and the datasource permission risk; make gate and docs-check pass
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [x] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [x] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add dashboard.InfinityAction and generator flags. 2. Tests. 3. Deploy PDC agent, Infinity datasource, restricted folder, dashboard. 4. Click-verify activate and clear. 5. Docs, gate, commit.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Grafana drops infinity actions unless vizActionsAuth is on (public/app/features/actions/utils.ts filter), which rendered empty action boards until the toggle went live. Click-verified on Grafana 13.3: activate set active_scenarios to the payment scenario and payment error logs rose from about 1 to 27 per minute; Clear all returned []. PDC agent on a region-named cell needed -region-format (HTTP 530 otherwise). No blueprint fields changed; blueprint-schema and dump not applicable. make gate and docs-check pass.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added server-side infinity action mode to the control dashboard generator and documented Grafana-dashboard control over Private Data Source Connect. Verified with generator tests, make gate, docs-check and a live click test that activated and cleared a scenario with the control Service still cluster-internal.
<!-- SECTION:FINAL_SUMMARY:END -->
