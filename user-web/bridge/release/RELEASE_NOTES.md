
## 固定扩展 ID（v1.0.0+）

- manifest 已内嵌固定 RSA key，扩展 ID 恒为 `odlameecjipmbmbejkplpemijjgpljce`
- 服务端 CORS 白名单已配 `CORS_ALLOWED_EXTENSIONS=chrome-extension://odlameecjipmbmbejkplpemijjgpljce`
- **重新加载插件后请确认 chrome://extensions 中 ID 与上述一致**，否则连接会因 CORS fail-closed 被拒
