CREATE TYPE tier AS ENUM (
  'student',
  'team',
  'community'
);

CREATE TABLE IF NOT EXISTS subscriptions (
  id TEXT PRIMARY KEY,
  backed_by UUID UNIQUE NOT NULL,
  tier tier,
  invitation_code TEXT,
  till TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  since TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subscription_members (
  user_id UUID UNIQUE NOT NULL,
  subscription_id TEXT NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
  since TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
