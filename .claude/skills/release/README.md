# Release Skill

A Claude Code skill for automating semantic versioning and GitHub releases with minimalist release notes.

## Purpose

This skill (the "releaser") automates the complete release workflow:

1. Analyzes commits since the last tag
2. Determines the appropriate semantic version bump
3. Builds executables for all platforms
4. Creates git tags
5. Generates minimalist release notes
6. Creates GitHub releases with binaries

## Usage

Simply ask Claude Code to cut a release:

```
Cut a new release
```

```
Create a release
```

```
Make a release with the latest changes
```

The skill will automatically activate and handle the entire release process.

## What It Does

### Analysis Phase
- Gets the current version from git tags
- Analyzes all commits since the last release
- Applies semantic versioning rules based on conventional commits
- Determines whether to bump MAJOR, MINOR, or PATCH version

### Build Phase
- Runs `make build-all` to compile binaries for:
  - macOS Intel (darwin-amd64)
  - macOS Apple Silicon (darwin-arm64)
- Verifies build artifacts exist

### Release Phase
- Creates a git tag with the new version
- Pushes the tag to origin
- Generates minimalist release notes (see philosophy below)
- Creates a GitHub release using `make release`
- Uploads binaries automatically

## Release Notes Philosophy

The releaser writes minimalist release notes inspired by tools like ripgrep, fd, and bat:

**Principles:**
- **Concise**: One line per change, no elaboration
- **User-focused**: What changed for users, not implementation details
- **No noise**: Skip refactoring, internal cleanup, dependency updates
- **No emotion**: Avoid superlatives like "exciting" or "awesome"
- **Direct**: Facts only, clear language

**Example:**

✅ Good:
```markdown
Add Cursor IDE support and fix timezone handling.

## Changes
- Add support for Cursor IDE theme switching

## Fixes
- Fix timezone handling for negative offsets
```

❌ Bad:
```markdown
We're excited to announce an amazing new release! 🎉

This release includes:
- We've added some really cool new Cursor IDE support that you'll love
- Fixed a critical bug that was causing major issues for some users
- Refactored the internal architecture for better performance
```

## Semantic Versioning Rules

The skill follows standard semantic versioning:

**MAJOR (x.0.0)** - Breaking changes:
- Commits with `BREAKING CHANGE:` in the body
- Commits with `!` after type (e.g., `feat!:`)
- API changes or removed features

**MINOR (0.x.0)** - New features:
- Commits starting with `feat:`
- New plugins or commands
- Backward-compatible additions

**PATCH (0.0.x)** - Bug fixes and improvements:
- Commits starting with `fix:`
- Documentation updates (`docs:`)
- Refactoring (`refactor:`)
- Performance improvements

**Default**: If commits don't follow conventions, defaults to PATCH bump.

## Prerequisites

The skill expects:
- Git repository with tags for versioning
- GitHub CLI (`gh`) installed and authenticated
- Clean working directory (no uncommitted changes)
- Makefile with `build-all` and `release` targets
- Conventional commit format (recommended)

## Examples

### Example 1: Feature Release

**Commits since v1.2.0:**
```
feat: add Cursor IDE plugin
fix: handle negative timezone offsets
docs: update README with Cursor example
```

**Skill actions:**
1. Detects `feat:` commit → MINOR bump
2. Calculates v1.2.0 → v1.3.0
3. Builds binaries with `make build-all`
4. Creates tag v1.3.0
5. Generates notes:
   ```
   Add Cursor IDE support and fix timezone handling.

   ## Changes
   - Add support for Cursor IDE theme switching

   ## Fixes
   - Fix timezone handling for negative offsets
   ```
6. Creates GitHub release

### Example 2: Bug Fix Release

**Commits since v1.3.0:**
```
fix: correct solar calculation for edge case
refactor: simplify plugin error handling
```

**Skill actions:**
1. Only `fix:` and `refactor:` → PATCH bump
2. Calculates v1.3.0 → v1.3.1
3. Generates brief notes:
   ```
   Fix solar calculations.

   ## Fixes
   - Correct solar calculation for edge case
   ```

### Example 3: Breaking Change Release

**Commits since v1.3.1:**
```
feat!: change plugin configuration format
docs: update migration guide
```

**Skill actions:**
1. Detects `!` → MAJOR bump
2. Calculates v1.3.1 → v2.0.0
3. Generates notes with breaking change warning:
   ```
   Breaking: Plugin configuration format changed.

   ## Breaking Changes
   - Plugin configuration now uses nested structure
   - See migration guide for update instructions

   ## Changes
   - Update plugin configuration format
   ```

## Workflow Integration

The skill uses your existing Makefile workflow:

1. **make build-all**: Compiles binaries for all platforms
2. **make release**: Creates GitHub release with binaries

This ensures consistency with your manual release process.

## Benefits

- **Automated**: No manual version calculation or changelog writing
- **Consistent**: Always follows semantic versioning rules
- **Quality**: Minimalist notes focus on user value
- **Fast**: Complete release in seconds
- **Safe**: Verifies builds before creating releases

## Customization

To adjust the skill behavior, edit `.claude/skills/release/SKILL.md`:
- Modify semantic versioning rules
- Adjust release note format
- Add additional validation steps
- Change build or release commands

## Troubleshooting

**"Cannot release 'dev' version"**
- Create a git tag first, or let the skill create one

**"GitHub CLI not installed"**
- Install with: `brew install gh`
- Authenticate with: `gh auth login`

**"No commits since last tag"**
- Make some changes and commit them first

**"Working directory not clean"**
- Commit or stash changes before releasing

## Files

- **SKILL.md** - Main skill definition and agent instructions
- **README.md** - This file, user-facing documentation

## Notes

- The skill is called "release" but acts as a "releaser" agent
- Release notes emphasize brevity and clarity
- Follows conventional commits when available
- Always builds before releasing
- Verifies artifacts exist before pushing

## Related Documentation

- [Makefile](../../../Makefile) - Build and release targets
- [README.md](../../../README.md) - Project overview
