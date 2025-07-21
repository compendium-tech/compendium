<img width="1640" height="1014" alt="image" src="https://github.com/user-attachments/assets/ca0dd747-8f15-468c-a748-de92a73d8ad6" />


# Table of Contents

-   **[Welcome!](#welcome)**
    -   [Why Compendium?](#why-compendium)
    -   [Our Growth Vision](#our-growth-vision)
-   **[Why Open Source? This is the Heart of It.](#why-open-source-this-is-the-heart-of-it)**
-   **[Making changes](#making-changes)**
-   **[System Architecture](#system-architecture)**
    -   [Setup guide](#setup-guide)
        -   [Paddle](#paddle)
        -   [Kafka](#kafka)
-   **[Microservice Backend](#microservice-backend)**
    -   [Running Locally](#running-locally)
    -   [Go Code Guidelines](#go-code-guidelines)
        -   [Logging](#logging)
        -   [Error Tracing](#error-tracing)
        -   [Error Handling](#error-handling)
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
-   **[Monolithic SPA Frontend](#monolithic-spa-frontend)**
    -   [Running Locally](#running-locally-1)

-----

# Welcome!
## Why Compendium?

**Think about it**: applying to university is a huge deal, right? But for so many aspiring students, it's a total maze. They're trying to find the perfect university, figure out what research opportunities exist and then get tangled in the mess of admissions criteria. Information is scattered everywhere, making it a super time-consuming and frustrating experience.

And then there's the shortlisting – trying to narrow down hundreds of universities to just a few that actually fit their goals, budget, and dreams? It's daunting. Many students end up missing out on ideal matches or just wasting tons of time on places that aren't right for them.

Then, the dreaded essays. How do you make your story stand out to an admissions committee that reads thousands of them? Most students don't have the tools or support to really refine their arguments and present their absolute best selves.

And let's not forget the international exams like SAT, GRE, IELTS. Prep for these often means expensive courses, private tutors, and a real lack of interactive learning. It's a huge financial burden and a barrier to success for many.

Basically, the existing solutions out there are fragmented and come with a hefty price tag, making comprehensive support inaccessible.

This is where Compendium steps in. We're bringing all these pieces together into one affordable, AI-powered platform. Our goal is to offer the "shortest way" for students to confidently research, plan, and achieve their higher education dreams, completely breaking down those walls of complexity and cost.

## Our Growth Vision

We're not just building a product; we're building a movement. Our growth strategy is all about getting Compendium into as many hands as possible, fast.

We're going all-in on **User Acquisition** with aggressive digital marketing, SEO, and creating killer content that genuinely helps students. We want to build a massive user base quickly because, frankly, the more students we help, the more impact we have.

When it comes to making money, we're using a **Freemium Model**. This means we'll offer a ton of value for free, but then have premium AI features and advanced exam prep courses for those who want to take their journey to the next level. This keeps Compendium accessible while still generating the revenue we need to keep innovating.

<!-- We're actively working with Startup Accelerators. This isn't just about money; it's about refining our entire business model, getting priceless mentorship, and tapping into networks that can truly propel us forward. This ties directly into our efforts to Pitch and Raise Investments from VCs and angel investors – securing the capital to really expand aggressively. -->

Of course, Strategic Partnerships are key. Teaming up with educational institutions and other Ed-Tech platforms will broaden our reach and bring even more valuable resources to our users.

Finally, we're thinking Global Market Entry from day one. This means targeted localization and marketing in key international student markets to ensure Compendium becomes the go-to platform for university research and admissions worldwide. We want to be there for every student, no matter where they are.

# Why Open Source? This is the Heart of It.

This is where the human element really shines through for us. Making Compendium open source isn't just a technical decision; it's deeply tied to our values and our business strategy.

First off, it's about community collaboration. We genuinely believe that by opening up our project, we tap into a global brain trust. Imagine developers, educators, and even students from all over the world contributing their ideas, their code, and their insights. This leads to more creative solutions, higher quality code, and a platform that truly serves a diverse user base. It's like having an army of passionate people helping you build something incredible.

Then there's transparency. In an era where trust is paramount, being open source means our codebase is public. There are no hidden agendas, no secret practices. Users can see exactly how their data is handled, how features work, and how we're evolving. This builds immense trust, which is invaluable for a platform dealing with something as critical as education.

Open sourcing also dramatically accelerates innovation. We're not limited by the resources of just our internal team. Contributions can come in from anywhere, leading to faster development cycles and quicker iterations on features. It means we can adapt to market needs and user feedback with incredible agility.

Ultimately, it helps us build a more robust and accessible platform. When more eyes are on the code, bugs are found faster, security vulnerabilities are patched sooner, and the overall quality of the software improves. And by being open, Compendium can be more easily adapted and integrated with other educational tools, creating a richer ecosystem for students.

In essence, we're open sourcing Compendium because we believe in the power of collective effort. It's about building a collaborative ecosystem around higher education guidance, making quality resources and tools not just accessible, but continuously improving, for everyone. We're inviting you to be a part of that journey.

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

# System Architecture

<img width="80%" alt="architecture" src="https://github.com/user-attachments/assets/bd038fe6-f1ef-4631-bc3e-9ae8e8bf9e3c" />

This section outlines the setup process for the various components of our system architecture. We leverage a mix of programming languages, databases, caching mechanisms, external APIs, and cloud services to build a robust and scalable application.
- **Python**: Primarily used for specialized microservice services, including interactions with the LLM API for the assistant service.
- **Go**: Chosen for the core microservices (application service, course service, user service, college service, subscription service, email delivery service).
- **Vite, Vue 3**: Provide a modern and efficient development environment for our Frontend SPA, allowing for rapid development and a highly interactive user experience. Vite's speed and Vue 3's reactive framework are key here.
- **PostgreSQL**: Serves as the primary persistent data store for various application domains: Course data, Application data, User data, and Subscription data. 
- **Redis**: Employed for high-speed data caching to manage authentication state and caching in some microservices.
- [**Paddle**](https://paddle.com/): An external platform integrated for handling Subscription management and payment processing.
- **Elasticsearch**: Utilized for the College database, specifically for its powerful search and analytics capabilities. This allows for efficient querying and retrieval of college-related information.
- **Nginx**: Acts as a high-performance reverse proxy and load balancer.
- **AWS S3, CloudFront**: Solution for static file storage and content delivery used in course microservice. 

## Setup Guide

### Paddle

This section will guide you through setting up your Paddle sandbox environment and configuring the subscription service.

If you don't already have one, you'll need a Paddle sandbox account for development and testing.

We need to create three specific subscription prices in your Paddle sandbox, which will represent our different subscription tiers.

- **Navigate to Catalog**: In your Paddle sandbox dashboard, go to Catalog > Products.
- **Create Products**:
  - Click New Product.
  - Create a product named "Student Subscription".
  - Create a product named "Team Subscription".
  - Create a product named "Community Subscription".
  - Create Prices for Each Product.
- **Configure Subscription Service**:
    The subscription service relies on environment variables to connect to Paddle and identify the correct products.
    Navigate to the `subscription-service/` directory:
    ```bash
    cd subscription-service/
    ```
    
    Add product IDs to `.env`: Open the `.env` file and add the following lines, replacing the placeholder values with the actual product IDs you copied from your Paddle sandbox.
    ```env
    PADDLE_STUDENT_SUBSCRIPTION_PRODUCT_ID=pro_...
    PADDLE_TEAM_SUBSCRIPTION_PRODUCT_ID=pro_...
    PADDLE_COMMUNITY_SUBSCRIPTION_PRODUCT_ID=pro_...
    ```
    
    Your subscription service is now configured with the correct Paddle product IDs and ready for development and testing!

### Kafka

This guide provides quick setup instructions for Apache Kafka, suitable for an email delivery service, using three common methods: manual download, Docker, and Homebrew.

### Prerequisites
One of these:
- Java 17+.
- Docker Desktop (for Docker method).

### Manual Download

- **Download Kafka Binary**:
  Download the latest release from https://kafka.apache.org/downloads.
- **Extract & Navigate**:
  ```bash
  tar -xzf kafka_....tgz # Use your version
  cd kafka_...
  ```
- **Generate Cluster ID & Format Logs**:
  ```bash
  KAFKA_CLUSTER_ID="$(bin/kafka-storage.sh random-uuid)"
  bin/kafka-storage.sh format --standalone -t $KAFKA_CLUSTER_ID -c config/server.properties
  ```
- **Start Kafka Server**:
  ```bash
  bin/kafka-server-start.sh config/server.properties
  ```
- **Create Topic**:
  ```bash
  bin/kafka-topics.sh --create --topic private.emaildelivery.emails --bootstrap-server localhost:9092
  ```

### Using Docker

- **Get the Docker image**:

  ```bash
  docker pull apache/kafka:4.0.0
  ```

- **Start the Kafka Docker container**:

  ```bash
  docker run -p 9092:9092 apache/kafka:4.0.0
  ```

- **Create Topic**:
  ```bash
  docker exec $(docker ps -q --filter ancestor=apache/kafka:4.0.0) /opt/kafka/bin/kafka-topics.sh --create --topic private.emaildelivery.emails --bootstrap-server localhost:9092
  ```

### Homebrew (macOS)

- **Install Zookeeper (required by Kafka)**:
  ```bash
  brew install zookeeper
  ```

- **Install Kafka**:
  ```bash
  brew install kafka
  ```

  This will install Kafka along with its dependencies, including Zookeeper.

- **Start Zookeeper**:
  Kafka depends on Zookeeper for distributed coordination. You need to start Zookeeper first. In your terminal, run:
  ```bash
  zkServer start
  ```
- **Start Kafka Server**:
  Now that Zookeeper is running, we need to start the Kafka server. In a new terminal window, run the following command:
  ```bash
  kafka-server-start /usr/local/etc/kafka/server.properties
  ```
  
  This starts the Kafka server on the default port (9092).

- **Create a Kafka Topic**
  ```bash
  kafka-topics --create --topic private.emaildelivery.emails --bootstrap-server localhost:9092
  ```

# Microservice backend

## Running locally

```bash
git clone github.com/compendium-tech/compendium

# User service
cd ./user-service
go test ./...      # test
go run cmd/main.go # run server

# Subscription service
cd ../subscription-service
go test ./...      # test
go run cmd/main.go # run server
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

# Monolithic SPA frontend

## Running Locally

```
git clone github.com/compendium-tech/compendium
cd frontend-spa
npm install
npm run dev
```
