---
id: doc-0004
title: Workshop deployment handoff
type: guide
created_date: '2026-09-17 05:28'
updated_date: '2026-09-17 05:39'
---
# Workshop deployment handoff

Updated: 2026-09-17. Status: implementation complete; image publication blocked; Kubernetes deployment not performed.

## Current state

- Work directly on main, in the user's SynthKit fork. Do not open an upstream PR yet.
- Implementation commit: 1563b19. Verification-record commit: 18e716caf0ca8522e7781e11e2bb60630c816424. Both were pushed.
- SKT-0013 is Done for implementation and local verification, not for live deployment.
- Working tree was clean before this handoff document was created.
- No live Kubernetes or Grafana resources were changed. Historical telemetry was not erased.
- No successful GHCR publication or deployable registry digest was obtained. The failed push may have uploaded unreferenced layers; do not mistake build output for a published image.

## Delivered

- charts/synthkit: digest-pinned image, external credential Secret, single replica with Recreate updates, retained PVC/existingClaim, nonroot/read-only-root container, private ClusterIP.
- blueprints/grafana-cloud-workshop.yaml: five synthetic shop services. Storefront, checkout and payment have metrics/logs/traces/CPU profiles; inventory has metrics only; shipping has metrics/logs.
- Instructor scenario checkout-regression increases checkout latency/errors/CPU. No automatic incident activation or Grafana fixture provisioning.
- docs/workshop.md: 90-minute shared-stack workshop. One instructor deploys and controls scenarios; trainees query Grafana with their own accounts.
- docs/kubernetes.md: build/push, private Secret, install, rotation, pause/resume, reset, retained-PVC reinstall and troubleshooting.
- Optional CLI flag -healthcheck-require-profiles prevents the chart reporting ready when profile credentials are empty. Ordinary deployments keep profiles optional.
- Helm tests wired into GitHub and Forgejo CI.

## Verification

- Final make gate passed: build, vet, unit/integration tests, race, RW2 compatibility, SPDX and forbidden-word checks.
- Ten Helm regression tests and docs-check passed.
- Selected workshop dry-run: 317 distinct metric names.
- Integration tests prove all five services' intended presence/absence, log/trace/profile correlation, scenario effects, and deactivation recovery. Weekend-midnight traffic is usable.
- Local production image built and ran with a read-only filesystem.
- Local control scenario survived container restart; control reset cleared it.
- New healthcheck failed as expected when profile credentials were missing.
- Review findings fixed: optional-profile readiness gap and waiting for Pod deletion before resume.
- Two full gate passes occurred: initial implementation and final reviewed readiness fix. Focused test-first failures were expected and fixed.
- Hosted CI for the final SHA and live cluster ingestion were not verified. No claim of signed release provenance or registry pullability.
- Local evidence retained in ignored codex/workshop-evidence-20260916/: final-gate.log, container-build.log, container-dump.log, synthkit-publish.i7QzNz-build.log. The container dump predates final private_link:false; final source dry-run/integration checks are the authority for the final 317-name inventory.

## Publication attempt and blockers

The explicit lowercase registry path in docs/kubernetes.md bypasses the inherited workflow's separate mixed-case repository-path bug.

A clean git archive of 18e716caf0ca8522e7781e11e2bb60630c816424 was built for linux/amd64 and linux/arm64. Intended tag: workshop-18e716caf0ca8522e7781e11e2bb60630c816424. Both platform builds completed, but registry export failed:

    denied: permission_denied: The token provided does not match expected scopes.

Docker login with the existing GitHub CLI credential succeeded, but that credential could not push packages. Authentication success is not proof of package-write permission. The temporary Docker GHCR login was removed immediately afterward; it had been stored unencrypted by Docker, so do not leave such temporary credentials behind.

The inherited reusable publishing workflow also previously failed because it used a mixed-case GitHub repository name in Docker output/cache paths. That workflow was not modified. Follow the explicit lowercase publishing procedure or implement and verify a narrowly scoped workflow fix later; do not rerun the unchanged failing workflow expecting a different result.

## Decisions and safety boundaries

- Keep work in the fork; upstream PR later.
- Kubernetes + Helm, one shared Grafana Cloud stack, 90-minute general workshop; no DOMM fixture dependency.
- Preserve historical telemetry. Reset controls for a healthy new session; ordinary restart restores persisted controls.
- Never delete the namespace/PVC/Grafana stack merely to repeat training.
- Wait for the old Pod to disappear before resuming from zero replicas.
- Dry-run intentionally remains NotReady; omit --wait for dry-run.
- Use an image built from the workshop commit or later; old upstream images lack the blueprint and required probe flag.
- No HPA or duplicate emitters sharing synthetic identities. PVC persists controls, not counters or a telemetry WAL.
- The current kubeconfig context was not approved for this deployment; do not deploy merely to whatever context happens to be active.
- No new product/design question remains. Publishing requires package-write authorization, followed by an explicit target cluster/namespace and externally managed runtime credentials.

## Cleanup

Verification logs were preserved in the ignored evidence directory before cleanup. Disposable source copies, scratch patches, local synthetic state and original push log are moved to the user's Trash under synthkit-workshop-session-20260917 and remain recoverable. The stopped synthkit-workshop-local-check container and synthkit:workshop-validation test image are removed; those are reproducible, not archived. Shared Docker builders/cache and unrelated repositories/resources are left alone. No publishing build or telemetry emitter remains running from this work.

## Required reading and next action

1. AGENTS.md and applicable contribution policies.
2. This handoff and SKT-0013.
3. docs/kubernetes.md, charts/synthkit/README.md, docs/workshop.md.

First obtain a credential authorized for GHCR package writes (or a correctly permissioned Actions publishing path). Do not print tokens or commit them. Authenticate with password-stdin, rebuild/push the exact reviewed revision using the lowercase repository from the deployment guide, and verify the registry index digest plus AMD64/ARM64 manifests. Remove temporary credentials afterward.

Then pin that verified digest in Helm values, establish the intended Kubernetes target and external Secret, deploy once, and verify fresh metrics/logs/traces/profiles plus participant access. Test pause/resume and scenario reset on that actual deployment before class. Record the digest and actual live results; local tests alone do not prove cluster storage, registry pulls or Grafana ingestion.
