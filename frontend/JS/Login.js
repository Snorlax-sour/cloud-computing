document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("Input_bar");
    const accountInput = document.getElementById("Account_input_box");
    const passwordInput = document.getElementById("Password_input_box");
    // 模擬帳號密碼驗證資料
    const VALID_ACCOUNT = "0";
    const VALID_PASSWORD = "0";

    form.addEventListener("submit", function (event) {
        event.preventDefault(); // prevent the form from automatically submit
        const account = accountInput.value.trim();
        const password = passwordInput.value.trim();
       
        if (account === VALID_ACCOUNT && password === VALID_PASSWORD) {
            // 驗證成功，跳轉到新頁面
             form.submit();
             //window.location.href = '../HTML/manage_home_page.html'; // 替換為你的目標頁面
        } else if (account === "" || password === ""){
            alert("請輸入帳號或密碼");
        } else {
            // 驗證失敗，顯示錯誤訊息
            alert("帳號或密碼錯誤");
        }
    });
});