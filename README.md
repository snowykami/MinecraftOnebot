## MinecraftOnebot

### 这是什么

- 一个基于OneBot标准的Minecraft服务器聊天机器人实现端，目标是兼容现有OneBot v11/12标准的Bot，使开发者无需改代码即可在MC中使用Bot

### 特性

- 在线模式服务器需要自己买正版账户
- 离线模式服务器不需要正版账户，但是无法使用正版皮肤
- 仅实现了部分OneBot标准的API，因为部分在MC中无法实现
- 支持同时多个服务器

### 配置

```yaml
# 这是一个示例配置文件，你可以根据自己的需求进行修改保存为config.yml
minecraft:
  servers:
    - address: "127.0.0.1:25565"
      id: 112121212 # 服务器id
      auth: "auth_2"
      ignore_self: false
      reconnect_interval: 5

  auths:
    - name: "auth_online"
      online: true

    - name: "auth_2"
      online: false
      player: LiteyukiOffline


onebot:
  implementations:
    - type: "reverse_websocket"  # ws rws http httppost
      address: "ws://127.0.0.1:20216/onebot/v11/ws"
      access_token: ""
      reconnect_interval: 5

    - type: "forward_websocket"
      host: "127.0.0.1"
      port: 20217
      access_token: ""

    - type: "http_post"
      address: "http://127.0.0.1:20218/onebot/v11"
      access_token: ""

    - type: "http"
      host: "127.0.0.1"
      port: 20219
      access_token: ""

redis:
  address: "127.0.0.1:6379"
  password: ""
  db: 0
```

### TODO

- [ ] 使用地图画插件来支持图片消息(需自备图片外链)
- [ ] 支持@player
- [x] `player_message_handler`自定义聊天样式的支持

### 常见问题

-
Bot也是玩家，死亡后是不能发送消息的，目前暂时不能自动复活，请确保你的Bot在游戏中是安全的，不要被刷出来的怪物杀死，若出现`Chat was disabled in client option`
大概率是你的bot死了。
- 某些情况下玩家消息会被接收为`system`，所以提供了`player_message_handler`选项来决定是否从系统消息中正则匹配提取玩家消息，
  当然这也会导致某些情况下系统消息被误判为玩家消息，本人能力有限，后续了解MC协议后会尝试修复。与`go-mc`
  的作者交流后推测应该是Mojang在`1.19.1`加入的聊天消息验证所致，该问题大概率出在非原版服务端上，
  目前测试发现`Paper` `Spigot` `purpur`均有此问题，其他服务端未测试，如果你有解决方案欢迎提出。
- 强烈建议开启正版验证，否则协议库可能会出现蜜汁问题
- 高版本出现`Chat message validation failure`，请安装`NoEncryption`(https://github.com/Doclic/NoEncryption)插件

### 其他

- 本项目旨在兼容现有的OneBot标准机器人，有诸多不稳定因素，如果你只是想用一些简单功能和消息互通，那么有一个更好的方案可以参考
  Nonebot的[Minecraft适配器](https://github.com/17TheWord/nonebot-adapter-minecraft)
  配合服务端插件使用，相比使用协议库，这个方案更加稳定且优雅，不过按照作者的说法可能部分Nonebot插件需要修改才能使用
- 协议库使用的是`go-mc`，一个构建仅支持一个Minecraft版本，如有其他Minecraft版本需求请修改`go.mod`自行构建
  如果你是腐竹，可以使用`ViaVersion`插件来支持多个版本

### 非常感谢以下开源项目及其作者

- Minecraft协议库：[go-mc](https://github.com/Tnze/go-mc)
- OneBot库参考：[go-libonebot](https://github.com/botuniverse/go-libonebot)
- 微软验证支持：[go-mc-ms-auth](https://github.com/maxsupermanhd/go-mc-ms-auth)
- OneBot标准为此项目和其他机器人交互提供了一套标准：[OneBot](https://onebot.dev/)