---
title: DOMM Qualification Demo
description: Run the bundled production-like application estate used to qualify all seven DOMM maturity levels.
---

# DOMM qualification demo

The bundled `domm-qualification` blueprint provides a ready-to-run, production-like application estate for demonstrating and qualifying the Dynamic Observability Maturity Model (DOMM). It models nine application workloads across production and staging, including deliberately uneven adoption of metrics, logs, traces, profiles, release attribution, and alertable failure signals.

Synthkit emits only ordinary application and Kubernetes telemetry. It does not emit DOMM facts, criteria, scores, levels, recording-rule output, dashboards, alerts, SLOs, or incident-management resources. Those remain independent evidence supplied and evaluated by DOMM. This separation keeps the demo representative of a real estate instead of teaching the scorer its expected answer.

## Offline proof

No credentials are required to load the exact bundled blueprint and inspect its inventory:

```bash
go build ./cmd/synthkit
DRY_RUN=true BLUEPRINT_NAMES=domm-qualification ./synthkit -once -dump
```

The command must report exactly one selected blueprint named `domm-qualification`. It pushes nothing.

## Live Grafana Cloud run

Create the private configuration file, select only the demo, and opt into live delivery:

```bash
install -m 600 .env.example .env
./plugins/synthkit/skills/initial-setup/scripts/set-env.sh BLUEPRINT_NAMES domm-qualification .env
./plugins/synthkit/skills/initial-setup/scripts/set-env.sh DRY_RUN false .env
```

Fill these values in `.env` without printing or committing them:

- `GC_TOKEN`
- `GC_PROM_RW` and `GC_PROM_USER`
- `GC_OTLP_ENDPOINT` and `GC_OTLP_USER`
- `GC_LOKI` and `GC_LOKI_USER`
- `GC_PROFILES_URL` and `GC_PROFILES_USER` for the complete profiles lane

Validate authentication without emitting telemetry, then start the continuous generator:

```bash
docker compose run --rm --no-deps synthkit -preflight
docker compose up -d --wait synthkit
curl -fsS http://127.0.0.1:8088/control/status | jq -e '.dry_run == false'
```

The long-running service emits continuously until stopped with `docker compose stop synthkit`. The `meaningful-alert-cycle` and `proactive-early-alert-flap` scenarios can be activated from the control plane when validating alert fire, recovery, and noise behavior.

## Qualifying DOMM

Deploy the DOMM engine and its qualification fixtures from the [DOMM repository](https://github.com/grafana-ps/domm). That harness adds the Grafana-side resources that telemetry alone cannot represent, such as meaningful alert rules, curated dashboards, and later-stage operational evidence, and then verifies the independently calculated L1 through L7 results.

The blueprint contains names such as `domm-reactive-early` and `domm-systematic-base` to make the intended test cases readable. Names are selectors only; they do not affect scoring. A level is valid only when DOMM derives it from the emitted signals and the actual Grafana resources.
