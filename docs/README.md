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
| [requirement](requirement/README.md) | 用户故事、目标、范围、功能与非功能需求 | [需求](templates/requirement.md) |
| [design](design/README.md) | 流程、架构、目录分层、代码落点与取舍 | [设计](templates/design.md) |
| [prd](prd/README.md) | 本增量业务规则与验收条件 | [产品规格](templates/prd.md) |
| [plan](plan/README.md) | SDD 任务、依赖、写入范围与验证方法 | [计划](templates/plan.md) |
| [research](research/README.md) | 有来源的分析结论、建议与局限 | [研究](templates/research.md) |

当前尚无产品主题。新增后在本页列出主题与各阶段链接，无需单独维护 INDEX。

## 使用约定

协作、敏捷/SDD 和文档规则统一见根 [AGENTS.md](../AGENTS.md)，工程分层与编码规范见 [PROJECT.md](../PROJECT.md)，视觉标准见 [DESIGN.md](../DESIGN.md)。这些文件在 vault 外，可由 agent 或代码编辑器读取。

在目标目录新建笔记，用内置 `Templates: Insert template` 插入对应模板，或直接复制并替换变量；删去不适用项。模板只填充标题和创建时日期，不自动分配编号或更新版本。

验收条件和验证方法留在文档，实际截图、日志、报告与中间产物留在库外。`.obsidian` 只共享基本配置；本地布局和缓存不进入 Git。旧文档不恢复。
