CREATE TABLE installations (
    name TEXT PRIMARY KEY,
    path TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY (type) REFERENCES installation_type (name)
);

CREATE TABLE installation_type (
  name TEXT PRIMARY KEY
);
