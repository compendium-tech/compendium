<p align="center">
<img width="35%" alt="download" src="https://github.com/user-attachments/assets/9e9dea7b-d67d-4c72-afd0-fa615f8e9063" /> <br />
<h4>Compendium provides comprehensive insights into university life, research opportunities, alumni networks, admission criteria and help you shortlist universities, plan admissions, assess essays and offers interactive courses for preparing to successfully pass international exams.</h4>
</p>
<img width="3168" height="792" alt="1584 (5)" src="https://github.com/user-attachments/assets/ea84f885-aaba-48bc-8c31-1edb22493b65" />


# Table of Contents

-   **[Non technical land](#non-technical-land)**
    -   [Problem](#problem)
    -   [Growth Strategy](#growth-strategy)
-   **[Making changes](#making-changes)**
-   **[Microservice Backend](#microservice-backend)**
    -   [Setting Up Locally](#setting-up-locally)
    -   [Go Code Guidelines](#go-code-guidelines)
        -   [Logging](#logging)
        -   [Error Tracing](#error-tracing)
        -   [Error Handling](#error-handling)
    -   [Authentication Flow Overview](#authentication-flow-overview)
        -   [1. Sign-Up](#1-sign-up)
            -   [1.1 Initial User Registration](#11-initial-user-registration)
            -   [1.2 Multi-Factor Authentication (MFA) Verification for Sign-Up](#12-multi-factor-authentication-mfa-verification-for-sign-up)
        -   [2. Sign-In](#2-sign-in)
            -   [2.1 Sign-In from a Known Device (Password Authentication)](#21-sign-in-from-a-known-device-password-authentication)
            -   [2.2 Sign-In from a New Device (MFA Required)](#22-sign-in-from-a-new-device-mfa-required)
        -   [3. Token Management and Security Note](#3-token-management-and-security-note)
-   **[Monolithic SPA Frontend](#monolithic-spa-frontend)**

-----

# Non technical land
## Problem

Navigating the journey to higher education is often fraught with significant challenges, and Compendium exists to address these pain points head-on. Many aspiring students face:

- **Overwhelming University Research**: The process of finding the right university is incredibly complex. Students struggle to gather comprehensive, up-to-date information on global institutions, their specialized programs, research opportunities, campus life, and alumni networks. This information is often scattered across countless websites, making it a time-consuming and frustrating endeavor.

- **Inefficient Shortlisting**: Sifting through hundreds or thousands of universities to create a tailored shortlist that aligns with personal academic goals, financial constraints, and career aspirations is a daunting task. Without personalized guidance, students often miss ideal matches or spend excessive time on unsuitable options.

- **Essay Writing Anxiety**: Crafting compelling and unique admission essays that stand out to admissions committees is a major hurdle. Many students lack the tools or support to effectively assess their essays, refine their arguments, and present their best selves.

- **Exam Preparation Stress & Cost**: Preparing for rigorous international exams (like SAT, GRE, GMAT, IELTS, TOEFL) often involves expensive courses, private tutors, and a lack of interactive, adaptable learning resources. This financial burden and limited accessibility to quality preparation can hinder a student's success.

- **High Costs of Existing Solutions**: While some businesses offer fragmented solutions for parts of this process, they often come with a hefty price tag, making comprehensive support inaccessible to a wide range of students.

Compendium consolidates these disparate needs into a single, affordable, and AI-powered platform. We provide the "shortest way" (as our name suggests) for students to confidently research, plan, and achieve their higher education goals, breaking down the barriers of complexity and cost.

## Growth Strategy

Compendium's growth strategy is focused on rapid scaling and market penetration. We will prioritize User Acquisition through aggressive digital marketing, SEO optimization, and strategic content creation to quickly build a substantial user base. Our Monetization Strategy will leverage a freemium model, offering premium AI features and advanced exam preparation courses to generate revenue while maintaining accessibility.

Crucially, we will actively engage with Startup Accelerators to refine our business model, gain mentorship, and access vital networks. This will directly support our efforts to Pitch and Raise Investments from venture capitalists and angel investors, securing the capital needed for aggressive expansion. Our Feature Expansion will be driven by user feedback and market demand, continuously enhancing our AI capabilities, integrating larger institutional databases, and fostering peer-to-peer networking. We will pursue Strategic Partnerships with educational institutions and Ed-Tech platforms to broaden our reach and integrate valuable resources.

Finally, Global Market Entry will be executed through targeted localization and marketing campaigns in key international student markets, ensuring Compendium becomes the leading platform for university research and admissions worldwide.

# Making changes

This section outlines the process for contributing changes to this GitHub repository. Following these guidelines ensures a smooth workflow, maintains code quality, and facilitates effective collaboration.

## Conventional Commits
We use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to standardize our commit messages. This provides a clear and concise history, automates changelog generation, and helps with semantic versioning.

A conventional commit message should follow this structure:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Common Types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc.)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies (e.g., npm, pip)
- `ci`: Changes to our CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files
- `revert`: Reverts a previous commit

Examples:
- `feat: add user authentication module`
- `fix(auth): correct password reset bug`
- `docs: update README with installation instructions`
- `refactor(api): simplify data fetching logic`

## Forking
While you can directly create branches in this private repository, for larger or more experimental changes, it's often cleaner to work from a fork. To fork the repository, click the "Fork" button in the top right corner on the main repository page. This will create a copy of the repository under your personal GitHub account.

Next, clone your forked repository to your local machine:

```bash
git clone https://github.com/YOUR_USERNAME/YOUR_REPOSITORY_NAME.git
cd YOUR_REPOSITORY_NAME
```

To easily pull updates, add the original repository as an `upstream` remote:

```bash
git remote add upstream https://github.com/ORGANIZATION_NAME/ORIGINAL_REPOSITORY_NAME.git
```

Finally, regularly pull changes from the upstream main (or master) branch to keep your fork up-to-date:

```bash
git checkout main
git pull upstream main
```

Always work on a new branch for your changes. This keeps the main branch clean and allows for easy review and integration. Before creating a new branch, ensure your local main branch is up-to-date with the remote repository:

```bash
git checkout main
git pull origin main # If working directly in the main repo

# OR
git pull upstream main # If working from a fork
```

Once your main branch is updated, create a new branch by choosing a descriptive name, typically prefixed with the type of change (e.g., feat/, fix/, docs/):

```bash
git checkout -b feat/add-user-profile
```

After creating the branch, implement your changes, commit them, and push them to your new branch.

## Making a Pull Request (PR)
Once your changes are complete and thoroughly tested on your branch, you can open a Pull Request to merge them into the main branch. First, ensure all your commits are pushed to your remote branch:

```bash
git push origin feat/add-user-profile # If working directly in the main repo

# OR
git push origin feat/add-user-profile # If working from a fork
```

To open a Pull Request, go to the GitHub repository page (either the original or your fork). GitHub will usually detect your newly pushed branch and prompt you to "Compare & pull request." Click this button, or alternatively, navigate to the "Pull requests" tab and click "New pull request."

When filling out the PR description, use a concise, descriptive title. In the description, provide a detailed explanation of your changes, including what problem the PR solves, how it solves it, any relevant context, design decisions, or trade-offs, screenshots or GIFs if applicable, and references to any related issues (e.g., Closes #123).

Request reviews from appropriate team members and add relevant labels or assign to a milestone if necessary. Be prepared to receive feedback and make further changes based on code reviews; pushing new commits to your branch will automatically update the PR. Once the PR has been approved by the required reviewers and all checks pass, it can be merged into the main branch.

# Microservice backend

## Setting up locally

```bash
git clone github.com/seacite-tech/compendium
cd user-service
go run cmd/main.go # run server
go test ./...      # test
```

## Go code Guidelines

### Logging

To ensure comprehensive logging with contextual information such as `requestId` and `userId`, we recommend enabling `common.pkg.middleware.LoggerMiddleware` globally. This middleware will instantiate and populate a reusable logger within `context.Context`, accessible via `common.pkg.log.L(ctx)`. This allows for consistent and enriched logging throughout the application.

```go
if university == nil {
    log.L(ctx).Errorf("University was not found")
    return ...
}
```

### Error tracing

Our code mostly uses [`tracerr`](https://github.com/ztrue/tracerr) for finding source of unexpected internal server errors (e.g. when database fails or external API doesn't respond). Errors must be wrapped into `tracerr.Error` (error with stacktrace) **only in the service dependency layer**.

This means `tracerr.Wrap()`, `tracerr.New()` and `tracerr.Errorf()` should only ever be applied at the point where an error originates from an external dependency (e.g., database, external API call, file system operation) and is first returned up the call stack. It should not be used repeatedly throughout the service or presentation layers.

```go
// In repository layer:
func (r *SomeRepo) GetItem(id string) error {
    err := db.Query("...")
    if err != nil {
        return tracerr.Wrap(fmt.Errorf("failed to query item %s from DB: %w", id, err))
    }
    return nil
}

// In service layer:
func (s *SomeService) ProcessItem(id string) error {
    err := s.repo.GetItem(id)
    if err != nil {
        return err
    }

    return nil
}
```

### Error handling

The presentation/API layer (e.g., HTTP handlers, gRPC servers) handles incoming requests, manages data serialization/deserialization, and invokes the service layer. Generic unexpected internal server errors originating from the service layer should be handled and logged at this presentation level, or by a dedicated error-handling middleware. This approach prevents code duplication across various handlers and the service layer itself.

Conversely, business logic-related errors (such as "not found" scenarios) should be returned as custom error types from the service layer. These custom errors should be logged with the contextual logger within the service layer itself, and then passed up to the presentation/API layer to be converted into an appropriate, user-friendly API response.

## Authentication Flow Overview
This part of the document outlines a typical authentication flow, covering user registration (sign-up) and user login (sign-in) scenarios, including considerations for multi-factor authentication (MFA) and token management.

### 1. Sign-Up
The sign-up process involves two primary steps: initial user registration and subsequent multi-factor authentication (MFA) verification.

#### 1.1 Initial User Registration
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

#### 1.2 Multi-Factor Authentication (MFA) Verification for Sign-Up
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

### 2. Sign-In
The sign-in process can vary depending on whether the user is logging in from the same device or a different device, particularly concerning MFA requirements.

#### 2.1 Sign-In from a Known Device (Password Authentication)
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
Set-Cookie: csrfToken=..., accessToken=..., refreshToken=...

{
  "isMfaRequired": false,
  "accessTokenExpiry": "...",
  "refreshTokenExpiry": "...",
}
```

#### 2.2 Sign-In from a New Device (MFA Required)
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
Set-Cookie: csrfToken=..., accessToken=..., refreshToken=...

{
  "accessTokenExpiry": "...",
  "refreshTokenExpiry": "...",
}
```

### 3. Token Management and Security Note
It's important to note the roles of the tokens and cookies in this flow:

- `accessToken` and `refreshToken`: These are critical for authentication and session management. They are stored as HTTP-only cookies, which means they cannot be accessed by client-side JavaScript. This significantly reduces the risk of Cross-Site Scripting (XSS) attacks stealing these tokens.
- `csrfToken`: This token is used to protect against Cross-Site Request Forgery (CSRF) attacks. While the `accessToken` and `refreshToken` are HTTP-only, the `csrfToken` is a non-HTTP-only cookie. This allows client-side JavaScript to read its value and include it in subsequent requests (e.g., in a custom header like `X-CSRF-Token`).

A key point regarding the `csrfToken`: Even if the `csrfToken` is also included as a claim within the `accessToken`, this design is acceptable. The primary protection for the `accessToken` and `refreshToken` comes from their HTTP-only nature, preventing client-side script access. The `csrfToken` serves its purpose by being readable by JavaScript for inclusion in requests, thereby validating that the request originated from the legitimate client application. So this is used in our current auth implementation in `user-service/`.

# Monolithic SPA frontend
Presentation/API Layer's (e.g., HTTP handlers, gRPC servers) purpose is to handle incoming requests, marshal/unmarshal data, and call the service layer. Generic unexpected internal server errors from the service layer should ideally be handled and logged at this level (or by a dedicated error handling middleware) to avoid code duplication across handlers and service layer. Business logic-related errors (e.g., "not found" scenarios) should be returned as custom error types from the service layer and then logged with the contextual logger and converted into an appropriate, user-friendly API response.
