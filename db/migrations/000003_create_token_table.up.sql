START TRANSACTION;

CREATE TABLE IF NOT EXISTS token(
    id varchar(36) NOT NULL,
    token varchar(100) NOT NULL,
    expireAt timestamp NOT NULL,
    isClaimed tinyint(1) NOT NULL DEFAULT 0,
    userId varchar(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE token
  ADD PRIMARY KEY (id),
  ADD UNIQUE KEY id (id),
  ADD UNIQUE KEY token (token);
COMMIT;