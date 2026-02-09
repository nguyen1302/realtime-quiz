CREATE TABLE IF NOT EXISTS questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    options JSONB NOT NULL, -- Array of strings
    correct_answer VARCHAR(255) NOT NULL,
    time_limit INTEGER NOT NULL DEFAULT 30, -- Seconds
    points INTEGER NOT NULL DEFAULT 100,
    item_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_questions_quiz_id ON questions(quiz_id);
