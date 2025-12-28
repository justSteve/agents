# {{project_name}}

{{description}}

## Getting Started

```bash
# Clone the repository
git clone {{repo_url}}
cd {{project_name}}

# View available work
bd ready
```

## Development

This project uses [beads](https://github.com/steveyegge/beads) for issue tracking.

### Common Commands

```bash
bd ready                    # Show issues ready to work on
bd show <id>                # View issue details
bd update <id> --status in_progress  # Claim an issue
bd close <id>               # Complete an issue
bd sync                     # Sync with remote
```

### Session Workflow

1. **Start**: `bd ready` → pick an issue → `bd update <id> --status in_progress`
2. **Work**: Make changes, commit regularly
3. **End**: `bd close <id>` → `bd sync` → `git push`

## Project Structure

```
{{project_name}}/
├── .beads/           # Issue tracking database
├── .claude/          # Claude Code configuration
│   └── CLAUDE.md     # Project-specific instructions
└── README.md         # This file
```

## Contributing

1. Check `bd ready` for available issues
2. Claim one issue at a time
3. Follow the session workflow above
4. Create new issues for discovered work: `bd create --title="..." --type=task`

## License

{{license}}
