ALTER TABLE user_subs
    ALTER COLUMN period            SET DEFAULT '',
    ALTER COLUMN next_payment_date DROP NOT NULL;
