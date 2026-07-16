# YUTANK 商城系统

[English version](README.md)

YUTANK 是一个面向企业员工的 Credits 商品兑换系统。项目包含管理员后台、员工商城和 GoFrame 后端服务，覆盖员工与 Credits 管理、商品与库存管理、订单处理以及管理员权限控制。

## 项目组成

| 模块 | 技术 | 默认地址 | 用途 |
| --- | --- | --- | --- |
| 后端服务 | GoFrame 2 + MySQL | `http://127.0.0.1:8000` | 提供认证、商城和管理接口 |
| 管理端 | Vue 2 + Element UI | `http://localhost:8080` | 管理员工、Credits、商品、库存、订单和权限 |
| 员工商城 | Vue 2 + Element UI | `http://localhost:9528` | 浏览商品、使用 Credits 兑换商品和管理个人订单 |

后端接口使用统一响应格式：

```json
{
  "code": 0,
  "message": "Success",
  "data": {}
}
```

登录后的请求通过以下请求头携带 token：

```text
Authorization: Bearer <token>
```

## 当前功能

### 管理端

- **Dashboard**：查看当前管理员和系统模块概览。
- **Employees**：查看详情、新增、编辑、启用或禁用、重置密码和软删除员工。
- **Credit Operations**：为员工增加或扣除 Credits，并查询员工 Credits 流水。
- **Categories**：查看商品分类。
- **Goods**：筛选、新增、编辑和查看商品，上传商品图片，以及控制上架或下架状态。
- **Stock**：增加或减少商品库存，并按商品、日期和变更类型查询库存流水。
- **Orders**：按订单号、员工、状态和日期查询订单，查看商品明细，完成或取消待处理订单。
- **Access Control**：管理管理员、角色和接口权限，支持状态控制、角色分配和权限分配。
- **Account Security**：修改当前管理员密码和退出登录。

### 员工商城

- 注册、登录和退出员工账号。
- 按分类或关键字浏览商品并查看商品详情。
- 将商品加入购物车、调整数量或移除商品。
- 使用 Credits 提交兑换订单。
- 查看订单列表和详情，并取消符合条件的订单。
- 查看个人 Credits 余额和流水。
- 修改当前员工账号密码。

## 技术栈

- Go `1.25.5`
- GoFrame CLI `2.10.0`
- GoFrame `2.9.5`
- MySQL `8.0.x`
- Vue `2.6.10`
- Element UI `2.13.2`
- Node.js `22.22.2`
- npm `10.9.7`

以上版本是当前项目已验证的本地开发环境。建议所有命令在 WSL Ubuntu 20.04 中执行。

## 目录结构

```text
BIT303_shop/
├── api/                         # 后端 API 请求与响应契约
│   ├── backend/                 # 管理端及员工账号接口契约
│   ├── frontend/                # 员工商城接口契约
│   └── base/                    # 公共接口契约
├── frontend_manage/             # 管理端 Vue 项目
├── frontend_web/                # 员工商城 Vue 项目
├── hack/init.sql                # 数据库表结构和基础数据
├── internal/                    # GoFrame controller、logic、service、dao 和 model
├── manifest/config/config.yaml  # 后端运行配置
├── resource/public/upload/      # 管理端上传的商品图片
└── main.go                      # 后端入口
```

## 环境准备

请先安装并确认以下工具可用：

```bash
go version
mysql --version
node --version
npm --version
```

如果 MySQL 尚未启动，可在 Ubuntu 中运行：

```bash
sudo service mysql start
```

## 配置数据库

默认配置文件为 `manifest/config/config.yaml`：

```yaml
database:
  default:
    type: mysql
    host: 127.0.0.1
    port: "3306"
    user: root
    pass: "123456"
    name: bit303_shop
```

注意！！！请根据目标环境修改数据库地址、账号和密码。默认配置只适合本地开发，不应直接用于生产环境。

## 初始化数据库

进入项目根目录并导入初始化 SQL：

```bash
cd /path/to/BIT303_shop
mysql -uroot -p123456 < hack/init.sql
```

`hack/init.sql` 会创建 `bit303_shop` 数据库、业务表、基础分类、权限数据，以及默认管理员和演示员工账号。脚本使用幂等写法，可在结构升级后重新执行。

## 启动项目

项目需要分别启动后端、管理端和员工商城。建议为三个服务分别打开终端。

### 1. 启动后端

```bash
cd /path/to/BIT303_shop
go mod download
go run main.go
```

启动后可访问：

- 健康检查：`http://127.0.0.1:8000/health`
- Swagger：`http://127.0.0.1:8000/swagger/`
- OpenAPI JSON：`http://127.0.0.1:8000/api.json`

### 2. 启动管理端

```bash
cd /path/to/BIT303_shop/frontend_manage
npm install
npm run dev
```

登录地址：`http://localhost:8080/#/login`

### 3. 启动员工商城

```bash
cd /path/to/BIT303_shop/frontend_web
npm install
npm run dev
```

商城地址：`http://localhost:9528`

两个前端默认请求 `http://127.0.0.1:8000`。如后端部署在其他地址，可通过 `VUE_APP_BASE_API` 调整请求地址。

## 测试账号

| 使用位置 | 身份 | 账号 | 密码 | 说明 |
| --- | --- | --- | --- | --- |
| 管理端 | 超级管理员 | `root` | `123456` | 显示名称为 `System Administrator`，拥有全部管理权限 |
| 员工商城 | 演示员工 | `root` | `123456` | 可用于登录和检查员工商城；也可以在商城注册新员工 |

管理员和员工是两套独立的认证身份，即使演示账号名称相同，token 也不能混用。交付或部署后请立即修改默认密码。

## 管理端操作指引

### 员工和 Credits

1. 使用超级管理员账号登录管理端。
2. 在 **Employees** 中创建或维护员工账号。
3. 在 **Credit Operations** 中选择员工，为其增加初始 Credits。
4. 可在同一页面查询该员工的 Credits 增减记录。

### 商品和库存

1. 在 **Categories** 中确认目标分类存在且已启用。
2. 进入 **Goods**，点击新增商品。
3. 选择分类，填写英文商品名称、Credits Price、库存和说明。
4. 使用上传控件上传商品图片，也可以手动填写图片 URL。
5. 保存后确认商品状态为上架，员工商城才会显示该商品。
6. 后续库存调整和库存流水统一在 **Stock** 中管理。

### 订单

1. 在 **Orders** 中按订单号、员工 ID、状态或日期筛选订单。
2. 打开详情查看订单商品快照和 Credits 合计。
3. 待处理订单可以执行 **Complete** 或 **Cancel**。
4. 已完成或已取消的订单不能重复处理。

### 管理员权限

1. 在 **Access Control / Permissions** 中维护接口权限。
2. 在 **Roles** 中创建角色并分配权限。
3. 在 **Admins** 中创建管理员并分配角色。
4. 超级管理员不受普通角色权限限制；普通管理员由后端权限中间件执行最终校验。

## 员工商城操作指引

1. 打开商城后可以直接浏览已上架商品。
2. 登录或注册员工账号后，可进入购物车、订单、My Credits 和 Account Security。
3. 将需要兑换的商品加入购物车并确认数量。
4. 提交订单时，系统会校验 Credits 余额和商品库存。
5. 订单创建成功后会扣除 Credits 和库存；符合条件的订单取消后会按业务规则恢复。
6. 在 **My Orders** 查看订单状态和明细，在 **My Credits** 查看余额和流水。

## API 边界

- `/backend/admin/...`：管理员登录、退出、信息和密码接口，使用管理员 token。
- `/backend/employee/...`：员工注册、登录、信息和密码接口；其中管理接口需要管理员 token 和对应权限。
- `/backend/points/...`：员工个人 Credits 与管理员 Credits 操作接口。
- `/backend/goods/...`、`/backend/stock/...`、`/backend/order/...`：管理端商品、库存和订单接口。
- `/backend/role/...`、`/backend/permission/...`：管理员 RBAC 接口。
- `/frontend/category/...`、`/frontend/goods/...`：商城公开浏览接口。
- `/frontend/cart/...`、`/frontend/order/...`：需要员工 token 的购物车和订单接口。

接口响应字段、参数和完整路径以 Swagger 为准。

## 构建与验证

### 后端测试

```bash
cd /path/to/BIT303_shop
go test ./...
```

### 管理端生产构建

```bash
cd /path/to/BIT303_shop/frontend_manage
NODE_OPTIONS=--openssl-legacy-provider npm run build:prod
```

构建结果位于 `frontend_manage/dist`。

### 员工商城生产构建

```bash
cd /path/to/BIT303_shop/frontend_web
npm run build
```

构建结果位于 `frontend_web/dist`。如果 Node.js 报 OpenSSL 兼容错误，请改用：

```bash
NODE_OPTIONS=--openssl-legacy-provider npm run build
```

## 数据与交付注意事项

- 商品、订单、员工、Credits、角色和权限等业务记录保存在 MySQL 中，不会随 Git 代码自动同步。
- 管理端上传的图片文件保存在 `resource/public/upload`，上传记录和商品图片 URL 保存在数据库中。
- 完整迁移环境时，需要同时备份 MySQL 数据库和 `resource/public/upload`。
- 新环境只导入 `hack/init.sql` 时会获得基础账号、分类和权限，不会自动获得本地管理员后来创建的正式商品。
- 正式商品应在目标环境中通过管理端上传图片并创建，或者通过经过确认的数据库备份进行迁移。
- `Credits` 是页面统一使用的业务术语；数据库和 JSON 中部分历史字段仍使用 `points` 命名，不影响前端显示。

## 常见问题

### 登录提示账号或密码错误

- 确认后端和 MySQL 已启动。
- 确认已导入最新的 `hack/init.sql`。
- 管理端必须使用管理员身份登录，员工商城必须使用员工身份登录。
- 如果账号已被禁用或软删除，需要由管理员恢复或重新创建。

### 接口返回 401 或 403

- `401` 通常表示 token 缺失、失效或使用了错误身份的 token。
- `403` 通常表示普通管理员没有对应接口权限。
- 退出当前账号并重新登录；如果仍然返回 `403`，由超级管理员检查角色和权限分配。

### 商城没有商品

- 初始化 SQL 不包含管理员后续创建的商品数据。
- 请在管理端创建商品、上传图片、设置库存，并确认商品已上架。

### 商品图片无法显示

- 确认后端服务正在运行，因为 `/upload/...` 图片由后端提供。
- 确认 `resource/public/upload` 中存在对应文件且运行用户具有读取权限。
- 确认数据库中的图片 URL 与上传接口返回的 URL 一致。

### 前端启动或构建出现 OpenSSL 错误

Vue CLI 4 在较新的 Node.js 中可能需要兼容选项：

```bash
export NODE_OPTIONS=--openssl-legacy-provider
npm run dev
```

### 端口被占用

检查默认端口：

```bash
lsof -i :8000
lsof -i :8080
lsof -i :9528
```

关闭占用进程或为对应服务指定其他端口。修改后端地址时，也要同步配置两个前端的 `VUE_APP_BASE_API`。
