+ **HTTP 是无状态的**： HTTP 协议本身是无状态的，这意味着服务器无法记住用户的登录状态。要实现登录状态保持，你需要使用 session。

+ **Session 的工作原理：**

    + 当用户登录成功时，服务器会为用户创建一个 session，并生成一个唯一的 session ID。

    + 服务器将 session ID 发送给客户端，通常是作为 HTTP cookie 存储在浏览器中。

    + 客户端后续的请求都会自动携带这个 session ID。

    + 服务器通过 session ID 来识别用户，并获取用户相关的 session 数据。

+ **Go 中的 Session 管理：** 你可以使用第三方库（如 gorilla/sessions）来管理 session。
+ install: go get github.com/gorilla/sessions