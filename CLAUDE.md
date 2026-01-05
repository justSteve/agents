# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Claude Code Plugins marketplace: A production-ready system combining **99 specialized AI agents**, **15 multi-agent workflow orchestrators**, **107 agent skills**, and **71 development tools** organized into **68 focused, single-purpose plugins**.

**Repository**: `justSteve/agents` (forked from `wshobson/agents`)

## Architecture

### Plugin Structure

Each plugin is isolated in `plugins/{plugin-name}/` with:
- `agents/` - Specialized agents (markdown with YAML frontmatter)
- `commands/` - Tools and workflows
- `skills/{skill-name}/SKILL.md` - Modular knowledge packages with progressive disclosure
- `plugin.json` - Plugin metadata

### Key Files

- `.claude-plugin/marketplace.json` - Central plugin registry (68 plugins)
- `.beads/` - AI-native issue tracking database
- `.claude/settings.local.json` - Claude Code permissions and allowlisted commands
- `.github/workflows/claude.yml` - GitHub Actions for Claude Code automation

### Agent Frontmatter

```yaml
---
name: agent-name
description: What the agent does
model: opus|sonnet|haiku|inherit
---
[System prompt]
```

### Skill Frontmatter

```yaml
---
name: skill-name
description: What the skill does. Use when [trigger].
---
[Progressive disclosure content]
```

## Issue Tracking with Beads

This repo uses **Beads** for AI-native issue tracking. Issues are stored in `.beads/issues.jsonl` and sync with git.

### Essential Commands

```bash
bd ready                              # Show unblocked, open issues
bd list --status=open                 # All open issues
bd create --title="..." --type=task   # Create issue
bd show <id>                          # View issue details
bd update <id> --status=in_progress   # Claim work
bd close <id>                         # Mark complete
bd dep add <issue> <depends-on>       # Add dependency
bd sync                               # Sync with git remote
bd stats                              # Project statistics
```

### Session Protocol

Before completing work, run:
```bash
git status && git add <files> && bd sync && git commit -m "..." && bd sync && git push
```

## Development Commands

### Plugin Operations

```bash
/plugin marketplace add justSteve/agents    # Add marketplace
/plugin install python-development          # Install specific plugin
/plugin                                      # List available commands
```

### Testing Changes

No build or test commands - this is a configuration-only repository (YAML, Markdown, JSON).

### Validation

- Ensure YAML frontmatter is valid in all agent/skill files
- Verify plugin.json references existing files
- Update marketplace.json when adding/modifying plugins

## Contributing Patterns

### Adding an Agent

1. Create `plugins/{plugin}/agents/{name}.md` with frontmatter
2. Write comprehensive system prompt
3. Add to plugin's `agents` array in plugin.json
4. Update marketplace.json

### Adding a Skill

1. Create `plugins/{plugin}/skills/{name}/SKILL.md`
2. Add YAML frontmatter with "Use when" trigger in description
3. Add to plugin's `skills` array in plugin.json

### Model Selection

| Model | Use Case |
|-------|----------|
| opus | Critical architecture, security, code review |
| sonnet | Complex tasks, support with intelligence |
| haiku | Fast operations, deployment, SEO |
| inherit | User-controlled via session default |

## Beads Workflows Plugin

The `beads-workflows` plugin (`plugins/beads-workflows/`) contains specialized agents and skills for workflow discipline:

- **Agents**: `beads-disciplinarian`, `beads-workflow-orchestrator`, `beads-issue-reviewer`
- **Commands**: session-start, session-end, issue-create, dependency-add
- **Skills**: session-rituals, dependency-thinking, description-quality
- **Go Library**: `lib/agent_tracking/` for database extensions
