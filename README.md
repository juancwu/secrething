# konbini

A service to manage secrets for your awesome projects.

This service should be use in conjunction with the [Bento CLI](https://github.com/juancwu/bento)

Feel free to fork/clone this repository and host the service on your own server.

## Database Setup

Make sure you have PostgreSQL (code works with v14) installed.

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

Now you can run the application and migrations will run automatically. `go run .`
