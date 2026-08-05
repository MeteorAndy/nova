# Mobile Workbench v2 Browser Acceptance

> Date: 2026-08-04
> Method: headless Chrome (CDP) against the local dev stack (Vite + Denova backend), driven by `Z:\nova\.acceptance\cdp-acceptance.mjs`. Layout is measured programmatically (viewport width vs. document scroll width) and key interactions are exercised with real input events.

## Results

| Width | Locale | Theme | Horizontal overflow | Primary navigation |
| --- | --- | --- | --- | --- |
| 320 | zh-CN | dark | none | 正文 / 项目 / Agent / 更多 |
| 360 | zh-CN | dark | none | 正文 / 项目 / Agent / 更多 |
| 390 | zh-CN | dark | none | 正文 / 项目 / Agent / 更多 |
| 430 | zh-CN | dark | none | 正文 / 项目 / Agent / 更多 |
| 768 | zh-CN | dark | none | 正文 / 项目 / Agent / 更多 |
| 1024 | zh-CN | dark | none | 正文 / 项目 / Agent / 更多 |
| 1280 | zh-CN | dark | none | desktop activity bar (no mobile nav) |
| 390 | en-US | light | none | Manuscript / Project / Agent / More |
| 1280 | en-US | light | none | desktop activity bar |
| 390 | zh-CN | light | none | 正文 / 项目 / Agent / 更多 |
| 768 | en-US | dark | none | Manuscript / Project / Agent / More |

## Interaction checks (390px, zh-CN, dark)

- Open the created book from the bookshelf: the writing manuscript surface becomes active.
- Open More: mode switch group and Task Center are visible.
- Open Task Center: `#mobile-task-center-title` is present.
- Press Escape: Task Center closes and More remains active (the unified back path).

## 200% zoom (390px, zh-CN, dark)

- Applying `zoom: 2` produces no horizontal overflow (`document.documentElement.scrollWidth` stays at viewport width).

## Landscape (2026-08-04)

| Viewport | Locale | Theme | Horizontal overflow | Primary navigation |
| --- | --- | --- | --- | --- |
| 844x390 | zh-CN | dark | none | 正文 / 项目 / Agent / 更多 |
| 932x430 | zh-CN | dark | none | 正文 / 项目 / Agent / 更多 |
| 844x390 | en-US | light | none | Manuscript / Project / Agent / More |
| 1180x820 | zh-CN | dark | none | desktop activity bar |
| 1280x800 | zh-CN | dark | none | desktop activity bar |

## Long English book name (2026-08-04)

| Width | Horizontal overflow | Switcher accessible name | Primary navigation |
| --- | --- | --- | --- |
| 320 | none | full long title preserved | 正文 / 项目 / Agent / 更多 |
| 390 | none | full long title preserved | 正文 / 项目 / Agent / 更多 |
| 430 | none | full long title preserved | 正文 / 项目 / Agent / 更多 |

## Notes and limitations

- The acceptance ran against a fresh local workspace (`acceptance-book`) with no chapters, so it covers shell/navigation/layout rather than long-form editor content.
- Pixel-level visual review was not possible in this session; layout invariants and interactions were verified programmatically via CDP.
- A 404 for the missing book cover image was observed; it is expected for a book without a cover and is not a layout defect.
- Safe-area emulation, touch targets, AX-tree, and notification plumbing are covered by `docs/verification/mobile-workbench-automated-audit.md`; real-device insets and the OS permission prompt still need a device/browser check.
