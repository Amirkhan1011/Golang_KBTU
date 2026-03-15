DROP TABLE IF EXISTS user_friends;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(150) NOT NULL UNIQUE,
    gender      VARCHAR(10)  NOT NULL,
    birth_date  DATE         NOT NULL
);

CREATE TABLE user_friends (
    user_id   INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT user_friends_no_self_friend
        CHECK (user_id <> friend_id)
);

