package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"encoding/json" // json need
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3" // Import the driver
)

type DB struct {
	db           *sql.DB
	filepath     string
	sessionStore *sessions.CookieStore
}

// CHANGED: added error return
func connect_sqlite() (*DB, error) {
	file := "./data/order_db.db" // change path because fit container of path 2025/06/10
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}
	// Enable foreign key constraints.
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Println("Error enabling foreign key constraints:", err)
		db.Close() // Close the database if we failed
		return nil, err
	}
	//  key should be an authentication key to encrypt the cookie
	//  make sure to replace it with your own key for your applications
	key := "your-secret-authentication-key"              // Replace with a strong key
	sessionStore := sessions.NewCookieStore([]byte(key)) // create a cookie store
	fmt.Println("Successfully connected to SQLite database with foreign keys enabled!")
	return &DB{db: db, filepath: file, sessionStore: sessionStore}, nil
}

// (db_name *DB)類似於class的部份的method
func (db_name *DB) insert_value_User(user_name string, user_password string) (*DB, bool) {
	if db_name == nil || user_password == "" || user_name == "" {
		log.Println("Invalid input: db_name is nil or user_password/user_name is empty")
		log.Println("User name:", user_name, "User password:", user_password)
		return nil, false
	}

	// 准備插入語句
	insertStmt, err := db_name.db.Prepare("INSERT INTO User (user_name, user_password, user_salt) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing insert statement:", err)
		return nil, false
	}
	defer insertStmt.Close()

	// 生成 Hash 和 Salt
	password, salt, err := hashPassword(user_password)
	if err != nil {
		log.Println("Error generating hash and salt:", err)
		return nil, false
	}
	log.Println("Generated Password:", password, "Generated Salt:", salt)

	// 執行插入
	_, err = insertStmt.Exec(user_name, password, salt)
	if err != nil {
		log.Println("Error executing insert statement:", err)
		return nil, false
	}

	return db_name, true
}

func (db *DB) show_User(username string) (string, bool) {
	if db == nil || username == "" {
		log.Println("error input db: ", db.filepath)
		log.Println("or error username: ", username)
		return "", false
	}
	// CHANGED: Removed .Open(), just use existing connection
	// db.db.Open("sqlite3", file) //Remove this as well

	searchStmt, err := db.db.Prepare("SELECT user_name FROM User WHERE user_name = ?")
	if err != nil {
		log.Println("Error preparing select statement", err)
		return "", false
	}
	defer func() {
		if err := searchStmt.Close(); err != nil {
			log.Println("Error closing prepared statement", err)
		}
	}()

	row := searchStmt.QueryRow(username)
	var userName string
	err = row.Scan(&userName)

	if err == sql.ErrNoRows {
		log.Println("User Not Found", err)
		return "", false
	}

	if err != nil {
		log.Println("Error executing select statement", err)
		return "", false
	}
	// CHANGED: DO NOT close the database here
	// db.db.Close()
	return userName, true
}
func (db *DB) show_all_User() ([]string, bool) {
	if db == nil {
		log.Println("db not exist")
		return nil, false
	}
	searchStmt, err := db.db.Prepare("SELECT user_name FROM User")
	if err != nil {
		log.Println("Error preparing select statement", err)
		return nil, false
	}
	defer func() {
		if err := searchStmt.Close(); err != nil {
			log.Println("Error closing prepared statement", err)
		}
	}()

	rows, err := searchStmt.Query()
	if err != nil {
		log.Println("Error executing select statement", err)
		return nil, false
	}
	defer rows.Close()

	var userNames []string
	for rows.Next() {
		var userName string
		err = rows.Scan(&userName)
		if err != nil {
			log.Println("Error scanning row", err)
			return nil, false
		}
		userNames = append(userNames, userName)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows", err)
		return nil, false
	}
	return userNames, true
}
func (db *DB) verify_User_password(user_name string, user_input_password string) bool {
	if db == nil {
		return fmt.Errorf("database connection is nil") == nil
	}
	if user_name == "" || user_input_password == "" {
		log.Println("empty input username: ", user_name, " or input password: ", user_input_password)
		return false
	}
	username, operation_sucessful := db.show_User(user_name)
	if (username == user_name) && operation_sucessful {
		query := "SELECT user_password, user_salt FROM User WHERE user_name = ?"

		row := db.db.QueryRow(query, username)
		var user_hash_password string
		var user_salt string

		// 提取結果
		err := row.Scan(&user_hash_password, &user_salt)

		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user not found") == nil
			}
			return fmt.Errorf("error querying database: %v", err) == nil
		}
		operation_sucessful = verifyPassword(user_input_password, user_hash_password, user_salt)
		return operation_sucessful
	}
	return false
}

type ResponseData struct {
	Username    string
	RedirectURL string
}

func (db *DB) submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Println("Received POST request with username:", username, "and password:", password)
	operation_successful := db.verify_User_password(username, password)
	var redirectURL string
	if operation_successful {
		// Set session data
		session, err := db.sessionStore.Get(r, "session-name") // the session name is "session-name"
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["username"] = username
		session.Save(r, w) // save session
		if username == "boss" {
			redirectURL = "/html/manage_home_page.html"
		} else {
			redirectURL = "/html/home_page.html"
		}

		tmpl := `
            <!DOCTYPE html>
            <html>
            <head>
            <meta charset="utf-8">
            <title>驗證成功</title>
            <meta http-equiv="refresh" content="2;url={{.RedirectURL}}" />
            </head>
            <body>
            <h1> Hello {{.Username}} </h1>
            </body>
            </html>
        `
		t := template.Must(template.New("response").Parse(tmpl))
		data := ResponseData{
			Username:    username,
			RedirectURL: redirectURL,
		}
		t.Execute(w, data)
		return
	} else {

		tmpl := `
            <!DOCTYPE html>
            <html>
            <head>
            <meta charset="utf-8">
            <title>驗證失敗</title>
            <meta http-equiv="refresh" content="2;url=/html/login.html" />
            </head>
            <body>
            <h1> 帳號密碼錯誤 </h1>
            </body>
            </html>
        `
		redirectURL = "/html/login.html"
		t := template.Must(template.New("response").Parse(tmpl))
		data := ResponseData{
			RedirectURL: redirectURL,
		}
		t.Execute(w, data)
		return
	}
}


// 在 backend/go/sqlite_connect.go 中

func (db *DB) startOrderHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("=====================================================")
    log.Println("startOrderHandler called")

    // 1. 檢查請求方法是否為 POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // 2. 獲取當前 session
    session, err := db.sessionStore.Get(r, "session-name")
    if err != nil {
        log.Printf("startOrderHandler: 無法獲取 session: %v。將以未登入用戶身份處理。", err)
        // 即使 session 獲取失敗，也應該讓用戶能點餐，所以我們不中斷流程
    }

    // 3. 檢查 session 中是否有 'username'
    username, ok := session.Values["username"].(string)
    
    var redirectURL string

    // 4. 根據是否為 'boss' 決定跳轉路徑
    // 這個邏輯與您最開始的設計保持一致
    if ok && username == "boss" {
        // 如果 session 中有 username 且值為 'boss'
        log.Println("startOrderHandler: User 'boss' is logged in. Redirecting to manage home page.")
        redirectURL = "/manage_home_page.html"
    } else {
        // 其他所有情況 (未登入、普通用戶登入等)
        log.Println("startOrderHandler: User is not 'boss' or not logged in. Redirecting to menu.")
        redirectURL = "/menu.html"
    }

    // 5. 【【【 核心修改 】】】
    // 使用 http.Redirect 函數發送一個標準的 HTTP 302 重導向。
    // 這會告訴瀏覽器立即跳轉到 redirectURL 指定的地址。
    // 這是最可靠、最高效的伺服器端跳轉方式。
    http.Redirect(w, r, redirectURL, http.StatusFound) // StatusFound 對應的狀態碼就是 302
}

// struct to store the query result of ingredient table
type IngredientData struct {
	IngredientName               string
	IngredientRemainingInventory int
	IngredientExpiryDate         string
	IngredientNote               string
}

// func to store all ingredient data
func (db *DB) getAllIngredientData() ([]IngredientData, error) {
	query := "SELECT ingredient_name, ingredient_remaining_inventory, ingredient_expiry_date, ingredient_delivery_date FROM Ingredient where ingredient_visable = 1;"
	rows, err := db.db.Query(query)
	if err != nil {
		log.Println("Error querying ingredient data", err)
		return nil, err
	}
	defer rows.Close()

	var allIngredientData []IngredientData
	for rows.Next() {
		var data IngredientData
		var ingredientDeliveryDate sql.NullString
		err := rows.Scan(&data.IngredientName, &data.IngredientRemainingInventory, &data.IngredientExpiryDate, &ingredientDeliveryDate)
		if err != nil {
			log.Println("Error scanning ingredient data", err)
			return nil, err
		}

		if ingredientDeliveryDate.Valid {
			data.IngredientNote = ingredientDeliveryDate.String
		}

		allIngredientData = append(allIngredientData, data)
	}
	if err := rows.Err(); err != nil {
		log.Println("error iterating ingredient rows:", err)
		return nil, err
	}

	return allIngredientData, nil
}
// in sqlite_connect.go

func (db *DB) manageIngredientHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("=====================================================")
    log.Println("manageIngredientHandler called")

    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    ingredientData, err := db.getAllIngredientData()
    if err != nil {
        http.Error(w, "Error querying ingredient data", http.StatusInternalServerError)
        log.Printf("manageIngredientHandler: [錯誤] 查詢食材資料失敗: %v", err)
        return
    }

    log.Printf("manageIngredientHandler: [成功] 查詢到 %d 筆食材資料。", len(ingredientData))
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(ingredientData); err != nil {
        log.Printf("manageIngredientHandler: [錯誤] 將食材資料編碼為 JSON 時失敗: %v", err)
    } else {
        log.Println("manageIngredientHandler: [成功] 已將食材資料以 JSON 格式回應。")
    }
}

type FinancialData struct {
    FinancialDate       string  `json:"financial_date"`
    FinancialActionCost float64 `json:"financial_action_cost"`
    FinancialActionType string  `json:"financial_action_type"`
    FinancialActionDescribe string `json:"financial_action_describe"`
}

func (db *DB) getAllFinancialData() ([]FinancialData, error) {
	query := "SELECT financial_date, financial_action_cost, financial_action_type, financial_action_describe FROM Financial_Management where financial_visable = 1;"
	rows, err := db.db.Query(query)
	if err != nil {
		log.Println("Error querying financial data", err)
		return nil, err
	}
	defer rows.Close()

	var allFinancialData []FinancialData
	for rows.Next() {
		var data FinancialData  
		err = rows.Scan(&data.FinancialDate, &data.FinancialActionCost, &data.FinancialActionType, &data.FinancialActionDescribe)
		if err != nil {
			log.Println("Error scanning financial data:", err)
			return nil, err
		}
		allFinancialData = append(allFinancialData, data)
	}
	if err = rows.Err(); err != nil {
		log.Println("error iterating financial rows", err)
		return nil, err
	}
	return allFinancialData, nil
}
func (db *DB) manageFinancialHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("=====================================================")
    log.Println("manageFinancialHandler called")

    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        log.Printf("manageFinancialHandler: Invalid method %s. Expected GET.", r.Method)
        return
    }

    financialData, err := db.getAllFinancialData()
    if err != nil {
        http.Error(w, "Error querying financial data", http.StatusInternalServerError)
        log.Printf("manageFinancialHandler: [錯誤] 查詢財務資料失敗: %v", err)
        return
    }

    // 【新增的除錯日誌】
    log.Printf("manageFinancialHandler: [成功] 查詢到 %d 筆財務資料。", len(financialData))
    log.Printf("manageFinancialHandler: [資料預覽] 準備回傳的財務資料：%+v", financialData)


    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(financialData); err != nil {
        log.Printf("manageFinancialHandler: [錯誤] 將資料編碼為 JSON 時失敗: %v", err)
    } else {
        log.Println("manageFinancialHandler: [成功] 已將財務資料以 JSON 格式回應。")
    }
}

func (db *DB) logoutHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("=====================================================")
    log.Println("logoutHandler called")

    // 1. 獲取當前的 session
    session, err := db.sessionStore.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 2. 將 session 的 MaxAge 設置為 -1，這會告訴瀏覽器立即刪除這個 Cookie
    session.Options.MaxAge = -1
	/*
	
	*/
    
    // 3. 儲存更改，將刪除指令發送給瀏覽器
    err = session.Save(r, w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 4. 重導向回首頁
    // 雖然前端 JS 會處理跳轉，但後端重導向是一個好的備用方案
    log.Println("Session cookie deleted. Redirecting to home page.")
    http.Redirect(w, r, "/home_page.html", http.StatusFound)
}

// 在 sqlite_connect.go 中

// =================================================================
// 菜單功能相關 (基於舊的資料庫結構)
// =================================================================

// MenuItemData struct 對應舊的 `Menu` 表
type MenuItemData struct {
    Sn          int    `json:"menu_sn"`
    Name        string `json:"menu_name"`
    Description string `json:"menu_describe_and_material"`
    Price       int    `json:"menu_money"`
    ImageURL    string `json:"menu_image"`
    // 舊表中的 menu_visable 和 menu_need_ingredient 我們在查詢時處理，不在 struct 中體現
}

// getAllMenuItems 從舊的 `Menu` 表獲取資料
func (db *DB) getAllMenuItems() ([]MenuItemData, error) {
    // 【關鍵】查詢語句使用舊的表名 `Menu` 和舊的欄位名
    query := "SELECT menu_sn, menu_name, menu_describe_and_material, menu_money, menu_image FROM Menu WHERE menu_visable = 1"
    rows, err := db.db.Query(query)
    if err != nil {
        log.Printf("[舊結構] 查詢 Menu 表時出錯: %v", err)
        return nil, err
    }
    defer rows.Close()

    var items []MenuItemData
    for rows.Next() {
        var item MenuItemData
        var description sql.NullString // 處理 menu_describe_and_material 可能為 NULL

        // 【關鍵】Scan 的變數數量和順序必須與 SELECT 的欄位完全對應
        err := rows.Scan(&item.Sn, &item.Name, &description, &item.Price, &item.ImageURL)
        if err != nil {
            log.Printf("[舊結構] 掃描 Menu 表的某一行時出錯: %v", err)
            continue 
        }
        item.Description = description.String // 處理 NULL 值
        
        items = append(items, item)
    }
    return items, nil
}

// getMenuItemsHandler 是處理 API 請求的函數 (這個函數通常不需要修改)
func (db *DB) getMenuItemsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("=====================================================")
    log.Println("getMenuItemsHandler (舊結構) 被呼叫")

    items, err := db.getAllMenuItems()
    if err != nil {
        http.Error(w, "無法獲取菜單項目", http.StatusInternalServerError)
        return
    }

    log.Printf("[舊結構] 成功從資料庫檢索到 %d 個菜單項目。", len(items))

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
}

// 在 sqlite_connect.go 中

// addMenuItemHandler 處理新增餐點的請求
func (db *DB) addMenuItemHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("=====================================================")
    log.Println("addMenuItemHandler (舊結構) 被呼叫")

    // 1. 檢查請求方法是否為 POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // 2. 解析表單數據
    // 如果未來要處理文件上傳，需要用 r.ParseMultipartForm()
    if err := r.ParseForm(); err != nil {
        http.Error(w, "無法解析表單", http.StatusBadRequest)
        return
    }

    // 3. 從表單中獲取值
    // r.FormValue() 會獲取與 <input name="..."> 匹配的值
    name := r.FormValue("menu_name")
    price := r.FormValue("menu_money")
    description := r.FormValue("menu_describe")
    imageURL := r.FormValue("menu_image")

    // 4. 簡單的後端驗證
    if name == "" || price == "" || imageURL == "" {
        http.Error(w, "餐點名稱、價格和圖片路徑為必填項", http.StatusBadRequest)
        return
    }
    
    log.Printf("收到新增餐點請求: 名稱=%s, 價格=%s, 描述=%s, 圖片=%s", name, price, description, imageURL)

    // 5. 準備 SQL 插入語句
    // 使用舊的 `Menu` 表和欄位名
    query := `
        INSERT INTO Menu (menu_name, menu_money, menu_describe_and_material, menu_image, menu_visable)
        VALUES (?, ?, ?, ?, 1)
    `
    stmt, err := db.db.Prepare(query)
    if err != nil {
        log.Printf("準備 SQL 語句時出錯: %v", err)
        http.Error(w, "資料庫內部錯誤", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    // 6. 執行 SQL 語句
    _, err = stmt.Exec(name, price, description, imageURL)
    if err != nil {
        // 處理可能的錯誤，例如 'UNIQUE constraint failed: Menu.menu_name'
        log.Printf("執行 SQL 插入時出錯: %v", err)
        http.Error(w, "無法新增餐點，可能是名稱重複或資料格式錯誤。", http.StatusInternalServerError)
        return
    }

    log.Printf("成功新增餐點: %s", name)

    // 7. 新增成功後，重導向回菜單管理頁面
    http.Redirect(w, r, "/manage_menu.html", http.StatusFound)
}