// =================================================================
// 這是 manage_financial.js 的最終、乾淨版本 (2025/06/10)
// =================================================================

console.log("!!!!!!!!!! 正在執行最新的 manage_financial.js，版本時間戳：", new Date().getTime(), "!!!!!!!!!!");

document.addEventListener('DOMContentLoaded', function() {

    /**
     * 負責將從後端獲取到的財務資料陣列，渲染到 HTML 頁面上。
     * @param {Array} data - 從後端 API 收到的財務資料物件陣列。
     */
    function displayFinancialData(data) {
        // 步驟 A：防禦性檢查，確保我們收到的是一個陣列
        if (!Array.isArray(data)) {
            console.error("收到的財務資料不是一個有效的陣列:", data);
            return; // 如果資料格式不對，直接停止執行，避免錯誤
        }
        
        const summaryElement = document.getElementById('financial-summary');
        const tableBody = document.getElementById('financial-table-body');

        // 步驟 B：在前端計算總收入、總支出和總利潤
        let totalRevenue = 0;
        let totalCost = 0;

        data.forEach(item => {
            // 使用後端 Go struct JSON 標籤定義的鍵名
            if (item.financial_action_type === '收入') {
                totalRevenue += item.financial_action_cost;
            } else if (item.financial_action_type === '支出') {
                totalCost += item.financial_action_cost;
            }
        });

        const totalProfit = totalRevenue - totalCost;

        // 步驟 C：更新財務總結區塊的 HTML
        summaryElement.innerHTML = `
            <p>總收入: ${totalRevenue.toFixed(2)} 元</p>
            <p>總支出: ${totalCost.toFixed(2)} 元</p>
            <p><strong>總利潤: ${totalProfit.toFixed(2)} 元</strong></p>
        `;

        // 步驟 D：清空舊的表格內容，並根據新資料建立表格的每一行
        tableBody.innerHTML = ''; 
        data.forEach(item => {
            const row = document.createElement('tr');
            const revenueCell = item.financial_action_type === '收入' ? item.financial_action_cost.toFixed(2) : '0.00';
            const costCell = item.financial_action_type === '支出' ? item.financial_action_cost.toFixed(2) : '0.00';
            
            row.innerHTML = `
                <td>${item.financial_date}</td>
                <td>${revenueCell}</td>
                <td>${costCell}</td>
                <td>${item.financial_action_describe}</td>
            `;
            tableBody.appendChild(row);
        });
    }

    /**
     * 負責向後端 API 發送請求，獲取財務資料。
     */
    async function fetchFinancialData() {
        try {
            const response = await fetch('/api/manageFinancial?_cache_bust=' + new Date().getTime());
            if (!response.ok) {
                throw new Error('Network response was not ok. Status: ' + response.status);
            }
            
            const data = await response.json();

            // 【【【 最關鍵的除錯步驟 】】】
            // 我們要在 Console 親眼看看後端到底回傳了什麼。
            console.log("從後端收到的原始資料:", data);

            // 將獲取到的資料交給顯示函數處理
            displayFinancialData(data);

        } catch (error) {
            console.error("獲取財務資料時發生錯誤:", error);
            const errorDisplay = document.getElementById('financial-summary');
            if (errorDisplay) {
                errorDisplay.innerHTML = '<p class="error">無法載入財務資料。</p>';
            }
        }
    }

    // 頁面載入完成後，立即執行獲取資料的動作
    fetchFinancialData();
});