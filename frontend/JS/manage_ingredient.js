document.addEventListener('DOMContentLoaded', function() {
    const tableBody = document.getElementById('ingredient-table-body');
    const searchInput = document.getElementById('searchInput');
    let allIngredients = []; // 新增一個變數來儲存所有從後端拿到的資料

    if (!tableBody || !searchInput) {
        return; 
    }

    function displayIngredients(ingredients) {
        tableBody.innerHTML = '';
        if (!Array.isArray(ingredients) || ingredients.length === 0) {
            tableBody.innerHTML = '<tr><td colspan="4">沒有符合條件的資料。</td></tr>';
            return;
        }

        ingredients.forEach(item => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${item.IngredientName || 'N/A'}</td>
                <td>${item.IngredientRemainingInventory || '0'}</td>
                <td>${item.IngredientExpiryDate || 'N/A'}</td>
                <td>${item.IngredientNote || ''}</td>
            `;
            tableBody.appendChild(row);
        });
    }

    async function fetchIngredients() {
        try {
            const response = await fetch('/api/manage_ingredient?_cache_bust=' + new Date().getTime());
            if (!response.ok) throw new Error('Network error');
            
            allIngredients = await response.json(); // 將資料存到全域變數
            console.log("已獲取所有食材資料:", allIngredients);
            displayIngredients(allIngredients); // 首次載入時顯示所有資料

        } catch (error) {
            console.error("獲取食材資料時發生錯誤:", error);
            tableBody.innerHTML = '<tr><td colspan="4" class="error">載入資料失敗！</td></tr>';
        }
    }

    // 【【【 新增搜尋功能 】】】
    searchInput.addEventListener('input', function(event) {
        const searchTerm = event.target.value.toLowerCase(); // 獲取輸入值並轉為小寫

        // 過濾 allIngredients 陣列
        const filteredIngredients = allIngredients.filter(item => {
            const name = (item.IngredientName || '').toLowerCase();
            const note = (item.IngredientNote || '').toLowerCase();
            // 如果名稱或備註包含搜尋詞，就回傳 true
            return name.includes(searchTerm) || note.includes(searchTerm);
        });

        // 用過濾後的資料重新渲染表格
        displayIngredients(filteredIngredients);
    });

    // 頁面載入後，立即執行獲取資料的動作
    fetchIngredients();
});