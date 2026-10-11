# HiveBridge 蜂桥 v1.0.0 发布包

- 产物：`release/hivebridge-1.0.0.zip`（205KB，解压后含 manifest.json，可直接「加载已解压的扩展程序」）
- 构建时间：2026-10-11

## 支持平台（网页私信 5/5 全覆盖）

| 平台 | 域名 | content script |
|---|---|---|
| 抖音 | douyin.com | content-douyin.js |
| 小红书 | xiaohongshu.com | content-xhs.js |
| TikTok | tiktok.com | content-tiktok.js |
| 闲鱼 | goofish.com / xianyu.com | content-xianyu.js |
| 快手 | kuaishou.com | content-kuaishou.js |

> 注：其余渠道（Telegram / WhatsApp / 企微 / 飞书 / 钉钉 / QQ / 邮件）均有官方 API / webhook，走 `user-server` 官方通道，无需浏览器插件。

## 包含能力
- 五平台网页私信桥接：抖音 / 小红书 / TikTok / 闲鱼 / 快手
- 实时上行：客户私信 -> user-server（触发 AI）
- 历史回填：页面加载/会话切换时存量消息仅落库，不触发 AI（防回环）
- 下行回写：user-server AI 回复 -> 回写网页私信输入框并发送
- 拟人限速风控：最小间隔 + 令牌桶 + 会话冷却 + 相同文案去重
- 多用户/多会话：按 (channel, account, conversation) 隔离路由与历史
- 协议：HTTP 三通道（/api/bridge/ingest 上行 / /api/bridge/outbox 下行 / /api/bridge/outbox/ack 状态）

## 安装（本机 Chrome）
1. 打开 chrome://extensions
2. 开启「开发者模式」
3. 点击「加载已解压的扩展程序」，选择解压后的目录
4. 打开 popup 配置 user-server 地址与凭据，登录对应平台网页即可桥接