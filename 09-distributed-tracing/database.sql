CREATE DATABASE IF NOT EXISTS dtdemo;

USE dtdemo;

CREATE TABLE IF NOT EXISTS dtdemo.people (
    name        VARCHAR(100),
    title       VARCHAR(10),
    description VARCHAR(100),
    PRIMARY KEY (name)
);

DELETE FROM dtdemo.people;

INSERT INTO dtdemo.people VALUES ('Gru', 'Felonius', 'Where are the minions?');
INSERT INTO dtdemo.people VALUES ('Nefario', 'Dr.', 'Why ... why are you so old?');
INSERT INTO dtdemo.people VALUES ('Agnes', '', 'Your unicorn is so fluffy!');
INSERT INTO dtdemo.people VALUES ('Edith', '', "Don't touch anything!");
INSERT INTO dtdemo.people VALUES ('Vector', '', 'Committing crimes with both direction and magnitude!');
INSERT INTO dtdemo.people VALUES ('Dave', 'Minion', 'Ngaaahaaa! Patalaki patalaku Big Boss!!');
