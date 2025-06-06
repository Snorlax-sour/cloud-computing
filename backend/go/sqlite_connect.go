package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

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
	file := "../sql/order_db.db"
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
			redirectURL = "/HTML/manage_home_page.html"
		} else {
			redirectURL = "/HTML/home_page.html"
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
            <meta http-equiv="refresh" content="2;url=/HTML/login.html" />
            </head>
            <body>
            <h1> 帳號密碼錯誤 </h1>
            </body>
            </html>
        `
		redirectURL = "/HTML/login.html"
		t := template.Must(template.New("response").Parse(tmpl))
		data := ResponseData{
			RedirectURL: redirectURL,
		}
		t.Execute(w, data)
		return
	}
}
func (db *DB) startOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("startOrderHandler called")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session data
	session, err := db.sessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username, ok := session.Values["username"].(string)
	var redirectURL string
	if !ok {
		redirectURL = "/HTML/home_page.html"
	} else if username == "boss" {
		redirectURL = "/HTML/manage_home_page.html"
	} else {
		redirectURL = "/HTML/menu.html" // 搞錯檔案
	}
	tmpl := `
        <!DOCTYPE html>
        <html>
        <head>
        <meta charset="utf-8">
        <title>跳转中...</title>
        <meta http-equiv="refresh" content="1;url={{.RedirectURL}}" />
        </head>
        <body>
          
        </body>
        </html>
    `

	t := template.Must(template.New("response").Parse(tmpl))
	data := ResponseData{
		RedirectURL: redirectURL,
	}
	t.Execute(w, data)
	return
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
func (db *DB) manageIngredientHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("manageIngredientHandler called")
	ingredientData, err := db.getAllIngredientData()
	if err != nil {
		http.Error(w, "Error querying ingredient data", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("../HTML/manage_ingredient.html")
	if err != nil {
		log.Println("Error parsing html files", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, ingredientData) // assign to err to check the error
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

type FinancialData struct {
	FinancialDate           string
	FinancialActionCost     int
	FinancialActionType     string
	FinancialActionDescribe string
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
	log.Println("manageFinancialHandler called")
	financialData, err := db.getAllFinancialData()
	if err != nil {
		http.Error(w, "Error querying financial data", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("../HTML/manage_financial.html")
	if err != nil {
		log.Println("Error parsing html files", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, financialData) // assign to err to check the error
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
