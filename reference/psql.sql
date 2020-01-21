CREATE DATABASE snippetbox ENCODING 'UTF8';

CREATE TABLE snippets (
  id SERIAL NOT NULL PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created TIMESTAMP NOT NULL,
  expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets(title, content, created, expires) VALUES (
    'My first snippet',
    'This is a fabulous life that one lives',
    now()::timestamp,
    now()::timestamp + '1 year'::interval
);

INSERT INTO snippets(title, content, created, expires) VALUES (
    'My second snippet',
    'Live life to the fullest',
    now()::timestamp,
    now()::timestamp + '1 year'::interval
);

INSERT INTO snippets(title, content, created, expires) VALUES (
    'My third snippet',
    'Do not be afraid of dealth',
    now()::timestamp,
    now()::timestamp + '1 year'::interval
);

CREATE USER web;

GRANT INSERT, SELECT ON TABLE snippets TO web;

ALTER USER web PASSWORD 'password';