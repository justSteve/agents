# {{PROJECT_NAME}}

{{PROJECT_DESCRIPTION}}

## Quick Start

```bash
# Clone the repository
git clone {{REPO_URL}}
cd {{PROJECT_SLUG}}

# Start working with beads
bd ready          # See available work
bd show <id>      # View issue details
```

## Project Structure

```
{{PROJECT_SLUG}}/
├── .beads/              # Issue tracking database
├── .claude/
│   └── CLAUDE.md        # Agent instructions
├── src/                 # Source code
├── tests/               # Test files
├── docs/                # Documentation
└── README.md
```

## Development Workflow

This project uses **Beads** for AI-native issue tracking. All work is tracked as beads issues.

### Starting a Session

```bash
bd ready                              # Find available work
bd show <id>                          # Review issue details
bd update <id> --status=in_progress   # Claim the work
```

### Completing Work

```bash
bd close <id>                         # Mark issue complete
bd sync                               # Sync with remote
git add . && git commit -m "..."      # Commit code changes
git push                              # Push to remote
```

### Creating Issues

```bash
bd create --title="..." --type=task   # Create a task
bd create --title="..." --type=bug    # Report a bug
bd create --title="..." --type=feature --priority=1  # High-priority feature
```

## Configuration

### Agent Plugins

This project is configured to use the following agent plugins:

{{#AGENT_PLUGINS}}
- `{{PLUGIN_NAME}}` - {{PLUGIN_DESCRIPTION}}
{{/AGENT_PLUGINS}}

## Contributing

1. Check `bd ready` for available work
2. Claim an issue with `bd update <id> --status=in_progress`
3. Make your changes
4. Close the issue with `bd close <id>`
5. Create a pull request

## License

{{LICENSE}}
