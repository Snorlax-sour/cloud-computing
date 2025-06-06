CREATE TABLE Ingredient (
    ingredient_sn INTEGER PRIMARY KEY AUTOINCREMENT, -- 食材編號，自動遞增的主鍵, autoincrement要在primary key後面
    ingredient_name TEXT UNIQUE NOT NULL,  -- 食材名稱，唯一且不可為空
    ingredient_remaining_inventory INTEGER NOT NULL, -- 食材剩餘庫存量，不可為空
    ingredient_visable INTEGER DEFAULT 1, -- 食材是否可見，預設為可見
    ingredient_expiry_date TEXT NOT NULL, -- 食材到期日，不可為空，使用 TEXT 儲存 "YYYY-MM-DD HH:MM:SS.SSS"
    ingredient_delivery_date TEXT,  -- 食材到貨日，使用 TEXT 儲存 "YYYY-MM-DD HH:MM:SS.SSS"
    ingredient_delivery_quantity INTEGER -- 食材進貨數量
);