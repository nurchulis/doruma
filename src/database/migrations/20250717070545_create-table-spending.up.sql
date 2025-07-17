CREATE TABLE spendings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_session_id UUID NOT NULL,
    category VARCHAR(255) NOT NULL,
    category_id UUID NULL,
    name VARCHAR(255) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    description TEXT,
    datetime TIMESTAMP WITH TIME ZONE NOT NULL,
    is_confirm BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
