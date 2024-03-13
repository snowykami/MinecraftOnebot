## MCOnebot

#### 本程序部分代码依靠bug运行，若无法正确运行可能是bug被修了
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
# 自行查看`config.example.yaml`文件
```

### 版本支持
协议库使用的是`go-mc`，一个构建仅支持一个Minecraft版本，如有其他Minecraft版本需求请修改`go.mod`自行构建

### 感谢以下项目
- [go-mc](https://github.com/Tnze/go-mc)
- [go-libonebot](https://github.com/botuniverse/go-libonebot)