# Mobile Workbench v2 Automated Audit

> Date: 2026-08-05 (initial run), 2026-08-06 and 2026-08-11 (post-fix regression runs)
> Method: headless Chrome (CDP) against the local dev stack (Vite + Denova backend), driven by temporary Node CDP scripts under `Z:\nova\.scratch-mobile-audit\` and `Z:\nova\.scratch-mobile-settings-audit\` (removed after each run). The scripts emulated mobile viewports and touch input, overrode CSS safe-area insets through `Emulation.setSafeAreaInsetsOverride`, dumped the AX tree, measured interactive element hit areas, and drove key mobile flows. Navigation clicks use DOM clicks so transient toasts cannot block the flow; Escape is dispatched as a real key event.

## Layout results

| Width | Theme | Horizontal overflow | Primary navigation | Nav placement |
| --- | --- | --- | --- | --- |
| 320 | dark | none | 正文 / 项目 / Agent / 更多 | bottom |
| 390 | dark | none | 正文 / 项目 / Agent / 更多 | bottom |
| 430 | dark | none | 正文 / 项目 / Agent / 更多 | bottom |
| 390 | light | none | 正文 / 项目 / Agent / 更多 | bottom |
| 768 | dark | none | 正文 / 项目 / Agent / 更多 | side |

Pixel sampling (deviceScaleFactor 2, 390x844): dark theme mean luminance 10, light theme 249, confirming both themes render and screenshots are non-blank.

## Safe-area simulation (390px dark)

With `Emulation.setSafeAreaInsetsOverride { top: 47, bottom: 34 }`:

| Check | Expected | Measured | Pass |
| --- | --- | --- | --- |
| Header top padding | 47px | 47px | yes |
| Bottom nav padding | 34px | 34px | yes |
| Nav bottom edge | touches viewport bottom | 844 / 844 | yes |
| Horizontal overflow | none | none | yes |

The page also ships `viewport-fit=cover` and `env(safe-area-inset-*)` declarations in the mobile shell, so the emulation exercised the real CSS path.

## Touch-target results

Surfaces that meet the 44px+ recommendation: primary nav buttons (48px), More menu items (48-56px), Task Center task rows (80px).

Measured after the 2026-08-06 fix (390px dark):

| Surface | Control | Before | After |
| --- | --- | --- | --- |
| Top bar | project switcher (切换书籍) | 32px | 44px |
| More | mode switch 写作模式 / 游戏模式 | 40px | 44px |
| Task Center | back button (返回更多) | 40x40 | 44x44 |
| First-run card | 和创作 Agent 聊灵感 | 34px | 44px |
| First-run card | 跳过引导 pill | 24px | 44px |
| First-run card | 去配置 | 24px | 44px |
| First-run card | close button | 24x24 | 44x44 |
| Toast | close button (coarse pointer) | 22x22 | 44x44 |

320 / 390 / 430px regression runs report zero visible interactive targets below 44px on the audited shell surfaces and no horizontal overflow.

### Settings touch-target regression (2026-08-11)

The mobile Settings page initially exposed 66 visible controls at 320x844; 64 were below the 44px recommendation. The close action measured 24x24, the pane trigger 40x40, native fields 34-36px high, shadcn inputs 32px, and compact Radix select triggers 28-40px.

After the fix, the same CDP audit reported zero failures at 320, 390, and 430px with coarse-pointer emulation and no horizontal overflow. The 320px English/light run also reported zero failures. Expanded Radix context-window options each measured 44px high. Fine-pointer desktop Settings retains its compact density.

## Accessibility tree results

- One named `navigation` landmark per shell (移动端工作台导航), with all four destinations named.
- Zero unnamed controls on the audited surfaces (bookshelf/onboarding, More, Task Center, Settings).
- No `<img>` without `alt`; headings are present for Command Palette and Task Center.
- Opening the Task Center moves focus to its back control; Escape and the header back control both return focus to the More button (verified in the 2026-08-06 run).

## System notification flow

- Settings page shows 启用系统通知 with the hint 默认关闭，应用不会在首次启动时请求权限.
- CDP permission flow: initial `default` -> grant -> `Notification` construction succeeds -> reset -> `requestPermission()` resolves (`denied` in headless) without throwing.
- The OS-level permission prompt itself still requires a real browser/device check.

## Interaction regression (390px dark)

- Dismiss first-run card, More opens with mode switch group and Task Center.
- Task Center opens (`#mobile-task-center-title`), Escape closes it and keeps More active.

## Findings for follow-up

1. Real-device notch insets, real notification permission prompt, screen-reader reading order, and final pixel-level human review remain external checks.
