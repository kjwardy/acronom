# Acronom

Placeholder for CircleCI and GitHub actions badges

[![GitHub Actions](https://github.com/kjwardy/acronom/actions/workflows/validate.yml/badge.svg?branch=main)](https://github.com/kjwardy/acronom/actions/workflows/validate.yml?query=branch%3Amain)

A simple collaborative acronym glossary in Go with a React frontend and Postgres database.

# Features

// TODO - list top-level features

# Getting Started

The recommended way to test and deploy is using Docker. You will need to run both the `acronom` app, and the Postgres database.

**Start Postgres**

```sh
docker run --name db -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=acronom -e POSTGRES_ADDR=db:5432 -d postgres:17
```

**Start App**

```sh
docker run -p 1323:1323 -e HOSTS=localhost -e APP_URI=http://localhost:1323 --link db acronom
```

Alternatively use the docker-compose file and run:

```sh
docker-compose up
```

## Development

Run PostgreSQL, the API, and the frontend separately during development. The React development server automatically reloads frontend changes; restart the Go process after backend changes.

**Start PostgreSQL**

The default API configuration works with this local Docker container:

```sh
docker run --name db -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=acronom -e POSTGRES_ADDR=db:5432 -d postgres:17
```

After creating the container, it can be stopped and started again with `docker stop db` and `docker start db`.

**Start the API**

```sh
cd api
go mod download
go run .
```

The API listens on `http://localhost:1323`. Check it with:

```sh
curl http://localhost:1323/health
curl http://localhost:1323/api/status
```

**Start the frontend**

In another terminal:

```sh
cd frontend
yarn install --frozen-lockfile
yarn start
```

The development UI is available at `http://localhost:3000/acronom`.

## Local Deployment

Build the frontend and copy it into the API's static file directory:

```sh
cd frontend
yarn install --frozen-lockfile
yarn build
cd ..
mkdir -p api/public
cp -R frontend/build/. api/public/
```

Then build and run the API from the `api` directory so it can find the `public` directory:

```sh
cd api
go build -o acronom .
./acronom
```

Open `http://localhost:1323`; the API redirects the root path to the React application at `/acronom`.

## Testing

```sh
cd api
go test ./...

cd ../frontend
yarn lint
yarn test --watchAll=false --passWithNoTests
yarn build
```

## Environment Configuration

| Variable | Required | Default | Example | Description |
| --- | --- | --- | --- | --- |
| `HOSTS` | Yes | | acronom.domain.com,acronom2.domain.com | List of comma separated hosts that the server will be able to be accessed from |
| `APP_URI` | Yes | | https://acronom.domain.com | Default URI of app - used to link back to app |
| `PORT` | No | `1323` | `8080` | Port used by the API server |
| `DEBUG` | No | `false` | `true` | Enable Echo debug mode |
| `POSTGRES_ADDR` | No | `localhost:5432` | `database:5432` | PostgreSQL server address in `host:port` format |
| `POSTGRES_DATABASE` | No | `acronom` | `acronom` | PostgreSQL database name |
| `POSTGRES_USER` | No | `postgres` | `acronom` | PostgreSQL user |
| `POSTGRES_PASS` | No | `password` | | PostgreSQL password |

See the [API documentation](./docs/api.md) for detailed endpoint usage and examples.

More advanced deployment guides can be found on the [Advanced Deployment](./docs/deployments/advanced-deployment.md) page.
