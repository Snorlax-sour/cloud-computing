// 这是一个简单的示例，后续可以添加更完善的设置和处理逻辑
function setRefreshInterval() {
    const refreshInterval = 60; // default 60 sec
    if(refreshInterval > 0) {
        setTimeout(function() {
            location.reload(); // reload the whole page.
        }, refreshInterval * 1000);
    }

}
setRefreshInterval();