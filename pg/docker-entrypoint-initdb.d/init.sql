create role lms_owner;

grant all privileges on database lms to lms_owner;

create schema draw;
grant all privileges on schema draw to lms_owner;
alter default privileges in schema draw grant all privileges on tables to lms_owner;

create schema payment;
grant all privileges on schema payment to lms_owner;
alter default privileges in schema payment grant all privileges on tables to lms_owner;

create schema ticket;
grant all privileges on schema ticket to lms_owner;
alter default privileges in schema ticket grant all privileges on tables to lms_owner;

create schema consumer;
grant all privileges on schema consumer to lms_owner;
alter default privileges in schema consumer grant all privileges on tables to lms_owner;