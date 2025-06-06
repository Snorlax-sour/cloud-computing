CREATE TABLE Menu (
    menu_sn INTEGER PRIMARY KEY AUTOINCREMENT,
    menu_name nvarchar(30) UNIQUE NOT NULL,
    menu_describe_and_material nvarchar(80),
    menu_need_ingredient INTEGER,
    menu_money INTEGER NOT NULL,
    menu_visable INTEGER DEFAULT 1 not null,  -- 1: 顯示, 0: 不顯示
    menu_image TEXT not null,
    FOREIGN KEY (menu_need_ingredient) REFERENCES Menu_need_ingredient(join_sn)
);