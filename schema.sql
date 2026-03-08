CREATE TYPE effort AS ENUM (
    'mins',
    'hours',
    'days',
    'weeks',
    'months',
);

CREATE TABLE todos (
    id SERIAL PRIMARY KEY,

    short TEXT NOT NULL,
    description TEXT,

    due_date DATE,

    cost_of_delay SMALLINT
        CHECK (cost_of_delay BETWEEN -2 AND 2)
        DEFAULT 0,

    effort effort DEFAULT 'hours',

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);