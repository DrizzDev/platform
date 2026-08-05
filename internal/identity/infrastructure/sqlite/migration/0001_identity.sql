CREATE TABLE epoch (
	id    INTEGER NOT NULL PRIMARY KEY CHECK (id = 1),
	value INTEGER NOT NULL CHECK (value >= 0)
) STRICT;

INSERT INTO epoch (id, value) VALUES (1, 0);

CREATE TABLE pointer (
	session  TEXT NOT NULL PRIMARY KEY,
	name     TEXT NOT NULL,
	revision INTEGER NOT NULL CHECK (revision > 0),
	epoch    INTEGER NOT NULL CHECK (epoch >= 0)
) STRICT;

CREATE TABLE attempt (
	session  TEXT NOT NULL,
	revision INTEGER NOT NULL CHECK (revision > 0),
	epoch    INTEGER NOT NULL CHECK (epoch >= 0),
	PRIMARY KEY (session, revision)
) STRICT;

CREATE TABLE lease (
	session TEXT NOT NULL PRIMARY KEY,
	owner   TEXT NOT NULL,
	expiry  INTEGER NOT NULL CHECK (expiry > 0)
) STRICT;

CREATE TABLE cleanup (
	name     TEXT NOT NULL PRIMARY KEY,
	reason   TEXT NOT NULL CHECK (reason IN ('LOGOUT', 'REJECTED', 'SUPERSEDED')),
	state    TEXT NOT NULL CHECK (state IN ('PENDING', 'BLOCKED')),
	attempts INTEGER NOT NULL CHECK (attempts >= 0),
	next     INTEGER NOT NULL CHECK (next >= 0),
	deadline INTEGER NOT NULL CHECK (deadline >= 0),
	created  INTEGER NOT NULL CHECK (created >= 0)
) STRICT;

CREATE INDEX due ON cleanup (state, next, created);
