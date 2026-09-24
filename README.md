# Acronom

Placeholder for CircleCI and GitHub actions badges

A simple collaborative acronym glossary in Go with a React frontend and Postgres database.

# Features

// TODO - list top-level features

# Getting Started

## Development

Run the API and frontend separately during development. The React development server automatically reloads frontend changes; restart the Go process after backend changes.

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
| `PORT` | No | `1323` | `8080` | Port used by the API server |
| `DEBUG` | No | `false` | `true` | Enable Echo debug mode |

More advanced deployment guides can be found on the [Advanced Deployment](./docs/deployments/advanced-deployment.md) page.
