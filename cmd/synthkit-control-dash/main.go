// SPDX-License-Identifier: AGPL-3.0-only

// Command synthkit-control-dash generates the CUSTOMER self-serve control dashboard: an
// Infinity-datasource-backed Grafana v2 dashboard exposing only the customer-safe knobs
// (master volume + incident scenarios) as read panels + native fetch-POST action buttons.
// Reads come from the synthkit control plane's GET routes (?audience=customer); writes POST
// to /control/load and /control/scenarios. The operator UI (/control/ui) is unaffected.
// Protected reads use the Infinity datasource's secure Basic auth. Browser-direct POST routes use
// a separate HTTP Basic challenge; no token is embedded in the dashboard.
package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/rknightion/synthkit/dashboard"
)

const (
	actionModeFetch    = "fetch"
	actionModeInfinity = "infinity"

	layoutCustomer     = "customer"
	layoutControlPlane = "control-plane"
)

type opts struct {
	writeBaseURL string
	dsName       string
	dsUID        string
	actionMode   string
	layout       string
	folder       string
	promUID      string
	lokiUID      string
	outDir       string
	blueprints   string
}

func (o opts) validate() error {
	switch o.layout {
	case "", layoutCustomer:
	case layoutControlPlane:
		if o.actionMode != actionModeInfinity {
			return errors.New("-layout control-plane requires -action-mode infinity")
		}
		if o.promUID == "" || o.lokiUID == "" {
			return errors.New("-layout control-plane requires -prom-uid and -loki-uid")
		}
	default:
		return errors.New("-layout must be customer or control-plane")
	}
	switch o.actionMode {
	case "", actionModeFetch:
		return nil
	case actionModeInfinity:
		if o.dsUID == "" || o.writeBaseURL == "" {
			return errors.New("-action-mode infinity requires -ds-uid and -write-base-url (the URL the datasource reaches)")
		}
		return nil
	default:
		return errors.New("-action-mode must be fetch or infinity")
	}
}

func main() {
	var o opts
	// Reads are RELATIVE paths resolved against the Infinity datasource's Base URL (no read base
	// here). Writes are browser fetches → need an absolute, browser-reachable base URL. The default
	// is the browser-trusted tailscale-serve endpoint; OVERRIDE per-deploy.
	flag.StringVar(&o.writeBaseURL, "write-base-url", "", "action-button POST base URL (absolute, HTTPS, browser-reachable; per-deploy)")
	flag.StringVar(&o.dsName, "ds-name", "", "Infinity datasource name (required)")
	flag.StringVar(&o.actionMode, "action-mode", actionModeFetch, "fetch (browser POST) or infinity (server-side via the datasource; needs Grafana toggle vizActionsAuth)")
	flag.StringVar(&o.dsUID, "ds-uid", "", "Infinity datasource UID (required with -action-mode infinity)")
	flag.StringVar(&o.layout, "layout", layoutCustomer, "customer (self-serve knobs) or control-plane (instructor console with scenario cards and impact charts)")
	flag.StringVar(&o.folder, "folder", "", "target Grafana folder UID (empty = General)")
	flag.StringVar(&o.promUID, "prom-uid", "", "Prometheus datasource UID for control-plane impact charts")
	flag.StringVar(&o.lokiUID, "loki-uid", "", "Loki datasource UID for control-plane impact charts")
	flag.StringVar(&o.outDir, "out", "", "output directory (required)")
	flag.StringVar(&o.blueprints, "blueprints", "./blueprints", "directory of *.yaml blueprints to enumerate scenarios from")
	flag.Parse()
	if o.dsName == "" || o.outDir == "" {
		log.Fatal("synthkit-control-dash: -ds-name and -out are required")
	}
	if err := generate(o); err != nil {
		log.Fatalf("synthkit-control-dash: %v", err)
	}
}

func generate(o opts) error {
	if err := o.validate(); err != nil {
		return err
	}
	build := buildControlDashboard
	if o.layout == layoutControlPlane {
		build = buildControlPlaneDashboard
	}
	d, err := build(o)
	if err != nil {
		return err
	}
	js, err := dashboard.Render(d)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(o.outDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(o.outDir, d.UID+".json")
	return os.WriteFile(path, js, 0o644)
}
