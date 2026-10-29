# WEARLOG 知识库

将 `wearlog-server/docs` 作为本地文件夹在 Obsidian 中打开即可。这里保存长期产品文档；任务认领和进度见根 [BACKLOG](../BACKLOG.md)，接手信息见 [HANDOVER](../HANDOVER.md)。

## 内容组织

主链路为 **requirement → design → prd → plan**，调研分析归 research。五个目录及各自 README 固定保留，各阶段按同一主题关联；有实际内容时新增主题文档。

```text
docs/
├── requirement/
├── design/
├── prd/
├── plan/
├── research/
├── templates/
└── README.md
```

| 目录 | 内容 | 模板 |
| --- | --- | --- |
| [requirement](requirement/README.md) | 用户故事、目标、范围、功能与非功能需求 | [需求](templates/01-需求模板.md) |
| [design](design/README.md) | 流程、架构、目录分层、代码落点与取舍 | [设计](templates/02-设计模板.md) |
| [prd](prd/README.md) | 本增量业务规则与验收条件 | [产品规格](templates/03-产品规格模板.md) |
| [plan](plan/README.md) | SDD 任务、依赖、写入范围与验证方法 | [计划](templates/04-执行计划模板.md) |
| [research](research/README.md) | 有来源的分析结论、建议与局限 | [研究](templates/05-调研分析模板.md) |

## 当前主题

| 主题 | 需求 → 设计 → PRD → 计划 | 调研 | 当前边界 |
| --- | --- | --- | --- |
| 01 人物环绕 | [REQ-01](requirement/01-人物环绕.md) → [DES-01](design/01-人物环绕.md) → [PRD-01](prd/01-人物环绕.md) → [PLAN-01](plan/01-人物环绕.md) | [RES-01](research/01-人物环绕.md) | 对标 Woo 1:1 复刻（revision 3）；锁定“纯 API 聚合网关 + 端侧 onnxruntime-web 零成本深度视差 + 2D 原生姿态换装 + 16 正交关键帧 WebP Sprite Sheet 闭环”架构，零自建 GPU 算力 |

本主题以前端 Vercel (Next.js) + shadcn + Radix + Tailwind v4 与后端 Go (Gin/Huma) 为工程基线，通过 `/openapi.json` 与 `npm run openapi` 强类型契约打通全链路，实现低延迟、低成本的 360° 立体数字衣橱。

## 文档编号

正式文档使用 `01-主题.md`、`02-主题.md`，按主题建立顺序递增。同一主题在 requirement → design → prd → plan 中使用相同序号，例如 `requirement/01-衣物录入.md` 对应 `design/01-衣物录入.md`；research 独立编号。

编号由协调者核对后分配，不重复、不因删减重排，各目录 README 按编号顺序列出文档。README 与规范文件不加编号，模板按链路顺序采用编号与中文名称；模板序号不替代主题序号。主题文档的 id 使用对应前缀，例如 `REQ-01`、`DES-01`、`PRD-01`、`PLAN-01`、`RES-01`。示例仅说明命名，不代表已录入产品需求。

## 使用约定

协作、敏捷/SDD 和文档规则统一见根 [AGENTS.md](../AGENTS.md)，工程分层与编码规范见 [PROJECT.md](../PROJECT.md)，视觉标准见 [DESIGN.md](../DESIGN.md)。这些文件在 vault 外，可由 agent 或代码编辑器读取。

在目标目录新建笔记，用内置 `Templates: Insert template` 插入对应模板，或直接复制并替换变量；删去不适用项。模板只填充标题和创建时日期，不自动分配编号或更新版本。

验收条件和验证方法留在文档，实际截图、日志、报告与中间产物留在库外。`.obsidian` 只共享基本配置；本地布局和缓存不进入 Git。旧文档不恢复。
