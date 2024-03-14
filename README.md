## MinecraftOnebot
### 这是什么
- 一个基于OneBot标准的Minecraft服务器聊天机器人实现端，目标是兼容现有OneBot v11/12标准的Bot，使开发者无需改代码即可在MC中使用Bot
- 还没有Release，快了快了...

### 特性
- 在线模式服务器需要自己买正版账户
- 离线模式服务器不需要正版账户，但是无法使用正版皮肤
- 仅实现了部分OneBot标准的API，因为部分在MC中无法实现
- 支持同时多个服务器

### 配置
```yaml
# 这是一个示例配置文件，你可以根据自己的需求进行修改保存为config.yml
common: # 通用配置
  join_interval: 3  # 加入服务器的间隔，单位秒，防止微软验证失败

servers: # 服务器列表
  - name: "server1" # 服务器名称，对应GroupID，v11下为数字，v12下为字符串
    address: "mc.example.top" # 服务器地址
    reconnect_interval: 5  # 重连间隔，单位秒，0为不重连，建议大一点不然被反作弊封号
    auth: "auth_example_online"  # 服务器的认证信息
    player_message_handler: "system"  # 这个是go-mc的一个bug，某些情况下玩家消息会被接收为system，这里可以设置为system来处理玩家消息

  - name: "Server2"
    address: "server.top"
    reconnect_interval: 5
    auth: "auth_example_offline"

auth: # 自定义认证信息
  auth_example_online: # 服务器的认证信息
    online: true  # 在线模式下使用交互式登录
    name:   # 在线模式下不生效

  auth_example_offline:
    online: false
    name: "Bot"

onebot: # OneBot配置
  reverse_ws: # 反向Websocket配置
    - address: "ws://127.0.0.1:8080/onebot/v12"
      reconnect_interval: 5  # 重连间隔，单位秒，0为不重连
      access_token:  # 反向Websocket的访问令牌(可选，下同)
      protocol_version: 12  # OneBot协议版本，v11或v12(可选，下同，默认12)
  forward_ws: # 正向Websocket配置
    - host: "127.0.0.1"
      port: 8080
      access_token:
      protocol_version: 12
  http_post: # 反向HTTP配置
    - address:
      access_token:
      protocol_version: 12
  http: # 正向HTTP配置
    - host: "127.0.0.1"
      port: 8080
      access_token:
      protocol_version: 12

  bot:
    heartbeat_interval: 20  # 心跳间隔，单位秒
    self_id: 114514  # 机器人的ID，v11下需使用数字形式，v12下使用字符串形式，要兼容的话请使用字符串数字形式
    player_id_type: "uuid_int"  # 为兼容onebotv11，用户的ID传输模式，uuid_int为UUID的整数形式(在某些语言下可能无法处理较大的整数，例如Javascript)，primary_key为使用数据库自增主键作为id(自动映射，但是在不同的bot上该值也不同)
```

### 版本支持
协议库使用的是`go-mc`，一个构建仅支持一个Minecraft版本，如有其他Minecraft版本需求请修改`go.mod`自行构建
如果你是腐竹，可以使用`ViaVersion`插件来支持多个版本

### TODO
- [ ] 使用地图画插件来支持图片消息(需自备图片外链)
- [ ] 支持@player
- [ ] `player_message_handler`自定义聊天样式的支持

### 已知Bug
- 某些情况下玩家消息会被接收为`system`，所以提供了`player_message_handler`选项来决定是否从系统消息中正则匹配提取玩家消息，
当然这也会导致某些情况下系统消息被误判为玩家消息，本人能力有限，后续了解MC协议后会尝试修复。与`go-mc`的作者交流后推测应该是Mojang在`1.19.1`加入的聊天消息验证所致，该问题大概率出在非原版服务端上，
目前测试发现`Paper` `Spigot` `purpur`均有此问题

### 感谢以下项目
- [go-mc](https://github.com/Tnze/go-mc)
- [go-libonebot](https://github.com/botuniverse/go-libonebot)