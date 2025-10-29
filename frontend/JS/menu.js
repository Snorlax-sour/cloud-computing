document.addEventListener('DOMContentLoaded', function() {
    const menuContainer = document.getElementById('menu');
    let allMenuItems = []; // 儲存從後端獲取的原始菜單資料

    // 【新增】客製化選項資料庫
        const customizationOptions = {
            '夏威夷鳳梨炒飯': ['不要鳳梨', '不要蝦仁', '加辣'],
            '經典瑪格麗特披薩': ['加倍起司', '薄皮', '厚皮'],
            '奶油培根義大利麵': ['麵加量', '醬多', '不要培根'],
            '黃金脆薯': ['去鹽', '附蜂蜜芥末醬'],
            '田園凱薩沙拉': ['醬另外放', '不要麵包丁'],
            '冰拿鐵咖啡': ['去冰', '半糖', '換燕麥奶 (+10元)'],
            '新鮮柳橙果汁': ['去冰', '半糖'],
            '重乳酪蛋糕': ['附巧克力醬'],
            '巧克力布朗尼冰淇淋': ['不要鮮奶油']
            // 未來可以從後端 API 獲取這些選項
        };

    // 【新增】獲取彈出視窗相關的元素
    const popup = document.getElementById("popup");
    const popupTitle = document.getElementById("popupTitle");
    const popupDescription = document.getElementById("popupDescription");
    const popupPrice = document.getElementById("popupPrice");
    const tagContainer = document.getElementById("tagContainer");
    const confirmOrderBtn = document.getElementById("confirmOrder");
    const closePopupBtn = document.getElementById("closePopup");

    if (!menuContainer || !popup) {
        console.error("找不到 'menu' 或 'popup' 容器！");
        return;
    }

    /**
     * 處理菜單項目的點擊事件，填充並顯示彈出視窗
     * @param {MouseEvent} event 
     */
    function handleMenuItemClick(event) {
        const clickedItem = event.currentTarget;
        
        // 從 data-* 屬性獲取資料
        const name = clickedItem.dataset.name;
        const price = clickedItem.dataset.price;
        const description = clickedItem.dataset.desc;
        const imageUrl = clickedItem.dataset.img;

        // 填充彈出視窗的基本內容
        popupTitle.textContent = name;
        popupDescription.textContent = description;
        popupPrice.textContent = price;

        // 【核心】動態生成客製化選項
        tagContainer.innerHTML = ''; // 先清空舊的選項
        const options = customizationOptions[name] || []; // 根據菜單名稱查找選項

        if (options.length > 0) {
            options.forEach(tag => {
                const label = document.createElement('label');
                label.className = 'custom-option';
                label.innerHTML = `
                    <input type="checkbox" value="${tag}">
                    <span>${tag}</span>
                `;
                tagContainer.appendChild(label);
            });
        } else {
            tagContainer.innerHTML = '<p>此餐點無客製化選項</p>';
        }
        
        // 顯示彈出視窗
        popup.style.display = "flex";
    }
    
    /**
     * 加入購物車的邏輯
     */
    function addToCart() {
        // 1. 從彈出視窗獲取當前餐點資訊
        const name = popupTitle.textContent;
        const price = parseFloat(popupPrice.textContent);
        
        // 2. 獲取所有被勾選的客製化選項
        const selectedTags = Array.from(tagContainer.querySelectorAll("input[type='checkbox']:checked"))
            .map(checkbox => checkbox.value);

        // 3. 從 localStorage 讀取現有的購物車，如果沒有就建立一個空陣列
        const cart = JSON.parse(localStorage.getItem('cart')) || [];

        // 4. 建立新的購物車項目物件
        const cartItem = {
            id: Date.now(), // 給一個簡單的唯一 ID，方便未來刪除
            name: name,
            price: price,
            tags: selectedTags,
            quantity: 1
            // 為了簡化，圖片路徑等其他資訊暫不存入
        };

        // 5. 將新項目加入購物車陣列
        cart.push(cartItem);

        // 6. 將更新後的購物車存回 localStorage
        localStorage.setItem('cart', 'JSON.stringify(cart)');

        // 7. 提示用戶並關閉彈出視窗
        alert(`「${name}」已成功加入購物車！`);
        closePopup();
    }

    /**
     * 關閉彈出視窗
     */
    function closePopup() {
        popup.style.display = "none";
    }

    // --- 以下是頁面載入和分類篩選的邏輯 (大部分不變) ---

    function renderMenuItems(itemsToRender) {
        menuContainer.innerHTML = '';
        if (!itemsToRender || itemsToRender.length === 0) {
            menuContainer.innerHTML = '<p>這個分類目前沒有餐點喔！</p>';
            return;
        }
        itemsToRender.forEach(item => {
            const name = item.menu_name || '未知餐點';
            const description = item.menu_describe_and_material || '';
            const price = item.menu_money || 0;
            const imageUrl = item.menu_image || '/static/IMAGE/default_food.png';
            let category = 'side';
            if (['飯', '披薩', '麵'].some(keyword => name.includes(keyword))) category = 'main';
            if (['咖啡', '果汁', '茶', '飲料'].some(keyword => name.includes(keyword))) category = 'drink';
            if (['蛋糕', '冰淇淋', '甜點'].some(keyword => name.includes(keyword))) category = 'dessert';
            const menuItemDiv = document.createElement('div');
            menuItemDiv.className = 'menu-item';
            menuItemDiv.setAttribute('data-category', category);
            menuItemDiv.setAttribute('data-name', name);
            menuItemDiv.setAttribute('data-price', price);
            menuItemDiv.setAttribute('data-img', imageUrl);
            menuItemDiv.setAttribute('data-desc', description);
            menuItemDiv.innerHTML = `
                <img src="${imageUrl}" alt="${name}">
                <div class="item-details">
                    <h3>${name}</h3>
                    <p>${description}</p>
                </div>
            `;
            menuItemDiv.addEventListener('click', handleMenuItemClick);
            menuContainer.appendChild(menuItemDiv);
        });
    }

    async function fetchMenuData() {
        try {
            menuContainer.innerHTML = '<p id="loading-message">正在從資料庫載入菜單...</p>';
            const response = await fetch('/api/menu_items');
            if (!response.ok) throw new Error(`伺服器回應錯誤: ${response.status}`);
            allMenuItems = await response.json();
            window.showCategory('all');
        } catch (error) {
            console.error('獲取菜單失敗:', error);
            menuContainer.innerHTML = '<p style="color: red;">無法載入菜單，請檢查後端服務或網路連線。</p>';
        }
    }

    window.showCategory = function(category) {
        document.querySelectorAll('.category').forEach(cat => {
            const catText = cat.textContent;
            let isActive = false;
            if (category === 'all' && catText === '全部') isActive = true;
            if (category === 'main' && catText === '主餐') isActive = true;
            if (category === 'side' && catText === '副餐') isActive = true;
            if (category === 'drink' && catText === '飲料') isActive = true;
            if (category === 'dessert' && catText === '甜點') isActive = true;
            cat.classList.toggle('active', isActive);
        });
        if (category === 'all') {
            renderMenuItems(allMenuItems);
        } else {
            const filteredItems = allMenuItems.filter(item => {
                let itemCategory = 'side';
                if (['飯', '披薩', '麵'].some(keyword => item.menu_name.includes(keyword))) itemCategory = 'main';
                if (['咖啡', '果汁', '茶', '飲料'].some(keyword => item.menu_name.includes(keyword))) itemCategory = 'drink';
                if (['蛋糕', '冰淇淋', '甜點'].some(keyword => item.menu_name.includes(keyword))) itemCategory = 'dessert';
                return itemCategory === category;
            });
            renderMenuItems(filteredItems);
        }
    };

    // 【新增】為彈出視窗的按鈕綁定事件
    confirmOrderBtn.addEventListener('click', addToCart);
    closePopupBtn.addEventListener('click', closePopup);
    // 點擊彈出視窗外部的半透明背景時，也關閉視窗
    popup.addEventListener('click', function(event) {
        if (event.target === popup) {
            closePopup();
        }
    });

    // 頁面載入後立即獲取菜單資料
    fetchMenuData();
});