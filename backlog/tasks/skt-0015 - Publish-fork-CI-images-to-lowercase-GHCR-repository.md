---
id: SKT-0015
title: Publish fork CI images to lowercase GHCR repository
status: In Progress
assignee:
  - '@codex'
created_date: '2026-09-20 22:44'
updated_date: '2026-09-20 22:56'
labels: []
dependencies: []
type: bug
ordinal: 38000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The inherited publisher uses mixed-case github.repository for registry output/cache and fails before pushing. The local fallback token also lacks package writes. Use Actions GITHUB_TOKEN to publish the fork workshop image without altering upstream release verification trust.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Fork main and manual edge builds use the explicit lowercase synthkit image name with package-write Actions authentication and native AMD64/ARM64 builds.
- [ ] #2 Existing upstream release verification is preserved; unsupported fork release publishing fails clearly rather than claiming verification.
- [ ] #3 Workflow validation and repository checks pass, and a hosted build publishes a verified multi-platform digest or records a concrete external blocker.
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [ ] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [ ] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
<!-- DOD:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Add a fork edge job using the immutable maintained reusable workflow with image-name support; preserve the upstream job and reject unsupported fork releases; validate static contracts and hosted publish, then record the digest and update deployment guidance.

Hosted Trivy scans block the existing dependencies: update grpc to 1.83.2 and x/crypto to 0.55.0 (with required x/net 0.58.0 and x/text 0.41.0), rerun the gate, then retry publication without suppressions.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Three regression tests failed before the fix and pass afterward; wired into make deploy-tests. actionlint v1.7.12, full make gate, docs-check and git diff --check pass. Read-only review found no blockers. Fork uses maintained reusable db2707e with image-name synthkit and no security bypasses. Upstream release trust unchanged. Await hosted build and digest verification.

Run 35542890980 built both platforms but failed prepublish Trivy: CVE-2026-84445 and CVE-2026-84304 in grpc 1.82.1, CVE-2026-56854 in x/crypto 0.54.0. GitHub advisories confirm grpc 1.83.2 fixes both; scanner identifies x/crypto 0.55.0. No image pushed and no security bypass applied.
<!-- SECTION:NOTES:END -->
