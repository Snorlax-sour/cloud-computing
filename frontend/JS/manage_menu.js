// script.js
function showCategory(category) {
    // 選取所有菜單項目
    const menuItems = document.querySelectorAll('.menu-item');
    menuItems.forEach(item => {
        if (item.getAttribute('data-category') === category) {
            item.style.display = 'flex'; // 顯示對應分類
        } else {
            item.style.display = 'none'; // 隱藏其他分類
        }
    });

    // 更新分類標籤的樣式
    const categories = document.querySelectorAll('.category');
    categories.forEach(cat => {
        if (cat.textContent === category) {
            cat.classList.add('active');
        } else {
            cat.classList.remove('active');
        }
    });
}

// 餐點資料
const menuData = {
  "義大利鳳梨炒飯": ["不要鳳梨", "不要蝦仁", "不要蔥"],
  "海鮮披薩": ["不要蝦仁", "不要魷魚", "不要起司"],
  "薯條": ["不要番茄醬", "不要芥末醬"],
  "沙拉": ["不要凱薩醬", "不要起司"],
  "咖啡": ["無糖", "微糖"],
  "果汁": ["正常冰", "半冰", "去冰"],
  "蛋糕": [],
  "冰淇淋": []
};

// 選取元素
const menuItems = document.querySelectorAll(".menu-item");
const popup = document.getElementById("popup");
const closePopupButton = document.getElementById("closePopup");
const tagContainer = document.getElementById("tagContainer");
const popupTitle = document.getElementById("popupTitle");
const popupDescription = document.getElementById("popupDescription");
const popupPrice = document.getElementById("popupPrice");
const confirmOrder = document.getElementById("confirmOrder");

// 監聽菜單項目點擊事件
document.querySelector('.menu').addEventListener('click', (event) => {
    const clickedItem = event.target.closest('.menu-item');
    if (clickedItem) {
        const name = clickedItem.getAttribute('data-name');
        const price = clickedItem.getAttribute('data-price');
        const description = clickedItem.querySelector('p').textContent;

        // 填充彈窗內容
        popupTitle.textContent = name;
        document.getElementById('editName').value = name;
        document.getElementById('editPrice').value = price;
        document.getElementById('editDescription').value = description;

        // 填充 Tag 選項
        tagContainer.innerHTML = "";
        const tags = menuData[name] || [];
        tags.forEach(tag => {
            const label = document.createElement("label");
            const checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.value = tag;
            checkbox.checked = true;
            label.appendChild(checkbox);
            label.appendChild(document.createTextNode(tag));
            tagContainer.appendChild(label);
        });

        popup.style.display = 'flex';
    }
});

// 更新菜單項目
document.getElementById('updateMenuItem').addEventListener('click', () => {
    const newName = document.getElementById('editName').value.trim();
    const newPrice = document.getElementById('editPrice').value.trim();
    const newDescription = document.getElementById('editDescription').value.trim();

    const oldName = popupTitle.textContent;
    const menuItem = document.querySelector(`.menu-item[data-name="${oldName}"]`);

    if (menuItem) {
        menuItem.setAttribute('data-name', newName);
        menuItem.setAttribute('data-price', newPrice);
        menuItem.querySelector('h3').textContent = newName;
        menuItem.querySelector('p').textContent = newDescription;
    }

    // 更新 Tag 選項
    const selectedTags = Array.from(tagContainer.querySelectorAll('input[type="checkbox"]:checked')).map(checkbox => checkbox.value);
    menuData[newName] = selectedTags;

    if (oldName !== newName) {
        delete menuData[oldName];
    }

    localStorage.setItem('menuData', JSON.stringify(menuData));
    popup.style.display = 'none';
    alert('餐點已成功更新！');
});

// 返回按鈕
document.getElementById('cancelUpdate').addEventListener('click', () => {
    popup.style.display = 'none';
});

// 關閉彈窗
document.getElementById('closePopup').addEventListener('click', () => {
    popup.style.display = 'none';
});

// 點擊外部關閉彈窗
window.addEventListener('click', (event) => {
    if (event.target === popup) {
        popup.style.display = 'none';
    }
});