# Ski-Verleihsystem (Prototyp)

**Autor:** Atabek Kudaikulov
**Matrikelnummer:** 12453349

## Beschreibung
Dieses Projekt ist ein Prototyp für ein In-Memory-Verleihsystem, das als Single Page Application (SPA) mit reinem HTML/JS im Frontend und einem Go-Backend umgesetzt wurde. Es nutzt keine externe Datenbank oder ORM-Frameworks.

## Umgesetzte Anforderungen (Aufgabe 1)
Das Projekt deckt die folgenden Teilaufgaben ab:
*   **a) Ski-Verleih & Dashboard:** Implementierung eines interaktiven Dashboards zur Verwaltung von Skiern, Personen und Ausleihen. Die Ausweise der Personen können beim Registrieren hochgeladen werden (Datei-Upload ins Verzeichnis `/static/uploads`).
*   **b) Flexibilität der Preise:** Die Verleihdauer und der Verleihbeginn werden stundenweise berechnet. Bei der Rückgabe wird die reale Zeitdifferenz ermittelt und der Endpreis auf Basis des Stundenpreises kalkuliert.

## Technologie-Stack
*   **Backend:** Go (Standardbibliothek `net/http`)
*   **Frontend:** Vanilla HTML, CSS, JavaScript
*   **Datenhaltung:** In-Memory Store (`sync.Mutex` für Thread-Sicherheit)

## Starten des Projekts
1. In das Projektverzeichnis wechseln.
2. Den Server starten:
   ```bash
   go run .