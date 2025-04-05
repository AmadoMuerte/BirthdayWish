## Migrations
Service runs migration up on start.
### To add new migration you need to
Install [migration tool](https://github.com/golang-migrate/migrate):
Run this to apply migrations:
```bash
migrate -path=internal/storage/migrations/ -database=postgres://postgres:postgres@localhost:5432/db_name?sslmode=disable up [<number>]
```
Create new migration:
```bash
migrate create -dir="internal/storage/migrations/" -ext="sql" <name-of-migration>
```
