<img width="1640" height="1014" alt="image" src="https://github.com/user-attachments/assets/ca0dd747-8f15-468c-a748-de92a73d8ad6" />


# Table of Contents

-   **[Welcome!](#welcome)**
-   **[Making changes](#making-changes)**
-   **[System Architecture](#system-architecture)**
    - [Port Usage](#port-usage)
    -   [Setup guide](#setup-guide)
        -   [Paddle](#paddle), [Kafka](#kafka), [Nginx](#nginx)
        -   [Protobuf](#protobuf), [Mockery](#mockery)
-   **[Microservice Backend](#microservice-backend)**
    -   [Running Locally](#running-locally)
    -   [Go Code Guidelines](#go-code-guidelines)
        -   [Logging](#logging)
        -   [Error Handling](#error-handling)
-   **[Monolithic SPA Frontend](#monolithic-spa-frontend)**
    -   [Running Locally](#running-locally-1)

-----

# Welcome!

**Think about it**: applying to university is a huge deal, right? But for so many aspiring students, it's a total maze. They're trying to find the perfect university, figure out what research opportunities exist and then get tangled in the mess of admissions criteria. Information is scattered everywhere, making it a super time-consuming and frustrating experience.

And then there's the shortlisting – trying to narrow down hundreds of universities to just a few that actually fit their goals, budget, and dreams? It's daunting. Many students end up missing out on ideal matches or just wasting tons of time on places that aren't right for them.

Then, the dreaded essays. How do you make your story stand out to an admissions committee that reads thousands of them? Most students don't have the tools or support to really refine their arguments and present their absolute best selves.

And let's not forget the international exams like SAT, GRE, IELTS. Prep for these often means expensive courses, private tutors, and a real lack of interactive learning. It's a huge financial burden and a barrier to success for many.

Basically, the existing solutions out there are fragmented and come with a hefty price tag, making comprehensive support inaccessible.

This is where Compendium steps in. We're bringing all these pieces together into one affordable, AI-powered platform. Our goal is to offer the "shortest way" for students to confidently research, plan, and achieve their higher education dreams, completely breaking down those walls of complexity and cost.

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

To open a PR, head to the GitHub repo page (either the original or your fork). GitHub usually spots your new branch and offers a "Compare & pull request" button. Click that, or just go to the "Pull requests" tab and hit "New pull request."

When you're filling out the PR, give it a short, clear title. In the description, really explain what you did: what problem does this PR fix? How does it fix it? Any context, design choices, or trade-offs? Screenshots or GIFs are always a plus! And if it relates to an issue, link it (e.g., Closes #123).

Ask the right folks to review your code, and add any labels or milestones if needed. Be ready for feedback – it's part of the process! Pushing new commits to your branch will automatically update the PR. Once enough people have approved and all our automated checks pass, your changes can be merged!

# System Architecture

We leverage a mix of programming languages, databases, caching mechanisms, external APIs, and cloud services:
- **Go**: Chosen for the core microservices (application service, course service, user service, college service, subscription service, email delivery service).
- **Vite, Vue 3**: Provide a modern and efficient development environment for our Frontend SPA, allowing for rapid development and a highly interactive user experience. 
- **PostgreSQL**: Serves as the primary persistent data store for various application domains: Course data, Application data, User data, and Subscription data.
- **Redis**: Employed for high-speed data caching to manage authentication state and caching in some microservices.
- [**Paddle**](https://paddle.com/): An external platform integrated for handling Subscription management and payment processing.
- **Elasticsearch**: Utilized for the College database, specifically for its powerful search and analytics capabilities. This allows for efficient querying and retrieval of college-related information.
- **Nginx**: Acts as a high-performance reverse proxy and load balancer.
- **AWS S3, CloudFront**: Solution for static file storage and content delivery used in course microservice.

## Port Usage

When you run everything on your local machine, here's which port each service hangs out on (by default):

| Service                 | Protocol | Port    | Description                                   |
| :---------------------- | :------- | :------ | :-------------------------------------------- |
| **User Service** | HTTP     | `1000`  | Handles user-related operations over HTTP.    |
| **User Service** | gRPC     | `2000`  | Provides user-related operations over gRPC.   |
| **Subscription Service**| HTTP     | `1001`  | Manages user subscriptions.                   |
| **API Gateway / Nginx** | HTTP     | `8080`  | Entry point for all external requests.        |
| **Vue Application** | HTTP     | `5173`  | Frontend application serving the UI.          |
| **Redis** | TCP     | `6379`  | Cache storage.          |
| **PostgreSQL** | TCP     | `5432`  | For testing purposes one db for all microservices.          |
| **Kafka** | TCP     | `9092`  | For email delivery with corresponding service.          |

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
  Download the latest release from [here](https://kafka.apache.org/downloads).
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

### Nginx

To get Nginx running with our configuration, simply point it to the `nginx/api.conf` file:
```
nginx -c /absolute/path/to/compendium/nginx/api.conf
```

### Mockery

Mockery generates mocks for Go interfaces, allowing us to test individual code components in isolation. This makes our unit tests faster, more reliable, and independent of real external services or complex dependencies. To install Mockery run:

```bash
go install github.com/vektra/mockery/v3@v3.5.1
```

If you ever change a Go interface and need to regenerate the mocks, just run this command from the root of the repository:

```bash
mockery
```

### Protobuf

You can install protobuf [here](https://github.com/protocolbuffers/protobuf/releases). To install protobuf plugin for Go use:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

To regenerate interfaces from our .proto files, run these commands from the _root of the repository_:

```bash
export PATH="$PATH:$(go env GOPATH)"/bin

make protoc
```

# Microservice backend

## Running locally

```bash
git clone github.com/compendium-tech/compendium

# Email delivery service
cd ./email-delivery-service
go run cmd/main.go # run kafka consumer

# User service
cd ../user-service
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
import (
    "github.com/compendium-tech/compendium/common/pkg/log"
)

func (s *service) GetUniversity() (*domain.University, error) {
    ...
    logger := log.L(ctx).WithField("universityId", universityID)
    
    if university == nil {
        logger.Errorf("University was not found")
        return ...
    }
    ...
}
```

### Error handling

Our code mostly uses [`tracerr`](https://github.com/ztrue/tracerr) for finding source of unexpected internal server errors (e.g. when database fails or external API doesn't respond). Errors must be wrapped into `tracerr.Error` (error with stacktrace) **only in the service dependency layer**.

This means `tracerr.Wrap()`, `tracerr.New()` and `tracerr.Errorf()` should only ever be applied at the point where an error originates from an external dependency (e.g., database, external API call, file system operation) and is first returned up the call stack. It should not be used repeatedly throughout the service or presentation layers.

```go
import (
    "github.com/ztrue/tracerr"
)

func (r *SomeRepo) GetItem(id string) error {
    err := db.Query("...")
    if err != nil {
        return tracerr.Errorf("failed to query item %s from DB: %v", id, err)
    }

    return nil
}

func (s *SomeService) ProcessItem(id string) error {
    err := s.repo.GetItem(id)
    if err != nil {
        return err
    }

    return nil
}
```

API layer (think HTTP handlers or gRPC servers) handles turning data into something our code can understand, and then calls our service layer. If there's a generic, unexpected internal server error from the service layer, it should be caught and logged at this presentation level, or by a special error-handling middleware. This way, we avoid writing the same error-handling code everywhere.

On the flip side, errors related to our business logic (like, "hey, we couldn't find that user!") should be sent back as custom error types from the service layer. These custom errors should be logged with our handy contextual logger within the service layer itself. Then, they get passed up to the presentation/API layer, which can turn them into a nice, polite message for the user.

# Monolithic SPA frontend

## Running Locally

Want to see the app in action on your machine? Run:

```
git clone github.com/compendium-tech/compendium
cd frontend-spa
npm install
npm run dev
```
