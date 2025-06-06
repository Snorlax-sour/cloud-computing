const cartItems = JSON.parse(localStorage.getItem('cart')) || [];
const cartContainer = document.getElementById('cartItems');
const totalAmountEl = document.getElementById('totalAmount');

// 計算總金額
function calculateTotal() {
    const total = cartItems.reduce((sum, item) => sum + parseFloat(item.price), 0);
    totalAmountEl.textContent = `總金額: ${total} 元`;
}

// 顯示購物車項目
if (cartItems.length > 0) {
    cartItems.forEach(item => {
        const itemDiv = document.createElement('div');
        itemDiv.className = 'cart-item';
        itemDiv.innerHTML = `
            <img src="${item.img}" alt="${item.name}" class="cart-item-img">
            <div>
                <h2>${item.name}</h2>
                <p>價格: ${item.price} 元</p>
                <p>客製化: ${item.tags.length > 0 ? item.tags.join(', ') : '無'}</p>
            </div>
        `;
        cartContainer.appendChild(itemDiv);
    });

    // 計算總金額
    calculateTotal();
} else {
    cartContainer.innerHTML = '<p>購物車為空</p>';
}

// 清空購物車
document.getElementById('clearCart').addEventListener('click', () => {
    localStorage.removeItem('cart');
    window.location.reload();
});

// 確認送出訂單
document.getElementById('submitOrder').addEventListener('click', () => {
    if (cartItems.length === 0) {
        alert('購物車為空，無法送出訂單！');
        return;
    }

    // 模擬送出訂單的處理邏輯
    alert('訂單已送出，感謝您的訂購！');
    localStorage.removeItem('cart');
    const orderId = generateOrderId(); // 生成訂單編號
    localStorage.setItem('orderId', orderId); // 將訂單編號存入 LocalStorage
    window.location.href = '../HTML/order_confirmation.html'; // 跳轉到確認頁面
});

// 生成 10 位隨機字母和數字
function generateOrderId() {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let orderId = '';
    for (let i = 0; i < 10; i++) {
        orderId += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    return orderId;
}