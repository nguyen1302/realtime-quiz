CREATE TABLE IF NOT EXISTS answers (
    id UUID PRIMARY KEY,
    quiz_id UUID NOT NULL,
    question_id UUID NOT NULL,
    user_id UUID NOT NULL,
    answer TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    points INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id),
    FOREIGN KEY (question_id) REFERENCES questions(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_answers_quiz_id ON answers(quiz_id);
CREATE INDEX idx_answers_question_id ON answers(question_id);
CREATE INDEX idx_answers_user_id ON answers(user_id);
