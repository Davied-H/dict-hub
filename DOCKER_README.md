# Dict-Hub

现代化的 MDX 词典查询平台，支持多词典管理、快速查词与历史记录追踪。

## 功能特性

- **快速查词** - 毫秒级响应，即时获取释义
- **多词典支持** - 同时加载多个 MDX 格式词典
- **词频统计** - 自动记录查词频率，了解学习进度
- **历史记录** - 完整查询历史追踪
- **现代 UI** - 简洁优雅的界面，支持深色模式

## 快速开始

### Docker Compose（推荐）

```yaml
version: '3.8'
services:
  dict-hub:
    image: dongdonghe122/dict-hub
    container_name: dict-hub
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - ./dicts:/app/dicts
    environment:
      - SERVER_MODE=production
      - MDX_AUTO_LOAD=true
      - TZ=Asia/Shanghai
```

```bash
# 创建目录
mkdir -p data dicts/source dicts/sound

# 启动服务
docker-compose up -d
```

### Docker Run

```bash
docker run -d \
  --name dict-hub \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/dicts:/app/dicts \
  -e SERVER_MODE=production \
  -e MDX_AUTO_LOAD=true \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  dongdonghe122/dict-hub
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `SERVER_PORT` | `8080` | 服务端口 |
| `SERVER_MODE` | `development` | 运行模式 (`development` / `production`) |
| `MDX_AUTO_LOAD` | `false` | 启动时是否自动加载词典 |
| `TZ` | `Asia/Shanghai` | 时区设置 |

## 数据卷

| 容器路径 | 说明 |
|----------|------|
| `/app/data` | SQLite 数据库存储目录 |
| `/app/dicts` | 词典文件根目录 |
| `/app/dicts/source` | MDX 词典源文件目录 |
| `/app/dicts/sound` | 音频文件目录 |

## 使用方法

1. 将 `.mdx` 和 `.mdd` 词典文件放入 `dicts/source` 目录
2. 将音频文件（如有）放入 `dicts/sound` 目录
3. 访问 `http://localhost:8080` 打开 Web 界面
4. 在设置中启用词典，开始查词

## 支持的架构

- `linux/amd64`
- `linux/arm64`

## 链接

- **GitHub**: [https://github.com/Davied-H/dict-hub](https://github.com/Davied-H/dict-hub)
- **问题反馈**: [https://github.com/Davied-H/dict-hub/issues](https://github.com/Davied-H/dict-hub/issues)

## 许可证

MIT License
