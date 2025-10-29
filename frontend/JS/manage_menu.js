// manage_menu.js (API Driven Version)

// ------------------- 分類判斷邏輯 -------------------
function getItemCategory(name) {
    if (['飯', '炒飯', '披薩', '麵', '義大利麵'].some(keyword => name.includes(keyword))) return 'main';
    if (['薯條', '沙拉', '湯'].some(keyword => name.includes(keyword))) return 'side';
    if (['咖啡', '果汁', '茶', '飲料'].some(keyword => name.includes(keyword))) return 'drink';
    if (['蛋糕', '冰淇淋', '甜點', '布朗尼'].some(keyword => name.includes(keyword))) return 'dessert';
    return 'main'; // 默認主餐
}

// ------------------- DOM & 渲染邏輯 -------------------
function showCategory(category) {
    const menuItems = document.querySelectorAll('.menu-item');
    menuItems.forEach(item => {
        if (item.getAttribute('data-category') === category) {
            item.style.display = 'flex';
        } else {
            item.style.display = 'none';
        }
    });

    const categories = document.querySelectorAll('.category');
    categories.forEach(cat => {
        if (cat.textContent.includes(category === 'main' ? '主餐' : category === 'side' ? '副餐' : category === 'drink' ? '飲料' : '甜點')) { // 修正分類標籤比對
            cat.classList.add('active');
        } else {
            cat.classList.remove('active');
        }
    });
}

function createMenuItemHTML(item) {
    // 修正：使用 getItemCategory 確保 data-category 正確
    const category = getItemCategory(item.menu_name); 
    
    return `
        <div class="menu-item" 
             data-category="${category}" 
             data-name="${item.menu_name}" 
             data-price="${item.menu_money}" 
             data-img="${item.menu_image}"
             data-sn="${item.menu_sn}"> <img src="${item.menu_image}" alt="${item.menu_name}">
            <div class="item-details">
                <h3>${item.menu_name}</h3>
                <p>${item.menu_describe_and_material}</p>
            </div>
            <button class="edit-btn" data-name="${item.menu_name}">編輯</button>
        </div>
    `;
}

function renderMenuItems(items) {
    const menuContainer = document.getElementById('menu');
    if (!menuContainer) return; // 安全檢查
    menuContainer.innerHTML = '';
    
    items.forEach(item => {
        menuContainer.innerHTML += createMenuItemHTML(item);
    });

    // 預設顯示主餐
    showCategory('main'); 
}

// ------------------- API 載入 (READ) -------------------
async function loadMenuItems() {
    try {
        const response = await fetch('/api/menu_items'); // READ API
        if (response.ok) {
            const menuItems = await response.json();
            renderMenuItems(menuItems); 
        } else {
            console.error('無法載入菜單:', response.statusText);
            alert('無法載入菜單，請檢查伺服器連線。');
        }
    } catch (error) {
        console.error('連線錯誤:', error);
        alert('伺服器連線失敗。');
    }
}

// ------------------- API 更新 (UPDATE) -------------------
async function updateMenuItemHandler() {
    const popup = document.getElementById("popup");
    const newName = document.getElementById('editName').value.trim();
    const newPrice = document.getElementById('editPrice').value.trim();
    const newDescription = document.getElementById('editDescription').value.trim();
    const oldName = document.getElementById("popupTitle").textContent;
    
    // 從 DOM 獲取必要的舊資訊
    const menuItem = document.querySelector(`.menu-item[data-name="${oldName}"]`);
    if (!menuItem) { alert('找不到原始餐點'); return; }

    const currentImageURL = menuItem.getAttribute('data-img') || '/static/IMAGE/default.jpg';
    // 假設你使用 data-sn 作為唯一識別碼
    // const menuSn = parseInt(menuItem.getAttribute('data-sn')); 

    const updateData = {
        old_name: oldName,
        new_name: newName,
        menu_image: currentImageURL, 
        menu_money: parseInt(newPrice),
        menu_describe_and_material: newDescription,
    };

    try {
        const response = await fetch('/api/update_menu', { 
            method: 'PUT', 
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(updateData)
        });

        if (response.ok) {
            popup.style.display = 'none';
            alert('✅ 餐點已成功更新，並已安全寫入資料庫！');
            // 成功後重新載入，確保前端數據最新
            loadMenuItems(); 
        } else {
            const errorText = await response.text();
            alert(`❌ 更新失敗: ${errorText}`);
            console.error('後端錯誤:', errorText);
        }
    } catch (error) {
        console.error('連線錯誤:', error);
        alert('❌ 連線錯誤，請檢查 Golang 伺服器是否已啟動。');
    }
}

// ------------------- 事件監聽器 (必須在 DOMContentLoaded 後註冊) -------------------
document.addEventListener('DOMContentLoaded', () => {
    // 頁面載入後立即載入數據
    loadMenuItems();
    
    const popup = document.getElementById("popup");
    const tagContainer = document.getElementById("tagContainer");
    
    // 註冊點擊菜單事件 (使用事件委派)
    document.querySelector('.menu').addEventListener('click', (event) => {
        const clickedItem = event.target.closest('.menu-item');
        if (clickedItem) {
            // 點擊編輯按鈕才彈窗
            if (!event.target.classList.contains('edit-btn')) return;

            const name = clickedItem.getAttribute('data-name');
            const price = clickedItem.getAttribute('data-price');
            const description = clickedItem.querySelector('p').textContent;

            document.getElementById("popupTitle").textContent = name;
            document.getElementById('editName').value = name;
            document.getElementById('editPrice').value = price;
            document.getElementById('editDescription').value = description;

            // 填充 Tag 選項 (靜態/舊邏輯已移除，避免 ReferenceError)
            tagContainer.innerHTML = "<p>客製化選項尚未從 API 載入，請手動更新 DB</p>"; 

            popup.style.display = 'flex';
        }
    });

    // 註冊彈窗按鈕事件
    document.getElementById('updateMenuItem').addEventListener('click', updateMenuItemHandler);
    document.getElementById('cancelUpdate').addEventListener('click', () => { popup.style.display = 'none'; });
    document.getElementById('closePopup').addEventListener('click', () => { popup.style.display = 'none'; });
    
    window.addEventListener('click', (event) => {
        if (event.target === popup) {
            popup.style.display = 'none';
        }
    });
});