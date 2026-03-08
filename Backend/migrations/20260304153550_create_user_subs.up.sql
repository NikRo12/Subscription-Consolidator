CREATE TABLE user_subs (
    user_id             INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sub_id              INT NOT NULL REFERENCES subs(id) ON DELETE CASCADE,
    period              VARCHAR(255) NOT NULL,
    price               NUMERIC(10, 2) NOT NULL,
	next_payment_date   DATE NOT NULL,
	is_active           BOOLEAN,       
    PRIMARY KEY         (user_id, sub_id)
);
