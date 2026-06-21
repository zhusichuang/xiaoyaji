# xiaoyaji

基于微信云托管 `wxcloudrun-golang` 模板改造的小芽记后端，使用 Gin + Gorm 提供小程序容器 API。

## 运行时

- CloudBase Run / 微信云托管容器
- Go `1.23.2`
- Gin
- Gorm + MySQL

## 环境变量

- `MYSQL_ADDRESS`
- `MYSQL_USERNAME`
- `MYSQL_PASSWORD`
- `MYSQL_DATABASE`，默认 `xiaoyaji`

## 鉴权约定

小程序通过 `wx.cloud.callContainer` 调用时，服务端优先从以下位置读取用户身份：

- `X-WX-OPENID`
- `X-Wx-Openid`
- `openid` 查询参数，仅用于本地调试

## 核心 API

- `POST /api/login`
- `GET /api/families`
- `POST /api/families`
- `GET /api/babies`
- `POST /api/babies`
- `GET /api/actions`
- `POST /api/actions`
- `POST /api/actions/batch`
- `GET /api/summary/today`
- `POST /api/ai/parse`
- `POST /api/ai/chat`

## 本地调试

```bash
MYSQL_ADDRESS=127.0.0.1:3306 \
MYSQL_USERNAME=root \
MYSQL_PASSWORD=your_password \
MYSQL_DATABASE=xiaoyaji \
/Users/zhusichuang/sdk/go1.23.2/bin/go run .
```
