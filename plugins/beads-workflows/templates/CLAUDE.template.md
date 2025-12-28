# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**{{PROJECT_NAME}}** - {{PROJECT_DESCRIPTION}}

## Issue Tracking with Beads

This project uses **Beads** for AI-native issue tracking. Issues are stored in `.beads/` and sync with git.

### Session Rituals

**Starting a session:**
```bash
bd ready                              # Find available work
bd show <id>                          # Review issue details
bd update <id> --status=in_progress   # Claim the work
```

**Completing work:**
```bash
bd close <id>                         # Mark issue complete
bd sync                               # Sync beads with remote
git add . && git commit -m "..."      # Commit code changes
git push                              # Push to remote
```

### Creating Issues

```bash
bd create --title="..." --type=task        # Create a task
bd create --title="..." --type=bug         # Report a bug
bd create --title="..." --type=feature     # New feature
bd create --title="..." --type=epic        # Large initiative
```

### Dependencies

```bash
bd dep add <issue> <depends-on>       # Add dependency
bd blocked                            # Show blocked issues
bd dep tree <id>                      # View dependency tree
```

### Project Health

```bash
bd stats                              # Project statistics
bd ready                              # Unblocked work
bd list --status=open                 # All open issues
```

## Agent Plugins

This project uses the following Claude Code plugins:

{{#AGENT_PLUGINS}}
- **{{PLUGIN_NAME}}**: {{PLUGIN_DESCRIPTION}}
{{/AGENT_PLUGINS}}

## Development Commands

{{#DEV_COMMANDS}}
```bash
{{COMMAND}}  # {{DESCRIPTION}}
```
{{/DEV_COMMANDS}}

## Project Structure

```
{{PROJECT_SLUG}}/
{{PROJECT_STRUCTURE}}
```

## Code Conventions

{{#CODE_CONVENTIONS}}
- {{CONVENTION}}
{{/CODE_CONVENTIONS}}

## Important Notes

- Always check `bd ready` before starting work
- Close issues immediately after completing work (don't batch)
- Run `bd sync` at session end to push beads changes
- All work should be tracked as beads issues
