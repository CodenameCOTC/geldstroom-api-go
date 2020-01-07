START TRANSACTION;

CREATE TABLE IF NOT EXISTS `transaction`(
  id varchar(36) NOT NULL,
  amount BIGINT NOT NULL,
  description varchar(256) NOT NULL DEFAULT "",
  category varchar(100) NOT NULL,
  type varchar(100) NOT NULL,
  createdAt timestamp NOT NULL DEFAULT current_timestamp(),
  updatedAt timestamp NOT NULL DEFAULT current_timestamp(),
  userId varchar(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `transaction`
  ADD PRIMARY KEY (id),
  ADD UNIQUE KEY id (id);

COMMIT;