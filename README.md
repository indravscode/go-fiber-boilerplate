# RESTful API Go Fiber Boilerplate

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
[![Go Report Card](https://goreportcard.com/badge/github.com/indravscode/go-fiber-boilerplate)](https://goreportcard.com/report/github.com/indravscode/go-fiber-boilerplate)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
![Repository size](https://img.shields.io/github/repo-size/indravscode/go-fiber-boilerplate?color=56BEB8)
![Build](https://github.com/indravscode/go-fiber-boilerplate/workflows/Build/badge.svg)
![Test](https://github.com/indravscode/go-fiber-boilerplate/workflows/Test/badge.svg)
![Linter](https://github.com/indravscode/go-fiber-boilerplate/workflows/Linter/badge.svg)

A boilerplate/starter project for quickly building RESTful APIs using Go, Fiber, and PostgreSQL. Inspired by the Express boilerplate.

The app comes with many built-in features, such as authentication using JWT and Google OAuth2, request validation, unit and integration tests, docker support, API documentation, pagination, etc. For more details, check the features list below.

## Quick Start

To create a project, simply run:

```bash
go mod init <project-name>
```

## Manual Installation

If you would still prefer to do the installation manually, follow these steps:

Clone the repo:

```bash
git clone --depth 1 https://github.com/indravscode/go-fiber-boilerplate.git
cd go-fiber-boilerplate
rm -rf ./.git
```

Install the dependencies:

```bash
go mod tidy
```

Set the environment variables:

```bash
cp .env.example .env

# open .env and modify the environment variables (if needed)
```

## Table of Contents

- [Features](#features)
- [Commands](#commands)
- [Environment Variables](#environment-variables)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
- [Error Handling](#error-handling)
- [Validation](#validation)
- [Authentication](#authentication)
- [Authorization](#authorization)
- [Logging](#logging)
- [Linting](#linting)
- [Contributing](#contributing)

## Features

- **SQL database**: [PostgreSQL](https://www.postgresql.org) Object Relation Mapping using [Gorm](https://gorm.io)
- **Database migrations**: with [golang-migrate](https://github.com/golang-migrate/migrate)
- **Validation**: request data validation using [Package validator](https://github.com/go-playground/validator)
- **Logging**: using [Logrus](https://github.com/sirupsen/logrus) and [Fiber-Logger](https://docs.gofiber.io/api/middleware/logger)
- **Testing**: unit and integration tests using [Testify](https://github.com/stretchr/testify) and formatted test output using [gotestsum](https://github.com/gotestyourself/gotestsum)
- **Error handling**: centralized error handling mechanism
- **API documentation**: with [Swag](https://github.com/swaggo/swag) and [Swagger](https://github.com/gofiber/swagger)
- **Sending email**: using [Gomail](https://github.com/go-gomail/gomail)
- **Environment variables**: using [Viper](https://github.com/spf13/viper)
- **Security**: set security HTTP headers using [Fiber-Helmet](https://docs.gofiber.io/api/middleware/helmet)
- **CORS**: Cross-Origin Resource-Sharing enabled using [Fiber-CORS](https://docs.gofiber.io/api/middleware/cors)
- **Compression**: gzip compression with [Fiber-Compress](https://docs.gofiber.io/api/middleware/compress)
- **Docker support**
- **Linting**: with [golangci-lint](https://golangci-lint.run)

## Commands

Running locally:

```bash
make start
```

Or running with live reload:

```bash
air
```

> [!NOTE]
> Make sure you have `Air` installed.\
> See 👉 [How to install Air](https://github.com/air-verse/air)

Testing:

```bash
# run all tests
make tests

# run all tests with gotestsum format
make testsum

# run test for the selected function name
make tests-TestUserModel
```

> [!IMPORTANT]
> Tests use a **separate test database**.
>
> By default, the test database name is defined in:
> `test/init.go`
>
> ```go
> DB = database.Connect("localhost", "testdb")
> ```
>
> Make sure the test database (`testdb`) **already exists** and all required
> tables (`users`, `tokens`, etc.) have been migrated before running the test commands.

Docker:

```bash
# run docker container
make docker

# run all tests in a docker container
make docker-test
```

Linting:

```bash
# run lint
make lint
```

Swagger:

```bash
# generate the swagger documentation
make swagger
```

Migration:

```bash
# Create migration
make migration-<table-name>

# Example for table users
make migration-users
```

```bash
# run migration up in local
make migrate-up

# run migration down in local
make migrate-down

# run migration up in docker container
make migrate-docker-up

# run migration down all in docker container
make migrate-docker-down
```

## Environment Variables

The environment variables can be found and modified in the `.env` file. They come with these default values:

```bash
# server configuration
# Env value : prod || dev
APP_ENV=dev
APP_HOST=0.0.0.0
APP_PORT=3000

# database configuration
DB_HOST=postgresdb
DB_USER=postgres
DB_PASSWORD=thisisasamplepassword
DB_NAME=fiberdb
DB_PORT=5432

# JWT
# JWT secret key
JWT_SECRET=thisisasamplesecret
# Number of minutes after which an access token expires
JWT_ACCESS_EXP_MINUTES=30
# Number of days after which a refresh token expires
JWT_REFRESH_EXP_DAYS=30
# Number of minutes after which a reset password token expires
JWT_RESET_PASSWORD_EXP_MINUTES=10
# Number of minutes after which a verify email token expires
JWT_VERIFY_EMAIL_EXP_MINUTES=10

# SMTP configuration options for the email service
SMTP_HOST=email-server
SMTP_PORT=587
SMTP_USERNAME=email-server-username
SMTP_PASSWORD=email-server-password
EMAIL_FROM=support@yourapp.com

# OAuth2 configuration
GOOGLE_CLIENT_ID=yourapps.googleusercontent.com
GOOGLE_CLIENT_SECRET=thisisasamplesecret
REDIRECT_URL=http://localhost:3000/v1/auth/google-callback
```

## Project Structure

```
src\
 |--config\         # Environment variables and configuration related things
 |--controller\     # Route controllers (controller layer)
 |--database\       # Database connection & migrations
 |--docs\           # Swagger files
 |--middleware\     # Custom fiber middlewares
 |--model\          # Postgres models (data layer)
 |--response\       # Response models
 |--router\         # Routes
 |--service\        # Business logic (service layer)
 |--utils\          # Utility classes and functions
 |--validation\     # Request data validation schemas
 |--main.go         # Fiber app
```

## API Documentation

To view the list of available APIs and their specifications, run the server and go to `http://localhost:3000/v1/docs` in your browser.

![Auth](https://indravscode.github.io/assets/images/swagger1.png)
![User](https://indravscode.github.io/assets/images/swagger2.png)

This documentation page is automatically generated using the [Swag](https://github.com/swaggo/swag) definitions written as comments in the controller files.

See 👉 [Declarative Comments Format.](https://github.com/swaggo/swag#declarative-comments-format)

## API Endpoints

List of available routes:

**Auth routes**:\
`POST /v1/auth/register` - register\
`POST /v1/auth/login` - login\
`POST /v1/auth/logout` - logout\
`POST /v1/auth/refresh-tokens` - refresh auth tokens\
`POST /v1/auth/forgot-password` - send reset password email\
`POST /v1/auth/reset-password` - reset password\
`POST /v1/auth/send-verification-email` - send verification email\
`POST /v1/auth/verify-email` - verify email\
`GET /v1/auth/google` - login with google account

**User routes**:\
`POST /v1/users` - create a user\
`GET /v1/users` - get all users\
`GET /v1/users/:userId` - get user\
`PATCH /v1/users/:userId` - update user\
`DELETE /v1/users/:userId` - delete user

## Error Handling

The app includes a custom error handling mechanism, which can be found in the `src/utils/error.go` file.

It also utilizes the `Fiber-Recover` middleware to gracefully recover from any panic that might occur in the handler stack, preventing the app from crashing unexpectedly.

The error handling process sends an error response in the following format:

```json
{
  "code": 404,
  "status": "error",
  "message": "Not found"
}
```

Fiber provides a custom error struct using `fiber.NewError()`, where you can specify a response code and a message. This error can then be returned from any part of your code, and Fiber's `ErrorHandler` will automatically catch it.

For example, if you are trying to retrieve a user from the database but the user is not found, and you want to return a 404 error, the code might look like this:

```go
func (s *userService) GetUserByID(c *fiber.Ctx, id string) {
	user := new(model.User)

	err := s.DB.WithContext(c.Context()).First(user, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}
}
```

## Validation

Request data is validated using [Package validator](https://github.com/go-playground/validator). Check the [documentation](https://pkg.go.dev/github.com/go-playground/validator/v10) for more details on how to write validations.

The validation schemas are defined in the `src/validation` directory and are used within the services by passing them to the validation logic. In this example, the CreateUser method in the userService uses the `validation.CreateUser` schema to validate incoming request data before processing it. The validation is handled by the `Validate.Struct` method, which checks the request data against the schema.

```go
import (
	"app/src/model"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

func (s *userService) CreateUser(c *fiber.Ctx, req validation.CreateUser) (*model.User, error) {
	if err := s.Validate.Struct(&req); err != nil {
		return nil, err
	}
}
```

## Authentication

To require authentication for certain routes, you can use the `Auth` middleware.

```go
import (
	"app/src/controllers"
	m "app/src/middleware"
	"app/src/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, u services.UserService, t services.TokenService) {
  userController := controllers.NewUserController(u, t)
	app.Post("/users", m.Auth(u), userController.CreateUser)
}
```

These routes require a valid JWT access token in the Authorization request header using the Bearer schema. If the request does not contain a valid access token, an Unauthorized (401) error is thrown.

**Generating Access Tokens**:

An access token can be generated by making a successful call to the register (`POST /v1/auth/register`) or login (`POST /v1/auth/login`) endpoints. The response of these endpoints also contains refresh tokens (explained below).

An access token is valid for 30 minutes. You can modify this expiration time by changing the `JWT_ACCESS_EXP_MINUTES` environment variable in the .env file.

**Refreshing Access Tokens**:

After the access token expires, a new access token can be generated, by making a call to the refresh token endpoint (`POST /v1/auth/refresh-tokens`) and sending along a valid refresh token in the request body. This call returns a new access token and a new refresh token.

A refresh token is valid for 30 days. You can modify this expiration time by changing the `JWT_REFRESH_EXP_DAYS` environment variable in the .env file.

## Authorization

The `Auth` middleware can also be used to require certain rights/permissions to access a route.

```go
import (
	"app/src/controllers"
	m "app/src/middleware"
	"app/src/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, u services.UserService, t services.TokenService) {
  userController := controllers.NewUserController(u, t)
	app.Post("/users", m.Auth(u, "manageUsers"), userController.CreateUser)
}
```

In the example above, an authenticated user can access this route only if that user has the `manageUsers` permission.

The permissions are role-based. You can view the permissions/rights of each role in the `src/config/roles.go` file.

If the user making the request does not have the required permissions to access this route, a Forbidden (403) error is thrown.

## Logging

Import the logger from `src/utils/logrus.go`. It is using the [Logrus](https://github.com/sirupsen/logrus) logging library.

Logging should be done according to the following severity levels (ascending order from most important to least important):

```go
import "app/src/utils"

utils.Log.Panic('message') // Calls panic() after logging
utils.Log.Fatal('message'); // Calls os.Exit(1) after logging
utils.Log.Error('message');
utils.Log.Warn('message');
utils.Log.Info('message');
utils.Log.Debug('message');
utils.Log.Trace('message');
```

> [!NOTE]
> API request information (request url, response code, timestamp, etc.) are also automatically logged (using [Fiber-Logger](https://docs.gofiber.io/api/middleware/logger)).

## Linting

Linting is done using [golangci-lint](https://golangci-lint.run)

See 👉 [How to install golangci-lint](https://golangci-lint.run/welcome/install)

To modify the golangci-lint configuration, update the `.golangci.yml` file.

## Contributing

Contributions are more than welcome! Please check out the [contributing guide](CONTRIBUTING.md).

If you find this boilerplate useful, consider giving it a star! ⭐

## Inspirations

- [hagopj13/node-express-boilerplate](https://github.com/hagopj13/node-express-boilerplate)
- [khannedy/golang-clean-architecture](https://github.com/khannedy/golang-clean-architecture)
- [zexoverz/express-prisma-template](https://github.com/zexoverz/express-prisma-template)

## License

[MIT](LICENSE)

## Contributors

[![Contributors](https://contrib.rocks/image?c=6&repo=indravscode/go-fiber-boilerplate)](https://github.com/indravscode/go-fiber-boilerplate/graphs/contributors)
