---
title: Blueprint Examples
description: A tour of every bundled blueprint — what each models, which constructs it exercises, and which to copy for your use case.
---

# Blueprint Examples

Synthkit ships ready-to-run blueprints in `blueprints/`. At startup it loads only the exact
runtime names selected by `BLUEPRINT_NAMES`; an empty selection emits no synthetic telemetry.
The `*` selector explicitly opts into the complete catalog and is not needed for these examples.

To author a variation, change the blueprint name and conflicting resource identities, validate
the YAML, then stage it through the [custom blueprint workflow](custom-blueprints.md). Select
the returned namespaced runtime identity before restarting. Uploading alone does not enable it.

## Practical use cases

A **blueprint** is the use-case configuration. A **scenario** is an incident definition inside
that blueprint which the instructor activates and deactivates. The YAML `profiles:` list means
reusable telemetry templates; actual continuous profiling uses the separate `pyroscope:` block.

| Use case | Select this blueprint | Exercise |
|---|---|---|
| General Grafana Cloud workshop | `grafana-cloud-workshop` | Follow a checkout request across metrics, logs, traces and CPU profiles |
| Instrumentation-gap training | `grafana-cloud-workshop` | Compare fully instrumented checkout with metrics-only inventory and metrics/logs shipping |
| Profiling specialist training | `profiling-demo` | Compare SDK-push, eBPF and Java-collected profiling shapes |
| DOMM qualification | `domm-qualification` | Generate uneven adoption and failures; let DOMM independently evaluate the evidence |

Run the commands below from the repository root. Stop the previous emitter before switching
use cases; do not create duplicate emitters for the same identities. Keep credentials in the
private `.env`, never in these commands. Live use requires the metrics/logs/traces credentials
and, for the profiling exercises, the synthetic `GC_PROFILES_*` settings described in
[Credentials](credentials.md). Do not overwrite an existing `.env`.

### 1. Workshop: investigate a slow checkout

Preview without sending telemetry:

```bash
DRY_RUN=true BLUEPRINT_NAMES=grafana-cloud-workshop \
  go run ./cmd/synthkit -once -dump
```

Once credentials are configured, start continuous generation:

```bash
DRY_RUN=false BLUEPRINT_NAMES=grafana-cloud-workshop go run ./cmd/synthkit
```

Allow at least five minutes of baseline data. In the instructor control UI, activate
`grafana-cloud-workshop/checkout-regression`. Compare checkout latency, error logs, request
spans and CPU profiles before/during the incident. Deactivate it after eight minutes, then
verify recovery using recent data. Activation does not auto-expire or create a Grafana incident.
The complete queries and timing are in the [90-minute workshop](workshop.md).

### 2. Training: distinguish missing instrumentation from an outage

Use the same running workshop, with no incident active. Have participants compare:

- `shop-checkout`: metrics, application logs, traces and CPU profiles.
- `shop-inventory`: application metrics only.
- `shop-shipping`: application metrics and logs, but no traces or profiles.

Ask: "Can you identify a failing request from metrics alone? Does a missing trace prove Tempo
is broken? What would you instrument next?" Expected result: the empty inventory log query and
empty shipping trace query are intentional adoption gaps, while checkout proves those sinks
are receiving data. Kubernetes infrastructure events are separate from application logs.

### 3. Specialist training: CPU and memory profiling

Stop the previous local process with Ctrl-C, then preview and run this separate estate:

```bash
DRY_RUN=true BLUEPRINT_NAMES=profiling-demo go run ./cmd/synthkit -once -dump
# After checking the preview and configuring the profile credentials:
DRY_RUN=false BLUEPRINT_NAMES=profiling-demo go run ./cmd/synthkit
```

Compare `profiling-api` (Go SDK-push) with the cluster-collected profiles in Pyroscope. In
the control UI, use the loaded schema's failure controls to enable `cpu_hotspot` on
`profiling-api`, record the interval, then disable it. Repeat separately with `memory_leak`.
Compare CPU and memory profile types rather than treating them as the same measurement.
Go span-profile correlation applies to CPU only. This blueprint has a business-hours traffic
curve, unlike the workshop's equal request-rate bounds; verify activity before class.

### 4. DOMM: alert fire and recovery evidence

```bash
DRY_RUN=true BLUEPRINT_NAMES=domm-qualification go run ./cmd/synthkit -once -dump
# After checking the preview and configuring credentials:
DRY_RUN=false BLUEPRINT_NAMES=domm-qualification go run ./cmd/synthkit
```

On a dedicated qualification stack, activate
`domm-qualification/meaningful-alert-cycle`, observe the generated failure signals, then
deactivate it and observe recovery. A second exercise uses
`domm-qualification/proactive-early-alert-flap`: activate and deactivate it at instructor-recorded
times to examine alert sensitivity and noise. The scenario name does not schedule automatic
flapping. Alert evaluation depends on the independently installed rules and their evaluation
windows; synthetic errors alone do not prove an alert fired.

Synthkit does not create DOMM scores, dashboards, alert rules, SLOs, teams or qualification
fixtures. Use the separate DOMM harness when those are required, as described in the
[DOMM qualification guide](domm-qualification.md). Do not run fixture teardown as part of
ordinary telemetry-generator startup or shutdown.

### Instructor API: activate and recover repeatably

For an authenticated local instance, the same workshop scenario can be controlled from another
terminal. On Kubernetes, first use the instructor port-forward from the deployment guide.
Curl prompts for the control password, keeping it out of the command line:

```bash
curl --fail --silent --show-error --user control \
  -H 'Content-Type: application/json' \
  -d '{"scenario":"grafana-cloud-workshop/checkout-regression"}' \
  http://127.0.0.1:8088/control/scenarios/activate

# End the exercise, including when the investigation finishes early:
curl --fail --silent --show-error --user control \
  -H 'Content-Type: application/json' \
  -d '{"scenario":"grafana-cloud-workshop/checkout-regression"}' \
  http://127.0.0.1:8088/control/scenarios/deactivate
```

For DOMM, replace the qualified scenario name with one of the DOMM names above and ensure that
blueprint is loaded. Do not change the route to an invented scenario name. Review
`/control/schema` or the UI for the actual available scenarios and failure targets.

Control changes persist across restart. Deactivate scenarios and disable ad-hoc failures before
the next class, or explicitly reset the dedicated instance's controls. Retain historical
telemetry and use fresh time windows; stopping generation does not erase past data.

### Kubernetes support today

The current [Helm chart](../charts/synthkit/README.md) selects `grafana-cloud-workshop` internally.
The workshop and instrumentation-gap exercises above work with that chart as shipped.
There is **no Helm value for blueprint selection yet**: setting `blueprintNames`, `profile`
or `extraEnv` is not supported. Use the local commands above for `profiling-demo` or DOMM until
chart selection is implemented. Uploading a custom blueprint also requires changing the startup
selection and restarting; the current chart does not expose that configuration.

Image publication remains a separate prerequisite. Follow the
[Kubernetes deployment guide](kubernetes.md) and pin an actually published digest, not a local
validation image or an assumed registry tag.

---

## Kubernetes

### [k8s-minimal.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/k8s-minimal.yaml)

Cheapest k8s-monitoring footprint: `cluster_metrics` only (KSM + cAdvisor + kubelet + node-exporter), no logs, events, profiling, OpenCost, or Kepler. Use this as the starting point for any Kubernetes blueprint. Exercises: `k8s_cluster`, `ec2`.

### [k8s-full-stack.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/k8s-full-stack.yaml)

Maximal k8s observability: every collector feature enabled — cluster metrics, events, pod logs, node logs, profiling, application observability — plus OpenCost cost allocation, Kepler energy monitoring, Fleet Management, control-plane deep monitoring, full addon set, Karpenter autoscaler, and Bottlerocket nodes. The reference for teams wanting everything at once. Exercises: `k8s_cluster`, `ec2`, `k8s_profiling`, `karpenter`, `cert_manager`, `coredns`, `vpc_cni`, `ebs_csi`, `argocd`, `envoy_gateway`, `external_dns`, `load_balancer_controller`, `fleetmgmt`.

### [k8s-cost-power.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/k8s-cost-power.yaml)

FinOps focus: OpenCost workload cost allocation + Kepler per-pod energy consumption on a standard EKS cluster. No logs or profiling overhead. Exercises: `k8s_cluster`, `k8s_profiling` (Kepler), OpenCost sub-family.

### [k8s-control-plane.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/k8s-control-plane.yaml)

EKS control-plane deep monitoring: all five control-plane component metric families (API server, kube-proxy, scheduler, controller-manager, kubelet probes) with cluster metrics. Reference for teams focused on k8s internals. Exercises: `k8s_cluster` control-plane sub-families.

### [k8s-logs-events.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/k8s-logs-events.yaml)

Logs-centric monitoring: pod logs + node logs + cluster events shipped to Loki via Alloy DaemonSet and singleton collectors. No profiling, OpenCost, or Kepler. Exercises: `k8s_cluster` logs + events features.

### [k8s-windows-mixed.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/k8s-windows-mixed.yaml)

Mixed Linux + Windows EKS node groups: exercises both the windows-exporter signal path (windows-pool) and the standard Linux node-exporter path (linux-general) with node-level log collection. Reference for teams running .NET or legacy workloads on Windows nodes. Exercises: `k8s_cluster` mixed-OS node groups.

---

## AWS / CloudWatch

### [cw-infra-aws.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/cw-infra-aws.yaml)

AWS CloudWatch infrastructure showcase: explicit sub-family toggles covering ALB/NLB/EBS/NAT/EKS/S3/Firehose/PrivateLink, plus RDS and ElastiCache CloudWatch lanes. Demonstrates every `cw_infra` sub-family switch. Exercises: `cw_infra`, `rds`, `elasticache`, `ec2`.

### [aws-cloud-services.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/aws-cloud-services.yaml)

AWS managed data/ETL services: OpenSearch Serverless (AOSS), Managed Workflows for Apache Airflow (MWAA), Glue ETL, DocumentDB, and Neptune. Focused on the cloud-service constructs that rarely appear in a basic k8s blueprint. Exercises: `aoss`, `mwaa`, `glue`, `docdb`, `neptune`.

---

## Databases

### [dbo11y-mysql.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/dbo11y-mysql.yaml)

Demonstrates the Database Observability MySQL lane: an RDS MySQL instance emitting `database_observability_*` + `mysql_*` metric families, log ops, replication (slave-status metrics), and the `query_data_locks` op (which appears only while a `lock_contention` incident is active). Pair with an incident targeting the db name. Exercises: `dbo11y_mysql`, `rds`.

---

## CSP Azure / GCP

### [csp-azure.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/csp-azure.yaml)

CSP Azure integration: `azure_microsoft_*` window-gauge metrics across compute, databases, storage, networking, messaging, and Event Hubs logs, via the serverless managed scraper or `azure_exporter` path. Demonstrates all `sub_signals` families and the `ingestion_path` discriminator. Exercises: `csp_azure`.

---

## AI / LLM

### [acme-ai-platform.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/acme-ai-platform.yaml)

The API-poll view of an AI assistant estate across multiple environments: AgentCore-vended AWS metrics (7 runtimes, no in-account Bedrock model inference), LLM gateway observed via the Portkey analytics poller (`portkey_api_*`) + LangSmith eval bridge, full traced estate across 8 EKS deployments and 4 request journeys. Exercises: `agentcore`, `portkeypoller`, `langsmithplatform`, `langsmitheval`, `rds` (Aurora PostgreSQL), `docdb`, `neptune`, `elasticache`, `k8s_cluster`, `app` workload with `gen_ai_client` + `agentic_flow`.

### [acme-ai-platform-eval.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/acme-ai-platform-eval.yaml)

The same assistant estate as `acme-ai-platform` but with the AI evaluation gateway exposed as a connected trace (Path-B gateway span), modelling the per-tenant gateway slice. Designed to run concurrently with `acme-ai-platform` using disjoint identities. Exercises: `portkeygateway` (connected trace), `bedrock`, `agentcore`, `app` workload with `gateway_export_log` + `gateway_native_scrape` profiles.

### [acme-ai-eval.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/acme-ai-eval.yaml)

The AI-gateway platform operator's view across an 8-cell estate (4 AWS account roles × 2 regions): gateway health, LangSmith platform health, per-cell AWS estate, multi-cloud LLM endpoints, edge, and the qualification pipeline. No single-tenant app traces. Exercises: `portkeygateway`, `portkeypoller`, `langsmithplatform`, `langsmitheval`, `qualificationpipeline`, `bedrock`, `csp_azure` (LLM-access footprint), `csp_gcp` (LLM-access footprint).

---

## Hosts

### [hostfleet.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/hostfleet.yaml)

Mixed-OS host fleet: Linux/Windows/macOS machines running Grafana Alloy's node/windows/macos exporter, plus optional Docker cAdvisor. Integration and full metric profiles. Exercises: `host` construct across all three OSes.

### [hosts-bare.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/hosts-bare.yaml)

Bare hosts with no container runtime: demonstrates the `docker: false` dimension and `observability.logs: false` (metrics-only, no log streams) across Linux/Windows/macOS. Exercises: `host`, logs-off configuration.

### [hosts-linux-docker.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/hosts-linux-docker.yaml)

Linux container hosts running node_exporter + Docker cAdvisor: container CPU/memory/network/filesystem metrics plus container log streams. Exercises: `host` with `docker: true` lane.

### [hosts-macos.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/hosts-macos.yaml)

macOS endpoint and developer fleet: macos-node exporter metrics (cpu/disk/net/fs + battery/power) on developer laptops and a CI runner. No Docker (unsupported on macOS in v1). Exercises: `host` macOS OS path.

### [hosts-windows.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/hosts-windows.yaml)

Windows Server estate: windows_exporter metrics + Application/System event log streams. Domain Controller, app server, and SQL Server roles on Server 2022/2025. No Docker. Exercises: `host` Windows path, event log streams.

---

## Network Topology

### [netobs-enterprise.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/netobs-enterprise.yaml)

"Average enterprise" archetype: one `network_topology` exporter watching a 2-spine / 6-leaf access fabric (mixed Arista + Cisco), standalone mode, LLDP/CDP/BGP, prod-realistic cold-start discovery churn. Exercises: `nettopo` construct, standalone sub-families.

### [netobs-global.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/netobs-global.yaml)

Maximal network topology: a federation HUB aggregating a 6-spine / 24-leaf Clos fabric across four vendors (Arista, Cisco, Juniper, Nokia), five spoke sites, all seven discovery protocols, OTLP push. Exercises: `nettopo`, federation sub-families (`federation_spoke_*`, `boundary_observation_info`), OTLP push families.

### [netobs-spoke.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/netobs-spoke.yaml)

A remote-site SPOKE `network_topology` exporter pushing its local graph to the federation hub (`netobs-global`). Exercises the spoke-side liveness families a hub/standalone deployment never emits. Exercises: `nettopo` spoke sub-families.

---

## Synthetic Monitoring / Fleet Management

### [synthetic-monitoring.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/synthetic-monitoring.yaml)

Grafana Synthetic Monitoring estate: HTTP probe checks (`probe_success` / `probe_duration_seconds` families) plus a Fleet Management collector roster across linux/windows/darwin. No cloud infrastructure — all telemetry flows from the `features:` block. Exercises: `sm`, `fleetmgmt`.

### [fleet-management.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/fleet-management.yaml)

Standalone Fleet Management showcase: a roster of synthetic Alloy collectors across linux/windows/darwin, emitting the Alloy self-metric set and registering with the Fleet Management API when `GC_FM_*` credentials are present. Exercises: `fleetmgmt`.

---

## Profiling

### [profiling-demo.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/profiling-demo.yaml)

End-to-end Pyroscope profiling: k8s eBPF process_cpu per pod, a `web_service` with SDK-push profiles + span profiles, and an `app` service-graph with per-node profiles. Exercises: `k8s_profiling`, `web_service` with `pyroscope:` block, `app` workload with per-node `pyroscope:` blocks.

---

## Native OTLP

### [otlp-native.yaml](https://github.com/rknightion/synthkit/blob/main/blueprints/otlp-native.yaml)

Native OTLP application-metrics showcase: two `web_service` workloads (one in `k8s_monitoring` mode, one in `naked` mode) emit `http.server.request.duration` and `http.server.active_requests` as OTLP/HTTP to `/v1/metrics`, letting the Grafana Cloud OTLP gateway own the Prometheus translation. Exercises: `web_service` with `otel.metrics: true`, both `mode: k8s_monitoring` and `mode: naked`.
