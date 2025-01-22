# 任务依赖管理系统

一个用于管理和可视化任务依赖关系的 Web 应用系统。

## 功能特点

- 📊 可视化任务依赖关系
  - 使用 D3.js 力导向图展示任务节点和依赖关系
  - 支持缩放、拖拽和交互操作
  - 直观展示任务之间的依赖链路

- 📝 任务管理
  - 创建、编辑和删除任务
  - 设置任务名称和描述
  - 管理任务间的依赖关系
  - 查看任务详细信息

- 🔄 依赖分析
  - 自动计算任务的前置依赖
  - 生成最优执行顺序
  - 识别循环依赖
  - 展示依赖路径

- 📋 执行计划
  - 选择多个目标任务
  - 自动包含所需的前置任务
  - 生成优化的执行顺序
  - 导出执行计划

## 技术栈

### 前端
- Vue 3 + TypeScript
- Element Plus UI 框架
- D3.js 可视化库
- Vite 构建工具

### 后端
- Go
- Gin Web 框架
- SQLite 数据库

## 快速开始

### 前端启动
```bash
# 进入前端目录
cd task-viewer

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 后端启动
```bash
# 进入后端目录
cd task-viewer/backend

# 启动服务
go run main.go
```

访问 http://localhost:5174 即可使用系统。

## 项目结构

```
task-viewer/
├── src/                # 前端源代码
│   ├── components/     # 组件
│   ├── views/         # 页面
│   ├── router/        # 路由配置
│   └── main.ts        # 入口文件
├── backend/           # 后端源代码
│   ├── api/          # API 接口
│   ├── model/        # 数据模型
│   ├── service/      # 业务逻辑
│   ├── db/           # 数据库操作
│   └── main.go       # 入口文件
└── docs/             # 文档
    ├── 1.0-功能规划.md
    ├── 1.1-前端页面设计-总览.md
    ├── 1.2-页面设计-*.md
    └── 1.3-总体架构设计.md
```

## API 接口

### 任务管理
- `GET /api/tasks` - 获取所有任务
- `POST /api/tasks` - 创建新任务
- `GET /api/tasks/:id` - 获取任务详情
- `PUT /api/tasks/:id` - 更新任务
- `DELETE /api/tasks/:id` - 删除任务

### 依赖管理
- `POST /api/tasks/:id/dependencies` - 添加依赖关系
- `DELETE /api/tasks/:id/dependencies/:depId` - 删除依赖关系

## 开发团队

- 前端开发：[Your Name]
- 后端开发：[Your Name]
- 设计：[Your Name]

## 许可证

MIT License 