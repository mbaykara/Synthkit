---
title: Grafana Cloud Workshop
description: Grafana Cloud Foundation in Action – drei Labs auf einem gemeinsamen Stack in 90 Minuten.
---

# Grafana Cloud Foundation in Action

**Ein gemeinsamer Stack, 90 Minuten, drei Labs:** ein Signal finden, eine Frage anhand von
zwei Signaltypen prüfen und den nützlichen Ablauf als persönlichen Quickstart wiederverwenden.
Die Labs gehören zu den Deck-Folien **9, 14 und 27**. Die Zeitaufteilung unten ist ein
Moderationsvorschlag; die Foliennummern stammen aus den Lab-Karten.

Die Kursleitung betreibt Synthkit auf Kubernetes. Alle Teilnehmenden verwenden denselben
Grafana-Cloud-Stack mit eigenem Grafana-Login und Browser. Kubernetes-Zugang, Ingest-Token und
Synthkit-Control-Passwort bleiben bei der Kursleitung. Der Shop ist simuliert: Es gibt keine
Shop-Webseite, deren Aufruf Requests erzeugt. Synthkit liefert die Telemetrie kontinuierlich.

## Das nehmen Sie mit

- Den Service, das Dashboard und den Zeitraum finden, die zählen.
- Eine Frage über Metriken und einen zweiten Signaltyp untersuchen; Traces und Profiles gezielt einordnen.
- Abfragen und Erkenntnisse validieren, bevor Sie handeln.
- Was funktioniert wiederverwenden: zunächst als persönlichen Assistant-Prompt, später mit geprüften Team-Vorgaben.

## Kursleitung: Vorbereitung vor den 90 Minuten

Das Blueprint `grafana-cloud-workshop` liefert die Daten. Helm erstellt **keine** Grafana-Dashboards,
Explore-Links, Assistant-Quickstarts, Benutzer, Teams, Alarmregeln oder DOMM-Fixtures und aktiviert
Assistant nicht. Die folgenden Stack-Artefakte müssen vor der Session vorbereitet und geprüft sein.

1. [Kubernetes-Deployment](kubernetes.md) mit allen vier Ingest-Signalen starten. Mindestens
   **30 Minuten Vorlauf** einplanen, bei der ersten Einrichtung länger. Die Abfragen weiter unten
   selbst ausführen und frische Daten in Metriken, Logs, Traces und Profiles prüfen; ein grüner Pod
   reicht nicht. Für den optionalen Vergleich mit gestern muss auch das gestrige Fenster Daten enthalten.
2. Mit der tatsächlichen Teilnehmerrolle Dashboard, Explore und alle vier Datenquellen testen.
   Assistant-Zugang und das Anlegen eines persönlichen Quickstarts separat prüfen. Falls nicht
   verfügbar, den Karten-Fallback für Lab 3 ankündigen; er ersetzt nicht den Nachweis eines gespeicherten Prompts.
3. Über den nur für die Kursleitung erreichbaren [Control-Plane-Zugang](control-plane.md) den
   [Übungszustand zurücksetzen](kubernetes.md#reset-the-exercise-not-the-telemetry): nur Workshop-Blueprint,
   Volumenmultiplikator `1`, keine aktiven Szenarien oder Ad-hoc-Fehler. Historische Telemetrie bleibt erhalten.
4. Ein Szenario-Dashboard anlegen, beispielsweise **Foundation – Checkout beobachten**. Mindestens
   ein Zeitreihen-Panel mit der Metrikabfrage aus Lab 2, Einheit Sekunden, Service `shop-checkout`.
   Das Panel als „Checkout: mittlere modellierte HTTP-Dauer“ benennen und als Einstieg hervorheben.
   Das Signal zeigen, aber keine Ursache vorwegnehmen. Ein vorhandenes Dashboard ist ebenfalls
   geeignet, wenn es exakt diese synthetischen Daten abfragt. Kein Dashboard-JSON wird hier mitgeliefert.
5. Gesunden Verlauf sammeln, dann `grafana-cloud-workshop/checkout-regression` in der
   Scenarios-Ansicht aktivieren. Aktivierungszeit notieren, etwa **acht Minuten** laufen lassen,
   anschließend ausdrücklich deaktivieren. Das Szenario läuft nicht automatisch ab. Es kombiniert
   Checkout-Latenz, Fehler und CPU-Hotspot; es verändert synthetische Daten, nicht echte Shop-Systeme
   oder die CPU-Auslastung des Kubernetes-Hosts, und erzeugt kein Grafana-Incident-Objekt.
6. Nach Ingest-Verzögerung die Veränderung tatsächlich prüfen. Ein **absolutes Fenster mit Datum
   und Zeitzone** wählen, das gesunden Verlauf und Veränderung enthält. Beispiel für den Ablauf:
   zehn Minuten gesund, acht Minuten Szenario, anschließend mindestens fünf Minuten Erholung.
   Der fünfminütige `rate()`-Bereich glättet Übergänge. Das Fenster in Dashboard und Explore speichern
   bzw. teilen; vor Lab 1 muss das Signal bereits sichtbar sein. Nicht erst während Lab 2 aktivieren.
7. Dashboard-Link und zwei vorbereitete Explore-Links aus der Grafana-Oberfläche kopieren:
   Metriken sowie Logs (oder Traces), jeweils derselbe Service und absolute Zeitraum. Mit einem
   Teilnehmerkonto öffnen und prüfen, dass Datenquelle, Filter und Zeitfenster erhalten bleiben.
   Keine erfundenen Data-Source-UIDs oder Beispiel-URLs an Teilnehmende verteilen.

### Ausgefüllten Session-Zettel verteilen

Die Kursleitung ersetzt die offenen Felder vor Beginn. `workshop-shop` ist der **synthetische**
Service-Namespace, nicht der Namespace des tatsächlichen Synthkit-Pods.

| Feld auf der Lab-Karte | Wert für diese Session |
|---|---|
| `{service}` | `shop-checkout` |
| `{namespace}` | `workshop-shop` |
| `{window}` | Tatsächlich geprüfte Start- und Endzeit mit Datum und Zeitzone eintragen |
| `{dashboard}` | Link zum geprüften Checkout-Dashboard eintragen |
| Datenquellen | Tatsächliche Namen von Prometheus, Loki, Tempo und Pyroscope eintragen |
| Explore-Fallback | Geprüfte Metrik- und Log-/Trace-Links für dasselbe Fenster eintragen |
| Assistant | Verfügbar und persönliche Quickstarts getestet, oder Karten-Fallback |

**Startkontrolle:** Ohne sichtbares Signal und geprüfte Dashboard-/Explore-Links ist Lab 1/2
nicht bereit. Bei einer Ingest-Störung ein vorher geprüftes, noch gespeichertes Fenster verwenden
und dessen Datum offen nennen. Nur eine Person verändert den Generator; Teilnehmende verändern
weder gemeinsame Dashboards noch Datenquellen oder Control-Zustand.

## Ablauf: 90 Minuten

| Minuten | Abschnitt | Ergebnis |
|---|---|---|
| 0–10 | Orientierung, Stack und Session-Zettel | Service, Dashboard, Datenquellen und Fenster gefunden |
| 10–25 | Lab 1, Folie 9: Ihr Signal finden | Eine untersuchbare Frage, noch keine Diagnose |
| 25–30 | Nachbesprechung | Auswirkung in Geschäftssprache, Unsicherheit benannt |
| 30–55 | Lab 2, Folie 14: Untersuchen und validieren | Zwei ausgeführte Abfragen über zwei Signaltypen |
| 55–65 | Evidenz besprechen, Trace-/Profile-Demo | Aussage und Grenze zusätzlicher Signale erklären |
| 65–85 | Lab 3, Folie 27: Persönlichen Quickstart erstellen | Persönlich speichern, ausführen, prüfen und verfeinern |
| 85–90 | Abschluss und Wiederverwendung | Nächster Einsatz und Grenzen benannt |

## Lab 1 · Deck-Folie 9 · Ihr Signal finden

### Szenario

Sie haben `shop-checkout` übernommen. Während des auf dem Session-Zettel genannten `{window}`
hat sich etwas verändert. Sie beheben es noch nicht und diagnostizieren es noch nicht: Sie
wählen eine Frage, der es sich zu folgen lohnt.

### Schritte

1. Öffnen Sie ein bekanntes Service-Dashboard, sofern es die Workshop-Daten zeigt; sonst `{dashboard}`.
2. Stellen Sie `{window}` ein. Nutzen Sie für die gemeinsame Untersuchung den absoluten Zeitraum
   vom Session-Zettel, nicht ein wanderndes „Last 15 minutes“.
3. Benennen Sie ein Signal und eine Frage. Halten Sie Service und Zeitraum für Lab 2 fest.

### Erwartetes Ergebnis und Erfolgskontrolle

Ein Satz mit Service, Zeitraum und Signal, beispielsweise:
„Bei `shop-checkout` steigt im markierten Zeitraum die mittlere modellierte HTTP-Dauer;
treten im selben Zeitraum auch Fehler auf?“ Ergänzen Sie die tatsächlichen Zeiten Ihrer Ansicht.
Sie können **Service, Zeitraum und Frage** benennen. Eine Ursache ist noch nicht festgestellt.

### In der Nachbesprechung

Formulieren Sie die mögliche geschäftliche Bedeutung: „Der Checkout war im beobachteten
Zeitraum langsamer; ob Kaufabschlüsse betroffen waren, müssen wir noch prüfen.“
Die Aussage „ein Fünftel der Kundschaft war zwanzig Minuten betroffen“ wäre ohne weitere
Belege unzulässig. Dieses Blueprint liefert weder echte Kundenzahlen noch Umsatzausfälle.
Beobachtete Dauer, mögliche Auswirkung und unbekannter Umfang bleiben getrennt.

### Wenn Sie nicht weiterkommen

Öffnen Sie das vorbereitete Szenario-Dashboard mit dem hervorgehobenen Latenz-Panel.
Beschreiben Sie zunächst nur, **was** sich **wann** verändert – nicht warum.

## Lab 2 · Deck-Folie 14 · Untersuchen und validieren

### Szenario und Schritte

Nehmen Sie Ihre Frage aus Lab 1 mit nach Explore. Halten Sie `shop-checkout` und `{window}`
durchgehend fest. Nicht gleichzeitig Service und Zeitraum wechseln.

1. Öffnen Sie das relevante Dashboard-Panel in Explore oder den vorbereiteten Metrik-Link.
2. Prüfen oder verfeinern Sie die eingegrenzte Abfrage. Erklären Sie, was sie tatsächlich misst.
3. **Wechseln Sie zwingend zu einem zweiten Signaltyp:** von Metriken zu Logs oder Traces.
   Query Builder und Assistant dürfen beim Entwurf helfen; erst die ausgeführte Abfrage samt
   Ergebnis ist ein Beleg. Vergleichen Sie denselben Service im selben Zeitfenster.

| Checkliste Untersuchung | Für diese Session festhalten |
|---|---|
| Service | `shop-checkout`, Namespace `workshop-shop` |
| Zeitraum | `{window}` einschließlich Datum und Zeitzone |
| Abfrage | Datenquelle, Filter, Messgröße und Einheit |
| Verlauf | Die zwei nützlichen Abfragen und ihre Ergebnisse bzw. Links behalten |

### Erste Abfrage: Metrik

In der Prometheus-Datenquelle ausführen:

```promql
sum by (service) (
  rate(http_server_request_duration_seconds_sum{blueprint="grafana-cloud-workshop",service="shop-checkout",namespace="workshop-shop"}[5m])
)
/
sum by (service) (
  rate(http_server_request_duration_seconds_count{blueprint="grafana-cloud-workshop",service="shop-checkout",namespace="workshop-shop"}[5m])
)
```

Das Verhältnis aus Dauer und Beobachtungszahl ergibt die **mittlere modellierte HTTP-Dauer
in Sekunden**, nicht p95, Fehlerrate oder die Dauer jedes einzelnen Requests. `[5m]` ist das
Berechnungsfenster an jedem Graph-Punkt, nicht der gesamte Explore-Zeitraum. Nach Neustart
braucht `rate()` genügend Samples. Die DSL-Histogrammzählung ist **kein echter Nutzer- oder
Request-Durchsatz**; daraus keine Kundenzahl oder Fehlerquote berechnen. Fehlende Daten sind nicht null.

### Zweite Abfrage: Logs

Zur Loki-Datenquelle wechseln, `{window}` unverändert lassen und ausführen:

```logql
{blueprint="grafana-cloud-workshop",source="app",service_name="shop-checkout",namespace="workshop-shop"} | json
```

Prüfen Sie `msg`, `route`, `status` und `outcome` im JSON sowie das Stream-Label `level`.
Für die Frage nach Fehlern anschließend eingrenzen:

```logql
{blueprint="grafana-cloud-workshop",source="app",service_name="shop-checkout",namespace="workshop-shop",level="error"} | json
```

Die Metrik zeigt den aggregierten Verlauf; Logs zeigen einzelne Ereignisse und deren Details.
Fehler im gleichen Fenster können eine Hypothese stützen, beweisen aber keine Ursache der Latenz.
Ein leeres Fehlerergebnis widerlegt die Hypothese erst dann sinnvoll, wenn die allgemeine
Logabfrage Daten liefert und Datenquelle, Service und Zeitfenster stimmen.

### Alternative oder Vertiefung: Traces

In Tempo bei unverändertem `{window}`:

```traceql
{ resource.service.name = "shop-checkout" && resource.service.namespace = "workshop-shop" }
```

Einen Treffer öffnen und Checkout-Spans im Request-Pfad
`shop-storefront → shop-checkout → shop-payment` betrachten. Für Fehlerspans zusätzlich
`&& status = error` innerhalb der Klammern ergänzen. Abhängigkeit und zeitliches Zusammentreffen
sind keine Ursachennachweise. Synthetische Metriken und Spans müssen nicht exakt numerisch übereinstimmen.

Logs enthalten `trace_id` und `span_id` als **strukturierte Metadaten**, nicht im JSON-Body
oder als indexierte Stream-Labels. Ohne vorkonfigurierte Links: ID aus dem Log kopieren und
in Tempos Trace-ID-Modus suchen. Zurück zu Loki, im gleichen Zeitraum:

```logql
{blueprint="grafana-cloud-workshop",source="app",service_name="shop-checkout",namespace="workshop-shop"} | trace_id="REPLACE_WITH_TRACE_ID"
```

`|= "TRACE_ID"` durchsucht den Body und ist hier kein Ersatz. Eine Trace-ID allein garantiert
nicht, dass der Trace tatsächlich ingestiert wurde.

### Erwartetes Ergebnis und Erfolgskontrolle

Zeigen Sie **zwei ausgeführte Abfragen über zwei Signaltypen**, beide für denselben Service
und dasselbe Fenster. Notieren Sie zu jeder Abfrage ihre Aussage und Grenze, zum Beispiel:
„Die Metrik zeigt erhöhte mittlere Dauer; die Logs zeigen im selben Fenster Fehlerereignisse.
Das stützt gleichzeitige Beeinträchtigung, erklärt aber noch nicht deren Ursache.“
Ein Widerspruch ist ein gutes Ergebnis: Sie haben Ihre Idee geprüft statt nur bestätigt.

### Wenn Sie nicht weiterkommen

Nutzen Sie die vorbereiteten Explore-Links und vergleichen Sie die beiden benannten Signale.
Lassen Sie Assistant die Abfrage erklären oder entwerfen, kontrollieren Sie aber Filter,
Zeitraum, Einheit und Ergebnis. Der Wechsel zum zweiten Signal ist kein optionaler Zusatz.

## 55–65: Zusätzliche Evidenz mit Traces und Profiles

Die Kursleitung zeigt kurz den Trace-Pfad und öffnet Profiles Drilldown oder Pyroscope in Explore.
Profiltyp `process_cpu:cpu:nanoseconds:cpu:nanoseconds`, gleiches `{window}`, Selektor:

```text
{service_name="shop-checkout",service_namespace="workshop-shop"}
```

Welche Funktionen tragen zur CPU-Zeit bei? Die Breite einer Flamegraph-Fläche steht für
CPU-Anteil, nicht für einen Request-Zeitstrahl. Ein CPU-Hotspot allein erklärt nicht jede
Latenz. Ohne getesteten Trace-to-Profile-Link vergleichen wir Service und Zeitraum, nicht
das Profil eines bestimmten Spans. Die Go-Span-Profile dieser Übung sind CPU-only.

Die Demo ergänzt Lab 2; sie ersetzt die zwei selbst ausgeführten Abfragen nicht.

## Lab 3 · Deck-Folie 27 · Einen persönlichen Quickstart erstellen

### Szenario und Schritte

Sie werden diese Untersuchung wiederholen. Speichern Sie einen Prompt, mit dem der nächste
Durchlauf bereits Kontext hat. Halten Sie ihn zunächst persönlich.

1. Wählen Sie eine wiederkehrende Aufgabe für `shop-checkout`, etwa einen morgendlichen
   Health-Check. Für die erste Validierung nutzen Sie das bekannte `{window}` aus Lab 2.
2. Schreiben Sie **Aktion + Service-Geltungsbereich + Zeitraum + @Kontext**. Wählen Sie nach
   Eingabe von `@` die tatsächlich vorhandene Datenquelle oder das Dashboard aus der Auswahl.
   `@workshop-shop-logs` ist nur gültig, wenn dieses Objekt wirklich existiert; der Namespace
   selbst erzeugt keine Datenquelle. Die Kontext-Platzhalter unten vor dem Speichern ersetzen.
3. Unter **Grafana Assistant → Settings → Quickstart prompts → Create Quickstart Prompt**
   Titel und Prompt eintragen, **Scope: Just me**, Enabled eingeschaltet lassen und speichern.
4. Den gespeicherten Quickstart ausführen. Mindestens eine erzeugte Abfrage und ihr Ergebnis
   prüfen: Service, Fenster, Datenquelle, Messgröße und tatsächlicher Beleg. Mit Lab 2 vergleichen.
   Den Prompt bei Bedarf präzisieren und erneut ausführen. Nicht auf **Everybody** umstellen.

Das persönliche Speichern und die Kontextauswahl folgen der
[offiziellen Assistant-Anleitung](https://grafana.com/docs/grafana-cloud/platform/grafana-assistant/introduction/).
Die Kursleitung prüft [Aktivierung und Zugriff](https://grafana.com/docs/grafana-cloud/platform/grafana-assistant/get-started/grafana-cloud/)
vorab; UI-Verfügbarkeit und Berechtigungen können je Stack abweichen.

### Zwei ausgearbeitete Prompts

**Untersuchung wiederholen – zunächst mit dem absoluten Lab-Fenster:**

> Prüfe shop-checkout im Namespace workshop-shop während {window} mit @PROMETHEUS und @LOKI.
> Zeige die mittlere modellierte HTTP-Dauer und passende Fehlerlogs. Verwende dieselben
> Service- und Zeitfilter. Zeige beide Abfragen und Ergebnisse; trenne Beobachtungen,
> Hypothesen und offene Fragen. Nimm keine Änderungen vor.

`@PROMETHEUS` und `@LOKI` sind Platzhalter für echte ausgewählte Kontextobjekte.
Nach erfolgreicher Prüfung kann eine persönliche Health-Check-Variante „letzte 30 Minuten“
verwenden. Sie zeigt nach dem Deaktivieren des Szenarios möglicherweise gesunde Daten –
das ist kein Fehlverhalten. Ein Quickstart ist kein automatisch laufender Monitor.

**Vergleichen – nur mit vorhandener Historie:**

> Vergleiche die mittlere modellierte HTTP-Dauer von shop-checkout im Namespace workshop-shop
> in der letzten abgeschlossenen Stunde mit derselben Stunde gestern anhand von @PROMETHEUS.
> Zeige die genauen Zeitfenster, Abfragen und Unterschiede. Weise auf fehlende Daten hin,
> statt sie als null zu behandeln. Behaupte keine Ursache ohne weiteren Beleg.

Diesen Vergleich nur verwenden, wenn beide Fenster tatsächlich Daten enthalten. Andernfalls
mit zwei geprüften gesunden/degradierten Zeitfenstern derselben Session vergleichen.
Eine „Fehlerrate“ nicht aus der DSL-Histogrammzählung ableiten; für den Einstieg Fehlerereignisse
untersuchen. Eine Quote aus Logs wäre höchstens eine Quote erzeugter Logereignisse, nicht Kundschaft.

### Erwartetes Ergebnis und Erfolgskontrolle

Ein **persönlich gespeicherter und ausgeführter Prompt**, mindestens ein geprüftes Ergebnis
und eine benannte Verfeinerung oder begründete Entscheidung, ihn so zu behalten.
Sie können erklären, wann Sie ihn wieder verwenden und welche Daten er braucht.
Teilen Sie ihn erst mit dem Team, wenn er zuverlässig funktioniert und jemand ihn verantwortet.

### Wenn Sie nicht weiterkommen

Entwerfen Sie den Prompt auf der Lab-Karte und gehen Sie Auswahl des Kontexts, erwartete
Abfrage und Ergebnisprüfung laut durch. Bei fehlendem Assistant-Zugang oder Speicherrecht
ist das ein bewusster Fallback, **kein** erfolgreich gespeicherter Quickstart. Zugang nicht
während des Labs durch neue gemeinsame Berechtigungen improvisieren.

## Abschluss und nächste Durchführung

Welche Frage haben Sie geprüft? Was konnte das zweite Signal ergänzen oder widerlegen?
Welche Aussage blieb unbelegt? Welchen persönlichen Quickstart würden Sie wiederverwenden?

Die Kursleitung bestätigt, dass das Szenario deaktiviert ist, und setzt den
[Übungszustand](kubernetes.md#reset-the-exercise-not-the-telemetry) zurück. Der Generator darf
weiterlaufen. Alte Telemetrie muss nicht gelöscht werden und ist für die nächste Durchführung
kein Hindernis. Für eine neue Klasse frische Fenster vorbereiten und alle Links neu prüfen;
ein Pod-Neustart allein ersetzt keinen Reset des persistenten Control-Zustands.

## Referenz: absichtlich ungleichmäßige Instrumentierung

Blueprint `grafana-cloud-workshop`, synthetischer Cluster `workshop-prod`, Namespace `workshop-shop`:

| Service | Metriken | App-Logs | Traces | Profiles |
|---|---|---|---|---|
| `shop-storefront` | Ja | Ja | Ja | Ja |
| `shop-checkout` | Ja | Ja | Ja | Ja |
| `shop-payment` | Ja | Ja | Ja | Ja |
| `shop-inventory` | Ja | Nein | Nein | Nein |
| `shop-shipping` | Ja | Ja | Nein | Nein |

Optionaler Transfer nach den Labs: Was könnten Sie über Inventory mit Metriken allein sagen?
Fehlende Anwendungssignale sind hier beabsichtigt. Infrastrukturtelemetrie beweist keine
App-Instrumentierung; eine Log-ID beweist keinen verfügbaren Trace.

### Optionale Verknüpfungen und Fehlerhilfe für die Kursleitung

Gemeinsame IDs konfigurieren keine Grafana-Links. Vor der Session testen oder die manuelle
Trace-ID-Suche aus Lab 2 verwenden:

- Loki → Tempo: Derived Field vom Typ Label, Matcher `^trace_id$`, interne Tempo-Verknüpfung,
  Query `${__value.raw}`. Die ID kommt aus strukturierten Metadaten.
- Tempo → Loki: `service.name` auf `service_name`, `service.namespace` auf `namespace` abbilden;
  Custom Query `{${__tags}} | trace_id="${__trace.traceId}"`. Kleinen Zeitpuffer erlauben.
  Kein zwingender `span_id`-Match: App-Logs tragen die Root-Span-ID des Requests.
- Tempo → Pyroscope: `service.name` auf `service_name`, `service.namespace` auf
  `service_namespace` abbilden, CPU-Profil auswählen und an einem Go-Server-Span prüfen.

Siehe [Trace-to-Logs](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/configure-trace-to-logs/)
und [Trace-to-Profiles](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/configure-trace-to-profiles/).
Provisionierte Datenquellen können schreibgeschützt sein; unterstützten Provisionierungs-/Klonweg
verwenden und nicht während der Session Produktionskonfiguration ändern.

| Beobachtung | Zuerst prüfen |
|---|---|
| Inventory ohne Logs oder Shipping ohne Traces | Erwartete Grenzen laut Tabelle |
| Alle Services ohne einen Signaltyp | Datenquelle und Fenster, dann Ingest-Status und Zugangsdaten |
| Alles leer | Dry-run, Blueprint-Auswahl, pausierter Pod und Readiness |
| Logs/Traces vorhanden, Links defekt | Manuelle ID-Suche, dann Datenquellen-Mappings |
| Nur Profiles fehlen | Profil-Zugangsdaten, Service/Typ und Zustellstatus |
| Alte Fehler nach Erholung sichtbar | Historisches Fenster zeigt weiterhin historische Fehler |

Quelle der Datenformen ist das
[`grafana-cloud-workshop`-Blueprint](https://github.com/mbaykara/Synthkit/blob/main/blueprints/grafana-cloud-workshop.yaml).
Weitere Diagnose: [Control Plane](control-plane.md) und [Troubleshooting](troubleshooting.md).
Teilnehmende melden Probleme der Kursleitung, statt Tokens, Helm-Werte oder gemeinsame Ressourcen zu ändern.
