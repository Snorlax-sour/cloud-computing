CREATE TABLE Menu_item_order (
    order_id nvarchar(10) PRIMARY KEY, -- 不是用sn是因為不是順序的數字 or 序號等
    menu_money INTEGER NOT NULL,
    additional_request text null, -- 需要應該要換成其他table來儲存
    total_quantity INTEGER not null,
    menu_sn INTEGER not null,
    FOREIGN KEY (menu_sn) REFERENCES Menu(menu_sn)
); 