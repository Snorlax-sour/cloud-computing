CREATE TABLE Menu_need_ingredient (
    join_sn INTEGER PRIMARY KEY AUTOINCREMENT,
    menu_sn INTEGER NOT NULL,  -- 不能直接寫在同一行！！REFERENCES
    ingredient_sn INTEGER NOT NULL,
    quantity INTEGER not null,
    FOREIGN KEY (menu_sn) REFERENCES Menu(menu_sn),
    FOREIGN KEY (ingredient_sn) REFERENCES Ingredient(ingredient_sn)
);