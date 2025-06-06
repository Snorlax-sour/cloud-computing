// 從 LocalStorage 中取得訂單編號
const orderId = localStorage.getItem('orderId') || '未知編號';
document.getElementById('orderId').textContent = orderId;