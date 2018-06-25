DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS restaurants;
DROP TABLE IF EXISTS visits;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS pictures;
DROP TABLE IF EXISTS settings;

CREATE TABLE users (
'username' VARCHAR(64) NULL,
'password' VARCHAR(64) NULL,
'token' VARCHAR(64) NULL,
'profilepicture' VARCHAR(64) NULL,
'firstname' VARCHAR(64) NULL,
'lastname' VARCHAR(64) NULL,
'address' VARCHAR(64) NULL,
'email' VARCHAR(64) NULL,
'id' VARCHAR(64) NULL,
'created' DATE NULL
);

CREATE TABLE events (
'category' VARCHAR(64) NULL,
'text' STRING NULL,
'time' int64 NULL
);

CREATE TABLE restaurants (
'id' VARCHAR(64) NULL,
'url' string NULL,
'addedby' VARCHAR(64) NULL,
'name' VARCHAR(64) NULL,
'addressid' VARCHAR(64) NULL
);

CREATE TABLE visits (
'id' VARCHAR(64) NULL,
'userid' VARCHAR(64) NOT NULL,
'restaurantid' string NOT NULL,
'time' int64 NULL,
'timecreated' int64 NULL,
'with' string NULL,
'comment' string NULL,
'ispublic' integer NOT NULL
);

CREATE TABLE addresses (
'street' string NULL,
'city' string NULL,
'state' string NULL,
'country' string NULL
);

CREATE TABLE pictures (
'visitid' string NULL,
'comment' string NULL,
'width' float64 NULL,
'height' float64 NULL,
'filename' string NOT NULL
);

CREATE TABLE settings (
'userid' VARCHAR(64) NOT NULL,
'setting' string NULL,
'value' string NULL
);



