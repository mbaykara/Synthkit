---
id: SKT-0015
title: Publish fork CI images to lowercase GHCR repository
status: Done
assignee:
  - '@codex'
created_date: '2026-09-20 22:44'
updated_date: '2026-09-26 19:06'
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
- [x] #1 Fork main and manual edge builds use the explicit lowercase synthkit image name with package-write Actions authentication and native AMD64/ARM64 builds.
- [x] #2 Existing upstream release verification is preserved; unsupported fork release publishing fails clearly rather than claiming verification.
- [x] #3 Workflow validation and repository checks pass, and a hosted build publishes a verified multi-platform digest or records a concrete external blocker.
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 make gate (build vet test race rw-proto-check spdx-check forbidden-words)
- [x] #2 make blueprint-schema (only if a blueprint field or construct/workload config struct changed)
- [x] #3 DRY_RUN=true go run ./cmd/synthkit -once -dump — inventory diffed against signals/
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

PAUSED at user request. Code fixes committed and pushed: 08d3740 (fork lowercase publisher, offline regression tests and docs), b3762b7488140cb9ad03de29fb9001c65d0bb8a2 (grpc 1.83.2, x/crypto 0.55.0 and required x/net/x/text updates). Full make gate docs-check, actionlint v1.7.12, go mod verify, git diff --check and offline grafana-cloud-workshop dry run passed. Latest publish run 35543218796 remains in_progress at pause: https://github.com/mbaykara/Synthkit/actions/runs/35543218796 . Leave external CI running; no local build process remains. Worktree was clean before this handoff. Next: gh run view 35543218796 --repo mbaykara/Synthkit --json status,conclusion,jobs . If failed inspect logs and Trivy findings; do not disable security gates. If successful inspect ghcr.io/mbaykara/synthkit:main, verify AMD64/ARM64 index digest and source labels, optionally offline-run workshop from exact digest, then record digest and finish task/handoff docs. No successful image publication claimed yet. Preserve upstream release signer/verifier; fork only supports main edge builds. Reusable checks out mutable branch ref: avoid source pushes while it builds. No Kubernetes deployment authorized or performed. Test cache moved to /private/tmp/synthkit-publish-ci-test-pycache; offline smoke log /private/tmp/synthkit-ci-deps-dryrun.log. No credentials printed or stored. Repository uses Backlog as sole handoff tracker, not parallel .planning files.

Resumed 2026-09-26. Run 35543218796 concluded success. Run 36263878498 for 2ad290d also succeeded and published index sha256:976de9dd9018c28edc61f732ecdad937f9dc767711047f0b20236647e7d9bf47 with linux/amd64 and linux/arm64 manifests, both labelled revision 2ad290d093660ed58335787fc59ecc27395bfd7b; anonymous pull returns 200. The digest was deployed with the Helm chart to the instructor cluster: pod ready with zero restarts, /app/synthkit -version reports 2ad290d, all four sinks delivering without failures, and workshop metrics, logs, traces and profiles queried in the target stack. Blueprint-schema not applicable (no field change); gate and dump evidence recorded in earlier notes.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Fork main builds now publish a lowercase, signed multi-platform synthkit image with GITHUB_TOKEN while upstream release verification is unchanged. Verified by the successful hosted publish of 2ad290d, index and per-platform revision inspection, anonymous pull, and a live Helm deployment delivering all four signals.
<!-- SECTION:FINAL_SUMMARY:END -->
