---
title: Grafana Assistant in Action
description: Grafana Assistant in Action – drei Labs auf einem gemeinsamen Stack in 60 Minuten.
---

# Grafana Assistant in Action

[English version](workshop-assistant-en.md) · Weitere Sessions:
[Foundation](workshop.md) ([English](workshop-foundation-en.md)) ·
[Alerting](workshop-alerting-de.md) ([English](workshop-alerting-en.md))

**Ein gemeinsamer Stack, 60 Minuten, drei Labs:** mit Assistant ein Dashboard bauen, ein Symptom
mit belegter Hypothese untersuchen und die eigene Methode als Rule, Skill und Automation
wiederholbar machen. Die Labs gehören zu den Deck-Folien **7, 11 und 17**. Die Zeitaufteilung
unten ist ein Moderationsvorschlag; die Foliennummern stammen aus den Lab-Karten.

Diese Session steht für sich. Wer Foundation verpasst hat, braucht nichts daraus: Die Karten
liefern Service und alle Prompts. Die Kursleitung betreibt Synthkit auf Kubernetes; alle
Teilnehmenden arbeiten im selben Grafana-Cloud-Stack mit eigenem Login. Der Shop ist simuliert.
Assistant entwirft und erkundet, die Ingenieurin oder der Ingenieur validiert: Jedes Lab enthält
einen Schritt, in dem Sie das Ergebnis selbst prüfen.

## Das nehmen Sie mit

- Ein Dashboard aus Prompts bauen, iterieren und die Abfragen dahinter prüfen.
- Eine Untersuchung mit Belegen durchführen und eine Erkenntnis selbst in Explore validieren.
- Die Auswirkung erst nach der Validierung beziffern und ihre Grenzen benennen.
- Eine Assistant-Rule, einen Skill und eine geplante Automation für die eigene Methode anlegen.
- Geteilte KI-Artefakte verantwortlich behandeln: zuerst persönlich, geteilt nur mit Owner.

## Kursleitung: Vorbereitung vor den 60 Minuten

Das Blueprint `grafana-cloud-workshop` liefert die Daten. Helm aktiviert Assistant nicht und legt
keine Dashboards, Rules, Skills, Automations, Benutzer oder Annotationen an. Die folgenden
Punkte spätestens am Vortag prüfen.

1. [Kubernetes-Deployment](kubernetes.md) mit allen vier Ingest-Signalen laufen lassen. Der
   Namespace braucht **24 Stunden** lebende Metriken, Logs und Traces: Den Generator mindestens
   einen Tag vor der Session starten. Synthkit füllt keine Historie nach.
2. Assistant ist für den Stack aktiviert und lizenziert. Einen Prompt selbst vollständig ausführen.
   Mit einem Teilnehmerkonto prüfen, dass **Settings → Rules, Skills und Automations** erreichbar
   sind; Lab 3 braucht alle drei.
3. **Wichtig:** Übrig gebliebene mandantenweite Rules (Scope **Everybody**) früherer Gruppen
   löschen. Sie verändern unbemerkt, wie Assistant sich für die nächste Gruppe verhält.
4. Prüfen, ob **Investigations** verfügbar ist. Falls nicht, den Konversations-Fallback auf der
   Lab-2-Karte ankündigen.
5. Über den nur für die Kursleitung erreichbaren [Control-Plane-Zugang](control-plane.md) den
   [Übungszustand zurücksetzen](kubernetes.md#reset-the-exercise-not-the-telemetry): nur
   Workshop-Blueprint, Volumenmultiplikator `1`, keine aktiven Szenarien. Historische Telemetrie
   bleibt erhalten.
6. Den Lab-2-Fehler vorbereiten: Nach gesundem Verlauf `grafana-cloud-workshop/payment-regression`
   in der Scenarios-Ansicht aktivieren. Die Aktivierungszeit mit **Datum und Zeitzone** als
   `{fault-time}` notieren, etwa **zehn Minuten** laufen lassen, danach ausdrücklich deaktivieren.
   Das Szenario läuft nicht automatisch ab. Es erhöht fehlgeschlagene Payment-Requests und die
   Antwortzeit von `shop-payment`; es verändert synthetische Daten, keine echten Systeme. Nicht
   während Foundation oder Alerting auf demselben Stack aktivieren, außer das ist gewollt.
7. Die Änderung `{change}` vorbereiten. Synthkit erzeugt **kein** Deployment- oder Change-Ereignis.
   Entweder nennt die Kursleitung die Änderung mündlich (etwa „Payment-Release um `{fault-time}`“)
   oder legt vorab eine Grafana-Annotation zu diesem Zeitpunkt an. Ohne das kann Assistant die
   Änderung in der Telemetrie nicht finden; sagen Sie das offen.
8. Nach Ingest-Verzögerung die Veränderung selbst prüfen: erhöhte Fehler in den Logs von
   `shop-payment` und eine höhere mittlere Dauer im selben Fenster.
9. Den Teilnehmenden die ausgefüllten Lab-Karten mit der Einladung schicken.

**Namenswarnung:** In Lab 3 bedeutet „Rule“ eine **Assistant-Verhaltensregel**, keine
Alarmregel. Sagen Sie das ausdrücklich; die Verwechslung ist häufig.

### Ausgefüllten Session-Zettel verteilen

| Feld auf der Lab-Karte | Wert für diese Session |
|---|---|
| `{service.namespace}` | `workshop-shop` |
| `{service}` | `shop-payment` |
| `{fault-time}` | Tatsächliche Aktivierungszeit mit Datum und Zeitzone eintragen |
| `{change}` | Mündlich genannte Änderung oder vorbereitete Annotation, z. B. Payment-Release |
| `{dashboard name - lab 1}` | `Shop - Application Health` plus eigene Initialen |
| Assistant | Aktiviert und getestet; Investigations verfügbar oder Konversations-Fallback |

`workshop-shop` ist der **synthetische** Service-Namespace, nicht der Namespace des tatsächlichen
Synthkit-Pods. Verfügbare Signale für `shop-payment`: Metriken
(`http_server_request_duration_seconds` mit den Labels `service`, `service_name`, `namespace`),
App-Logs mit `status` und `outcome` sowie `trace_id` als strukturierte Metadaten, Traces im Pfad
`shop-storefront → shop-checkout → shop-payment` und CPU-Profiles.

## Ablauf: 60 Minuten

| Minuten | Abschnitt | Ergebnis |
|---|---|---|
| 0–5 | Assistant beschleunigt eine bekannte Schleife; was Assistant ist | Kontextbewusster Agent in Grafana, unter bestehendem RBAC |
| 5–20 | Lab 1, Folie 7: Mit Assistant bauen | Gespeichertes Dashboard, jede Abfrage erklärbar |
| 20–27 | Investigations und Prompt-Handwerk | Symptom, Scope, Zeitraum und Änderung benennen |
| 27–42 | Lab 2, Folie 11: Ein Symptom untersuchen | Belegte Hypothese und eine Aussage zur Auswirkung |
| 42–47 | Rules, Skills, Automations, MCP | Wann was passt; Governance |
| 47–57 | Lab 3, Folie 17: Wiederholbar machen | Rule, Skill und Automation, soweit die Zeit reicht |
| 57–60 | Abschluss | Nächster Einsatz benannt |

## Lab 1 · Deck-Folie 7 · Mit Assistant bauen

### Szenario

Verwenden Sie einen Service, den Sie auf diesem Stack kennen, oder die vorbereiteten Anwendungen
in `workshop-shop`. Behandeln Sie sie wie einen Service, den Sie gerade übernommen haben.

### Schritte

1. Fragen Sie Assistant nach dem Zustand Ihres Service: Fehler, Latenz und Sättigung der letzten Stunde.
2. Lassen Sie sich die verwendeten Abfragen zeigen und daraus ein Dashboard erstellen.
3. Lassen Sie eine Abfrage erklären und bestätigen Sie Service-Scope und Zeitraum.

Prompts in dieser Reihenfolge:

> Erzähl mir etwas über die Performance der Anwendungen in workshop-shop und des Clusters, auf
> dem sie laufen. Zeig mir die Abfragen, die du verwendet hast.

> Erstelle aus diesen Abfragen ein Dashboard mit Zeilen für Fehler, Latenz und Sättigung.

> Erkläre die Latenzabfrage. Was genau misst sie, und über welchen Zeitraum?

**Iterieren, nicht neu starten:** Bleiben Sie in derselben Unterhaltung und verfeinern Sie, etwa
ein Panel ergänzen, eine Visualisierung ändern oder auf einen Pod filtern. Ein Neustart verliert
den Kontext. Speichern Sie das Dashboard als **Shop - Application Health** mit Ihren Initialen,
damit niemand das Dashboard einer anderen Person überschreibt.

Zur Einordnung: Für Sättigung gibt es in diesem Stack Kubernetes-Metriken (cAdvisor und
kube-state-metrics) für alle Pods in `workshop-shop` auf dem Cluster `workshop-prod`. Die
Metrik `process_cpu_usage` existiert nur für `shop-catalog`. Welche Panels Assistant vorschlägt,
ist nicht festgelegt; prüfen Sie jedes selbst.

### Erwartetes Ergebnis und Erfolgskontrolle

Ein Dashboard mit Panels für Fehler, Latenz und Sättigung, bei dem Sie für jedes Panel sagen
können, was es misst und über welches Fenster. Ein oder zwei Panels brauchen meist einen
Folge-Prompt; die Datenquelle mit `@` zu benennen reduziert leere Panels.
**Ein Dashboard existiert, und Sie können jede Abfrage darauf erklären.**

### Früher fertig?

> Erstelle eine Alarmregel, die feuert, wenn die CPU-Auslastung eines Service in workshop-shop
> fünf Minuten lang über 90 % liegt.

Prüfen Sie den Entwurf, bevor Sie ihn speichern: Metrik, Filter, Schwelle und Wartezeit.
Speichern ist für diese Session nicht nötig.

### Wenn Sie nicht weiterkommen

Verwenden Sie die Prompts oben wörtlich für `workshop-shop`. Liefert ein Panel keine Daten,
benennen Sie die Datenquelle ausdrücklich mit `@` und fragen Sie erneut.

## Lab 2 · Deck-Folie 11 · Ein Symptom untersuchen

### Szenario

Seit `{fault-time}` zeigt `shop-payment` in `workshop-shop` eine erhöhte Fehlerrate. Zu diesem
Zeitpunkt ging eine Änderung live: `{change}`. Ihre Aufgabe ist nicht die Behebung, sondern eine
Hypothese, die Sie verteidigen können, und danach eine Aussage darüber, was es gekostet hat.

### Schritte

1. Beschreiben Sie ein echtes Symptom mit Service, Zeitraum und Änderung.
2. Starten Sie die Untersuchung und lesen Sie Befunde und die Belege dahinter.
3. Validieren Sie einen Befund selbst in Explore, bevor Sie ihn akzeptieren.
4. Lassen Sie Assistant die Auswirkung beziffern: wie viele Requests betroffen waren und wie lange.

Prompts:

> Untersuche die erhöhte Fehlerrate von shop-payment in workshop-shop in den letzten 3 Stunden.
> Sie begann nach {fault-time}.

> Zeig mir die Abfrage hinter diesem Befund.

Führen Sie diese Abfrage selbst in Explore aus und prüfen Sie Scope, Datenquelle und Zeitraum.

> Wie viele Nutzer oder Requests waren betroffen, und wie lange?

`{fault-time}` ersetzt die Kursleitung vorab durch die tatsächliche Zeit mit Datum und Zeitzone.

### Falls Investigations nicht verfügbar ist: als Unterhaltung

> Zeige die Fehlerrate von shop-payment in der letzten Stunde. Wann genau hat sie sich verändert?

> Zeige die Fehlerlogs von shop-payment in diesem Fenster. Welche Muster siehst du?

> Finde mit Traces heraus, welcher Aufruf fehlschlägt.

Liegt `{fault-time}` länger als eine Stunde zurück, ersetzen Sie „letzte Stunde“ durch das
tatsächliche Fenster.

### Zur Auswirkung: erst validieren, dann beziffern

Die Daten sind synthetisch: Es gibt keine echten Nutzer und keinen Umsatz. Gezählte
Fehlerlog-Ereignisse oder Fehler-Spans sind in dieser Übung simulierte Requests. Die Zählung des
Histogramms sind modellierte Beobachtungen, **kein** Request-Durchsatz; leiten Sie daraus keine
Nutzerzahl ab. Eine unvalidierte Zahl zur Auswirkung wird in Management-Berichten noch lange
wiederholt, nachdem die Details vergessen sind. Deshalb die Reihenfolge: erst den Befund prüfen,
dann die Auswirkung benennen und ihre Grundlage dazuschreiben.

### Erwartetes Ergebnis und Erfolgskontrolle

Eine Ursachenhypothese, die Sie verteidigen oder verwerfen können, mit einer Abfrage, die Sie
selbst ausgeführt haben, plus eine einzeilige Aussage zur Auswirkung. Eine Hypothese zu verwerfen
ist ein gültiges und wertvolles Ergebnis. Abhängigkeit und zeitliches Zusammentreffen im
Trace-Pfad sind noch kein Ursachennachweis.
**Sie können die Hypothese mit Belegen verteidigen und die Auswirkung benennen.**

### Wenn Sie nicht weiterkommen

Untersuchen Sie das vorbereitete Szenario und seinen bekannten Fehler im Fenster ab
`{fault-time}`. Fehlt die Änderung in den Befunden, liegt das daran, dass sie nur mündlich
oder als Annotation existiert, nicht an Ihrem Prompt.

## Lab 3 · Deck-Folie 17 · Wiederholbar machen

### Szenario

Drei kurze Schritte machen aus der heutigen Methode etwas, das Ihr Team behält: eine Rule für
immer, ein Skill für eine bestimmte Aufgabe, eine Automation für den richtigen Zeitpunkt.

### Schritt 1: Eine Rule anlegen

**Assistant → Settings → Create rule.** Konkret formulieren, Scope heute **Just me**.

> Prüfe bei der Untersuchung von Problemen immer zuerst Traces, dann Logs. Beginne mit Tempo und
> gehe dann in Loki.

> Wenn ein kritisches Problem gefunden wird, nenne immer das verantwortliche Team und schlage vor,
> über die IRM-Integration einen Incident anzulegen.

Das ist eine Assistant-Verhaltensregel, keine Alarmregel.

### Schritt 2: Einen Skill anlegen

Eigenen Skill schreiben oder aus einer Vorlage erstellen. Beispiel:

> Name: payment triage
>
> Verwenden bei: Alarmen oder Fragen zum Zustand von shop-payment
>
> Anweisungen:
>
> 1. Öffne das Dashboard Shop - Application Health (mit meinen Initialen) aus Lab 1.
> 2. Vergleiche die Fehlerrate der letzten Stunde mit der vorherigen Stunde.
> 3. Prüfe auf Deployments, Konfigurations- oder Fehlerratenänderungen im Fenster.
> 4. Korrelieren die Fehler mit einer Änderung oder einem Anstieg, fasse zusammen: Symptom,
>    Zeitpunkt, vermutete Änderung, Belege.
> 5. Ist der Schweregrad kritisch, nenne das Team, das alarmiert werden soll.

### Schritt 3: Als Automation planen

Den Slash-Befehl des Skills aktivieren und einmal ausführen. Dann **Assistant → Settings →
Automations → New**: Zeitplan und Prompt hinzufügen oder den eben gebauten Skill aufrufen.
Einmal manuell ausführen, um die Funktion zu prüfen.

### Erwartetes Ergebnis und Erfolgskontrolle

Eine Rule, die sichtbar verändert, wie sich eine neue Unterhaltung verhält; ein Skill, der den
ersten Schritt liefert, den Sie tatsächlich gehen würden, keine generische Checkliste; und eine
geplante Automation mit einem erfolgreichen manuellen Lauf. Vermutlich schaffen Sie nicht alle
drei; das ist erwartet, den Rest können Sie danach auf Ihrem eigenen Stack erledigen.
**Eine Rule prägt jede Unterhaltung, und ein geplanter Skill liefert Ihren ersten Schritt.**

### Wenn Sie nicht weiterkommen

Legen Sie die Rule an und passen Sie den Skill oben mit dem vorbereiteten Szenario-Dashboard und
`shop-payment` an. Bei fehlendem Zugriff auf Settings ist das ein Fallback, **kein** erfolgreich
gespeichertes Artefakt; Berechtigungen nicht während des Labs improvisieren.

Aktivierung und Zugriff prüft die Kursleitung vorab anhand der
[offiziellen Assistant-Anleitung](https://grafana.com/docs/grafana-cloud/platform/grafana-assistant/introduction/)
und [Aktivierung und Zugriff](https://grafana.com/docs/grafana-cloud/platform/grafana-assistant/get-started/grafana-cloud/);
UI-Verfügbarkeit und Berechtigungen können je Stack abweichen.

## Abschluss und nächste Durchführung

Welche Abfrage auf Ihrem Dashboard konnten Sie erklären, welche nicht? Welche Hypothese haben
Sie validiert oder verworfen, und worauf beruht Ihre Aussage zur Auswirkung? Welche Rule, welchen
Skill würden Sie mit dem Team teilen, und wer verantwortet ihn?

Die Kursleitung bestätigt, dass `payment-regression` deaktiviert ist, und setzt den
[Übungszustand](kubernetes.md#reset-the-exercise-not-the-telemetry) zurück. Vor der nächsten
Gruppe erneut alle Rules mit Scope **Everybody** löschen und ein frisches `{fault-time}`
vorbereiten. Alte Telemetrie muss nicht gelöscht werden.

## Fehlerhilfe für die Kursleitung

| Beobachtung | Zuerst prüfen |
|---|---|
| Assistant antwortet anders als erwartet | Übrig gebliebene Everybody-Rules, persönliche Rules der Person |
| Panels ohne Daten | Datenquelle mit `@` benennen, Zeitraum und Filter prüfen |
| Keine erhöhten Fehler bei `shop-payment` | Szenario-Fenster, Ingest-Verzögerung, `{fault-time}` samt Zeitzone |
| Änderung erscheint nicht in den Befunden | Erwartet ohne Annotation; mündlich nennen |
| Investigations fehlt | Konversations-Fallback auf der Lab-2-Karte |
| Settings nicht erreichbar | Rolle und Assistant-Zugriff des Teilnehmerkontos |
| Alles leer | Dry-run, Blueprint-Auswahl, pausierter Pod und Readiness |

Weitere Diagnose: [Control Plane](control-plane.md) und [Troubleshooting](troubleshooting.md).
Teilnehmende melden Probleme der Kursleitung, statt Tokens, Helm-Werte oder gemeinsame
Ressourcen zu ändern.
