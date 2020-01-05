START TRANSACTION;

CREATE TABLE IF NOT EXISTS user(
  id varchar(36) NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  isActive tinyint(1) NOT NULL,
  joinDate timestamp NOT NULL DEFAULT current_timestamp(),
  lastActivity timestamp NOT NULL DEFAULT current_timestamp(),
  isEmailVerified tinyint(1) NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE user
  ADD PRIMARY KEY (id),
  ADD UNIQUE KEY id (id),
  ADD UNIQUE KEY email (email) USING HASH;
COMMIT;

