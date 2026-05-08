# Bug/Incident Tracker (single-tenant) — kontekst projektu (Go + PostgreSQL)

## Cel biznesowy
Aplikacja służy do rejestrowania i obsługi zgłoszeń błędów oraz incydentów w zespole (QA/Support/Dev), tak aby:
- zgłoszenia nie ginęły w mailach/komunikatorach,
- było jasne ownership (kto odpowiada za temat),
- można było priorytetyzować pracę (P1–P4, środowisko),
- mieć historię zmian i w przyszłości liczyć metryki (MTTA/MTTR, re-open rate).

## Zakres MVP (co ma działać na początku)
- REST API do tworzenia i obsługi incydentów.
- PostgreSQL jako baza danych.
- Migracje bazy (wersjonowanie schematu).
- Historia zdarzeń (audit trail) od samego początku: każda istotna zmiana incydentu zapisuje event.
- Na start brak uwierzytelniania (bez JWT) — użytkownicy są w DB, ale endpointy są „otwarte” (auth jako iteracja 2).
- Paginacja list (limit/offset na MVP).
- Prosty endpoint healthcheck sprawdzający połączenie z DB.

## Model domeny (MVP)
### Statusy incydentu
- `OPEN`
- `IN_PROGRESS`
- `RESOLVED`
- `CLOSED`
(Później opcjonalnie: `REOPENED`)

### Priorytety
- `P1`, `P2`, `P3`, `P4`

### Środowiska
- `prod`, `stage`, `dev`

## Minimalny model danych (PostgreSQL)
Tabele (MVP):
1. `users`
   - `id` (UUID)
   - `email` (unique)
   - `name`
   - `created_at`

2. `incidents`
   - `id` (UUID)
   - `title`
   - `description`
   - `environment` (enum-like: tekst + walidacja po stronie API, ewentualnie CHECK w DB)
   - `priority` (P1–P4)
   - `status` (OPEN/IN_PROGRESS/RESOLVED/CLOSED)
   - `reporter_id` (FK -> users)
   - `assignee_id` (FK -> users, nullable)
   - `created_at`, `updated_at`

3. `incident_comments`
   - `id` (UUID)
   - `incident_id` (FK -> incidents)
   - `author_id` (FK -> users)
   - `body`
   - `created_at`

4. `incident_events` (audit / historia)
   - `id` (UUID)
   - `incident_id` (FK -> incidents)
   - `actor_id` (FK -> users)
   - `type` (np. `CREATED`, `STATUS_CHANGED`, `ASSIGNED`, `PRIORITY_CHANGED`, `COMMENT_ADDED`)
   - `description` 
   - `created_at`

Ważne: tworzenie incydentu i każda aktualizacja ma zapisywać event w tej samej transakcji.

## Endpointy REST (MVP)
### Health
- `GET /healthz` — status serwera + test połączenia do DB

### Users
- `POST /users` — utworzenie użytkownika (email, name)
- `GET /users` — lista użytkowników

### Incidents
- `POST /incidents` — tworzenie incydentu
- `GET /incidents` — lista z filtrami i paginacją:
  - query params: `status`, `priority`, `environment`, `assignee_id`, `q` (search po tytule/ILIKE)
  - paginacja: `limit`, `offset`
- `GET /incidents/{id}` — szczegóły incydentu
- `PATCH /incidents/{id}` — zmiana wybranych pól (np. `status`, `priority`, `assignee_id`)
  - każda zmiana zapisuje odpowiedni event do `incident_events`

### Comments
- `POST /incidents/{id}/comments` — dodanie komentarza
  - zapis komentarza + event `COMMENT_ADDED`
- `GET /incidents/{id}/comments` — lista komentarzy

## Wymagania niefunkcjonalne (MVP)
- Czytelna obsługa błędów (walidacja inputu, sensowne kody HTTP).
- Spójność danych poprzez transakcje dla operacji „zmień incydent + dopisz event”.
- Konfiguracja przez env (DB URL/host/port/user/pass/name, port aplikacji).

## Planowane technologie (proponowany stack)
- Język: Go (1.22+ lub najnowszy zainstalowany)
- HTTP router: `chi` (`github.com/go-chi/chi/v5`)
- PostgreSQL driver: `pgx` (`github.com/jackc/pgx/v5`) + `pgxpool`
- Migracje: `golang-migrate/migrate` (CLI lub biblioteka) **albo** `goose` (do wyboru)
- JSON: standard `encoding/json`
- Logowanie: na start standard `log/slog` (Go) lub później `zap`
- Kontenery: Docker + docker-compose (app + postgres)

## Struktura projektu (propozycja)
- `cmd/api/main.go` — uruchomienie serwera
- `internal/config/` — wczytywanie env/config
- `internal/http/` — router, handlers, middleware, DTO
- `internal/service/` — logika biznesowa (transakcje, reguły, eventy)
- `internal/repo/` — dostęp do DB (SQL)
- `internal/model/` — typy domenowe
- `internal/db/migrations/` — migracje SQL
- `docker-compose.yml` — postgres (i ewentualnie app)

## Iteracje rozwoju (po MVP)
1. Auth (JWT) + autoryzacja (np. user może edytować własne zgłoszenia / role).
2. Cursor-based pagination.
3. SLA + job/worker: wykrywanie incydentów “po SLA” i powiadomienia.
4. Raporty/metryki: MTTA/MTTR z `incident_events`.
5. OpenAPI/Swagger + generowanie klienta.
6. Integracje (webhook/Slack/email) — opcjonalnie.

## Notatki implementacyjne
- UUID: na start generować w Go (np. `github.com/google/uuid`) i zapisywać do DB jako `uuid`.
- Eventy (audit): trzymać minimalny `type` + `data` w JSONB; to daje elastyczność i pozwala liczyć metryki później.
- Lista incydentów: startowo sortowanie `created_at desc`.
