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

// 點擊餐點顯示彈窗
menuItems.forEach((item) => {
  item.addEventListener("click", () => {
      const name = item.getAttribute("data-name");
      const price = item.getAttribute("data-price");
      const description = item.getAttribute("data-desc");

      // 更新彈窗內容
      popupTitle.textContent = name;
      popupDescription.textContent = description;
      popupPrice.textContent = price;

      // 生成 tag 選項
      tagContainer.innerHTML = ""; // 清空舊的 tag
      const tags = menuData[name] || [];
      tags.forEach(tag => {
          const label = document.createElement("label");
          const checkbox = document.createElement("input");
          checkbox.type = "checkbox";
          checkbox.value = tag;
          label.appendChild(checkbox);
          label.appendChild(document.createTextNode(tag));
          tagContainer.appendChild(label);
      });

      // 顯示彈窗
      popup.style.display = "flex";
  });
});

// 關閉彈窗
closePopupButton.addEventListener("click", () => {
  popup.style.display = "none";
});

// 點擊外部關閉彈窗
window.addEventListener("click", (event) => {
  if (event.target === popup) {
      popup.style.display = "none";
  }
});

// 確認訂單按鈕
confirmOrder.addEventListener("click", () => {
  const selectedTags = Array.from(tagContainer.querySelectorAll("input[type='checkbox']:checked"))
      .map(checkbox => checkbox.value);
  alert(`已加入購物車！\n客製化選項: ${selectedTags.join(", ")}`);
  popup.style.display = "none";
});

// 主頁面 JavaScript（menu.js）
const cart = JSON.parse(localStorage.getItem('cart')) || []; // 讀取 LocalStorage

// 點擊彈窗的確認訂單按鈕
confirmOrder.addEventListener("click", () => {
    const selectedTags = Array.from(tagContainer.querySelectorAll("input[type='checkbox']:checked"))
        .map(checkbox => checkbox.value);

    const name = popupTitle.textContent; // 餐點名稱
    const price = popupPrice.textContent; // 價格
    const img = document.querySelector(`.menu-item[data-name="${name}"]`).getAttribute('data-img'); // 圖片

    const menuItem = { name, price, tags: selectedTags, img };

    // 加入到 LocalStorage 的購物車
    cart.push(menuItem);
    localStorage.setItem('cart', JSON.stringify(cart)); // 儲存資料

    alert(`餐點 ${name} 已加入購物車！\n客製化選項: ${selectedTags.join(", ")}`);
    popup.style.display = "none"; // 關閉彈窗
});