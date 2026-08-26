// SPDX-License-Identifier: AGPL-3.0-only

package blueprint

import (
	"os"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestBundledDOMMQualificationContract(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("../../blueprints/domm-qualification.yaml")
	if err != nil {
		t.Fatalf("read bundled DOMM qualification blueprint: %v", err)
	}

	var demo struct {
		Name      string `yaml:"name"`
		Workloads []struct {
			Name     string `yaml:"name"`
			Services []struct {
				Name      string `yaml:"name"`
				Namespace string `yaml:"namespace"`
				Metrics   []struct {
					Name string `yaml:"name"`
				} `yaml:"metrics"`
			} `yaml:"services"`
		} `yaml:"workloads"`
	}
	if err := yaml.Unmarshal(raw, &demo); err != nil {
		t.Fatalf("parse bundled DOMM qualification blueprint: %v", err)
	}

	if demo.Name != "domm-qualification" {
		t.Fatalf("blueprint name = %q, want domm-qualification", demo.Name)
	}

	wantWorkloads := []string{
		"northstar-checkout-platform",
		"northstar-commerce-platform",
		"northstar-digital-storefront",
		"northstar-fulfillment-platform",
		"northstar-marketplace",
		"northstar-order-intake",
		"northstar-partner-portal",
		"northstar-preview-store",
		"northstar-product-catalog",
	}
	gotWorkloads := make([]string, 0, len(demo.Workloads))
	for _, workload := range demo.Workloads {
		gotWorkloads = append(gotWorkloads, workload.Name)
		for _, service := range workload.Services {
			if service.Namespace != workload.Name {
				t.Errorf("service %q namespace = %q, want workload scope %q", service.Name, service.Namespace, workload.Name)
			}
			for _, metric := range service.Metrics {
				if strings.HasPrefix(metric.Name, "domm") {
					t.Errorf("service %q emits forbidden DOMM-native metric %q", service.Name, metric.Name)
				}
			}
		}
	}
	sort.Strings(gotWorkloads)
	if strings.Join(gotWorkloads, "\n") != strings.Join(wantWorkloads, "\n") {
		t.Fatalf("bundled workloads = %v, want %v", gotWorkloads, wantWorkloads)
	}

	for _, forbidden := range []string{"rubric_version", "criterion_key", "maturity_score", "maturity_level", "domm_team_", "domm:"} {
		if strings.Contains(string(raw), forbidden) {
			t.Errorf("blueprint contains forbidden scorer-native token %q", forbidden)
		}
	}

	for _, maturityScope := range []string{
		"domm-reactive-early",
		"domm-reactive-base",
		"domm-proactive-early",
		"domm-proactive-base",
		"domm-proactive-boundary",
		"domm-proactive-late",
		"domm-systematic-early",
		"domm-systematic-base",
		"domm-nonprod-only",
	} {
		if strings.Contains(string(raw), maturityScope) {
			t.Errorf("blueprint contains maturity-revealing scope %q", maturityScope)
		}
	}
}
