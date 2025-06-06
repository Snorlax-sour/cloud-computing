package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting database connection...")
	db, err := connect_sqlite()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err) // Use Fatalf for errors
	}
	defer db.db.Close() // Close the connection in defer, after it is used

	// fmt.Println("We now have a database connection and can use it")
	// operation_sucessful := db.verify_User_password("front", "fu06t;3bp6m06")
	// if operation_sucessful {
	// 	fmt.Println("Sucessfully verified user")
	// } else { // else 必須和 if 的右大括號在同一行：
	// 	fmt.Println("verified user failed")
	// }

	// Call show all users here
	// allUsernames, ok := db.show_all_User()
	// if ok {
	// 	fmt.Println("printing usernames:")
	// 	for _, userName := range allUsernames {
	// 		log.Println(userName)
	// 	}
	// }

	// Serve Static Files
	http.Handle("/CSS/", http.StripPrefix("/CSS/", http.FileServer(http.Dir("../CSS"))))
	http.Handle("/JS/", http.StripPrefix("/JS/", http.FileServer(http.Dir("../JS"))))
	http.Handle("/IMAGE/", http.StripPrefix("/IMAGE/", http.FileServer(http.Dir("../IMAGE"))))
	http.Handle("/HTML/", http.StripPrefix("/HTML/", http.FileServer(http.Dir("../HTML"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/HTML/home_page.html", http.StatusFound)
	})

	// 業務邏輯： 登入、開始點餐按鈕按下去的處理
	http.HandleFunc("/login", db.submitHandler)
	// start order
	// change to /start_order
	http.HandleFunc("/start_order", db.startOrderHandler)
	// manage ingredient
	http.HandleFunc("/manage_ingredient", db.manageIngredientHandler)
	// manage financial
	http.HandleFunc("/manage_financial", db.manageFinancialHandler)
	// Start Server
	log.Println("Server is listening on: http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
