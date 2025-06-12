  # FIDO-TEST3

A Go-based backend service using PostgreSQL, Docker, and database migrations via `golang-migrate`.
---


### 1. Clone the Repository

```bash
git pull https://github.com/Dzhodddi/FIDO-TEST3.git
cd FIDO-TEST3(if needed)
```
### 2.Build the Go Application

```bash
go build -o ./bin/main ./cmd/api
```
### 3. Run PostgreSQL via Docker
```bash
 docker run -d --name mypostgres -p 5432:5432 -e POSTGRES_PASSWORD=yourpassword postgres
```
### 4. Install migrations
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path internal/db/migrations -database "postgresql://postgres:yourpassword@localhost:5432/postgres?sslmode=disable" up
```
### 5. Run binary
```bash
./bin/main
```
Database is empty, seed script isn't here, need to add data manually to table quotes
