use benjerry_dev;

CREATE TABLE IF NOT EXISTS `icecreams` (
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `name`        varchar(255)  COLLATE utf8_unicode_ci not null,
    `description` varchar(512)  COLLATE utf8_unicode_ci,
    `story`       varchar(1024) COLLATE utf8_unicode_ci,

    `image_open`     varchar(1024)  COLLATE utf8_unicode_ci,
    `image_closed`   varchar(1024)  COLLATE utf8_unicode_ci,
    `allergy_info`   varchar(1024)  COLLATE utf8_unicode_ci,
    `dietary_certifications` varchar(1024) COLLATE utf8_unicode_ci,

    `status`      varchar(20)  COLLATE utf8_unicode_ci not null,
    `created_dt`  datetime,
    `modified_dt` datetime,
    `created_by`     int(11) ,
    `modified_by`    int(11) ,
    UNIQUE  KEY idx_uniq_icecream(`name`),
    KEY idx_auth_icecream(`created_by`,`modified_by`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `sourcing_values` (
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `name`        varchar(255)  COLLATE utf8_unicode_ci not null,
    `icecream_id` int(11) NOT NULL,
    `created_dt`  datetime,
    `modified_dt` datetime,
    UNIQUE  KEY idx_uniq_sourcing(`icecream_id`,`name`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `ingredients` (
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `name`        varchar(255)  COLLATE utf8_unicode_ci not null,
    `icecream_id` int(11) NOT NULL,
    `created_dt`  datetime,
    `modified_dt` datetime,
    UNIQUE  KEY idx_uniq_ingredients(`icecream_id`,`name`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS `users` (
    `id`            int(11) NOT NULL AUTO_INCREMENT,
    `user`          varchar(255)  COLLATE utf8_unicode_ci not null,
    `pass`          varchar(255)  COLLATE utf8_unicode_ci not null,
    `status`        varchar(20)   COLLATE utf8_unicode_ci not null,
    `otp`           varchar(20)   COLLATE utf8_unicode_ci,
    `otp_exp`       datetime,
    `logged`        int(1) default 0,
    `token`         varchar(512),
    `token_exp`     datetime,
    `created_dt`  datetime,
    `modified_dt` datetime,
    UNIQUE  KEY idx_user(`user`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;


