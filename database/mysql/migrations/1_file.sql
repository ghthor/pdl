-- +goose Up
CREATE TABLE IF NOT EXISTS `file` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `filename` text NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +goose Down
DROP TABLE IF EXISTS `file`;
