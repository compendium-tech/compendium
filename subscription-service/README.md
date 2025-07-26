# Subscription Service

Hey there! This Subscription service is basically the brain behind managing all our user subscriptions.  
At its core, this service is responsible for:
- **Keeping track of subscriptions**: It knows exactly who is subscribed, what tier they're on (e.g., Student, Team).
- **Talking to Paddle**: It's got a direct line to Paddle, our payment processor. This means it can initiate new checkouts, tell Paddle to cancel old subscriptions, and receive updates (via webhooks) when something important happens on Paddle's side (like a successful payment or a cancellation).
- **Controlling access**: It works hand-in-hand with other services (like Service X) via gRPC to make sure that only users with active subscriptions can access certain features or operations.
- **Managing subscription changes**: Need to upgrade? Downgrade? Add more team members in collective subscriptions? This service handles all those changes.

## How does it work with other parts of our system?

<img width="1260" height="871" alt="Снимок экрана 2025-07-26 в 13 08 20" src="https://github.com/user-attachments/assets/20c52047-8eef-4648-9127-7c4cb01432e1" />

You can see a few other players in the diagram, and the Subscription service interacts with them quite a bit:
- **User**: When they're checking out to buy a subscription, they'll go through Paddle.js which then talks to Paddle, and Paddle talks back to our Subscription service via a webhook.
- **User service**: This service likely handles core user data. The Subscription service talks to it to get account details via gRPC, to show them when fetching collective subscription members.
- **Paddle (our payment service)**: As mentioned, Paddle is where the actual money transactions processing happen. The Subscription service listens to events from Paddle and adds changes to its subscriptions storage.
- **Subscriptions storage & cache**: This is where all the nitty-gritty details of subscriptions are kept. There's a main Subscriptions storage (think of it as the authoritative record) and a Subscriptions storage cache which is there to make things super fast when we need to check a user's subscription status frequently.
- **Service X**: This is an example of another service in our system (could be anything, like a feature access gate, an API, etc.).
Service X asks the Subscription service (via gRPC) if a user has an active subscription and what tier they're on before allowing them to do something.
