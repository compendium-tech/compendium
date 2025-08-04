DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'tier' AND typnamespace = 'public'::regnamespace) THEN
        CREATE TYPE tier AS ENUM (
            'student',
            'team',
            'community'
        );
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS subscriptions (
  id TEXT PRIMARY KEY,
  backed_by UUID UNIQUE NOT NULL,
  tier tier,
  invitation_code TEXT,
  till TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  since TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subscription_members (
  user_id UUID NOT NULL,
  subscription_id TEXT NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
  since TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  PRIMARY KEY (subscription_id, user_id)
);
