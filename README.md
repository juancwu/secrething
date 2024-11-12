# Konbini

A service to manage secrets for your awesome projects.
You can use any software like [curl](https://github.com/curl/curl)
or [bruno](https://github.com/usebruno/bruno) to access the service because its like any other
REST API.

Try [Mi CLI](https://github.com/juancwu/mi), the official CLI for Konbini.

Feel free to fork/clone this repository and host the service on your own server.
Read the [documentation](/.github/docs/DOCUMENTATION.md) for more help.

This service does not really do any encryption on secrets since it won't be storing/generating
any type of crypto keys. This ensures that if a breach were to happen, no secret stored in the database,
are leaked and vulnerable to be decrypted.

> As of now, the service does not encrypt any of the user information that is collected during
> sign-up, which include email and name. There are no enforcement to use a real name, so it is preferred
> for users to input a nickname that teammates or other parties can recognize.
> The purpose of getting a name is to be able to address that user via emails, such as verification codes
> or for teammates to know who did what.

## Table of Content

- [Self-host Konbini](#self-host-konbini)
- [Development](#development)

## Self-host Konbini

To self-host Konbini on your private machine, you can build a Docker image by cloning this repository
or forking it to modify/extend its features for your needs.

If you do not wish to use Docker, you can always just build a binary and run it.

You will also need to get yourself a [Resend](https://resend.com) API key to send emails.

Finally, make sure to have a usable PostgreSQL database.

Here is a list of environment variables and their usage that you will need:

```
# The connection string for the PostgreSQL database.
DB_URL=
# The port the server will listen to. If you are using a Docker image, make sure
# port forwarding points to this PORT for the Docker port.
PORT=
# This variable can have two possible values "production" or "development".
APP_ENV=
# The name of the database.
DB_NAME=
# Resend API key to send emails.
RESEND_API_KEY=
# This is used to fill email templates so that the BE knows its own URL.
# For example: https://my.konbini.com
SERVER_URL=
# Currently not used, so it can be left empty.
# This can be used to encrypt data in the database.
PGP_SYM_KEY=
# The hashing algorithm for passwords.
PASS_ENCRYPT_ALGO=
# The do not reply email address used by the service.
DONOTREPLY_EMAIL=
# Make sure it is a secure string that will be used for access tokens.
JWT_ACCESS_TOKEN_SECRET=
# Make sure it is a secure string that will be used for refresh tokens.
JWT_REFRESH_TOKEN_SECRET=
# This is just a string that identifies the type of jwt for access tokens.
JWT_ACCESS_TOKEN_TYPE=
# This is just a string that identifies the type of jwt for refresh tokens.
JWT_REFRESH_TOKEN_TYPE=
# A name to put in the issuer field for jwt.
JWT_ISSUER=
```

## Development

This section highlights pre-requisites and how to get started for development.

### Pre-requisites

- [Docker](https://docs.docker.com/get-started/get-docker/)
- [Go](https://go.dev/dl/)
- [Air](https://github.com/air-verse/air)
- `make` command

### Get Started

After meeting all the pre-requisites, we can now start the development environment.

The Konbini PostgreSQL database runs on a docker container which makes it easy to just start/stop.

Use the following `make` commands to initialize databas, run migrations and start development server.
Read the Makefile for more commands.

```
make init-dev-db
make up
make dev
```

### Testing

To run all tests just use the command `make run-tests`. This command will do the following:

- Start a temporary testing database using docker.
- Run all tests.
- Stop/Remove temporary testing database.
