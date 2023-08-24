# Transaction Proxy Backend


### Database migrations

To get the [goose][golang-goose-cli], use your package manager or run:

```sh
go install github.com/pressly/goose/v3/cmd/goose
```

Ensure that you have `${GOPATH}/bin` in your `$PATH`. If you don't have GOPATH, then it's usually `${HOME}/go`

To create a migration, use:
```sh
goose -dir ./internal/pkg/db/migrations/sql create create_tokens_table sql
```

To run migrations manually, use:
```sh
goose -dir ./internal/pkg/db/migrations/sql up
```

To run migrations manually and apply missing migrations (out of order), use:
```sh
goose -allow-missing -dir ./internal/pkg/db/migrations/sql up
```

[golang-goose-cli]: https://github.com/pressly/goose
