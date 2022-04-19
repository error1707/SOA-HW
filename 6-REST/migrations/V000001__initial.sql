/* pgmigrate-encoding: utf-8 */

CREATE TYPE GENDER AS ENUM ('m', 'f');

CREATE TABLE user_profiles (
    username VARCHAR(128) UNIQUE PRIMARY KEY,
    email varchar(64) UNIQUE NOT NULL,
    gender GENDER NOT NULL,
    win_count INT DEFAULT 0,
    lose_count INT DEFAULT 0,
    time_in_game INT DEFAULT 0
)