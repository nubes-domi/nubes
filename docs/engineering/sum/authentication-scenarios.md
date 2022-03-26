# Signin flowchart

## Case: no user is signed in

```mermaid
flowchart TB
    doany([Any page visit])-->redirect
    redirect[Redirect to signin page]-->typecred[Write credentials]-->confirm[Signin]

    confirm-->valid{Credentials are valid?}

    valid-- yes -->back([Redirect back])
    valid-- no -->tryagain[Try again]
    tryagain-->typecred
```

## Case: one user already signed in

```mermaid
flowchart TB
    switch([Click on switch user])-->typecred[Write credentials]-->confirm[Signin]

    confirm-->valid{Credentials are valid?}

    valid-- no -->tryagain[Try again]
    valid-- yes -->duplicate{Already logged in?}
    duplicate-- no -->newsess[Start session]-->back([Redirect back])
    duplicate-- yes --> back
    tryagain-->typecred
```

## Case: multiple users already signed in

```mermaid
flowchart TB
    switch([Click on switch user])-->showswitcher[Show account switcher]

    showswitcher-->pick[Pick existing session]
    showswitcher-->newsession[Start new session]

    pick-->forcereauth{Settings force reauth?}

    forcereauth-- no -->expired{Session expired?}
    expired-- no -->oidcmaxage{OIDC reauth?}
    oidcmaxage-- no -->switchsession[Change current user]-->back

    forcereauth-- yes -->reauthenticate[Reauthenticate]
    expired-- yes -->reauthenticate
    oidcmaxage-- yes -->reauthenticate

    reauthenticate-->pass[Type password / Webauthn]-->confirmreauth[Confirm]
    confirmreauth-->validreauth{Accepted?}
    validreauth-- yes -->switchsession
    validreauth-- no -->tryagainreauth[Try again]
    tryagainreauth-->pass

    newsession-->typecred[Write credentials]-->confirm[Signin]
    confirm-->valid{Credentials are valid?}
    valid-- yes -->duplicate{Already logged in?}
    valid-- no -->tryagain[Try again]
    duplicate-- no -->newsess[Start session]-->back([Redirect back])
    duplicate-- yes --> back
    tryagain-->typecred
```

## Case: OIDC, no session

```mermaid
sequenceDiagram
    participant RP as RP
    participant U as User/Browser
    participant E as Experior
    participant S as Sum

    RP->>U: Initiate OIDC Authorization Flow
    U->>E: Visit OIDC Authorization
    E->>S: Store new Authorization Flow
    S->>E: Auth flow ID
    E->>U: Request credentials
    U->>E: Provide credentials
    E->>S: Start session
    S->>E: Authentication token
    E->>U: Have a cookie!
    E->>U: Redirect to grants
    U->>E: Confirm grants
    E->>S: Save grants
    E->>U: Redirect to RP
    U->>RP: Continue
```

## Case: OIDC, one session

```mermaid
sequenceDiagram
    participant RP as RP
    participant U as User/Browser
    participant E as Experior
    participant S as Sum

    RP->>U: Initiate OIDC Authorization Flow
    U->>E: Visit OIDC Authorization
    E->>S: Store new Authorization Flow
    S->>E: Auth flow ID
    alt reauthentication required
        U->>E: Reauthenticate
        E->>S: Refresh session
        S->>E: Authentication token
        E->>U: Have a cookie!
    end
    E->>U: Redirect to grants
    alt User wants to start new session
        U->>E: Provide credentials
        E->>S: Start session
        S->>E: Authentication token
        E->>U: Have a cookie!
        E->>U: Redirect to grants
        U->>E: Confirm grants
    else
        U->>E: Confirm grants
    end
    E->>S: Save grants
    E->>U: Redirect to RP
    U->>RP: Continue
```

## Case: OIDC, multiple session

```mermaid
sequenceDiagram
    participant RP as RP
    participant U as User/Browser
    participant E as Experior
    participant S as Sum

    RP->>U: Initiate OIDC Authorization Flow
    U->>E: Visit OIDC Authorization
    E->>S: Store new Authorization Flow
    S->>E: Auth flow ID
    E->>U: Redirect to user picker
    U->>E: Pick user
    alt reauthentication required
        U->>E: Reauthenticate (if needed)
        E->>S: Refresh session
        S->>E: Authentication token
        E->>U: Have a cookie!
    end
    E->>U: Redirect to grants
    U->>E: Confirm grants
    E->>S: Save grants
    E->>U: Redirect to RP
    U->>RP: Continue
```