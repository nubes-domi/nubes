_sum_ (_I am_) is the IdP in Nubes.

# Purposes

- Store user identifiers, credentials, roles and details.
- Authenticate users to internal applications and to external apps that supports so.
- Allow administrators to manage such users

# Protocols

## OpenID Connect

**Required**. The latest and shiniest Web Authentication protocol. All new apps should be using this.

## SAML

**Maybe**. The enterprise weirdo. Speaks XML (👻)

## LDAP, Kerberos, RADIUS

**Unlikely**. Potentially authenticate to computers (Windows Logon, SSH) and networks (WPA2/3 Enterprise)

## Authentication

- Traditional username & password.
- TOTP
- Yubi?
- TouchID/FaceID?
