# Authoring APIs

The main Web app should just be very nice facade on top of a number of apps.

The frontend can only interact with the apps through HTTP API requests. Whether these API should be exposing a JSON RESTful or GraphQL interface is still undecided.

It's not excluded that some services might decide to expose both protocols, and different apps can choose which one fits best.

## REST

### The good and the bad

- ✔️ Very common, well known
- ✔️ Good tooling
- 🔴 Hard or impossible to get exactly what's needed (and only that) in a single request.
- 🔴 No real standardisation of request/response formats except for jsonapi.org, which might just be too strongly opinionated.

## GraphQL 

- ✔️ Surgical precision in fetching data
- 🔴 Not well known
- 🔴 Hard to implement

## Stick with what's out there

Interoperability is king. What's the saying? Embrace, extend, extinguish 😉.

The object storage layer should expose an S3-like API. The identity later an OpenID Connect API.

Is there a standard API for browsing photos? Let's try and use it.
