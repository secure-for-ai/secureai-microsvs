-- DROP SCHEMA test;

CREATE SCHEMA test;

-- test.test_user definition

-- Drop table

-- DROP TABLE test.test_user;

CREATE TABLE test.test_user (
	uid int8 NOT NULL,
	username varchar(100) NULL,
	"password" varchar(100) NOT NULL,
	CONSTRAINT test_user_pk PRIMARY KEY (uid)
);

-- Insert test data
INSERT INTO test.test_user (uid, username, "password") VALUES(0, 'user1', 'password');
INSERT INTO test.test_user (uid, username, "password") VALUES(1, 'hello world', 'p@ssword');
INSERT INTO test.test_user (uid, username, "password") VALUES(2, 'user2', 'test_pass');
