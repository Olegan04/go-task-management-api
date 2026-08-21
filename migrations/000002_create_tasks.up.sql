CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    discription TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    create_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    update_ap TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);