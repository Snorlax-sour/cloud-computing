CREATE TABLE User (
    sn INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT UNIQUE NOT NULL,
    privilege INTEGER DEFAULT 1, -- 1: powerful
    user_password TEXT NOT NULL, -- hashed password
    user_salt TEXT NOT NULL -- hash salt value store
);
-- 至少新增user_name跟user_password, user_salt, privilege默認是1