-- =================================================================
-- 角色表 (Roles) - 用於權限管理
-- =================================================================
CREATE TABLE Roles (
    role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_name TEXT UNIQUE NOT NULL -- 'admin', 'kitchen', 'front_desk', 'customer'
);

-- =================================================================
-- 用戶表 (Users) - 關聯到角色表
-- =================================================================
CREATE TABLE Users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT UNIQUE NOT NULL,
    user_password TEXT NOT NULL,
    user_salt TEXT NOT NULL,
    user_role_id INTEGER NOT NULL, -- 外鍵
    FOREIGN KEY (user_role_id) REFERENCES Roles(role_id)
);

-- =================================================================
-- 食材主檔 (Ingredients) - 儲存食材的基本資訊
-- =================================================================
CREATE TABLE Ingredients (
    ingredient_id INTEGER PRIMARY KEY AUTOINCREMENT,
    ingredient_name TEXT UNIQUE NOT NULL, -- 例如: '番茄', '麵粉', '起司絲'
    unit TEXT NOT NULL, -- 單位 (例如: '克', '個', '毫升', '包')
    -- 安全庫存量，當總庫存低於此值時可以發出警告
    safe_stock_level REAL NOT NULL DEFAULT 0,
    is_visible INTEGER NOT NULL DEFAULT 1
);

-- =================================================================
-- 食材批次庫存表 (Ingredient_Batches) - 追蹤每一批進貨
-- =================================================================
CREATE TABLE Ingredient_Batches (
    batch_id INTEGER PRIMARY KEY AUTOINCREMENT,
    ingredient_id INTEGER NOT NULL, -- 關聯到是哪一種食材
    
    -- 這一批的剩餘數量
    remaining_quantity REAL NOT NULL, 
    
    -- 【您提到的關鍵欄位】
    purchase_date TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 進貨日期
    expiry_date TEXT NOT NULL, -- 有效期限

    -- 可選的追蹤欄位
    supplier TEXT, -- 供應商
    purchase_cost REAL, -- 該批次的進貨成本

    FOREIGN KEY (ingredient_id) REFERENCES Ingredients(ingredient_id)
);

-- =================================================================
-- 菜單表 (Menu_Items) - 保持獨立
-- =================================================================
CREATE TABLE Menu_Items (
    menu_item_id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_name TEXT UNIQUE NOT NULL,
    description TEXT,
    price REAL NOT NULL, -- 使用 REAL 類型來處理金額
    image_url TEXT,
    is_visible INTEGER NOT NULL DEFAULT 1
);

-- =================================================================
-- 菜單所需食材 (連接表) - 實現 Menu 和 Ingredient 的多對多關係
-- =================================================================
CREATE TABLE Recipe_Ingredients (
    recipe_ingredient_id INTEGER PRIMARY KEY AUTOINCREMENT,
    menu_item_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    quantity_needed REAL NOT NULL, -- 所需數量
    FOREIGN KEY (menu_item_id) REFERENCES Menu_Items(menu_item_id),
    FOREIGN KEY (ingredient_id) REFERENCES Ingredients(ingredient_id)
);

-- =================================================================
-- 訂單主表 (Orders) - 儲存訂單總體資訊
-- =================================================================
CREATE TABLE Orders (
    order_id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_uuid TEXT UNIQUE NOT NULL, -- 給客戶看的10位亂碼訂單號
    order_date TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    total_amount REAL NOT NULL,
    order_status TEXT NOT NULL DEFAULT 'pending' -- 'pending', 'preparing', 'ready', 'completed', 'cancelled'
);

-- =================================================================
-- 訂單項目表 (Order_Items) - 儲存訂單內的每個餐點
-- =================================================================
CREATE TABLE Order_Items (
    order_item_id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    menu_item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price REAL NOT NULL, -- 當時的單價
    special_requests TEXT, -- 客製化選項 (例如 '不要鳳梨, 少冰')
    FOREIGN KEY (order_id) REFERENCES Orders(order_id),
    FOREIGN KEY (menu_item_id) REFERENCES Menu_Items(menu_item_id)
);

-- =================================================================
-- 財務紀錄表 (Financial_Transactions)
-- =================================================================
CREATE TABLE Financial_Transactions (
    transaction_id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_date TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    transaction_type TEXT NOT NULL, -- 'income', 'expense'
    amount REAL NOT NULL,
    description TEXT NOT NULL,
    related_order_id INTEGER, -- 可選，關聯到訂單
    is_visible INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (related_order_id) REFERENCES Orders(order_id)
);