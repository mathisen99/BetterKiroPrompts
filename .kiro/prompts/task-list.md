List current work items from Kiro spec tasks.

INPUTS:
- Read every file matching: .kiro/specs/*/tasks.md
- If none exist, stop and say no specs/tasks were found.

OUTPUT REQUIREMENTS:
- Show tasks grouped by spec slug.
- Prefer tasks that are NOT completed.
- If the file uses checkboxes:
  - treat "- [ ]" as open and "- [x]" as done
- If the file does not use checkboxes:
  - treat lines containing "(done)" or "[x]" as done, otherwise open
- Show at most the top 10 open tasks across all specs, but include counts:
  - open count per spec
  - done count per spec

ALSO INCLUDE:
- "Suggested next task" chosen by:
  1) earliest task in the list that is open
  2) if dependencies are explicitly stated, pick the first unblocked task
- A one-paragraph reminder of the workflow:
  1) implement work
  2) /task-complete to mark done
  3) /commit to commit changes

Do not modify any files.
