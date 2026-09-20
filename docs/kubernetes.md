---
title: Kubernetes and Helm
description: Continuously run the shared Grafana Cloud workshop with a digest-pinned image, external credentials, and persistent operator state.
---

# Kubernetes and Helm

The chart in `charts/synthkit` runs one continuous emitter for the
[Grafana Cloud workshop](workshop.md). It selects the builtin `grafana-cloud-workshop` blueprint.
It does not deploy DOMM, provision Grafana fixtures, or run Terraform. Restarting a Pod never
triggers cloud-resource cleanup. Historical metrics, logs, traces, and profiles remain in Grafana
Cloud under the stack's normal retention policy.

This is a single-instance deployment: `replicaCount` is `1` (running) or `0` (paused), and updates
use `Recreate`. Do not attach an HPA or run another emitter for the same synthetic identities.
The PVC retains operator settings and staged blueprints, not a telemetry queue or cumulative
metric counters. A restart can lose in-memory queued data and resets synthetic counters; it is
not exactly-once delivery or an HA deployment.

## Prerequisites

- A reviewed Kubernetes context, Helm, kubectl, and a persistent-volume provisioner supporting
  the selected storage class. This chart creates no cluster-wide RBAC.
- A Grafana Cloud stack reserved for synthetic data, with metrics, logs, traces, and profiles
  ingest access. Learners use Grafana accounts; only the instructor has Kubernetes/control access.
- A container registry writable by the image publisher and readable by the cluster. Configure
  `imagePullSecrets` if the package is private.
- Outbound HTTPS/DNS to the selected Grafana Cloud endpoints. Budget and monitor ingestion for
  a continuously running generator; no retention or cost limits are imposed by this chart.

The chart uses a nonroot process (UID/GID `65532`), a read-only root filesystem, and a writable
state volume. The storage driver must support that volume ownership. Control access is ClusterIP
only, with a required password and the explicit `trusted-network` acknowledgement. Use an
isolated trusted cluster network and instructor port-forward; do not expose plaintext Basic
authentication through a public Service or ingress. NetworkPolicy/TLS and tenant isolation are
cluster responsibilities, not features installed by this chart.

## 1. Publish the correct image

Use your fork and a reviewed commit that **contains**
`blueprints/grafana-cloud-workshop.yaml`. An older upstream image does not contain this blueprint.
The chart requires a digest, not `main`, `latest`, or a mutable tag. Publishing an image and
installing a chart are separate operations.

From that clean source checkout, with registry authentication already configured, an explicit
build/push workflow is:

```bash
git status --short
test -f blueprints/grafana-cloud-workshop.yaml
WORKSHOP_REVISION="$(git rev-parse HEAD)"
WORKSHOP_IMAGE=ghcr.io/mbaykara/synthkit
WORKSHOP_BUILD_DIR="$(mktemp -d)"
docker buildx build --platform linux/amd64,linux/arm64 \
  --build-arg "VERSION=workshop-${WORKSHOP_REVISION}" \
  --build-arg "REVISION=${WORKSHOP_REVISION}" \
  --tag "${WORKSHOP_IMAGE}:workshop-${WORKSHOP_REVISION}" \
  --metadata-file "${WORKSHOP_BUILD_DIR}/metadata.json" --push .
WORKSHOP_DIGEST="$(jq -er '."containerimage.digest"' "${WORKSHOP_BUILD_DIR}/metadata.json")"
```

Review `git status` before building; do not stamp a commit ID onto unreviewed working-tree changes.
Use a registry path you control. A single-platform build is also valid if it matches every eligible
node; adjust `--platform` deliberately. Alternatively, use your fork's publishing workflow and
obtain the published index digest from its successful build output. Record the revision and
digest together. These instructions do not imply a workshop image is already published.

The fork's `publish` workflow builds on pushes to `main`, or can be dispatched manually on
`main` with `release_tag` empty. Its `image-fork` job uses the explicit lowercase image name
`synthkit` and GitHub's automatic `GITHUB_TOKEN` with package-write permission; no personal
token or repository secret is needed. It builds native AMD64 and ARM64 images, scans before
pushing, then merges, signs and attests the index. Security failures must be resolved, not bypassed.
The image tags include `ghcr.io/mbaykara/synthkit:main` and a `main-<short-sha>` tag. After a
successful run, inspect the published image and pin its index digest in Helm, never the mutable
`main` tag. The reusable builder checks out the branch ref; avoid advancing `main` during a
publication and verify the reported source revision. Fork release tags are deliberately rejected:
the existing release verification trust policy remains specific to upstream.

```bash
gh workflow run publish.yml --repo mbaykara/Synthkit --ref main
gh run list --repo mbaykara/Synthkit --workflow publish.yml --limit 5
docker buildx imagetools inspect ghcr.io/mbaykara/synthkit:main
```

For a private package, grant the cluster pull access and configure `imagePullSecrets`.
A successful registry login alone does not establish package-write permission for the manual
build/push fallback above. Do not treat a failed workflow or an older image as the workshop artifact.

## 2. Supply credentials outside Helm

Confirm `kubectl config current-context` is the intended cluster. The examples use release and
namespace `synthkit`; the synthetic application namespace remains `workshop-shop`.

```bash
kubectl create namespace synthkit --dry-run=client -o yaml | kubectl apply -f -
```

Prepare a mode-0600 file outside the repository using your private editor or secret manager.
Use plain `KEY=value` entries, no `export`, shell substitutions, quotes around values, or inline
comments. Populate these exact keys (see [credentials](credentials.md) for endpoint formats):

```dotenv
CONTROL_TOKEN=
GC_TOKEN=
GC_PROM_RW=
GC_PROM_USER=
GC_OTLP_ENDPOINT=
GC_OTLP_USER=
GC_LOKI=
GC_LOKI_USER=
GC_PROFILES_URL=
GC_PROFILES_USER=
```

`CONTROL_TOKEN` must be a strong, unique, nonempty password even in dry-run. `GC_TOKEN` needs
`metrics:write`, `logs:write`, `traces:write`, and `profiles:write` restricted to the selected stack.
Use each sink's own user/instance ID; they are not interchangeable. The chart requires the
profiles keys for live workshops as well as metrics/logs/traces. These are synthetic profile
credentials, not the separate `GC_PYROSCOPE_*` self-profiling settings. Self-observability is off.

Create or update the externally managed Secret without putting values in command arguments,
Helm values, Git, or rendered terminal output:

```bash
WORKSHOP_ENV=/absolute/private/path/synthkit-workshop.env
chmod 600 "$WORKSHOP_ENV"
kubectl -n synthkit create secret generic synthkit-runtime \
  --from-env-file="$WORKSHOP_ENV" --dry-run=client -o yaml \
  | kubectl -n synthkit apply --server-side --field-manager=synthkit-instructor -f -
```

Do not add `tee`, shell tracing (`set -x`), or print the Secret. Kubernetes Secret storage is
not automatically encrypted by base64 encoding; apply cluster encryption/RBAC policy. A secret
manager can own `synthkit-runtime` instead; do not give two managers conflicting ownership.
The chart references only its documented credential keys and never stores their values in the
Helm release. It does not create or delete this Secret.

## 3. Install and verify

Create a nonsecret deployment values file outside this checkout (for example
`/absolute/private/path/workshop-values.yaml`), substituting your actual published digest:

```yaml
image:
  repository: ghcr.io/mbaykara/synthkit
  digest: sha256:REPLACE_WITH_PUBLISHED_DIGEST
existingSecret: synthkit-runtime
dryRun: false
replicaCount: 1
persistence:
  size: 1Gi
  # storageClass: your-storage-class
```

The literal placeholder fails validation. Use `WORKSHOP_DIGEST` from the build or the verified
publishing output. Keep this file as the source of truth for every lifecycle command.

```bash
WORKSHOP_VALUES=/absolute/private/path/workshop-values.yaml
helm lint charts/synthkit -f "$WORKSHOP_VALUES"
helm template synthkit charts/synthkit -n synthkit -f "$WORKSHOP_VALUES"
helm upgrade --install synthkit charts/synthkit -n synthkit \
  -f "$WORKSHOP_VALUES" --wait --timeout 10m
kubectl -n synthkit get pods,pvc
kubectl -n synthkit logs deployment/synthkit --tail=100
```

Rendering contains Secret references, not credentials. Live startup validates configuration;
readiness waits for fresh successful delivery of the intended lanes and writable persisted
state. Initial delivery takes time. A failed `--wait` does not delete the release or PVC and
does not mean the process has stopped; inspect it before retrying the same command.

For an offline smoke run, set `dryRun: true` and install **without `--wait`**. It loads the
blueprint and emits no telemetry. Dry-run remains NotReady because delivery-aware readiness
requires live delivery; that is expected, not a restart-loop. Change back to `false` and repeat
the live upgrade command when credentials are ready.

Open the instructor UI through the Deployment so port-forward also works for a NotReady Pod:

```bash
kubectl -n synthkit port-forward --address 127.0.0.1 deployment/synthkit 8088:8088
```

Visit `http://127.0.0.1:8088/control/ui` and authenticate as `control` with your private
`CONTROL_TOKEN`. Do not retrieve or display the Secret as a workaround. Check:

- Exactly `grafana-cloud-workshop` is loaded/enabled, dry-run is off, and diagnostics are clear.
- All expected sink lanes show recent successful deliveries and persisted state is writable.
- The four [workshop queries](workshop.md) return **fresh** data in the intended Grafana stack.

The unauthenticated `/healthz` liveness endpoint answers process/HTTP health; it does not prove
delivery. The readiness command `/app/synthkit -healthcheck -healthcheck-require-profiles` is
delivery-aware and rejects missing or empty synthetic profile credentials. A sink outage
should make the Pod NotReady, not trigger repeated liveness restarts. Reconnect port-forward
after any replacement Pod. See [control-plane diagnostics](control-plane.md) for detail.

## Operations

### Start again, upgrade, or rotate credentials

Repeat `helm upgrade --install` with the same release, namespace, values, and digest to reconcile
the deployment without duplicating the emitter. It does not clear persisted operator changes.
For an image upgrade, record the old digest, snapshot the PVC using your storage system, and
update the values file to the tested new digest. `Recreate` deliberately introduces a short gap.

For credential rotation, update the private file/secret manager, reconcile the same Secret, then:

```bash
kubectl -n synthkit rollout restart deployment/synthkit
kubectl -n synthkit rollout status deployment/synthkit --timeout=10m
```

Environment-based credentials change only in a new Pod. Verify all four sinks and Grafana queries
again before retiring old credentials. Helm does not detect changes to an externally owned
Secret. For image rollback, select the recorded digest and reuse compatible persisted state;
if state compatibility is uncertain, restore the reviewed pre-upgrade snapshot first.

### Pause and resume

Use the same values file to preserve the image/Secret/storage identity:

```bash
helm upgrade --install synthkit charts/synthkit -n synthkit \
  -f "$WORKSHOP_VALUES" --set replicaCount=0 --wait --timeout 10m
kubectl -n synthkit wait --for=delete pod \
  -l app.kubernetes.io/instance=synthkit,app.kubernetes.io/name=synthkit --timeout=120s
```

Wait for Pod deletion to finish before resuming: Helm's zero-replica readiness alone does not
prove that the old process has finished shutdown. This stops new generation; already sent data
remains queryable. Resume:

```bash
helm upgrade --install synthkit charts/synthkit -n synthkit \
  -f "$WORKSHOP_VALUES" --set replicaCount=1 --wait --timeout 10m
```

Repeated pause/resume commands are safe for this release. In a GitOps deployment, change the
declarative replica count there instead; an external reconciler may undo CLI overrides. A resume
restores persisted scenarios/settings. Reset the exercise explicitly if a healthy baseline is
required. Counter-rate queries need fresh samples after restart; historical data need not be erased.

### Reset the exercise, not the telemetry

With port-forward active, the instructor can deactivate `checkout-regression` in the UI's
Scenarios view. To restore **all** control defaults for the dedicated workshop instance, use:

```bash
curl --fail --silent --show-error --user control --request POST \
  http://127.0.0.1:8088/control/reset
```

Curl prompts for the control password; it is not part of the command line. This reset clears
active scenarios, ad-hoc failures, scaling/load overrides, and disabled controls. It affects
every loaded blueprint, so use it only when this dedicated instance is under your ownership.
Confirm the default state and fresh healthy telemetry in the UI/Grafana. It does not erase
historical telemetry, rewrite the blueprint, or reset in-memory counters. Reset before each
class if desired; ordinary restarts intentionally preserve operator state rather than silently
changing it.

### Uninstall and reinstall

To stop generation and remove this release's Deployment/Service, while retaining state:

```bash
helm uninstall synthkit -n synthkit --ignore-not-found --wait --timeout 10m
kubectl -n synthkit get pvc synthkit-state
```

Repeated uninstall is safe. The chart marks its PVC `helm.sh/resource-policy: keep`; it survives
uninstall. The external Secret also remains. There are no deletion hooks and no calls to delete
Grafana resources or telemetry. Do **not** delete the namespace, PVC, or stack to repeat a class.

For reinstall, add `persistence.existingClaim: synthkit-state` to the same values file, retaining
the other persistence settings. Then use the normal install command. This tells Helm to use the
retained claim rather than trying to recreate it. Verify the claim exists and contains the state
you intend to reuse; never attach it to two active emitters. Keep `existingClaim` in subsequent
upgrades. With your own pre-existing PVC, set this value from the first install; you own its
retention and backups.

Deleting a retained PVC is a separate destructive storage action, not part of workshop reset.
Only do that after reviewing its exact contents, backup, and storage-class reclaim policy. It
does not remove data already ingested into Grafana Cloud.

## Troubleshooting boundary

| Symptom | Instructor checks |
|---|---|
| Pod Pending | PVC binding/storage class, scheduling, resource requests |
| ImagePullBackOff | Digest exists, supported node architecture, registry pull credentials |
| CreateContainerConfigError | External Secret and all required keys exist; never print values |
| State not writable | PVC ownership and CSI support for UID/GID `65532` |
| Running but NotReady | Dry-run intentional? Then sink errors, credentials/endpoints, egress, state health |
| Healthy process but missing profiles | Profile endpoint/user/token scope and per-lane delivery, then Grafana time range |
| Old incident returns on restart | Expected PVC state restoration; deactivate/reset explicitly |

Chart rendering and local tests cannot prove cluster storage compatibility, registry access,
Grafana permissions, or successful live ingestion. Treat those as deployment acceptance checks;
the workshop is ready only after the instructor has verified the actual four-signal queries.
