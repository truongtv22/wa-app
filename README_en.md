# wa-app

`wa-app` is a WA application-chain service that provides account management, number probing, registration, login-state checks, long-connection sessions, and message handling, with a built-in management dashboard.

> [!CAUTION]
> By using this project, you agree to all terms in [NOTICE](./NOTICE). This project is limited to protocol modeling, educational demos, authorized security research, and internal non-commercial validation. Commercial use, unauthorized targets, or scenarios that violate third-party service terms are prohibited.

## Features

- Account management: maintain WAAccount records, client profiles, registration records, and login-state projections.
- Number and registration: support number probing, SMS probing, registration requests, OTP submission, and login-state checks.
- Connection and messaging: support long-connection sessions, message receiving, message ack, 1:1 text sending, and conversation viewing.
- Data extraction: extract OTP/Flag candidates from messages and store references or redacted projections under sensitive-data rules.
- Admin UI: provide a dashboard for accounts, contacts, messages, connection status, and account-profile operations.

## Deployment

The recommended way to start the service is with the Docker Compose file included in this repository:

```sh
cp .env.example .env
docker compose pull
docker compose up -d
```

For a quick local start, `docker compose up -d` also works directly. The service can start without creating `.env`; missing values will use the compose defaults.

Default ports (fixed):

- Dashboard: `http://127.0.0.1:8080` (`docker-compose.yml` mapping)
- gRPC: `127.0.0.1:50091`

If you really need different host-mapped ports, edit the `ports` lines in `docker-compose.yml` directly. They are not exposed as config options.

### Configuration

Only a small set of required runtime settings is kept in `.env`:

- `WA_APP_IMAGE_TAG`: image tag. In production, pin a fixed version.
- `WA_APP_AUTH_PASSWORD`: optional single-password dashboard login. Leave empty to disable auth.
- `WA_APP_DATA_DIR`: persistent directory inside the container. Default: `/var/lib/wa-app`.
- `WA_APP_PG_DSN`: optional PostgreSQL DSN. Leave empty to use built-in SQLite persistence.
- `WA_APP_REDIS_URL`: optional Redis URL. Leave empty to use built-in SQLite runtime-state storage.
- `WA_COMMON_PROXY`: optional default outbound WA proxy. Leave empty for direct connection.
- `WA_NUMBER_PROBE_PROXY`: optional number/SMS probing proxy. Falls back to `WA_COMMON_PROXY`; if both are empty, direct connection is used.
- `WA_REGISTRATION_PROXY`: optional registration and OTP-submit proxy. Falls back to `WA_COMMON_PROXY`; if both are empty, direct connection is used.

PostgreSQL and Redis are both optional. To enable them, uncomment the related services in `docker-compose.yml` and set `WA_APP_PG_DSN` / `WA_APP_REDIS_URL` in `.env`.

### Build image from source

The `Dockerfile` supports building directly inside the `wa-app` repository without requiring a `common-lib` build context:

```sh
docker build -t wa-app-service:local .
```

It also works when building from the `byte-v-forge` aggregate directory:

```sh
docker build -f wa-app/Dockerfile -t wa-app-service:local .
```

## Links

- [LINUX DO - A new ideal community](https://linux.do/)
