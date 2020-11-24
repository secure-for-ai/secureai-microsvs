-- DROP SCHEMA test;

CREATE SCHEMA test;

ALTER SCHEMA test owner to test;

-- test.testUser definition

-- Drop table

-- DROP TABLE test.test_user;

CREATE TABLE test.test_user
(
	uid bigint not null
		CONSTRAINT test_user_pk
			PRIMARY KEY,
	username varchar(100),
	password varchar(100) not null
);

ALTER TABLE test.test_user owner to test;

-- Insert test data
INSERT INTO test.test_user (uid, username, "password") VALUES (0, 'user1', 'password');
INSERT INTO test.test_user (uid, username, "password") VALUES (1, 'hello world', 'p@ssword');
INSERT INTO test.test_user (uid, username, "password") VALUES (2, 'user2', 'test_pass');

-- test.test_session definition

-- Drop table

-- DROP TABLE test.test_session;

CREATE TABLE test.test_session
(
	uid bigint,
	sid bigint,
	nonce bytea,
	data jsonb,
	ip inet,
	user_agent varchar,
	create_time bigint,
	update_time bigint,
	expire_time bigint,
	CONSTRAINT test_session_pk
		UNIQUE (uid, sid)
);

ALTER TABLE test.test_session owner to test;

-- DROP SCHEMA demo_pg;

CREATE SCHEMA demo_pg;

ALTER SCHEMA demo_pg owner to test;

-- demo_pg.session definition

-- Drop table

-- DROP TABLE demo_pg.session;

CREATE TABLE demo_pg.session
(
	uid bigint,
	sid bigint,
	nonce bytea,
	data jsonb,
	ip inet,
	user_agent varchar,
	create_time bigint,
	update_time bigint,
	expire_time bigint,
	CONSTRAINT test_session_pk
		UNIQUE (uid, sid)
);

ALTER TABLE demo_pg.session owner to test;

-- demo_pg.user definition

-- Drop table

-- DROP TABLE demo_pg.user;

CREATE TABLE demo_pg.user
(
	uid bigint not null
		CONSTRAINT user_pk_uid
			PRIMARY KEY,
	username varchar(100)
		CONSTRAINT user_pk_username
			UNIQUE,
	nickname varchar(100),
	email varchar(100),
	create_time bigint,
	update_time bigint

);

ALTER TABLE demo_pg.user owner to test;
