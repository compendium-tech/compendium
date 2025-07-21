CREATE TYPE subscription_level AS ENUM (
  'student',
  'team',
  'community'
);

CREATE TABLE IF NOT EXISTS subscriptions (
  user_id UUID UNIQUE NOT NULL,
  subscription_level subscription_level,
  till TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  since TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
