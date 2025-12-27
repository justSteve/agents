---
name: description-quality
description: Ensure high-quality issue descriptions that enable effective work. Use when creating or updating issues to provide sufficient context, clear problem statements, and verifiable acceptance criteria.
---

# Description Quality

Write issue descriptions that provide enough context for any agent to understand the problem, implement a solution, and verify completion without needing to ask clarifying questions.

## When to Use This Skill

- Creating new issues with `bd create`
- Updating existing issues with `bd update`
- Reviewing issues before claiming them
- Filing discovered work during implementation
- Breaking down epics into actionable tasks

## The Quality Standard

A high-quality issue description answers these questions:

1. **What?** - What needs to be done (clear problem statement)
2. **Why?** - Why does this matter (context and motivation)
3. **Where?** - Where in the codebase (specific locations)
4. **How?** - How to verify completion (acceptance criteria)
5. **Notes?** - Any technical considerations (constraints, edge cases)

## Required Sections

### 1. Context

**Purpose**: Explain the background and motivation.

**Include**:
- Why this work is needed
- What triggered this issue
- How it relates to larger goals
- Any relevant history

**Example**:
```markdown
## Context

The authentication system currently uses a 1-hour token expiry, but users
are reporting frequent session timeouts during long form submissions.
This is causing data loss and user frustration. Product has prioritized
this for the Q1 release.
```

### 2. Problem Statement

**Purpose**: Clearly define what needs to change.

**Include**:
- Current behavior (what's wrong or missing)
- Desired behavior (what should happen)
- Scope boundaries (what's NOT included)

**Example**:
```markdown
## Problem

**Current**: Auth tokens expire after 1 hour with no refresh mechanism.
Users lose unsaved work if their session expires during form submission.

**Desired**: Implement token refresh that extends sessions during active
use. Tokens should refresh transparently without interrupting user workflow.

**Scope**: This issue covers the token refresh mechanism only. UI session
warnings are tracked in issue `session-warning-ui`.
```

### 3. Acceptance Criteria

**Purpose**: Define verifiable conditions for completion.

**Format**: Use checkboxes for clear verification.

**Include**:
- Functional requirements (what it must do)
- Edge cases (what happens in unusual conditions)
- Non-functional requirements (performance, security)

**Example**:
```markdown
## Acceptance Criteria

- [ ] Tokens refresh automatically when less than 15 minutes remain
- [ ] Refresh happens on any authenticated API call
- [ ] Failed refresh attempts redirect to login
- [ ] Refresh tokens are stored securely (httpOnly cookie)
- [ ] Token refresh completes in < 100ms
- [ ] Existing sessions remain valid after deployment
```

### 4. Technical Notes

**Purpose**: Provide implementation guidance and constraints.

**Include**:
- Relevant code locations
- Technical constraints
- Dependencies or related systems
- Known edge cases
- Security considerations

**Example**:
```markdown
## Technical Notes

- Token logic is in `src/auth/token-manager.ts`
- Refresh endpoint should be `POST /api/auth/refresh`
- Must maintain backward compatibility with mobile clients (v2.1+)
- Consider rate limiting refresh requests (max 1 per minute)
- Edge case: handle clock skew between client and server
```

## The Minimum Viable Description

For quick tasks, at minimum include:

```markdown
## Context
[1-2 sentences on why this matters]

## Problem
[What needs to change]

## Acceptance Criteria
- [ ] [At least one verifiable criterion]

## Location
[File path or component name]
```

**Example minimal description**:
```markdown
## Context
The login button color doesn't match the new brand guidelines.

## Problem
Change login button from blue (#0066CC) to brand green (#00AA55).

## Acceptance Criteria
- [ ] Login button uses color #00AA55
- [ ] Hover state uses #008844

## Location
src/components/LoginButton.tsx
```

## Good vs Bad Descriptions

### Example 1: Bug Fix

**Bad description**:
```
Fix the login bug
```

**Why it's bad**:
- No context (what bug?)
- No steps to reproduce
- No acceptance criteria
- Could mean anything

**Good description**:
```markdown
## Context
Users on Safari 16+ are unable to log in. The authentication endpoint
returns 401 but credentials are correct. This affects approximately 15%
of our user base.

## Problem
**Current**: Login fails on Safari 16+ with error "Invalid credentials"
even with correct username/password.

**Reproduction**:
1. Open Safari 16 or later
2. Navigate to /login
3. Enter valid credentials
4. Click "Sign In"
5. Error: "Invalid credentials"

**Expected**: Login succeeds and redirects to dashboard.

## Acceptance Criteria
- [ ] Login works on Safari 16, 17 (latest)
- [ ] Login works on Chrome, Firefox, Edge (regression test)
- [ ] Error logging captures browser info for future debugging

## Technical Notes
- Suspect SameSite cookie issue introduced in Safari 16
- Auth cookies are set in `src/api/auth-middleware.ts`
- Related: Apple's ITP changes in WebKit
```

### Example 2: New Feature

**Bad description**:
```
Add dark mode
```

**Why it's bad**:
- Massive scope (where? how?)
- No acceptance criteria
- No technical direction
- Could take hours or weeks

**Good description**:
```markdown
## Context
User research shows 73% of users prefer dark mode for evening use.
This is Phase 1 of the theming initiative, focusing on the core
dashboard only.

## Problem
**Current**: Application only supports light theme.

**Desired**: Users can toggle between light and dark themes on the
dashboard. Theme preference persists across sessions.

**Scope**: Dashboard only. Settings page and login are Phase 2
(tracked in `theme-phase-2`).

## Acceptance Criteria
- [ ] Theme toggle visible in dashboard header
- [ ] Dark theme applies to all dashboard components
- [ ] Theme preference saved to localStorage
- [ ] Theme loads correctly on page refresh
- [ ] No flash of wrong theme on load (FOUC prevention)
- [ ] Respects system preference if no user preference set

## Technical Notes
- Use CSS custom properties for theme values
- Theme context: `src/contexts/ThemeContext.tsx`
- Dashboard components in `src/components/dashboard/`
- Consider reduced motion for theme transition
- Test contrast ratios for accessibility (WCAG AA)
```

### Example 3: Discovered Work

**Bad description**:
```
Found a bug while working on auth
```

**Why it's bad**:
- No detail on the bug
- No reproduction steps
- No connection to original work
- Future agent can't act on this

**Good description**:
```markdown
## Context
Discovered while implementing token refresh (issue `token-refresh`).
Not blocking current work but should be addressed.

## Problem
**Current**: Password reset tokens are stored in plaintext in the
database. This is a security vulnerability.

**Reproduction**: Query `password_reset_tokens` table - tokens are
visible in plaintext.

**Expected**: Tokens should be hashed before storage.

## Acceptance Criteria
- [ ] Reset tokens stored as bcrypt hash
- [ ] Token verification works with hashed storage
- [ ] Existing tokens invalidated (force new reset requests)
- [ ] Migration handles existing token cleanup

## Technical Notes
- Token storage in `src/models/PasswordReset.ts`
- Follow same pattern as auth tokens in `src/models/AuthToken.ts`
- Security priority: should be P1

## Discovery Context
Found during: `token-refresh`
Related code: `src/auth/password-reset.ts:47`
```

## Common Patterns

### Pattern 1: Epic Breakdown

When breaking an epic into tasks:

```markdown
## Epic: User Authentication System

### Child Issues:

1. **user-model**: Create User database model
   - Context: Foundation for auth system
   - Criteria: Model with email, password_hash, created_at
   - Location: src/models/User.ts

2. **auth-register**: Implement registration endpoint
   - Context: First user-facing auth feature
   - Criteria: POST /register, validation, duplicate check
   - Location: src/api/auth/register.ts
   - Depends: user-model

3. **auth-login**: Implement login endpoint
   - Context: Core authentication flow
   - Criteria: POST /login, token generation, rate limiting
   - Location: src/api/auth/login.ts
   - Depends: user-model
```

### Pattern 2: Refactoring Task

```markdown
## Context
The `UserService` class has grown to 1200 lines with mixed
responsibilities. This violates SRP and makes testing difficult.

## Problem
**Current**: Monolithic UserService handles auth, profile,
preferences, and notifications.

**Desired**: Split into focused services: AuthService,
ProfileService, PreferencesService, NotificationService.

## Acceptance Criteria
- [ ] Each service < 300 lines
- [ ] Single responsibility per service
- [ ] Existing tests pass after refactor
- [ ] No public API changes (internal refactor only)
- [ ] Dependency injection maintained

## Technical Notes
- Current file: src/services/UserService.ts
- Target files: src/services/auth/, src/services/profile/, etc.
- Use strangler pattern: create new, delegate, remove old
- Watch for circular dependencies
```

### Pattern 3: Performance Issue

```markdown
## Context
Dashboard load time increased from 1.2s to 4.8s after the recent
data grid update. Users are complaining about sluggish performance.

## Problem
**Current**: Dashboard takes 4.8s to load.

**Desired**: Dashboard loads in < 2s.

**Measurement**: Use Chrome DevTools Performance tab, measure
Largest Contentful Paint (LCP).

## Acceptance Criteria
- [ ] Dashboard LCP < 2000ms
- [ ] Data grid renders visible rows only (virtualization)
- [ ] API response cached for 5 minutes
- [ ] No regression in functionality

## Technical Notes
- Profiling shows data grid rendering all 10,000 rows
- Consider react-virtualized or react-window
- API endpoint: GET /api/dashboard/data
- Current component: src/components/Dashboard/DataGrid.tsx
```

## Anti-Patterns to Avoid

### Anti-Pattern 1: Vague One-Liners

**Problem**: Description provides no actionable information.

```
Fix the thing
Update styles
Refactor code
```

**Solution**: Add context, specific problem, and acceptance criteria.

### Anti-Pattern 2: Missing Context

**Problem**: What needs to be done is clear, but not why.

```
Change button color to red
```

**Solution**: Explain the motivation:
```
Change error button color to red (#FF0000) to align with
accessibility guidelines. Current gray doesn't meet contrast
requirements.
```

### Anti-Pattern 3: No Acceptance Criteria

**Problem**: No way to verify completion.

```
Improve search performance
```

**Solution**: Add measurable criteria:
```
Acceptance Criteria:
- [ ] Search returns results in < 200ms for queries up to 1000 chars
- [ ] Memory usage stays below 512MB during search
- [ ] Typo tolerance works for 1-2 character errors
```

### Anti-Pattern 4: Scope Creep in Description

**Problem**: Description includes multiple unrelated tasks.

```
Fix login bug, also update the header styling, and refactor
the auth module while we're at it
```

**Solution**: One issue per concern. File separate issues:
- `fix-login-bug`: Fix the specific login bug
- `update-header-styles`: Header styling changes
- `refactor-auth`: Auth module refactoring

### Anti-Pattern 5: Implementation Prescription

**Problem**: Description dictates implementation instead of outcomes.

```
Use a HashMap with O(1) lookup to cache user sessions in Redis
with a 3600 second TTL using the ioredis library.
```

**Solution**: Describe the problem and criteria, let implementer choose approach:
```
Problem: User session lookups are slow (200ms average).
Criteria: Session lookups complete in < 10ms.
Notes: Consider caching solutions. Redis available in infra.
```

## Quality Checklist

Before submitting an issue, verify:

- [ ] **Context**: Is it clear why this matters?
- [ ] **Problem**: Is the current vs desired state explained?
- [ ] **Scope**: Are boundaries clear (what's NOT included)?
- [ ] **Criteria**: Can someone verify completion objectively?
- [ ] **Location**: Are relevant files/components mentioned?
- [ ] **Dependencies**: Are blockers and related issues linked?
- [ ] **Standalone**: Can someone work on this without asking questions?

## Quick Reference Card

**Minimum description structure**:
```markdown
## Context
[Why this matters - 1-2 sentences]

## Problem
Current: [What's wrong]
Desired: [What should happen]

## Acceptance Criteria
- [ ] [Verifiable criterion 1]
- [ ] [Verifiable criterion 2]

## Technical Notes
- Location: [file paths]
- [Any constraints or considerations]
```

**Quality test**:
> "Could a new team member work on this without asking clarifying questions?"

**Golden rule**:
> Write descriptions for your future self who has forgotten all context.

## Summary

**Core principle**: Good descriptions enable autonomous work.

**Key elements**:
1. Context - Why does this matter?
2. Problem - What needs to change?
3. Acceptance Criteria - How do we verify completion?
4. Technical Notes - What should the implementer know?

**Remember**:
- Vague descriptions waste everyone's time
- Acceptance criteria prevent "is this done?" confusion
- Scope boundaries prevent creep
- Technical notes accelerate implementation
- One issue = one concern (no bundling)

Master description quality and you'll create issues that any agent can pick up and complete successfully!
