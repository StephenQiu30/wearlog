# WEARLOG Handover

更新：2026-10-29。原地维护当前快照；任务状态见 [BACKLOG](BACKLOG.md)，文档入口见 [docs](docs/README.md)。

## 当前接手

- 协调者 / 集成人：本聊天主 agent；没有已分派的产品实现任务。
- 核对时仓库为 `wearlog-server`，分支 `main`；接手先重新核对 Git 状态与基线。
- 本轮文档升级：针对 Woo 1:1 复刻、零自建 GPU 算力约束与明确的前后端技术栈（Next.js + shadcn + Radix + Tailwind v4 + Axios + umiopenapi 与 Go /backend），完成 [01 人物环绕](docs/README.md) 全套规格（REQ/DES/PRD/PLAN/RES-01）向 **revision 3** 升级。
  - 核心架构锁定：纯 API 聚合网关（/backend Go 封装调用 Gemini 2.0 Flash、Fal.ai、FASHN.ai）+ 端侧 onnxruntime-web 零成本 2.5D 深度视差 + 2D 原生姿态换装 + 16 关键帧 WebP Sprite Sheet 闭环环绕；
  - 接口契约：Go Huma 导出 `/openapi.json`，前端通过 `npm run openapi` 自动生成全类型 Axios 客户端。
- 验证：本地链接无缺失，id/status/revision/updated 保持一致，计划与 AC 闭环对齐。
- 下一动作：执行 T-01/T-02 搭建 `/backend` Orbit 路由，导出 `/openapi.json` 并通过 `npm run openapi` 生成前端 SDK。

## 有开发任务时补充

只记录任务编号与规格版本、实际 owner、worktree/分支/SHA、独占和未提交文件、已完成/剩余行为、验证摘要与库外位置、阻塞及第一条接手动作。共享文件由指定写入者更新；交接明确转移所有权，不追加聊天流水账。
