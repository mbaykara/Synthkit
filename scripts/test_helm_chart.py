#!/usr/bin/env python3
"""Offline chart regression tests; requires Helm 3+, Python standard library only."""

import pathlib
import subprocess
import unittest


ROOT = pathlib.Path(__file__).resolve().parents[1]
CHART = ROOT / "charts" / "synthkit"
DIGEST = "sha256:" + "a" * 64


class HelmChartTests(unittest.TestCase):
    def helm(self, *args, ok=True):
        result = subprocess.run(
            ["helm", *args], cwd=ROOT, text=True, capture_output=True, check=False
        )
        if ok:
            self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        else:
            self.assertNotEqual(result.returncode, 0, result.stdout)
        return result.stdout + result.stderr

    def render(self, *args, ok=True):
        return self.helm(
            "template", "workshop", str(CHART), "--namespace", "training",
            "--set", f"image.digest={DIGEST}",
            "--set", "existingSecret=workshop-credentials", *args, ok=ok
        )

    def test_lint(self):
        self.helm("lint", str(CHART), "--strict", "--set-string",
                  f"image.digest={DIGEST}", "--set", "existingSecret=workshop-credentials")

    def test_default_safety_and_runtime_contract(self):
        rendered = self.render()
        for expected in [
            "replicas: 1", "type: Recreate", "type: ClusterIP",
            'image: "ghcr.io/mbaykara/synthkit@' + DIGEST + '"',
            "automountServiceAccountToken: false", "runAsNonRoot: true",
            "runAsUser: 65532", "runAsGroup: 65532", "fsGroup: 65532",
            "readOnlyRootFilesystem: true", "allowPrivilegeEscalation: false",
            "type: RuntimeDefault", "- ALL", 'value: "grafana-cloud-workshop"',
            'value: "/data/control-state.json"', 'value: "/data/blueprints"',
            'value: "trusted-network"', 'value: "false"',
            'name: DRY_RUN\n              value: "true"', "path: /healthz",
            "- /app/synthkit", "- -healthcheck", "- -healthcheck-require-profiles", 'key: CONTROL_TOKEN',
            'helm.sh/resource-policy: keep', "mountPath: /data",
        ]:
            self.assertIn(expected, rendered)
        for forbidden in ["kind: Secret", "helm.sh/hook", "envFrom:", "hostPort:", "kind: Ingress"]:
            self.assertNotIn(forbidden, rendered)
        self.assertEqual(rendered.count("kind: Deployment"), 1)
        self.assertEqual(rendered.count("kind: PersistentVolumeClaim"), 1)

    def test_live_requires_all_sink_keys(self):
        rendered = self.render("--set", "dryRun=false")
        self.assertIn('name: DRY_RUN\n              value: "false"', rendered)
        for key in ["GC_TOKEN", "GC_PROM_RW", "GC_PROM_USER", "GC_OTLP_ENDPOINT",
                    "GC_OTLP_USER", "GC_LOKI", "GC_LOKI_USER", "GC_PROFILES_URL", "GC_PROFILES_USER"]:
            self.assertIn(f"key: {key}\n                  optional: false", rendered)

    def test_pause_preserves_state(self):
        rendered = self.render("--set", "replicaCount=0")
        self.assertIn("replicas: 0", rendered)
        self.assertIn("kind: PersistentVolumeClaim", rendered)

    def test_existing_claim_not_owned(self):
        rendered = self.render("--set", "persistence.existingClaim=retained-state")
        self.assertNotIn("kind: PersistentVolumeClaim", rendered)
        self.assertIn('claimName: "retained-state"', rendered)

    def test_storage_class_modes(self):
        self.assertNotIn("storageClassName:", self.render())
        self.assertIn('storageClassName: "fast"', self.render("--set", "persistence.storageClass=fast"))
        self.assertIn('storageClassName: ""', self.render("--set-string", "persistence.storageClass="))

    def test_missing_inputs_fail(self):
        self.helm("template", "workshop", str(CHART), ok=False)
        self.render("--set-string", "existingSecret=", ok=False)
        self.render("--set-string", "image.digest=", ok=False)

    def test_invalid_inputs_fail(self):
        for setting in ["replicaCount=2", "replicaCount=-1", "image.digest=latest",
                        "image.tag=latest", "service.type=LoadBalancer", "persistence.enabled=false",
                        "existingSecret=invalid/name", "image.repository=example.com/image:main"]:
            with self.subTest(setting=setting):
                self.render("--set", setting, ok=False)

    def test_release_name_and_namespace(self):
        rendered = self.render()
        self.assertIn("name: workshop-synthkit", rendered)
        self.assertIn("namespace: training", rendered)

    def test_live_example(self):
        self.render("-f", str(CHART / "values-live.example.yaml"))


if __name__ == "__main__":
    unittest.main()
