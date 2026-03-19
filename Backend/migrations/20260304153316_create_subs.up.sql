CREATE TABLE subs (
    id          SERIAL PRIMARY KEY,
	title       VARCHAR(255) UNIQUE NOT NULL,
	currency    VARCHAR(255) NOT NULL,
	category    VARCHAR(255) NOT NULL DEFAULT '',
	icon_url     VARCHAR(255) NOT NULL DEFAULT '',
	brand_color  VARCHAR(255) NOT NULL DEFAULT '',
	description VARCHAR(255) NOT NULL DEFAULT ''
);