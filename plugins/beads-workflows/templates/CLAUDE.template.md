# CLAUDE.md

This file provides guidance to Claude Code when working with this project.

## Project Overview

**Name**: {{project_name}}
**Description**: {{description}}
**Created**: {{created_date}}

## Beads Issue Tracking

This project uses [beads](https://github.com/steveyegge/beads) for AI-native issue tracking.

### Quick Reference

```bash
# Find work
bd ready                              # Show unblocked issues
bd list --status=open                 # All open issues
bd show <id>                          # View issue details

# Work on issues
bd update <id> --status=in_progress   # Claim an issue
bd close <id>                         # Complete an issue

# Create issues
bd create --title="..." --type=task   # Create new issue
bd dep add <issue> <depends-on>       # Add dependency

# Sync
bd sync                               # Sync with git remote
```

### Session Workflow

**Start of session:**
1. Run `bd ready` to find available work
2. Review issue with `bd show <id>`
3. Claim with `bd update <id> --status=in_progress`

**End of session:**
1. Run `git status` to check changes
2. Commit your work
3. Run `bd sync` to sync beads
4. Run `git push` to push changes
5. Close completed issues with `bd close <id>`

### Issue Creation Guidelines

When creating issues, include:
- **Why**: Context and motivation
- **What**: Specific problem or feature
- **How**: Acceptance criteria (checkboxes)

Example:
```bash
bd create --title="Add user authentication" --type=feature \
  --description="## Context
We need to secure the API endpoints.

## Requirements
- [ ] Implement JWT tokens
- [ ] Add login endpoint
- [ ] Add logout endpoint

## Acceptance Criteria
- Users can log in and receive a token
- Protected endpoints reject invalid tokens"
```

## Agent Plugins

This project uses the following agent plugins:
{{#agents}}
- **{{name}}**: {{description}}
{{/agents}}

## Project-Specific Rules

{{#rules}}
- {{.}}
{{/rules}}
{{^rules}}
<!-- Add project-specific rules here -->
{{/rules}}

## Technology Stack

{{#tech_stack}}
- {{.}}
{{/tech_stack}}
{{^tech_stack}}
<!-- Document your technology stack here -->
{{/tech_stack}}

## Development Commands

{{#commands}}
### {{name}}
```bash
{{command}}
```
{{description}}

{{/commands}}
{{^commands}}
<!-- Document common development commands here -->
{{/commands}}

## Architecture Notes

{{#architecture}}
{{.}}
{{/architecture}}
{{^architecture}}
<!-- Add architecture notes as the project evolves -->
{{/architecture}}
