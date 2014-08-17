CREATE TABLE IF NOT EXISTS txns (
	id serial,
	hash varchar(128),
	date date,
	description varchar(255),
	debit_cents integer,
	credit_cents integer,
	balance_cents integer,
	PRIMARY KEY (id)
);

CREATE UNIQUE INDEX txns_idx_hash ON txns (hash);
