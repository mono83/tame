CREATE TABLE `original` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sourceId` int(11) unsigned DEFAULT NULL COMMENT 'Optional source ID',
  `createdAt` int(11) unsigned NOT NULL COMMENT 'Unix timestamp',
  `hash` int(11) unsigned NOT NULL COMMENT 'CRC32 checksum of URL',
  `url` varchar(1024) NOT NULL COMMENT 'Page URL',
  `headers` text NOT NULL COMMENT 'Page headers',
  `body` mediumtext NOT NULL COMMENT 'Page content',
  `timeDns` mediumint(8) unsigned NOT NULL COMMENT 'DNS request time in ms',
  `timeConnection` mediumint(8) unsigned NOT NULL COMMENT 'TCP connection time in ms',
  `timeRequestSent` mediumint(8) unsigned NOT NULL COMMENT 'Request sent time in ms',
  `timeTotal` mediumint(9) unsigned NOT NULL COMMENT 'Total time in ms',
  PRIMARY KEY (`id`),
  KEY `urlSearchIdx` (`hash`,`createdAt`,`sourceId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;