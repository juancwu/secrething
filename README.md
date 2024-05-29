# Konbini

A service to manage secrets for your awesome projects.
You can use any software like [curl](https://github.com/curl/curl)
or [bruno](https://github.com/usebruno/bruno) to access the service because its like any other
REST API.

Try [Konbini CLI](https://github.com/juancwu/konbini-cli), the official CLI for Konbini.

Feel free to fork/clone this repository and host the service on your own server.

## Table of Content
- [Database Setup](#database-setup)

## Database Setup

Make sure you have [PostgreSQL](https://www.postgresql.org/download/) (code works with v14) installed.

After installing POstgreSQL, run the followin commands to create the database and user.

Change user to postgres and connect to database server:
```
sudo -i -u postgres
psql
```

Crate a new database named `konbini`:
```
CREATE DATABASE konbini;
```

Create a new user named `cashier`:
```
CREATE USER cashier WITH PASSWORD 'mypassword';
```

Grant privileges to the user `cashier`:
```
GRANT ALL PRIVILEGES ON DATABASE konbini TO cashier;
```
