CREATE TABLE user_subs (
    user_id         INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sub_id INT NOT NULL REFERENCES subs(id) ON DELETE CASCADE,
    start_at   TIMESTAMP DEFAULT NOW(),
    end_at   TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, sub_id)
);