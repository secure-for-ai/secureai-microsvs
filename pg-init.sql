-- DROP SCHEMA test;

CREATE SCHEMA test;

-- test.testUser definition

-- Drop table

-- DROP TABLE test.testUser;

CREATE TABLE test.test_user (
	uid int8 NOT NULL,
	username varchar(100) NULL,
	"password" varchar(100) NOT NULL,
	CONSTRAINT test_user_pk PRIMARY KEY (uid)
);

-- Insert test data
INSERT INTO test.testUser (uid, username, "password") VALUES(0, 'user1', 'password');
INSERT INTO test.testUser (uid, username, "password") VALUES(1, 'hello world', 'p@ssword');
INSERT INTO test.testUser (uid, username, "password") VALUES(2, 'user2', 'test_pass');

create table "test_session"
(
	uid int,
	sid int,
	nonce bytea,
	data jsonb,
	ip inet,
	"userAgent" varchar,
	"createTime" int,
	"updateTime" int,
	"expireTime" int,
	constraint testsession_pk
		unique (uid, sid)
);
