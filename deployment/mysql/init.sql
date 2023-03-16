-- Active: 1678972307496@@oldgeneral.top@33306@oldgeneral

Create Table
    `user` (
        `id` varchar(40) NOT NULL,
        `username` varchar(255) NOT NULL,
        `password` varchar(255) NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `user_info` (
        `id` varchar(40) NOT NULL,
        `name` varchar(255) NOT NULL,
        `signature` varchar(255) NOT NULL,
        `gender` varchar(10) NOT NULL,
        `avatar` text NOT NULL,
        `created_at` DATETIME,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;