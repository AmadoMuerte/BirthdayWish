module github.com/AmadoMuerte/BirthdayWish/API/ms/wishlist

go 1.24.3

replace github.com/AmadoMuerte/BirthdayWish/API => ../../

require (
	github.com/AmadoMuerte/BirthdayWish/API v0.0.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/render v1.0.3
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/uptrace/bun v1.2.11
	github.com/uptrace/bun/dialect/pgdialect v1.2.11
	github.com/uptrace/bun/driver/pgdriver v1.2.11
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	mellium.im/sasl v0.3.2 // indirect
)
