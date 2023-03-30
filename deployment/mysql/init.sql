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
        `created_at` BIGINT,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `flag_info` (
        `id` varchar(40) NOT NULL,
        `user_id` varchar(40) NOT NULL,
        `name` varchar(255) NOT NULL,
        `status` varchar(10) NOT NULL,
        `created_at` BIGINT,
        `total_time` INT NOT NULL,
        `start_time` BIGINT,
        `challenge_num` INT NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `sign_in` (
        `id` varchar(40) NOT NULL,
        `flag_id` varchar(40) NOT NULL,
        `user_id` varchar(40) NOT NULL,
        `created_at` BIGINT,
        `current_time` INT NOT NULL,
        `total_time` INT NOT NULL,
        `content` varchar(255) NOT NULL,
        `picture_url` text NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `Wallet` (
        `id` varchar(40) NOT NULL,
        `user_id` varchar(40) NOT NULL,
        `gold_num` BIGINT NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `water_flow` (
        `id` varchar(40) NOT NULL,
        `user_id` varchar(40) NOT NULL,
        `gold_num` BIGINT NOT NULL,
        `created_at` BIGINT,
        `content` varchar(255) NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `vip` (
        `id` varchar(40) NOT NULL,
        `user_id` varchar(40) NOT NULL,
        `start_time` BIGINT,
        `end_time` BIGINT,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `siege` (
        `id` varchar(40) NOT NULL,
        `user_id` varchar(40) NOT NULL,
        `flag_id` varchar(40) NOT NULL,
        `created_at` BIGINT,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `comment` (
        `id` varchar(40) NOT NULL,
        `user_id` varchar(40) NOT NULL,
        `signin_id` varchar(40) NOT NULL,
        `content` text NOT NULL,
        `created_at` BIGINT,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

Create Table
    `props` (
        `id` varchar(40) NOT NULL,
        `flag_id` varchar(40) NOT NULL,
        `type` INT NOT NULL,
        `use_at` BIGINT,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;