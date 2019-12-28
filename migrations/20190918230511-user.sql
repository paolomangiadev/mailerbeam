-- +migrate Up
CREATE TABLE users(
  id UUID NOT NULL PRIMARY KEY,
  [name] VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL, 
  email VARCHAR(100) NOT NULL UNIQUE, 
  role VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE users;