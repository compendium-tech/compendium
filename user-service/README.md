# Table of Contents

 -   [Authentication Flow Overview](#authentication-flow-overview)
        -   [Sign-Up](#sign-up)
            -   [Initiating Sign-Up](#initiating-sign-up)
            -   [Multi-Factor Authentication (MFA) for Sign-Up](#multi-factor-authentication-mfa-for-sign-up)
        -   [Sign-In](#sign-in)
            -   [Sign-In from a Known Device (Password Authentication)](#sign-in-from-a-known-device-password-authentication)
            -   [Sign-In from a New Device (MFA Required)](#sign-in-from-a-new-device-mfa-required)
        -   [Password Resets](#password-resets)
            -   [Initiating Password Resets](#initiating-password-resets)
            -   [MFA for Password Resets](#mfa-for-password-resets)
        -   [Refresh](#refresh)
        -   [Token Management and Security Note](#token-management-and-security-note)

-----

## Authentication Flow Overview
This part of the document outlines a typical authentication flow, covering user registration (sign-up) and user login (sign-in) scenarios, including considerations for multi-factor authentication (MFA) and token management.

### Sign-Up
The sign-up process involves two primary steps: initial user registration and subsequent multi-factor authentication (MFA) verification.

#### Initiating Sign-Up
When a new user wishes to create an account, they send a POST request to the `/api/v1/users` endpoint. This request typically includes their chosen name, email address, and a password.

```http
POST /api/v1/users

{
  "name": "Adi",
  "email": "adisalimgereev@gmail.com",
  "password": "Qwerty12345!!"
}
```

Upon successful receipt of this information, the server creates the user's account. At this stage, it's common for the server to initiate an MFA process, such as sending a One-Time Password (OTP) to the provided email address or phone number.

#### Multi-Factor Authentication (MFA) for Sign-Up
After the initial registration, the user needs to verify their identity using the OTP they received. This is done by sending a POST request to the `/api/v1/sessions` endpoint with the `flow` parameter set to `mfa`. The request body includes their email and the otp received.

```http
POST /api/v1/sessions?flow=mfa

{
  "email": "adisalimgereev@gmail.com",
  "otp": "12345"
}
```

If the OTP is valid, the server responds by setting several HTTP-only cookies: `csrfToken`, `accessToken`, and `refreshToken`. These cookies are crucial for maintaining the user's session and security. The response body also provides the expiry times for the accessToken and refreshToken, allowing the client to manage their session validity.

```http
Set-Cookie: csrfToken=..., accessToken=..., refreshToken=...

{
  "accessTokenExpiry": "...",
  "refreshTokenExpiry": "...",
}
```

### Sign-In
The sign-in process can vary depending on whether the user is logging in from the same device or a different device, particularly concerning MFA requirements.

#### Sign-In from a Known Device (Password Authentication)
When a user attempts to sign in from a device they have previously used, the system might allow a password authentication flow without immediate MFA if the device is recognized or trusted. The user sends a POST request to `/api/v1/sessions` with `flow` set to `password`, providing their email and password.

```http
POST /api/v1/sessions?flow=password

{
  "email": "adisalimgereev@gmail.com",
  "password": "Qwerty12345!!"
}
```

If the credentials are correct and MFA is not required for this specific login attempt (e.g., the device is remembered), the server responds similarly to the MFA verification step during sign-up. It sets the csrfToken, accessToken, and refreshToken as HTTP-only cookies and provides their expiry times in the response body. The `isMfaRequired` flag will be false.

```http
201 Created

Set-Cookie: csrfToken=..., accessToken=..., refreshToken=...

{
  "isMfaRequired": false,
  "accessTokenExpiry": "...",
  "refreshTokenExpiry": "...",
}
```

#### Sign-In from a New Device (MFA Required)
If a user tries to sign in from a new or unrecognized device, the system typically enforces MFA for added security. The initial POST request to `/api/v1/sessions` with `flow` set to `password` will still be made with the email and password.

```http
POST /api/v1/sessions?flow=password

{
  "email": "adisalimgereev@gmail.com",
  "password": "Qwerty12345!!"
}
```

In this scenario, the server's response will indicate that MFA is required by setting `isMfaRequired` to `true`. It will not issue access or refresh tokens at this stage.

```http
202 Accepted

{
  "isMfaRequired": true,
}
```

Following this, the user must complete the MFA step by providing the OTP, similar to the sign-up MFA process. They send a POST request to `/api/v1/sessions` with `flow` set to `mfa`, including their email and the otp.

```http
POST /api/v1/sessions?flow=mfa

{
  "email": "adisalimgereev@gmail.com",
  "otp": "12345"
}
```

Upon successful OTP verification, the server will then issue the `csrfToken`, `accessToken`, and `refreshToken` as HTTP-only cookies, along with their expiry times in the response body.

```http
201 Created

Set-Cookie: csrfToken=..., accessToken=..., refreshToken=...

{
  "accessTokenExpiry": "...",
  "refreshTokenExpiry": "...",
}
```

### Password Resets
The password reset process typically involves two main steps: initiating the reset request and then confirming the new password using a verification mechanism like an OTP.

#### Initiating Password Resets
When a user forgets their password, they can initiate a password reset by sending a PUT request to the `/api/v1/password` endpoint with `flow` set to `init`. The request body usually contains the user's email address. This informs the server to begin the reset process and, for security, send an OTP to the user's registered email or phone number.

```http
PUT /api/v1/password?flow=init

{
  "email": "adisalimgereev@gmail.com"
}
```

Upon successful receipt of this request, the server confirms that a password reset has been initiated and that an OTP has been sent. No tokens are issued at this stage.

```http
202 Accepted
```

#### MFA for Password Resets
After receiving the OTP, the user can then set a new password. This is done by sending a PUT request to the `/api/v1/password` endpoint with `flow` set to `finish`. The request body must include their email, the received OTP, and their new chosen password.

```http
PUT /api/v1/password?flow=finish

{
  "email": "adisalimgereev@gmail.com",
  "otp": "67890",
  "password": "NewStrongPassword!1"
}
```

If the OTP is valid and the new password meets the system's security requirements, the server updates the user's password.

### Refresh

The refresh mechanism allows users to obtain a new access token without needing to re-authenticate with their credentials. To refresh an expired access token, the client sends a POST request to the `/api/v1/sessions` endpoint with the `flow` parameter set to `refresh`. This request must include refresh token in the Cookie header. The server uses the refresh token to verify the user's session and issue a new, valid access token.

```http
POST /api/v1/sessions?flow=refresh

Cookie: accessToken=..., refreshToken=...
```

If the provided refresh token is valid and has not expired, the server will respond by issuing a new access token and refresh token (to implement refresh token rotation). These tokens are set as HTTP-only cookies.

```http
201 Created

Set-Cookie: csrfToken=..., accessToken=..., refreshToken=...

{
  "accessTokenExpiry": "...",
  "refreshTokenExpiry": "...",
}
```

### Token Management and Security Note
It's important to note the roles of the tokens and cookies in this flow:

- `accessToken` and `refreshToken`: These are critical for authentication and session management. They are stored as HTTP-only cookies, which means they cannot be accessed by client-side JavaScript. This significantly reduces the risk of Cross-Site Scripting (XSS) attacks stealing these tokens.
- `csrfToken`: This token is used to protect against Cross-Site Request Forgery (CSRF) attacks. While the `accessToken` and `refreshToken` are HTTP-only, the `csrfToken` is a non-HTTP-only cookie. This allows client-side JavaScript to read its value and include it in subsequent requests (e.g., in a custom header like `X-CSRF-Token`).

A key point regarding the `csrfToken`: Even if the `csrfToken` is also included as a claim within the `accessToken`, this design is acceptable. The primary protection for the `accessToken` and `refreshToken` comes from their HTTP-only nature, preventing client-side script access. The `csrfToken` serves its purpose by being readable by JavaScript for inclusion in requests, thereby validating that the request originated from the legitimate client application. So this is used in our current auth implementation in `user-service/`.
