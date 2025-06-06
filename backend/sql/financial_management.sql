CREATE TABLE Financial_Management (
    financial_sn INTEGER PRIMARY KEY AUTOINCREMENT,
    financial_date text DEFAULT CURRENT_TIMESTAMP, 
    -- 注意默認是當下新增日期、時間，當然可以自己指定，不過有個問題時間是UTC而不是其他的
    financial_action_type nvarchar(30) not null, -- in (Income) or out (Expenses)
    financial_action_cost integer not null, -- Expense money or income money  
    financial_action_describe nvarchar(80) not null, -- describe the action and motivation
    financial_visable integer DEFAULT 1 -- 1: visable, 0: not visable
);