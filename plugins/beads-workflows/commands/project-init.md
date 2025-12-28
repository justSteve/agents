# Project Init

Initialize a new beads-native project with proper structure, templates, and initial issues.

## Overview

This command scaffolds a complete project setup:
1. Directory setup (create/verify, git init)
2. Beads initialization (bd init, hooks, config)
3. Template application (files, CLAUDE.md)
4. Initial issues (setup epic with tasks)
5. First session (optional beads-session-start)

## Arguments

**Required**:
- `<project-name>` - Name for the project (used for directory and module)

**Optional**:
- `--prefix <prefix>` - Issue ID prefix (default: derived from project name)
- `--template <type>` - Project template: standard, python, typescript, go (default: standard)
- `--agents <list>` - Comma-separated agent plugins to enable (default: beads-workflows)
- `--path <dir>` - Parent directory for project (default: current directory)
- `--description <text>` - Project description for README and CLAUDE.md
- `--no-session` - Skip running beads-session-start after setup

## Phase 1: Directory Setup

Create or verify project directory with git initialization.

**Steps**:
```bash
# Create project directory
mkdir -p <path>/<project-name>
cd <path>/<project-name>

# Initialize git if not already a repo
git init
```

**Validation**:
- Directory created successfully
- Git initialized (or already a git repo)
- No conflicting .beads directory

**Output**:
```
✓ Directory created: <path>/<project-name>
✓ Git initialized
```

## Phase 2: Beads Initialization

Initialize beads issue tracking with hooks and agent configuration.

**Commands**:
```bash
# Initialize beads with prefix
bd init --prefix <prefix>

# Install git hooks
bd hooks install

# Configure agent marketplace (if available)
bd config set agents.marketplace.repo "<agents-repo-path>"
bd config set agents.marketplace.enabled "<agents-list>"
bd config set agents.default_model "sonnet"
```

**Validation**:
- .beads directory created
- beads.db exists
- Hooks installed
- Configuration set

**Output**:
```
✓ Beads initialized with prefix: <prefix>
✓ Git hooks installed
✓ Agent plugins enabled: <list>
```

## Phase 3: Template Application

Apply project template and generate CLAUDE.md.

**Template Selection**:
| Template | Files Created |
|----------|---------------|
| standard | .gitignore, README.md |
| python | + pyproject.toml |
| typescript | + package.json, tsconfig.json |
| go | + go.mod, Makefile |

**Steps**:
1. Copy template files from `templates/<type>/`
2. Substitute variables in templates:
   - `{{project_name}}` - Project name
   - `{{description}}` - Project description
   - `{{module_path}}` - Go module path (for go template)
   - `{{created_date}}` - Current date
3. Generate `.claude/CLAUDE.md` from template
4. Commit initial structure

**Variable Substitution**:
```
{{project_name}} → my-awesome-app
{{description}} → A new beads-native project
{{created_date}} → 2025-01-15
{{repo_url}} → https://github.com/<user>/<project>
{{license}} → MIT
```

**Commands**:
```bash
# Create .claude directory
mkdir -p .claude

# Copy and process template files
# (Agent performs file operations)

# Initial commit
git add .
git commit -m "chore: Initialize project with beads"
```

**Output**:
```
✓ Template applied: <type>
✓ CLAUDE.md generated
✓ Initial commit created
```

## Phase 4: Initial Issues

Create setup epic and child tasks from template.

**Source**: `templates/<type>/initial-issues.json`

**Steps**:
1. Parse initial-issues.json
2. Create epic issue
3. Create child task issues
4. Add dependencies (tasks depend on epic)

**Commands**:
```bash
# Create epic
bd create --title="<epic.title>" --description="<epic.description>" \
  --type=epic --priority=<epic.priority>

# Create tasks (for each task in tasks[])
bd create --title="<task.title>" --description="<task.description>" \
  --type=<task.type> --priority=<task.priority>

# Add dependencies
bd dep add <task-id> <epic-id>
```

**Output**:
```
✓ Created epic: <prefix>-xxx - <epic.title>
✓ Created task: <prefix>-yyy - <task.title>
✓ Created task: <prefix>-zzz - <task.title>
...
✓ Dependencies configured
```

## Phase 5: First Session (Optional)

Start the first work session unless --no-session specified.

**Command**:
```bash
# Run session start
/beads-session-start
```

**Output**:
```
✓ Project ready!

📋 Ready work:
1. [P1] <prefix>-xxx: <first-task>
...

Run 'bd update <id> --status in_progress' to claim an issue.
```

## Complete Example

```bash
# Create a Python project
/project-init my-api --template python --prefix api \
  --description "REST API for my application" \
  --agents beads-workflows

# Output:
✓ Directory created: ./my-api
✓ Git initialized
✓ Beads initialized with prefix: api
✓ Git hooks installed
✓ Agent plugins enabled: beads-workflows
✓ Template applied: python
✓ CLAUDE.md generated
✓ Initial commit created
✓ Created epic: api-xxx - Python Project Setup
✓ Created task: api-yyy - Configure Python virtual environment
✓ Created task: api-zzz - Set up pytest testing
✓ Created task: api-aaa - Configure linting and formatting
✓ Created task: api-bbb - Add type checking
✓ Dependencies configured
✓ Project ready!

📋 Ready work (4 issues):
1. [P1] api-yyy: Configure Python virtual environment
2. [P1] api-zzz: Set up pytest testing
3. [P2] api-aaa: Configure linting and formatting
4. [P2] api-bbb: Add type checking
```

## Error Handling

| Error | Resolution |
|-------|------------|
| Directory exists with .beads | Ask user to use different name or --force |
| Git init fails | Check permissions, directory state |
| Template not found | Default to 'standard' with warning |
| bd init fails | Check beads installation |

## Files Created

After running `/project-init my-app --template python`:

```
my-app/
├── .beads/
│   ├── beads.db
│   ├── config.yaml
│   └── issues.jsonl
├── .claude/
│   └── CLAUDE.md
├── .git/
│   └── hooks/
│       └── post-commit (beads hook)
├── .gitignore
├── pyproject.toml
└── README.md
```

## Integration with beads-workflows

This command is part of the beads-workflows plugin and follows these patterns:

- **Session Rituals**: Sets up hooks for automatic sync
- **Description Quality**: Initial issues have proper Why/What/How structure
- **Dependency Thinking**: Tasks correctly depend on epic
- **Single-Issue Discipline**: Ready issues shown for focused work
