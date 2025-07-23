CREATE TYPE subscription_level AS ENUM (
  'student',
  'team',
  'community'
);

CREATE TABLE IF NOT EXISTS subscriptions (
  id TEXT UNIQUE NOT NULL,
  backed_by UUID UNIQUE NOT NULL,
  subscription_level subscription_level,
  till TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  since TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subscription_members (
  user_id UUID UNIQUE NOT NULL,
  subscription_id TEXT NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
  since TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
);
