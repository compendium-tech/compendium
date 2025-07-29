-- Custom types --

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'activity_category') THEN
        CREATE TYPE activity_category AS ENUM (
            'academic', 'art', 'athletics', 'career_oriented', 'community_service',
            'cultural', 'debate/speech', 'environmental', 'family_responsibilities',
            'journalism_publication', 'music', 'religious', 'research', 'robotics',
            'school_spirit', 'student_government', 'theatre_drama', 'work', 'other'
        );
    END IF;

    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'grade') THEN
        CREATE TYPE grade AS ENUM ('9', '10', '11', '12', 'post_graduate');
    END IF;

    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'honor_level') THEN
        CREATE TYPE honor_level AS ENUM ('school', 'regional', 'national', 'international');
    END IF;
END $$;

-- Tables --

CREATE TABLE IF NOT EXISTS applications (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS activities (
  application_id UUID PRIMARY KEY REFERENCES applications (id),
  name TEXT NOT NULL,
  role TEXT NOT NULL,
  description TEXT,
  hours_per_week INTEGER NOT NULL,
  weeks_per_year INTEGER NOT NULL,
  category activity_category NOT NULL,
  grades grade[] NOT NULL
);

CREATE TABLE IF NOT EXISTS honors (
  application_id UUID PRIMARY KEY REFERENCES applications (id),
  title TEXT NOT NULL,
  description TEXT,
  level honor_level NOT NULL,
  grade grade NOT NULL
);
