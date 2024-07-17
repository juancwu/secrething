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

### Prepare bento

It prepares a new bento and stores it in the database. You will need an access token.

```
POST /bento/prepare HTTP/1.1
Host: konbini.juancwu.dev
Content-Type: application/json
Authorization: Bearer <token>

JSON Body:
    name: required,min=3,max=50,ascii
    pub_key: required
    ingridients?: []{ name: string, value: string }

200 OK: Bento prepared but failed to add ingridients (if provided)
    Content-Type: application/json
    JSON Body:
        message: string
        bento_id: string

201 Created: Bento prepared and ingridients added (if provided)
    Content-Type: application/json
    JSON Body:
        message: string
        bento_id: string

400 Bad Request:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
        errors?: []string

403 Unauthorized:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```

### Order bento

This will get you the bento with all the ingridients in it. There must be an existing bento before ordering it.

```
GET /bento/order/:bento_id HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <token>

Query:
    signature: signature using the RSA private key for the bento.
    challenge: a random message to sign with the RSA private key.

200 OK:
    Content-Type: application/json
    JSON Body:
        message: string
        ingridients: []{ name: string, value: string }

400 Bad Request:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

404 Not Found:
    No content

500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```
