CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username       TEXT    UNIQUE NOT NULL,
  password_hash  TEXT    NOT NULL,
  role           TEXT CHECK( role IN ('user', 'agent', 'admin') ) NOT NULL,
  session        TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS ports (
  id   INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT    UNIQUE NOT NULL,
  country TEXT NOT NULL,
  city TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS quotations (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  agent_id    INTEGER NOT NULL REFERENCES users(id),
  port_id     INTEGER NOT NULL REFERENCES ports(id),
  rate        REAL    NOT NULL,
  valid_until DATE    NOT NULL,
  updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(agent_id, port_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS agent_port
  ON quotations(agent_id, port_id);