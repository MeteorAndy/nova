# 移动工作台状态恢复设计（故事 24 / 33）

## 背景

故事 24 要求“项目”工作面在切换目的地或模式后恢复搜索、筛选、滚动与选中项；故事 33 要求每个 Agent 会话保留输入草稿、滚动位置与运行状态。输入草稿与运行状态已分别通过会话级 `draftKey` 和任务恢复实现，剩余缺口是侧栏搜索/筛选/滚动与消息列表滚动。

2026-08-11 曾尝试让移动工作台目的地内容常驻挂载，结果破坏共享主面板与 portal 的装配，导致“保存失败阻止项目切换”的应用级保护失效，已回退。因此本设计不采用常驻挂载，而是在状态所有者（ModeRouter）持久化可序列化状态。

## 目标

- 切换“正文 / 项目 / Agent / 更多”后再返回，项目侧栏的当前视图（作品目录/项目文件/全局搜索）、搜索词、筛选状态、滚动位置和选中文件保持原样。
- 切换 Agent 会话后，消息列表滚动位置按会话恢复；运行状态与输入草稿保持现有行为。
- 不改变 MobileWorkbenchShell 的主面板装配，不触碰共享 portal。

## 方案

### 1. 项目侧栏状态（故事 24）

在 ModeRouter 增加按 `workspace` 分桶的内存状态（可选同步到 `localStorage`）：

```ts
interface ProjectSurfaceState {
  view: 'outline' | 'files' | 'search'
  searchQuery: string
  expandedKeys: string[]
  scrollTop: number
  selectedFile: string | null
}

const projectSurfaceState = new Map<string, ProjectSurfaceState>()
```

- FileTree 的搜索词、展开键与 `scrollTop` 通过受控 props 或 ref 上报到 ModeRouter。
- 切换目的地/模式时不销毁该状态；切回“项目”时按 workspace 恢复。
- 状态只保存在内存（单次会话），不做跨会话持久化，避免与版本恢复/离线草稿语义混淆。
- 落地点：`ModeRouter` 持有 Map，`ProjectSurface`（FileTree/SearchPanel 容器）通过 props 读写。

### 2. Agent 消息滚动（故事 33）

在 `useAgentChat` 增加按会话的滚动位置表：

```ts
const sessionScroll = new Map<string, number>()
```

- `MessageList` 暴露 `onScrollCapture(scrollTop)` 与 `scrollRestoreTop`。
- 切换会话前记录旧会话 `scrollTop`；新会话历史加载完成后，若该会话有记录则恢复，否则保持底部锁定语义。
- 滚动恢复只发生在用户离开过该会话的场景；新消息到达仍遵循现有底部跟随逻辑。

## 验收

- 项目侧栏输入搜索词、滚动到中部、选中文件后，切到“正文”再返回，搜索词/滚动/选中项一致。
- Agent 会话 A 滚动到中部后切到会话 B，再切回 A，滚动位置一致；B 的新消息仍按底部跟随策略。
- 覆盖矩阵故事 24、33 的滚动部分由新集成测试证明后转 ✅。

## 不做

- 不做目的地常驻挂载（已证伪）。
- 不做跨会话持久化搜索/滚动（避免语义混淆）。
