CREATE TABLE IF NOT EXISTS weibo
(
    id              varchar(20) NOT NULL,
    bid             varchar(12) NOT NULL,
    user_id         varchar(20),
    screen_name     varchar(30),
    text            varchar(2000),
    article_url     varchar(100),
    topics          varchar(200),
    at_users        varchar(1000),
    pics            varchar(3000),
    video_url       varchar(1000),
    location        varchar(100),
    created_at      DATETIME,
    source          varchar(30),
    attitudes_count INT,
    comments_count  INT,
    reposts_count   INT,
    retweet_id      varchar(20),
    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS user
(
    id                varchar(20) NOT NULL,
    screen_name       varchar(30),
    gender            varchar(10),
    statuses_count    INT,
    followers_count   INT,
    follow_count      INT,
    registration_time varchar(20),
    sunshine          varchar(20),
    birthday          varchar(40),
    location          varchar(200),
    education         varchar(200),
    company           varchar(200),
    description       varchar(400),
    profile_url       varchar(200),
    profile_image_url varchar(200),
    avatar_hd         varchar(200),
    urank             INT,
    mbrank            INT,
    verified          BOOLEAN DEFAULT 0,
    verified_type     INT,
    verified_reason   varchar(140),
    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;