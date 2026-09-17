# SynthKit Helm chart

Runs one continuous emitter with the builtin `grafana-cloud-workshop` blueprint. This chart
does not install Grafana, provision cloud fixtures, run DOMM Terraform, or erase telemetry.
Build the image from this revision or newer: older images do not contain the workshop blueprint.
Use the [Kubernetes deployment guide](../../docs/kubernetes.md) for image publication, private
credential preparation, rollout troubleshooting, and the full operating procedure.
For concrete workshop, instrumentation-gap, profiling and DOMM exercises, see
[practical use-case examples](../../docs/blueprint-examples.md#practical-use-cases).
This chart currently fixes the workshop blueprint; it does not yet support selecting the
other examples through Helm values.

## Install

Requires Helm 3+, Kubernetes 1.25+, a writable persistent volume supporting ownership by
UID/GID 65532, and an externally managed Secret in the release namespace. Select a supported
Kubernetes version for your environment; the version constraint is a manifest compatibility floor.

The Secret must contain `CONTROL_TOKEN` (HTTP Basic password for username `control`). Live mode
also requires `GC_TOKEN`, `GC_PROM_RW`, `GC_PROM_USER`, `GC_OTLP_ENDPOINT`, `GC_OTLP_USER`,
`GC_LOKI`, `GC_LOKI_USER`, `GC_PROFILES_URL`, and `GC_PROFILES_USER`. Supply a least-privilege
Grafana Cloud access-policy token permitting metrics, logs, traces, and profiles writes to the
training stack. Keep credentials out of values files and `--set` arguments: Helm stores values
in release history. The chart reads only these explicit keys, never all Secret entries.

From the repository root, set `SYNTHKIT_DIGEST` to the published image's real `sha256:...` digest,
then install a safe preview:

```sh
helm upgrade --install synthkit ./charts/synthkit \
  --namespace synthkit --create-namespace \
  --set-string image.digest="$SYNTHKIT_DIGEST" \
  --set existingSecret=synthkit-runtime
kubectl -n synthkit logs deployment/synthkit
kubectl -n synthkit port-forward deployment/synthkit 8088:8088
```

Open `http://127.0.0.1:8088/control/ui`. Dry-run intentionally stays **NotReady**, because readiness
proves live delivery. Do not use `--wait` in preview mode. Use deployment port-forward rather than
the Service while unready. Dry-run still persists instructor control changes to the state volume.

After adding all live keys to the Secret:

```sh
helm upgrade --install synthkit ./charts/synthkit \
  --namespace synthkit --create-namespace \
  -f charts/synthkit/values-live.example.yaml \
  --set-string image.digest="$SYNTHKIT_DIGEST" \
  --wait --timeout 5m
```

Allow at least a minute for the first metric tick and queue flush. A timeout is a diagnostic
signal, not permission to erase state. Inspect Pod events, PVC binding, application logs and
authenticated `/control/status`. Failed delivery affects readiness but not liveness; a downstream
outage must not cause restart loops. Empty or invalid Secret values may still fail runtime
validation; Helm can validate references, not external Secret contents or endpoint reachability.

## Values

| Value | Default | Purpose |
| --- | --- | --- |
| `image.repository` | `ghcr.io/mbaykara/synthkit` | Registry/repository without tag or digest |
| `image.digest` | empty, required | Immutable `sha256:` plus 64 lowercase hex characters |
| `existingSecret` | empty, required | Same-namespace external runtime Secret |
| `imagePullSecrets` | `[]` | Existing registry Secret references, for example `[{name: ghcr-pull}]` |
| `dryRun` | `true` | Disable synthetic ingestion; readiness remains false |
| `replicaCount` | `1` | Only `0` (pause) or `1` (emit) accepted |
| `persistence.existingClaim` | empty | Use an existing claim instead of creating one |
| `persistence.size` | `1Gi` | Size of chart-created state claim |
| `persistence.storageClass` | `null` | Cluster default; empty string selects no class |
| `persistence.accessModes` | `[ReadWriteOnce]` | One mode; `ReadWriteOncePod` also accepted if supported by the cluster and CSI driver |
| `resources` | 100m/256Mi requested; 1 CPU/1Gi limit | Starting budget, not a benchmark; tune after observing the workload |
| `nodeSelector` | `{}` | Pod placement |
| `tolerations` | `[]` | Pod scheduling tolerations |
| `affinity` | `{}` | Pod affinity/anti-affinity |

For release `synthkit`, Deployment and Service names are `synthkit`, claim `synthkit-state`.
Other releases use `<release>-synthkit` and `<release>-synthkit-state` (truncated to DNS label
limits). Unknown top-level values fail validation rather than silently changing nothing.

## Operate without destructive cleanup

Pause and resume using the same release; no telemetry erasure is needed:

```sh
helm upgrade synthkit ./charts/synthkit -n synthkit --reuse-values --set replicaCount=0
kubectl -n synthkit wait --for=delete pod -l app.kubernetes.io/instance=synthkit --timeout=120s
helm upgrade synthkit ./charts/synthkit -n synthkit --reuse-values --set replicaCount=1
```

Control settings, including active scenarios, persist across restarts; pausing does not reset
them. Reset instructor changes through the control UI before the next session. Counters restart
in memory, which `rate()` handles as a counter reset. No telemetry WAL exists or is required.

Rotate credentials in the external Secret, then run
`kubectl -n synthkit rollout restart deployment/synthkit`; environment variables are loaded only
at process startup. Follow with `kubectl -n synthkit rollout status deployment/synthkit --timeout=5m`
in live mode. Do not expose the token to participants.

`helm uninstall synthkit -n synthkit` removes only the emitter Deployment and Service. The
chart-created PVC carries `helm.sh/resource-policy: keep`; the external Secret is never owned
by this chart. Historical Grafana Cloud data expires according to the stack's retention policy.
To reinstall using the retained claim, use the install command above with
`--set persistence.existingClaim=synthkit-state`. Do not delete the namespace or PVC as a routine
restart step: namespace deletion bypasses chart retention. Volume deletion is an explicit,
separate destructive operation, not part of this workflow.

## Safety boundaries

The chart fixes `Recreate`, one replica, an internal `ClusterIP` Service, nonroot execution,
read-only root filesystem, disabled service-account token mount and all Linux capabilities
dropped. State is mounted as a directory at `/data` for atomic saves; `/tmp` is a bounded
ephemeral volume. The pod does not need Kubernetes API permissions. There are no hooks, init
cleanup jobs, ingress, autoscaling, or fixture provisioning commands.

`CONTROL_EXPOSURE_ACK=trusted-network` acknowledges the pod's non-loopback HTTP listener.
ClusterIP is **not** network isolation or TLS: deploy only on a trusted cluster/network, restrict
network access and Kubernetes RBAC separately, and use localhost port-forward for the instructor.
Participants need Grafana access only. Self-observability is explicitly disabled for this chart;
synthetic profiles use `GC_PROFILES_*`, not the process-profiling `GC_PYROSCOPE_*` credentials.

Recreate prevents overlapping pods during normal Deployment updates; it is not a distributed
leader lock. Never create a second release/emitter with these same synthetic identities against
the same stack, force-delete a pod on an unreachable node, or manually scale above one.

Run local, network-free chart regression checks with:

```sh
python3 scripts/test_helm_chart.py
```
