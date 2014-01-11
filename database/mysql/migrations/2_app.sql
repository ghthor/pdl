-- +goose Up
CREATE TABLE IF NOT EXISTS `app` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `pkgId` int unsigned,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +goose Down
DROP TABLE IF EXISTS `app`;
