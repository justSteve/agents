# Project Init

Initialize a new beads-enabled project with templates, issues, and agent configuration.

## Overview

This command orchestrates complete new project scaffolding:
1. Directory setup (create/verify project directory, init git)
2. Beads initialization (bd init, install hooks)
3. Template application (copy template files)
4. Documentation generation (create CLAUDE.md)
5. Initial issues creation (from template)
6. First session start

**Goal**: Initialize new projects with beads + agents integration in <5 minutes.

## Arguments

**Required**:
- `<project-name>` - Human-readable project name (e.g., "My Awesome App")

**Optional**:
- `--template <type>` - Template type: `standard` (default), `python`, `typescript`, `go`
- `--prefix <prefix>` - Beads issue prefix (default: derived from project name)
- `--path <path>` - Project directory path (default: current directory / project-slug)
- `--agents <plugins>` - Comma-separated agent plugins to configure
- `--description <desc>` - Project description for README and CLAUDE.md
- `--skip-issues` - Skip initial issue creation
- `--skip-session` - Skip starting first session

## Phase 1: Validate Arguments

Parse and validate all provided arguments.

**Derive defaults**:
- `project-slug`: Convert project name to lowercase, replace spaces with hyphens
- `prefix`: First 3-4 characters of project-slug
- `path`: If not provided, use `./<project-slug>`

**Validation**:
- Project name is not empty
- Template type is valid (standard, python, typescript, go)
- Path doesn't exist OR is empty directory OR user confirms overwrite

**Output**:
```
Project Configuration:
  Name: <project-name>
  Slug: <project-slug>
  Template: <template-type>
  Prefix: <prefix>
  Path: <path>
  Plugins: <plugin-list>

Proceed? [Y/n]
```

## Phase 2: Directory Setup

Create or verify project directory and initialize git.

**Commands to run**:
```bash
mkdir -p <path>
cd <path>
git init
```

**If directory exists and is not empty**:
```
⚠️  Directory <path> is not empty

Contents:
<ls output>

Options:
[1] Continue anyway (may overwrite files)
[2] Choose different path
[3] Abort

Your choice: _
```

**Validation**:
- Directory exists and is accessible
- Git init succeeds
- Write permissions verified

**Output**:
```
✓ Directory setup complete
  Path: <full-path>
  Git: Initialized
```

## Phase 3: Beads Initialization

Initialize beads and install git hooks.

**Commands to run**:
```bash
bd init --prefix <prefix>
bd hooks install
```

**Validation**:
- bd command available
- .beads/ directory created
- hooks installed successfully

**Output**:
```
✓ Beads initialized
  Prefix: <prefix>
  Database: .beads/beads.db
  Hooks: Installed
```

## Phase 4: Template Application

Copy template files based on selected template type.

**Template locations**:
- Standard: `plugins/beads-workflows/templates/standard/`
- Python: `plugins/beads-workflows/templates/python/`
- TypeScript: `plugins/beads-workflows/templates/typescript/`
- Go: `plugins/beads-workflows/templates/go/`

**Files to copy** (varies by template):

### Standard Template
- `.gitignore` → `.gitignore`
- `README.template.md` → `README.md` (with variable substitution)

### Python Template
- `.gitignore` → `.gitignore`
- `pyproject.toml.template` → `pyproject.toml` (with variable substitution)
- Standard `README.template.md` → `README.md`

### TypeScript Template
- `.gitignore` → `.gitignore`
- `package.json.template` → `package.json` (with variable substitution)
- `tsconfig.json.template` → `tsconfig.json`
- Standard `README.template.md` → `README.md`

### Go Template
- `.gitignore` → `.gitignore`
- `go.mod.template` → `go.mod` (with variable substitution)
- `Makefile.template` → `Makefile` (with variable substitution)
- Standard `README.template.md` → `README.md`

**Variable substitution**:
Replace template variables in copied files:
- `{{PROJECT_NAME}}` → project name
- `{{PROJECT_SLUG}}` → project slug
- `{{PROJECT_DESCRIPTION}}` → description or default
- `{{REPO_URL}}` → GitHub URL if derivable, or placeholder
- `{{AUTHOR_NAME}}` → git config user.name or placeholder
- `{{AUTHOR_EMAIL}}` → git config user.email or placeholder
- `{{MODULE_PATH}}` → Go module path (for Go template)
- `{{GO_VERSION}}` → Current Go version (for Go template)

**Output**:
```
✓ Template applied: <template-type>
  Files created:
  - .gitignore
  - README.md
  - <template-specific files>
```

## Phase 5: CLAUDE.md Generation

Generate .claude/CLAUDE.md from template.

**Commands to run**:
```bash
mkdir -p .claude
```

**Template source**: `plugins/beads-workflows/templates/CLAUDE.template.md`

**Variable substitution**:
- `{{PROJECT_NAME}}` → project name
- `{{PROJECT_SLUG}}` → project slug
- `{{PROJECT_DESCRIPTION}}` → description or "TODO: Add project description"
- `{{#AGENT_PLUGINS}}...{{/AGENT_PLUGINS}}` → Expand for each configured plugin
- `{{#DEV_COMMANDS}}...{{/DEV_COMMANDS}}` → Template-specific dev commands
- `{{PROJECT_STRUCTURE}}` → Template-specific directory structure

**Template-specific DEV_COMMANDS**:

Python:
```
uv sync --extra dev    # Install dependencies
pytest                 # Run tests
ruff check .           # Lint code
mypy src/              # Type check
```

TypeScript:
```
npm install            # Install dependencies
npm test               # Run tests
npm run build          # Build project
npm run lint           # Lint code
```

Go:
```
go mod tidy            # Sync dependencies
go test ./...          # Run tests
go build ./...         # Build project
golangci-lint run      # Lint code
```

**Output**:
```
✓ CLAUDE.md generated
  Path: .claude/CLAUDE.md
  Plugins: <plugin-count> configured
```

## Phase 6: Initial Issues Creation

Create initial issues from template's initial-issues.json.

**Skip if**: `--skip-issues` flag provided

**Template source**: `plugins/beads-workflows/templates/<template>/initial-issues.json`

**Process**:
1. Read initial-issues.json
2. Substitute variables in titles and descriptions
3. Create issues in dependency order (epics first, then dependents)
4. Add dependencies between issues

**Commands to run** (for each issue):
```bash
bd create --title="<title>" --type=<type> --priority=<priority> --description="<description>"
```

**For dependencies**:
```bash
bd dep add <issue-id> <depends-on-id>
```

**ID mapping**:
- Track mapping from template `id` to actual beads ID
- Use mapping when adding dependencies

**Output**:
```
✓ Initial issues created
  Epic: <epic-id> - <epic-title>
  Tasks: <count> created
  Dependencies: <count> configured
```

## Phase 7: First Session (Optional)

Start first beads session if not skipped.

**Skip if**: `--skip-session` flag provided

**Use beads-session-start command**:

The session start will:
- Verify environment
- Sync database
- Show ready work
- Help select first issue
- Load context

**Output**:
```
✓ Ready to start first session

Run: /beads-session-start
Or: bd ready (to see available work)
```

## Phase 8: Summary

Provide complete project initialization summary.

**Output**:
```
✅ Project initialized successfully

Project: <project-name>
Location: <full-path>
Template: <template-type>

Structure:
├── .beads/              # Issue tracking
├── .claude/
│   └── CLAUDE.md        # Agent configuration
├── .git/                # Version control
├── .gitignore
├── README.md
└── <template-specific files>

Beads:
• Prefix: <prefix>
• Issues: <count> created
• Hooks: Installed

Next steps:
1. cd <path>
2. Review .claude/CLAUDE.md
3. Run: /beads-session-start (or bd ready)
4. Start building!

Documentation:
• Beads: https://github.com/steveyegge/beads
• Claude Code: https://claude.ai/code
```

## Error Handling

### Error: bd not installed

**Detection**: `bd` command not found

**Output**:
```
❌ Beads CLI not installed

This command requires beads (https://github.com/steveyegge/beads)

Installation:
1. Clone: git clone https://github.com/steveyegge/beads
2. Build: cd beads && go build -o bd ./cmd/bd
3. Add to PATH

Project init aborted.
```

### Error: Template not found

**Detection**: Template directory doesn't exist

**Output**:
```
❌ Template not found: <template-type>

Available templates:
• standard - Generic project template
• python - Python project with uv/pytest
• typescript - TypeScript project with npm/vitest
• go - Go project with modules

Use: /project-init "My Project" --template <type>
```

### Error: Directory not writable

**Detection**: Cannot create files in target directory

**Output**:
```
❌ Cannot write to directory: <path>

Please check:
• Directory permissions
• Disk space
• Path validity

Project init aborted.
```

### Error: Git not available

**Detection**: git command not found

**Output**:
```
❌ Git not installed

This command requires git for version control.

Installation:
• macOS: xcode-select --install
• Linux: apt install git
• Windows: https://git-scm.com/download/win

Project init aborted.
```

## Examples

### Basic usage
```
/project-init "My Awesome App"
```
Creates `./my-awesome-app/` with standard template.

### Python project
```
/project-init "Data Pipeline" --template python --prefix dp
```
Creates Python project with pyproject.toml, pytest setup.

### TypeScript project with custom path
```
/project-init "API Server" --template typescript --path ./backend
```
Creates TypeScript project in `./backend/`.

### Go project with full options
```
/project-init "CLI Tool" --template go --prefix cli --description "A command-line utility" --agents beads-workflows,developer-essentials
```
Creates Go project with custom description and multiple agent plugins.

### Skip automatic session start
```
/project-init "Quick Test" --skip-session
```
Creates project without starting first session.

## Success Criteria

Project init is successful when:
- ✅ Directory created and accessible
- ✅ Git repository initialized
- ✅ Beads initialized with prefix
- ✅ Template files copied and variables substituted
- ✅ CLAUDE.md generated with plugin configuration
- ✅ Initial issues created (unless skipped)
- ✅ Ready for first session

## Notes

**Performance**: Expect 30-60 seconds for complete initialization due to:
- File system operations
- Git initialization
- Multiple bd create commands for issues
- Template variable substitution

**Templates are extensible**: Additional templates can be added to `plugins/beads-workflows/templates/<name>/` with:
- `.gitignore`
- `initial-issues.json`
- Any `*.template` files for variable substitution

**Integration**: This command works with:
- `beads-session-start` for beginning first work session
- `beads-workflow-orchestrator` for intelligent issue selection
- All beads CLI commands for issue management
