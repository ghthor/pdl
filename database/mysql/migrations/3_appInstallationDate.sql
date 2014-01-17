-- +goose Up
ALTER TABLE `app` ADD COLUMN (
  `installedAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
ALTER TABLE `app` DROP COLUMN `installedAt`;
