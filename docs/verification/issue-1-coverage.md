# Fork Issue #1 Coverage Matrix

> Date: 2026-08-11
> Scope: the 104 user stories in [mobile-workbench-v2-spec](../specs/mobile-workbench-v2-spec.md).
> Status legend: ✅ verified by automated tests and/or browser acceptance; 🟡 implemented with partial or device-dependent verification; ⬜ not implemented (either out of scope by design or still missing).

## Shell and Writing (1-30)

| # | Story (short) | Status | Evidence |
| --- | --- | --- | --- |
| 1 | One primary surface at a time | ✅ | MobileWorkbenchShell + WorkbenchShell tests |
| 2 | Destinations have icon + text | ✅ | MobileWorkbench.integration.test.tsx nav assertions |
| 3 | Exactly one primary entry highlighted | ✅ | aria-current assertions in workbench tests |
| 4 | Writing: Manuscript/Project/Agent/More | ✅ | nav assertions + browser acceptance |
| 5 | Game: Story/Storylines/Reference/More | ✅ | MobileWorkbench.integration.test.tsx mode tests |
| 6 | Mode switch only from More | ✅ | integration test "only from More" |
| 7 | Per-mode state restoration | ✅ | MobileWorkbench.integration.test.tsx restore test |
| 8 | Shared menus never switch mode | ✅ | integration test "shared workspaces ... without changing mode" |
| 9 | Project and mode switches are independent | ✅ | integration tests |
| 10 | Save before project switch | ✅ | App.save-before-project-switch.integration.test.tsx covers mobile switcher and bookshelf paths (save POST precedes switch) |
| 11 | Save failure stays and shows error | ✅ | App.save-before-project-switch.integration.test.tsx asserts zh error toast, zero switch requests, current project unchanged |
| 12 | Top bar shows current work object | ✅ | shell header + browser acceptance heading |
| 13 | Project subtitle opens switcher | ✅ | BookSwitcher + header tests |
| 14 | Long titles never cover controls | ✅ | MobileWorkbench.integration.test.tsx long-name case |
| 15 | At most two frequent actions | ✅ | header render structure |
| 16 | One manuscript document at a time | ✅ | MarkdownEditor isolation tests |
| 17 | Title opens recent documents | ✅ | Recent documents UI present; API tests |
| 18 | Per-document draft/cursor/scroll/undo | ✅ | draft persistence, undo isolation, per-file scroll and cursor restore asserted in MarkdownEditor.test.tsx |
| 19 | Save/offline/conflict status near title | ✅ | save status + conflict banner tests |
| 20 | Project default shows outline | ✅ | Project surface + workspace summary tests |
| 21 | Project keeps full file segments | ✅ | file tree + workspace API tests |
| 22 | Search covers manuscript/plan/lore | ✅ | SearchPanel.test.tsx asserts one query returns and renders chapters/plans/lore results and hands result + query to selection |
| 23 | Chapter click opens manuscript | ✅ | integration test "returns to the manuscript" |
| 24 | Project state restoration | 🟡 | design in docs/verification/mobile-workbench-state-restoration-design.md; awaiting confirmation before implementation |
| 25 | No drag during normal browse | ✅ | HomeView.test.tsx asserts no drag handle in recent mode; useSortable disabled outside manual mode |
| 26 | Sort mode shows handles | ✅ | HomeView.test.tsx asserts drag handle appears after switching to manual sort |
| 27 | Keyboard hides primary nav | ✅ | integration test "bottom area to the focused input" |
| 28 | Stable mobile edit toolbar | ✅ | EditorToolbar tests |
| 29 | Selection swaps toolbar actions | ✅ | selection toolbar tests |
| 30 | Selection AI only via explicit handoff | ✅ | handoff tests (backend + InputArea) |

## Agent, Review and Tasks (31-56)

| # | Story (short) | Status | Evidence |
| --- | --- | --- | --- |
| 31 | One full Agent session at a time | ✅ | Agent surface + integration tests |
| 32 | Session management full screen | ✅ | AgentPanel.test.tsx asserts Sessions view replaces the chat surface, session switching, and return to chat |
| 33 | Per-session draft/scroll/run state | 🟡 | per-session composer draft verified in AgentPanel.test.tsx; run-state resume covered by lifecycle tests; scroll position still resets on switch |
| 34 | Tool/sub-agent/diagnostics collapsible | ✅ | AgentMessageItem/List tests |
| 35 | No automatic attachment on switch | ✅ | handoff design + tests |
| 36 | Explicit "Ask Agent"/"Discuss chapter" | ✅ | quote selection flow + tests |
| 37 | Show source/purpose/version/size + remove | ✅ | InputArea + backend handoff tests |
| 38 | Configurable high context cap | ✅ | agent_context_handoff_limit_kb + API tests |
| 39 | Runs continue after leaving Agent | ✅ | handler_chat_task_lifecycle_test.go |
| 40 | Recover persisted run state | ✅ | task lifecycle + stream replay tests |
| 41 | Only explicit stop ends a run | ✅ | lifecycle tests (no timeout/route cancel) |
| 42 | File changes show summaries first | ✅ | AgentChangeSummaryCard tests |
| 43 | Full-screen single-column diff review | ✅ | ChangeReviewWorkspace tests |
| 44 | Bounded sourced review feedback | ✅ | ReviewFeedbackTray + context tests |
| 45 | Accept/undo/continue from review | ✅ | review mutation tests |
| 46 | Return to Agent after review | ✅ | use-writing-change-review.test.tsx close case |
| 47 | Runs bound to source project/session | ✅ | API lifecycle tests |
| 48 | Project switch never kills runs | ✅ | version lease + task registry tests |
| 49 | Open cross-project run restores source | ✅ | task recovery navigation tests |
| 50 | Unified task center for long work | ✅ | task center API + integration tests |
| 51 | Task shows project/type/status/time/recovery | ✅ | task center API assertions |
| 52 | Badge counts waiting/failed only | ✅ | WorkbenchShell badge tests |
| 53 | Completion never steals focus | ✅ | task center design + polling tests |
| 54 | In-app task center always available | ✅ | More menu + integration tests |
| 55 | Notifications off by default, no first-run request | ✅ | settings default + module tests + live-preference workbench integration test |
| 56 | Notifications show type/project only | ✅ | task-notifications tests + locale/privacy workbench integration test |

## Interactive Story (57-81)

| # | Story (short) | Status | Evidence |
| --- | --- | --- | --- |
| 57 | Story opens stage first | ✅ | StoryStage default tests |
| 58 | Full-width narrative reader | ✅ | StoryStage tests |
| 59 | Speaker + avatar for dialogue | ✅ | MessageItem dialogue tests |
| 60 | Compact player-action blocks | ✅ | StoryStage action block tests |
| 61 | Rolls/state changes collapsible cards | ✅ | roll/state card tests |
| 62 | Scene images full screen when present | ✅ | MessageItem ImagePreviewDialog tests |
| 63 | Full text experience without images | ✅ | no-image rendering tests |
| 64 | Max three suggestions + More | ✅ | StoryStage hot choices tests |
| 65 | Suggestion fills input only | ✅ | StoryStage hot choices tests |
| 66 | Keep draft, block duplicate submit | ✅ | streaming input tests |
| 67 | Enter/return director workspace | ✅ | DirectorBackstage tests |
| 68 | Turn locks submission-time snapshot | ✅ | interactive snapshot tests |
| 69 | New director rules apply next turn | ✅ | director plan tests |
| 70 | Stop/regenerate to apply new config | ✅ | director run controls |
| 71 | Latest-turn revision records | ✅ | turn narrative revision tests |
| 72 | Historical edits create branches | ✅ | branch creation tests |
| 73 | Original branch preserved | ✅ | branch history tests |
| 74 | Storylines list-first on phone | ✅ | StorylinesView rendered for mobile timeline in InteractiveLayout.test.tsx; list-first tests in StorylinesView.test.tsx |
| 75 | Branch item shows divergence/summary/count/time | ✅ | StorylinesView.test.tsx asserts current-first order, divergence, parent, latest summary, turn count, updated time |
| 76 | Branch detail timeline operations | ✅ | StorylinesView.test.tsx asserts timeline, continue, switch, delete, rename; handler_interactive_branch_rename_test.go + InteractiveLayout.test.tsx cover API and wiring |
| 77 | Relation graph secondary on wide | ✅ | mobile graph-overview toggle in StorylinesView.test.tsx; desktop timeline keeps BranchTimeline in InteractiveLayout.test.tsx |
| 78 | Reference separates state and lore | ✅ | story-state model tests |
| 79 | Current state read-only by default | ✅ | state display tests |
| 80 | State edit requires full-screen revision | ✅ | CurrentStateRevisionWorkspace tests |
| 81 | Revision scoped to branch, undo/restore | ✅ | state revision API tests |

## Management, Offline and Accessibility (82-104)

| # | Story (short) | Status | Evidence |
| --- | --- | --- | --- |
| 82 | List-to-full-detail management pages | ✅ | AdaptiveSurface drawer flows asserted (Agents narrow + Skills mobile tests); kept-mounted routes |
| 83 | Restore search/filter/scroll/selection | ✅ | MobilePaneHost keeps drawer panes mounted; adaptive-surface + SkillsView tests assert search survives close/reopen |
| 84 | Executable config staged + validated | ✅ | settings/config-manager tests |
| 85 | Unsaved config return protection | 🟡 | design in docs/verification/executable-config-return-guard-design.md; awaiting confirmation before implementation |
| 86 | Continuous draft autosave + conflict display | ✅ | useEditorDraftPersistence tests |
| 87 | Immediate preferences autosave | ✅ | settings auto-save tests |
| 88 | Short offline editing of open docs | ✅ | pending-drafts tests |
| 89 | Offline edits bound to project + base version | ✅ | pending-drafts + editor tests |
| 90 | Remote changes enter conflict review | ✅ | offline conflict-path test |
| 91 | No submit/switch/migration while offline | 🟡 | offline switch, Agent composer submit, and structure migration gates verified (App integration + InputArea tests); story-stage submit remains natural-failure only |
| 92 | Top/Android/browser share one back stack | ✅ | popstate + Escape tests |
| 93 | Back priority: keyboard/surface/unsaved/full-screen | ✅ | Escape + dialog-first tests |
| 94 | No left/right edge gestures | ✅ | no edge-swipe code in mobile shell |
| 95 | Tablet two-pane when space allows | ✅ | AgentsView.test.tsx asserts list+detail visible on wide viewport and drawer flow on narrow screens; AdaptiveSurface tests cover collapse/expand |
| 96 | Desktop keeps multi-pane + tabs | ✅ | desktop workbench retained; responsive tests |
| 97 | Rotation keeps editor/session mounted | ✅ | adaptive-surface + stable host tests |
| 98 | Coarse pointer keeps touch targets | ✅ | CDP measurement covers shell, task center, onboarding, toast, Settings controls, and Radix options at 320/390/430px |
| 99 | Complete Simplified Chinese copy | ✅ | i18n alignment + zh browser run |
| 100 | English long-text no truncation/overlap | ✅ | long-English book name browser case (320/390/430) |
| 101 | Screen-reader semantics | 🟡 | aria labels present; SR audit not run |
| 102 | Keyboard core flows | ✅ | Escape/Tab focus tests + keyboard nav tests |
| 103 | 200% font zoom accessibility | ✅ | zoom acceptance (CSS zoom 2x, no overflow) |
| 104 | Reduce-motion respected | ✅ | motion intensity setting + provider |

## Remaining gaps

- Device-only verification: safe-area insets and the real system-notification permission flow on hardware.
- Human/vision review remains incomplete across all surfaces; Settings was visually inspected at 320px in zh-CN/dark and en-US/light.
- Explicit tests for a few 🟡 rows (long-title truncation, project-switch save failure, mobile two-pane at wide tablet, screen-reader pass).
