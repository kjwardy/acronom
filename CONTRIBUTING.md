# Contributing

## Development setup

The application requires Node.js, Yarn, Go, and PostgreSQL. On macOS, the required development tools can be installed with:

```sh
brew install nvm yarn go
nvm install 22
nvm alias default 22
```

### One-time setup

These commands only need to be run once to set up your development environment:

**Start PostgreSQL container**

```sh
docker run --name db -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=acronom -e POSTGRES_ADDR=db:5432 -d postgres:17
```

**Install Go dependencies**

```sh
cd api
go mod download
```

**Install frontend dependencies**

```sh
cd frontend
yarn install
```

### Daily development

These commands are run each time you want to work on the project:

**Start the frontend**

```sh
cd frontend
yarn start
```

**Start the API**

In a separate terminal:

```sh
cd api
POSTGRES_PASS=password HOSTS=localhost APP_URI=http://localhost:3000 go run .
```

Open http://localhost:3000/acronom after the services have started.

## Verification

Run the Go test suite from the API directory:

```sh
cd api
go test ./...
go vet ./...
```

Run the frontend formatting, lint, and production build checks from the frontend directory:

```sh
cd frontend
yarn prettier:check
yarn lint
yarn build
```

Please resolve new warnings or failures caused by your changes before submitting a pull request.

GitHub Actions runs these API and frontend checks for every pull request targeting `master` and every push to `master`. CircleCI remains responsible for validating tagged releases and publishing Docker images.

## Pull request review

Maintainers may request changes or ask questions during review. Please keep discussion constructive and update the pull request when feedback is addressed. Once approved, a maintainer will merge the change.
