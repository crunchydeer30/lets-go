CREATE TABLE
  snippets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
  );

CREATE INDEX idx_snippets_created ON snippets (created);

INSERT INTO
  snippets (title, content, created, expires)
VALUES
  (
    'An old silent pond',
    'An old silent pond...
A frog jumps into the pond,
splash! Silence again.
– Matsuo Bashō',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + INTERVAL '1 year'
  );

INSERT INTO
  snippets (title, content, created, expires)
VALUES
  (
    'Over the wintry forest',
    'Over the wintry
forest, winds howl in rage
with no leaves to blow.
– Natsume Soseki',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + INTERVAL '1 year'
  );

INSERT INTO
  snippets (title, content, created, expires)
VALUES
  (
    'First autumn morning',
    'First autumn morning
the mirror I stare into
shows my father''s face.
– Murakami Kijo',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + INTERVAL '7 days'
  );

CREATE TABLE
  sessions (
    token CHAR(43) PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
  );

CREATE INDEX sessions_expiry_idx ON sessions (expiry);