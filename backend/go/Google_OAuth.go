package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
) // 專門用於 Google 的 Endpoints
// $\text{Web}$ 應用程式授權流程 (最常見且推薦)
// !!! 警告 !!! 新增後不加入實體蜜鑰
var googleOauthConfig = &oauth2.Config{
	ClientID:     "YOUR_CLIENT_ID",
	ClientSecret: "YOUR_CLIENT_SECRET",
	RedirectURL:  "http://localhost:80/go/callback", // 必須與 Google Cloud Console 中設定的完全一致
	Scopes:       []string{"profile", "email"},
	Endpoint:     google.Endpoint,
}

// 必須定義這個 Session 名稱，否則編譯失敗
const sessionName = "oauth_state"

// generateState generates a secure random string for the OAuth state parameter.
func generateState() (string, error) {
	// 建立一個 32 位元組的緩衝區 (byte buffer)
	b := make([]byte, 32)

	// 使用 crypto/rand 填充安全隨機位元組
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// 將位元組編碼為 URL 安全的字串
	return base64.URLEncoding.EncodeToString(b), nil
}
func Google_login(w http.ResponseWriter, r *http.Request) {
	state, error_ := generateState()
	if error_ != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}
	// 2. 儲存 state 到 Cookie (模擬 Session)
	// 注意：在正式環境中，請使用安全的 Session 管理庫！
	http.SetCookie(w, &http.Cookie{
		Name:     sessionName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   300, // 5分鐘有效期
	})
	url := googleOauthConfig.AuthCodeURL(state)
	// 將用戶導向這個 URL
	// 4. 將用戶導向 Google 進行授權
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	// 獲取 token
	// token, err := googleOauthConfig.Exchange(ctx, code)
	// // ... 處理錯誤
	// if err != nil {
	// 	log.Printf("Token exchange token: %v", token)

	// }
	// 成功獲取 token 後，您就可以使用它來驗證用戶身份了
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// 1. 驗證 state 參數 (防止 CSRF 攻擊)

	// 從 URL 參數中獲取 state
	receivedState := r.FormValue("state")

	// 從 Cookie (Session) 中獲取儲存的原始 state
	cookie, err := r.Cookie(sessionName)
	if err != nil {
		http.Error(w, "State cookie not found or expired", http.StatusBadRequest)
		return
	}
	storedState := cookie.Value

	// 刪除一次性 state cookie (銷毀 state)
	http.SetCookie(w, &http.Cookie{Name: sessionName, MaxAge: -1, Path: "/"})

	// 嚴格比對
	if receivedState != storedState {
		log.Printf("Invalid state: received %s, expected %s", receivedState, storedState)
		http.Error(w, "Invalid state parameter", http.StatusUnauthorized)
		return
	}

	// 2. 獲取授權碼 (code)
	code := r.FormValue("code")
	if code == "" {
		// 處理 Google 返回的錯誤 (例如用戶拒絕授權)
		http.Error(w, "Authorization code not provided", http.StatusBadRequest)
		return
	}

	// 3. 使用 code 換取 Token (核心 Exchange 步驟)
	// 創建一個空的、非取消的（$\text{non-cancellable}$）
	// $\text{Context}$ 對象，作為這個操作的根節點
	ctx := context.Background()

	// googleOauthConfig.Exchange 將授權碼 code 與 clientID/clientSecret 一起發送到 Google
	// 換取 Access Token 和 ID Token
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("Token exchange error: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// 至此，token 已經成功獲取。
	// token.AccessToken 是 Access Token
	// token.Extra("id_token") 包含了 ID Token (用於 OpenID Connect 驗證身份)

	// 4. 驗證身份（OIDC）：這裡使用 ID Token 進行身份驗證並獲取用戶數據
	// 這是實際登入的步驟，通常需要額外處理 ID Token

	// 成功! 返回成功訊息或重定向到應用程式主頁
	fmt.Fprintf(w, "Login successful! Access Token: %s", token.AccessToken)
	// 實際應用中，會在這裡建立用戶的 Session
}
