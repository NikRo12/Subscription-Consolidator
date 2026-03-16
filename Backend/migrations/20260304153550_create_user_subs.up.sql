CREATE TABLE user_subs (
    user_id             INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sub_id              INT NOT NULL REFERENCES subs(id) ON DELETE CASCADE,
    price               NUMERIC(10, 2) NOT NULL,
    period              VARCHAR(255) NOT NULL,
	next_payment_date   DATE NOT NULL,
	is_active           BOOLEAN,       
    PRIMARY KEY         (user_id, sub_id)
);
