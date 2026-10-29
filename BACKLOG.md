# WEARLOG Backlog

更新：2026-10-29。只保留未完成任务；认领与协作规则见 [AGENTS](AGENTS.md)，接手见 [HANDOVER](HANDOVER.md)。由协调者串行维护，共享清单不提供并发锁。

| 任务 | 目标 / 规格 | 优先级 | 状态 | 依赖 / 阻塞 | 所有者 | 写入范围 |
| --- | --- | --- | --- | --- | --- | --- |
| PLAN-01 / T-07 | 人物环绕与立体衣橱端到端全链路串联与三套衣物回归；[计划](docs/plan/01-人物环绕.md) revision 4 | 主线 | in-review | T-01~T-06 已全部实现并通过双规范代码审查与全量测试 | 主 agent | `frontend/`、`backend/` |

状态：needs-input / ready / in-progress / in-review / blocked；完成且经集成人确认后移出清单，摘要归 CHANGELOG。ready 须目标、验收与依赖清楚；认领须写明实际 agent、独占路径和共享文件写入者，重叠任务由协调者安排串行。产品任务用 plan 中的任务编号，不复制规格正文。
