#!/usr/bin/env python3
"""Offline publishing contract checks; workflow syntax is checked by actionlint."""

import pathlib
import re
import unittest

ROOT = pathlib.Path(__file__).resolve().parents[1]


class PublishWorkflowTests(unittest.TestCase):
    def setUp(self):
        self.workflow = (ROOT / ".github/workflows/publish.yml").read_text()

    def job(self, name):
        match = re.search(r"^  " + re.escape(name) + r":\n(.*?)(?=^  [a-z][a-z-]*:|\Z)",
                          self.workflow, re.M | re.S)
        self.assertIsNotNone(match, name)
        return match.group(1)

    def test_fork_publisher_is_pinned_and_lowercase(self):
        job = self.job("image-fork")
        self.assertIn("github.repository == 'mbaykara/Synthkit'", job)
        self.assertIn("github.ref == 'refs/heads/main'", job)
        self.assertIn("needs: identity", job)
        self.assertIn("container-publish.yml@db2707eeb563706efc1b92ec969df0aba0193e02", job)
        self.assertIn("image-name: synthkit", job)
        self.assertIn("packages: write", job)
        self.assertIn("REVISION=${{ needs.identity.outputs.revision }}", job)
        self.assertNotIn("secrets:", job)  # automatic GITHUB_TOKEN, no PAT or inherited secrets
        for feature in ("sign", "sbom", "provenance", "trivy"):
            self.assertNotIn(feature + ": false", job)

    def test_upstream_release_trust_is_unchanged(self):
        job = self.job("image")
        self.assertIn("github.repository == 'rknightion/synthkit'", job)
        self.assertIn("container-publish.yml@f31690684f4292d1fe8e528618f7c8306fe27d9a", job)
        self.assertIn("github.repository == 'rknightion/synthkit'", self.job("notices"))
        self.assertIn("needs: [identity, image]", self.job("verify-release"))

    def test_fork_releases_fail_explicitly(self):
        identity = self.job("identity")
        self.assertIn('"${GITHUB_REPOSITORY}" != "rknightion/synthkit"', identity)
        self.assertIn('fork publishing supports main edge builds only', identity)


if __name__ == "__main__":
    unittest.main()
