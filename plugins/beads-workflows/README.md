# Beads Workflows Plugin

Integration of [beads](https://github.com/steveyegge/beads) issue tracking with the agents marketplace, providing session management, pattern enforcement, and workflow orchestration for beads-native development.

## Overview

This plugin makes agents committed beads users by teaching beads patterns through skills, enforcing discipline through agents, and orchestrating workflows through commands.

## Components

### Agents

- **beads-disciplinarian** - Enforces beads workflow discipline (session rituals, dependency direction, description quality, single-issue focus)
- **beads-workflow-orchestrator** - Manages session lifecycle and work selection
- **beads-issue-reviewer** - Validates issue quality before creation

### Commands

- **beads-session-start** - Initialize a beads work session (sync, ready, claim, context load)
- **beads-session-end** - Clean up and sync session (bd sync, git push, file discovered work)
- **beads-issue-create** - Guided issue creation with quality validation
- **beads-dependency-add** - Guided dependency creation with causal reasoning validation

### Skills

- **session-rituals** - Session start/end patterns and checklists
- **dependency-thinking** - Causal vs temporal reasoning (avoid "Phase 1 before Phase 2" trap)
- **description-quality** - Context-rich descriptions (Why/What/How structure)
- **single-issue-discipline** - One issue at a time workflow
- **beads-cli-reference** - CLI command reference and patterns

## Installation

```bash
/plugin marketplace add justSteve/agents
/plugin install beads-workflows
```

## Prerequisites

- [beads](https://github.com/steveyegge/beads) CLI installed (`bd` command available)
- Git repository with beads initialized (`bd init`)

## Usage

### Starting a Session

```bash
/beads-session-start
```

This will:
1. Verify environment (`pwd && git status`)
2. Sync database (`bd sync`)
3. Find ready work (`bd ready --json`)
4. Help you select ONE issue
5. Claim the issue (`bd update <id> --status in_progress`)
6. Load relevant context and skills

### Ending a Session

```bash
/beads-session-end
```

This will:
1. Ensure all work is filed (`bd create` for discovered issues)
2. Sync database (`bd sync`)
3. Push to remote (`git push`)
4. Verify session cleanup

### Creating Issues

```bash
/beads-issue-create
```

Guided issue creation with description quality validation ensuring Why/What/How structure.

### Adding Dependencies

```bash
/beads-dependency-add
```

Guided dependency addition with causal reasoning validation to avoid temporal thinking traps.

## Pattern Enforcement

This plugin enforces beads discipline through multiple layers:

1. **Agent System Prompts** - All agents include beads workflow requirements
2. **Command Orchestration** - Commands wrap beads operations with validation checkpoints
3. **Skill Knowledge** - Skills teach patterns through concrete examples
4. **Soft Validation** - Helpful warnings without blocking progress

## Key Patterns

### Session Rituals

**Start**:
- Run `bd ready` to find available work
- Choose ONE issue to work on
- Run `bd update <id> --status in_progress`

**End**:
- Run `bd sync` to flush changes
- Run `git push` to sync remote
- File discovered work with `--deps discovered-from:<parent-id>`

### Dependency Thinking

- Think CAUSALLY: "Y needs X" → `bd dep add Y X`
- NOT temporally: "X before Y"
- Verify with `bd blocked` - tasks should be blocked by prerequisites

### Description Quality

Every issue MUST answer:
- **Why**: Problem statement or need
- **What**: Planned scope and approach
- **How discovered**: Context if filed during work

### Single-Issue Discipline

- Work on ONE issue at a time
- File discovered work separately
- Use `discovered-from` to link related work

## Configuration

Project-level beads configuration (stored in `.beads/beads.db`):

```bash
bd config set agents.marketplace.repo "/path/to/agents"
bd config set agents.marketplace.enabled "beads-workflows,developer-essentials"
bd config set agents.session.auto_prime true
bd config set agents.enforcement.description_min_length 50
```

## Integration with Other Plugins

Recommended plugin combinations:

- **developer-essentials** - Git, testing, debugging workflows
- **git-pr-workflows** - PR creation and review with beads tracking
- **python-development** - Python development with beads task management

## Development Status

**Version**: 1.0.0 (In Development)

This plugin is currently being built as part of the beads-agents integration project. See the [integration plan](https://github.com/justSteve/agents/blob/main/.claude/plans/beads-agents-integration-plan.md) for details.

## License

MIT
