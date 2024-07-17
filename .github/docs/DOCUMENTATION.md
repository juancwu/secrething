# Konbini Documentation

This is a simple documentation for Konbini. Things might not be 1-1 with the current
implementation in `main`.

## Table of Content
- [Routes](#routes)
   * [Sign up / Create an account](#sign-up--create-an-account)
   * [Sign in / Get access and refresh tokens](#sign-in--get-access-and-refresh-tokens)
   * [Get new access token](#get-new-access-token)
   * [Verify Email](#verify-email)
   * [Resend verification email](#resend-verification-email)

## Routes

Documentation on the available routes, depracated routes, and upcoming routes.
Here you will find the route path and method along with the different request bodies
and response bodies.

### Sign up / Create an account

This route handles requests to create a new account. An account is required to prepare new bentos
and manage existing bentos.

```
POST /auth/signup HTTP/1.1
Host: konbini.juancwu.dev
Content-Type: application/json

{
    "email": "your_email@mail.com",
    "password": "strong password",
    "name": "Your Name"
}
```

### Sign in / Get access and refresh tokens

This route handles requests to sign into an account. This route will response with access and refrehs tokens.
```
POST /auth/signin HTTP/1.1
Host: konbini.juancwu.dev
Content-Type: application/json

{
    "email": "your_email@mail.com",
    "password": "strong password"
}
```

### Get new access token

In case the access token had expired, you can get a new one using this route as long as the refresh token is still valid.

```
PATCH /auth/refresh HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <access_token>
```

### Verify email

Verify the email of an account using the code that was sent to the email.
This route's method is `GET` because this route is also sent to the email to facilitate one click verify.

```
GET /auth/email/verify?code HTTP/1.1
Host: konbini.juancwu.dev

Query:
    code: required, len=20
```

### Resend verification email

It resends an email to verify if the account using the given email has not been verified yet.

```
POST /auth/email/resend HTTP/1.1
Host: konbini.juancwu.dev
Content-Type: application/json

{
    "email": "your@mail.com"
}
```
