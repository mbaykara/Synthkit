// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/grafana/grafana-foundation-sdk/go/dashboardv2"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"

	"github.com/rknightion/synthkit/dashboard"
)

// Status colors (good / serious / critical) and the first two categorical slots, light and dark,
// from the validated reference palette. Status colors carry state only; series colors carry data.
const (
	colorGood     = "#0ca30c"
	colorSerious  = "#ec835a"
	colorCritical = "#d03b3b"
	colorNeutral  = "#6e6e6a"
)

var (
	seriesLatency = [2]string{"#2a78d6", "#3987e5"} // slot 1 blue: light, dark
	seriesErrors  = [2]string{"#eb6834", "#d95926"} // slot 2 orange: light, dark
)

// buildControlPlaneDashboard is the instructor console: status tiles, one card per scenario
// with Start/Stop buttons, live impact charts per affected service, and a collapsed row for
// fleet-wide controls. It requires -action-mode infinity so buttons run server-side.
func buildControlPlaneDashboard(o opts) (dashboard.Dashboard, error) {
	d, err := dashboard.NewDashboard("control-plane", "Control Plane")
	if err != nil {
		return dashboard.Dashboard{}, err
	}
	d.Folder = o.folder
	d.Builder.TimeSettings(dashboardv2.NewTimeSettingsBuilder().From("now-3h").To("now").AutoRefresh("30s"))

	scenarios, err := loadScenarios(o.blueprints)
	if err != nil {
		return dashboard.Dashboard{}, err
	}
	act := o.action()
	colored := func(title, path, body, color string) *dashboardv2.ActionBuilder {
		return dashboard.WithActionColor(act(title, path, body), color)
	}

	dashboard.AddPanel(&d, "header", dashboard.TextPanel("",
		"Start a scenario before its session and stop it afterwards: scenarios never stop on their own. "+
			"Buttons apply within seconds; the impact charts follow about a minute later."))

	dashboard.AddPanel(&d, "status-active", dashboard.StatusTile("Active scenarios",
		dashboard.InfinityTarget("A", "/control/state", o.dsName, "",
			dashboard.Col("active_scenarios", "Active scenarios", "string")),
		colorSerious,
		dashboard.TextMapping{Value: "[]", Text: "None", Color: colorGood}))
	dashboard.AddPanel(&d, "status-volume", dashboard.StatTile("Volume multiplier", "none",
		dashboard.InfinityTarget("A", "/control/state", o.dsName, "",
			dashboard.Col("volume_multiplier", "Volume", "number"))))
	dashboard.AddPanel(&d, "status-delivery", dashboard.StatusTile("Telemetry delivery",
		dashboard.InfinityTarget("A", "/control/readiness", o.dsName, "",
			dashboard.Col("live_ready", "Delivery", "string")),
		"",
		dashboard.TextMapping{Value: "true", Text: "Healthy", Color: colorGood},
		dashboard.TextMapping{Value: "false", Text: "Degraded", Color: colorCritical}))

	var cards []dashboard.Cell
	var services []string
	seen := map[string]bool{}
	for i, s := range scenarios {
		info := fmt.Sprintf("card-%d", i)
		buttons := fmt.Sprintf("buttons-%d", i)
		affects := make([]string, len(s.Targets))
		for j, t := range s.Targets {
			affects[j] = "`" + t + "`"
			if !seen[t] {
				seen[t] = true
				services = append(services, t)
			}
		}
		body := `{"scenario":"` + s.id() + `"}`
		dashboard.AddPanel(&d, info, dashboard.TextPanel(s.Title, strings.Join([]string{
			s.Summary,
			"",
			"Affects: " + strings.Join(affects, ", "),
		}, "\n")))
		dashboard.AddPanel(&d, buttons, dashboard.ActionBoardPanel("", o.dsName,
			colored("Start "+s.Title, "/control/scenarios/activate", body, colorSerious),
			colored("Stop", "/control/scenarios/deactivate", body, colorGood)))
		cards = append(cards, dashboard.At(info, 16, 4), dashboard.At(buttons, 8, 4))
	}

	var charts []dashboard.Cell
	width := 24
	if n := 2 * len(services); n > 0 {
		width = max(24/n, 6)
	}
	for _, svc := range services {
		lat, errs := "latency-"+svc, "errors-"+svc
		dashboard.AddPanel(&d, lat, dashboard.EChartsPanel(svc+": mean response time", "s",
			lineChartCode(seriesLatency, "ms"),
			promTarget(o.promUID, fmt.Sprintf(
				`sum(rate(http_server_request_duration_seconds_sum{blueprint=%q,service=%q}[5m])) / sum(rate(http_server_request_duration_seconds_count{blueprint=%q,service=%q}[5m]))`,
				scenarioBlueprint(scenarios, svc), svc, scenarioBlueprint(scenarios, svc), svc))))
		dashboard.AddPanel(&d, errs, dashboard.EChartsPanel(svc+": error log events per minute", "short",
			lineChartCode(seriesErrors, "/min"),
			lokiTarget(o.lokiUID, fmt.Sprintf(
				`sum(count_over_time({blueprint=%q,service_name=%q,level="error"}[1m]))`,
				scenarioBlueprint(scenarios, svc), svc))))
		charts = append(charts, dashboard.At(lat, width, 8), dashboard.At(errs, width, 8))
	}

	dashboard.AddPanel(&d, "fleet", dashboard.ActionBoardPanel("", o.dsName,
		colored("Stop all scenarios", "/control/scenarios", `{"active_scenarios":[]}`, colorGood),
		colored("Normal volume (1x)", "/control/load", `{"volume_multiplier":1}`, colorNeutral),
		colored("Peak volume (3x)", "/control/load", `{"volume_multiplier":3}`, colorSerious)))

	dashboard.WithRows(&d,
		dashboard.Section("Status", dashboard.At("header", 24, 2),
			dashboard.At("status-active", 12, 4), dashboard.At("status-volume", 6, 4), dashboard.At("status-delivery", 6, 4)),
		dashboard.Section("Scenarios", cards...),
		dashboard.Section("Live impact", charts...),
		dashboard.CollapsedSection("Fleet controls", dashboard.At("fleet", 24, 3)),
	)
	return d, nil
}

func scenarioBlueprint(scenarios []scenario, service string) string {
	for _, s := range scenarios {
		if slices.Contains(s.Targets, service) {
			return s.Blueprint
		}
	}
	return ""
}

func promTarget(uid, expr string) *dashboardv2.TargetBuilder {
	return dashboardv2.NewTargetBuilder().RefId("A").
		Query(prometheus.NewQueryV2Builder().Expr(expr).Range(true).
			Datasource(dashboardv2.NewDashboardv2DataQueryKindDatasourceBuilder().Name(uid)))
}

func lokiTarget(uid, expr string) *dashboardv2.TargetBuilder {
	return dashboardv2.NewTargetBuilder().RefId("A").
		Query(loki.NewQueryV2Builder().Expr(expr).
			Datasource(dashboardv2.NewDashboardv2DataQueryKindDatasourceBuilder().Name(uid)))
}

// lineChartCode is the Business Charts getOption body for one series: a thin line with a light
// area fill, a crosshair tooltip, a zero-based single y-axis and recessive grid, in the given
// categorical color stepped for the current theme. suffix "ms" converts seconds to milliseconds.
func lineChartCode(color [2]string, suffix string) string {
	format := `(v) => (v == null ? '-' : Math.round(v) + ' ` + suffix + `')`
	if suffix == "ms" {
		format = `(v) => (v == null ? '-' : Math.round(v * 1000) + ' ms')`
	}
	return strings.Join([]string{
		`const dark = context.grafana.theme.isDark;`,
		`const color = dark ? '` + color[1] + `' : '` + color[0] + `';`,
		`const ink = dark ? '#c3c2b7' : '#52514e';`,
		`const fmt = ` + format + `;`,
		`const frame = context.panel.data.series[0];`,
		`const time = frame && frame.fields.find((f) => f.type === 'time');`,
		`const value = frame && frame.fields.find((f) => f.type === 'number');`,
		`if (!time || !value) {`,
		`  return { title: { text: 'No data in this time range', left: 'center', top: 'middle', textStyle: { color: ink, fontSize: 13, fontWeight: 'normal' } } };`,
		`}`,
		`const t = Array.from(time.values);`,
		`const v = Array.from(value.values);`,
		`return {`,
		`  backgroundColor: 'transparent',`,
		`  grid: { left: 8, right: 16, top: 12, bottom: 4, containLabel: true },`,
		`  tooltip: { trigger: 'axis', axisPointer: { type: 'line' }, valueFormatter: fmt },`,
		`  xAxis: { type: 'time', axisLabel: { color: ink }, axisLine: { lineStyle: { color: ink, opacity: 0.3 } }, splitLine: { show: false } },`,
		`  yAxis: { type: 'value', min: 0, axisLabel: { color: ink, formatter: fmt }, splitLine: { lineStyle: { color: ink, opacity: 0.12 } } },`,
		`  series: [{ type: 'line', data: t.map((x, i) => [x, v[i]]), showSymbol: false, smooth: 0.3,`,
		`    lineStyle: { width: 2, color: color }, itemStyle: { color: color }, areaStyle: { color: color, opacity: 0.12 } }],`,
		`};`,
	}, "\n")
}
