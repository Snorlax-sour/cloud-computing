document.addEventListener('DOMContentLoaded', function() {
    const summaryContainer = document.getElementById('financial-summary');
    const tableBody = document.getElementById('financial-table-body');
    const searchInput = document.getElementById('searchInput');
    let allFinancialData = []; // 儲存所有資料

    if (!summaryContainer || !tableBody || !searchInput) {
        return; 
    }

    function displayFinancialData(data) {
        // ... (原有的 displayFinancialData 函數內容保持不變，但需要注意：
        // 總結區塊(summary)應該只根據 allFinancialData 計算一次，
        // 或者根據過濾後的 data 重新計算，這裡我們先保持簡單，總結不變)
        let totalRevenue = 0;
        let totalCost = 0;

        // 總結永遠根據全部資料計算
        allFinancialData.forEach(item => {
            if (item.financial_action_type === '收入') totalRevenue += item.financial_action_cost;
            else if (item.financial_action_type === '支出') totalCost += item.financial_action_cost;
        });
        
        const totalProfit = totalRevenue - totalCost;
        const profitClass = totalProfit >= 0 ? 'profit' : 'loss';

        summaryContainer.innerHTML = `
            <p>總收入<br><strong class="income">${totalRevenue.toFixed(2)} 元</strong></p>
            <p>總支出<br><strong class="expense">${totalCost.toFixed(2)} 元</strong></p>
            <p>總利潤<br><strong class="${profitClass}">${totalProfit.toFixed(2)} 元</strong></p>
        `;

        // 表格只顯示過濾後的資料
        tableBody.innerHTML = '';
        if (!Array.isArray(data) || data.length === 0) {
            tableBody.innerHTML = '<tr><td colspan="4">沒有符合條件的資料。</td></tr>';
            return;
        }

        data.forEach(item => {
            const row = document.createElement('tr');
            const revenueCell = item.financial_action_type === '收入' ? item.financial_action_cost.toFixed(2) : '0.00';
            const costCell = item.financial_action_type === '支出' ? item.financial_action_cost.toFixed(2) : '0.00';
            const revenueClass = item.financial_action_type === '收入' ? 'class="income"' : '';
            const costClass = item.financial_action_type === '支出' ? 'class="expense"' : '';

            row.innerHTML = `
                <td>${item.financial_date}</td>
                <td ${revenueClass}>${revenueCell}</td>
                <td ${costClass}>${costCell}</td>
                <td>${item.financial_action_describe}</td>
            `;
            tableBody.appendChild(row);
        });
    }

    async function fetchFinancialData() {
        try {
            const response = await fetch('/api/manageFinancial?_cache_bust=' + new Date().getTime());
            if (!response.ok) throw new Error('Network error');
            
            allFinancialData = await response.json();
            console.log("已獲取所有財務資料:", allFinancialData);
            displayFinancialData(allFinancialData); // 首次載入

        } catch (error) {
            console.error("獲取財務資料時發生錯誤:", error);
            summaryContainer.innerHTML = '<p class="error">無法載入財務資料。</p>';
        }
    }

    // 【【【 新增搜尋功能 】】】
    searchInput.addEventListener('input', function(event) {
        const searchTerm = event.target.value.toLowerCase();
        
        const filteredData = allFinancialData.filter(item => {
            const description = (item.financial_action_describe || '').toLowerCase();
            return description.includes(searchTerm);
        });

        // 只更新表格，不更新總結
        displayFinancialData(filteredData);
    });

    fetchFinancialData();
});