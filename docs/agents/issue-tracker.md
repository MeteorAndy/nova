# Issue tracker: GitHub

Issues and PRDs for this repository live in GitHub Issues. Use the `gh` CLI
for all operations and infer the repository from `git remote -v`.

## Conventions

- Create: `gh issue create --title "..." --body "..."`
- Read: `gh issue view <number> --comments`
- List: `gh issue list` with suitable state, label, and JSON filters
- Comment: `gh issue comment <number> --body "..."`
- Label: `gh issue edit <number> --add-label "..."` or `--remove-label "..."`
- Close: `gh issue close <number> --comment "..."`

## Pull requests as a triage surface

**PRs as a request surface: no.**

GitHub issues and pull requests share one number space. Resolve ambiguous
references with `gh pr view <number>`, falling back to `gh issue view <number>`.

## Skill conventions

"Publish to the issue tracker" means creating a GitHub issue.
"Fetch the relevant ticket" means reading the issue and its comments.

For `/wayfinder`, use a `wayfinder:map` issue and child issues. Prefer native
GitHub sub-issues and dependencies; fall back to task lists and explicit
`Blocked by: #<number>` references when those features are unavailable.
