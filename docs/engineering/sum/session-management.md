# Sum - Session management

Every user can sign in on multiple devices, and multiple users can sign in on
the same device.

To do so, we use two cookies:

`current_session`: an unique Session ID that identifies the Current Active Session
on this particular device.
`sessions`: a pipe ("|") separated array of signed JWTs, each one identifying
a session.

Each JWT in the `session` cookie:
- must refer to a different user;
- must have a `jti` field that is a valid Session ID;
- can have its own expiration time.

## Ensuring validity

Each JWT in the `sessions` cookie:
- is signed using the RS256 algorithm;
- references an entry in the database.

If the JWT seems tampered (signature does not match) or the `jti` field does 
not correspond to an active session, the JWT must be discarded and removed from 
`sessions` cookie.

If the `current_session` cookie was referring to a deleted JWT, that cookie has
to be deleted too.

## Session duration

At this stage we don't see the need to place restrictions on sessions duration.

The session related cookies are set to expire 10 years in the future.

## Switching users

When multiple users are signed in on a device (the `sessions` cookie has multiple 
JWTs), users can switch from one to the other by changing the `current_session` 
cookie.

### Reauthenticating

When switching from a user to another, it may be desireable to enforce 
reauthentication, especially on shared devices.

This should be left as an option for the user, defaulting to "do not 
reauthenticate".

Reauthentication should be "lightweight" and fast, avoiding having to re-type a 
password. Any sort of Biometrics, or push notification would be acceptable.

## Terminating sessions

If a User decides to terminate its session on the current device:
- the `current_session` cookie is cleared;
- the JWT for the session is removed from the `sessions` cookie array;
- the session is removed from the database.

If a User decides to terminate its session on another device:
- the session is removed from the database;

On the following Web request on that device, the service will check against the
database and discard the session JWT.
