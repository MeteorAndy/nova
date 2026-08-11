# 可执行配置返回保护设计（故事 85）

## 背景

故事 85 要求：可执行配置（Agents、Skills、自动化、模型、导演规则）存在未保存修改时，返回操作必须提供“继续编辑”或“放弃修改”，不能静默应用半成品或丢失草稿。

当前各配置面草稿机制不一致：

- Agents：`useLayeredSettingsDraft` + 自动保存通道，草稿随编辑自动排队保存；
- Skills：`useResourceAutosave`，无效草稿暂停保存但保留；
- Automations：独立定义草稿，显式保存；
- 游戏模式 SettingPanel：`usePresetResourceAutosave`，无效草稿保留。

直接在各面复制“返回确认”会产生不一致行为，因此需要一个共享守卫。

## 目标

- 在移动端“更多 → Skills/Agents/自动化”以及游戏模式 SettingPanel 返回时，若当前面存在未保存/待保存的可执行配置草稿，弹出统一对话框：
  - “继续编辑”：留在当前面，草稿保留；
  - “放弃修改”：丢弃草稿并继续返回。
- 各面负责向守卫注册“是否有草稿”与“如何放弃”，守卫不关心具体保存语义。
- 桌面端关闭面板/切换模式时复用同一守卫。

## 方案

### 1. 共享注册表

新增 `web/src/features/config-guard/executable-draft-guard.ts`（Zustand store）：

```ts
interface ExecutableDraftGuardEntry {
  hasPending: boolean
  discard: () => void | Promise<void>
}

interface ExecutableDraftGuardState {
  entries: Record<string, ExecutableDraftGuardEntry>
  register: (key: string, entry: ExecutableDraftGuardEntry) => void
  unregister: (key: string) => void
}
```

### 2. 统一对话框

新增 `UnsavedConfigGuardDialog`：读取注册表 `hasPending`，提供“继续编辑 / 放弃修改”两个动作；放弃后调用对应 `discard` 并放行导航。

### 3. 接线

- `WorkbenchShell` 的返回/Esc/模式切换：离开 `skills`、`agents`、`automations` 前先询问守卫。
- 游戏模式 `SettingPanel` 返回故事/导演台前询问守卫。
- Agents/Skills/Automations/SettingPanel 各在自己的草稿 hook 中 `register/unregister`，并给出各自的 `discard`（Agents 调 `reload()`、Skills 调内容草稿重置、Automations 清空未保存定义等）。

### 4. 各面接入顺序

1. Automations（草稿语义最明确，先做红测）；
2. Skills；
3. 游戏模式 SettingPanel；
4. Agents（需先明确其自动保存语义与“未保存”定义）。

## 验收

- 各面编辑草稿后按返回：出现“继续编辑 / 放弃修改”，选择“继续编辑”留在原面且草稿保留；选择“放弃修改”后草稿消失并完成返回。
- 无草稿时不弹窗，返回行为不变。
- 中英文文案齐全；覆盖矩阵故事 85 转 ✅。

## 待确认

- Agents 的“未保存”定义：是仅指保存中/保存失败/无效草稿，还是也包含自动保存队列中的草稿？
- 是否只保护移动端返回，还是桌面端关闭/切换模式也统一拦截？
