---
title: Grafana Alerting in Action
description: Grafana Alerting in Action – drei Labs auf einem gemeinsamen Stack in 60 Minuten.
---

# Grafana Alerting in Action

[English version](workshop-alerting-en.md) · Weitere Sessions: [Foundation](workshop.md) ·
[Assistant](workshop-assistant-de.md)

**Ein gemeinsamer Stack, 60 Minuten, drei Labs:** eine Alarmregel aus einer erklärbaren Abfrage
erstellen, sie mit Labels an genau ein Team routen und Benachrichtigungsrauschen reduzieren. Die
Labs gehören zu den Deck-Folien **8, 12 und 16**. Die Zeitaufteilung unten ist ein
Moderationsvorschlag; die Foliennummern stammen aus den Lab-Karten.

Die Kursleitung betreibt Synthkit auf Kubernetes. Alle Teilnehmenden verwenden denselben
Grafana-Cloud-Stack mit eigenem Grafana-Login. Kubernetes-Zugang, Ingest-Token und
Synthkit-Control-Passwort bleiben bei der Kursleitung. Diese Session steht für sich: Wer
Foundation verpasst hat, braucht nichts daraus.

Die Abfragen der Lab-Karten verwenden Micrometer-/Spring-Metriknamen. Auf diesem Stack liefert sie
der synthetische JVM-Service `shop-catalog`. Die Werte sind modelliert, nicht aus einer echten
Anwendung gemessen.

## Das nehmen Sie mit

- Eine Grafana-verwaltete Alarmregel mit begründbarem Schwellwert und Pending-Zeitraum erstellen.
- Den Lebenszyklus einer Regel erklären, einschließlich der Zustände No data und Error.
- Mit wenigen Labels routen, sodass eine Benachrichtigung genau einen Contact Point erreicht.
- Rauschen mit Gruppierung, Silences und Mute Timings reduzieren.

## Kursleitung: Vorbereitung vor den 60 Minuten

Das Blueprint `grafana-cloud-workshop` liefert die Daten. Synthkit und Helm erstellen **keine**
Ordner, Alarmregeln, Contact Points, Notification Policies, Silences, Mute Timings, Benutzer oder
Teams. Diese Stack-Artefakte müssen vor der Session in Grafana angelegt und geprüft sein.

1. **Mindestens sieben Tage Vorlauf.** Synthkit kann keine Historie nachträglich erzeugen. Lab 1
   verlangt einen Blick auf die letzten sieben Tage; der Generator muss also mindestens sieben
   Tage vor der Session mit dem [Kubernetes-Deployment](kubernetes.md) laufen. Die Lab-1-Abfrage
   selbst über `Last 7 days` ausführen und lückenlose Daten prüfen.
2. Über den nur für die Kursleitung erreichbaren [Control-Plane-Zugang](control-plane.md) den
   [Übungszustand zurücksetzen](kubernetes.md#reset-the-exercise-not-the-telemetry): nur
   Workshop-Blueprint, Volumenmultiplikator `1`, keine aktiven Szenarien. Diese Session braucht
   kein Synthkit-Szenario.
3. `{folder}` anlegen und den Teilnehmenden Editor-Rechte darauf geben.
4. Den Contact Point `workshop-webhook` auf `{webhook-url}` anlegen und testen. Das Ziel muss auf
   dem geteilten Bildschirm sichtbar sein, damit alle das Eintreffen einer Nachricht sehen.
   Zusätzlich `workshop-email` (Catch-all) und `workshop-servicenow` anlegen.
5. Den vorbereiteten Policy-Baum aus Lab 2 anlegen und stehen lassen.
6. Die verrauschte Regel aus Lab 3 anlegen und **pausiert** lassen. Sie wird im Lab auf ein
   Stichwort aktiviert.
7. Das Lab-Karten-Paket mit der Einladung verschicken.

### Ausgefüllten Session-Zettel verteilen

| Feld auf der Lab-Karte | Wert für diese Session |
|---|---|
| `{namespace}` | `workshop-shop` |
| `{service}` | `shop-catalog` (JVM-Service mit 12 Pods, synthetischer Cluster `workshop-prod`) |
| `{team}` | Ein Routing-Wert, den die Kursleitung festlegt, zum Beispiel `workshop-a` |
| `{folder}` | Tatsächlichen Regelordner mit Schreibrechten für Teilnehmende eintragen |
| `{webhook-url}` | Tatsächliches, auf dem Bildschirm sichtbares Webhook-Ziel eintragen |
| `{runbook-link}` | Beliebige URL; ein Platzhalter-Runbook genügt |

**Startkontrolle:** Ohne sieben Tage Daten für die Lab-1-Abfrage, einen getesteten
`workshop-webhook` und die pausierte Lab-3-Regel ist die Session nicht bereit. Teilnehmende ändern
weder den Policy-Baum der anderen noch Contact Points oder den Control-Zustand.

## Ablauf: 60 Minuten

| Minuten | Abschnitt | Ergebnis |
|---|---|---|
| 0–8 | Eine Regel ist eine validierte Abfrage mit Entscheidung; Schwellwert, Error Budget, adaptiv; Aufbau und Lebenszyklus | Begriffe für Lab 1 geklärt |
| 8–22 | Lab 1, Folie 8: Alarmregel erstellen | Regel wertet aus, Schwellwert begründet |
| 22–27 | Labelstrategie, von der Regel zur Benachrichtigung | Vier bis fünf Labels, Policy-Baum verstanden |
| 27–39 | Lab 2, Folie 12: Labeln und routen | Genau eine Benachrichtigung am richtigen Ziel |
| 39–44 | Rauschen, Eskalation, Betrieb im großen Maßstab | Gruppierung, Silences, Alerts as Code eingeordnet |
| 44–56 | Lab 3, Folie 16: Das Rauschen dämpfen | Eine Benachrichtigung statt zwölf |
| 56–60 | Abschluss und nächste Schritte | Was um 3 Uhr nachts pagt, ist benannt |

## Lab 1 · Deck-Folie 8 · Alarmregel erstellen

### Szenario

Alarmieren Sie, wenn ein Service in `workshop-shop` auf der CPU heiß läuft. Sehen Sie sich zuerst
die letzten sieben Tage an und wählen Sie dann einen Schwellwert, den Sie einer Kollegin oder
einem Kollegen gegenüber vertreten können.

### Startabfrage

Ihre eigene Abfrage, ein Panel, das Sie erklären können, oder diese:

```promql
avg by (service_name) (process_cpu_usage{namespace="workshop-shop"})
```

Kubernetes-Alternative (cAdvisor, alle Shop-Pods):

```promql
sum by (pod) (rate(container_cpu_usage_seconds_total{namespace="workshop-shop"}[5m]))
```

### Schritte

1. Startabfrage wählen.
2. Eine Bedingung aus dem Verlauf ableiten und einen Pending-Zeitraum, der normales Rauschen
   überdauert.
3. Eine Summary-Annotation mit Runbook-Link sowie die Labels `severity` und `team` ergänzen.
4. Zusatz: Assistant dieselbe Regel entwerfen lassen und dessen Schwellwert mit Ihrem vergleichen.

### Vorgeschlagene Einstellungen – am Verlauf anpassen

- Bedingung: zunächst über `0.8`, dann begründen oder verschieben.
- Auswertung jede Minute (`1m`), Pending-Zeitraum `5m`.
- No data und Error Handling: bewusst wählen, Standardwerte nicht ungelesen übernehmen.

### Labels und Annotationen

```text
labels:      team={team} severity=warning
             service_name=shop-catalog namespace=workshop-shop
annotations: summary     = CPU high on shop-catalog - check recent deploys and load
             runbook_url = {runbook-link}
```

### Erwartetes Ergebnis und Erfolgskontrolle

Eine Regel im Zustand **Normal**, die planmäßig auswertet, und ein Satz, der Schwellwert und
Pending-Zeitraum begründet. Wer den Zusatz gemacht hat, sieht einen vernünftigen, aber generischen
Assistant-Schwellwert: Assistant kennt Ihre Incident-Historie nicht. Dieser Kontrast ist der Punkt.
Sie können **genau erklären, wann die Regel feuert und warum**.

### Für die Kursleitung

`process_cpu_usage` liegt pro Pod stabil zwischen 0,30 und 0,72; der Durchschnitt über die zwölf
Pods liegt nahe 0,5. Der Startwert `0.8` feuert deshalb nie, und der Zustand Normal ist beabsichtigt.
Die Teilnehmenden müssen den Schwellwert aus dem Verlauf begründen. Die Werte hängen an keinem
Szenario.

### Wenn Sie nicht weiterkommen

Verwenden Sie Abfrage, Schwellwert und Labels genau wie oben abgedruckt; sie funktionieren auf dem
gemeinsamen Stack.

## Lab 2 · Deck-Folie 12 · Labeln und routen

### Vorbereiteter Policy-Baum – bereits auf dem Stack

```text
Default policy -> workshop-email (catch-all)
|- team = {team}
|    -> workshop-webhook
|    group_by: [service_name]
|
|- team = {team}, severity = critical
     -> workshop-servicenow
```

### Schritte

1. Prüfen Sie, dass Ihre Regel die vereinbarten Labels trägt: mindestens `team` und `severity`.
2. Bauen Sie den Policy-Zweig: zuerst auf `team` matchen, darin auf `severity`.
   **Alerting > Notification policies > new child policy**.
3. Senden Sie eine Testbenachrichtigung an den gemeinsamen Contact Point und prüfen Sie, dass sie
   genau einmal ankommt: **Test** an `workshop-webhook` drücken und `{webhook-url}` beobachten.

### Erwartetes Ergebnis und Erfolgskontrolle

Genau eine Benachrichtigung erreicht `{webhook-url}` mit Ihren Labels, und nichts erreicht die
Catch-all-E-Mail. Kommt sie doppelt an, matchen zwei Zweige: `continue matching` am Elternknoten
prüfen. Landet sie nur im Catch-all, stimmen Labels und Matcher nicht überein.
Die Benachrichtigung **erreicht den vorgesehenen Contact Point und keinen anderen**.

### Hinweis zu Contact Points

Der kritische Zweig zeigt auf einen ServiceNow-Contact-Point, weil viele Teams ihre Incidents
dorthin senden. Grafana Alerting kann alles benachrichtigen, was einen Webhook annimmt; derselbe
Baum funktioniert mit Teams, E-Mail, Grafana IRM oder eigenen Werkzeugen. Entscheidend ist die
Routing-Entscheidung, nicht das Ziel.

### Wenn Sie nicht weiterkommen

Verwenden Sie den vorbereiteten Policy-Baum und den gemeinsamen Webhook-Contact-Point, statt einen
eigenen Baum zu entwirren.

## Lab 3 · Deck-Folie 16 · Das Rauschen dämpfen

### Szenario

Die Kursleitung aktiviert eine vorbereitete, verrauschte Regel. Sie wertet pro Pod aus; ein
einziges Speicherproblem erzeugt deshalb zwölf Alarminstanzen. Ungruppiert sind das zwölf
Benachrichtigungen für eine Geschichte.

### Die verrauschte Regel – bereits angelegt, auf Stichwort aktiviert

```promql
(jvm_memory_used_bytes{namespace="workshop-shop", area="heap"}
  / jvm_memory_max_bytes{namespace="workshop-shop", area="heap"})
  > 0.7
```

### Schritte

1. Gruppierung am Policy-Zweig setzen, damit zusammengehörige Instanzen als eine Benachrichtigung
   ankommen: `group_by [team, service_name]`, `group_wait 30s`, `group_interval 5m`,
   `repeat_interval 4h`.
2. Eine Silence oder ein Mute Timing für eine geplante Änderung ergänzen: Matcher
   `service_name = shop-catalog`, Dauer `2h`, Kommentar `planned change`.
3. Den Alarm Ende zu Ende verfolgen: Regel, Policy-Zweig, Contact Point und weiter zu dem, was die
   Eskalation übernimmt.

### Erwartetes Ergebnis und Erfolgskontrolle

Eine Benachrichtigung, die alle betroffenen Pods auflistet, statt einer pro Pod. Notieren Sie die
Anzahl vorher und nachher; sie wird in der Nachbesprechung abgefragt. Nach der Silence feuert die
Regel weiterhin und bleibt in der Oberfläche sichtbar, aber es wird keine Benachrichtigung
gesendet. Diese Unterscheidung ist wichtig. **Eine Benachrichtigung erzählt die ganze Geschichte,
und Sie können sagen, was um 3 Uhr nachts pagt.**

### Für die Kursleitung

Das „Stichwort“ ist das Aktivieren der pausierten Regel in Grafana, kein Synthkit-Szenario.
`shop-catalog` meldet auf allen zwölf Pods dauerhaft 74 bis 88 Prozent von 512 MiB Heap; die Regel
feuert daher sofort zwölf Instanzen. Alle zwölf teilen `service_name="shop-catalog"`, deshalb
fasst `group_by [team, service_name]` sie zu einer Benachrichtigung zusammen. Erwartet: vorher 12,
nachher 1.

### Wenn Sie nicht weiterkommen

Nutzen Sie die vorbereitete laute Regel, sobald die Kursleitung sie aktiviert hat, und lösen Sie das Problem allein mit Gruppierung. Die Gruppierung
bringt den sichtbaren Gewinn.

## Abschluss und nächste Durchführung

Welche Regel würden Sie morgen behalten? Besteht sie die vier Tests: umsetzbar, symptombasiert,
verantwortet, begründet? Welche Labels braucht Ihr Team für Routing, Triage und Auswertung?

Die Kursleitung pausiert die Lab-3-Regel wieder, entfernt Silences der Teilnehmenden und löscht
bei Bedarf deren Regeln aus `{folder}`. Der Generator darf weiterlaufen; die Historie bleibt für die
nächste Klasse erhalten.

## Referenz: Alerting-Daten dieses Stacks

Blueprint `grafana-cloud-workshop`, synthetischer Cluster `workshop-prod`, Namespace `workshop-shop`:

| Metrik | Labels | Modell |
|---|---|---|
| `process_cpu_usage` | `service_name`, `namespace`, `pod` (12) | pro Pod 0,30–0,72 |
| `jvm_memory_used_bytes` | `area="heap"`, `service_name`, `namespace`, `pod` (12) | 74–88 % von 512 MiB |
| `jvm_memory_max_bytes` | `area="heap"`, `service_name`, `namespace`, `pod` (12) | konstant 512 MiB |
| `container_cpu_usage_seconds_total` | `namespace`, `pod` (alle Shop-Pods) | cAdvisor |

| Beobachtung | Zuerst prüfen |
|---|---|
| Lab-1-Abfrage liefert weniger als sieben Tage | Laufzeit des Generators; es gibt keine nachträgliche Historie |
| Lab-1-Regel feuert bei `0.8` nicht | Erwartet; der Durchschnitt liegt nahe 0,5 |
| Lab-3-Regel feuert nicht | Regel noch pausiert, Namespace oder `area="heap"` falsch |
| Doppelte Benachrichtigung | Zwei matchende Zweige, `continue matching` am Elternknoten |
| Nur Catch-all erreicht | Labels der Regel und Matcher des Zweigs stimmen nicht überein |
| Alles leer | Dry-run, Blueprint-Auswahl, pausierter Pod und Readiness |

Weitere Diagnose: [Control Plane](control-plane.md) und [Troubleshooting](troubleshooting.md).
Teilnehmende melden Probleme der Kursleitung, statt Tokens, Helm-Werte oder gemeinsame Ressourcen
zu ändern.
