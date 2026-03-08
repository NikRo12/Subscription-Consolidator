CREATE TABLE subs (
    id          SERIAL PRIMARY KEY,
	title       VARCHAR(255) NOT NULL,
	currency    VARCHAR(255) NOT NULL,
	category    VARCHAR(255) NOT NULL,
	icon_url     VARCHAR(255) NOT NULL,
	brand_color  VARCHAR(255) NOT NULL,
	description VARCHAR(255) NOT NULL 
);