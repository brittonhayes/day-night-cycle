---
name: release
description: Analyze changes, determine semantic version bump, build executables, and cut a GitHub release with Claude-written minimalist notes.
allowed-tools: Read, Bash, Grep, Glob
---

# Release

Automatically analyze changes, determine semantic versioning, build binaries, and create GitHub releases with custom Claude-written release notes.

## Overview

This skill automates the complete release process:

1. Analyze git commits since last tag
2. Determine semantic version bump (major.minor.patch)
3. Build executables for all platforms using `make build-all`
4. Create git tag with new version
5. **Generate minimalist release notes (always written by Claude)**
6. Write release notes to temporary file
7. Create GitHub release using `make release RELEASE_NOTES=/tmp/release-notes.md`

## The Releaser Agent

The releaser agent writes minimalist, actionable release notes that:
- Focus on user impact, not implementation details
- Use clear, concise language
- Group related changes logically
- Avoid redundancy and noise
- Skip obvious or trivial changes

**Tone:** Professional, direct, minimal. Think release notes from tools like ripgrep or fd.

## Semantic Versioning Rules

Analyze commits since the last tag to determine the version bump:

**MAJOR (x.0.0)** - Breaking changes:
- Commits with `BREAKING CHANGE:` in body or footer
- Commits with `!` after type (e.g., `feat!:`, `fix!:`)
- Removed features or changed APIs

**MINOR (0.x.0)** - New features:
- Commits starting with `feat:`
- New plugins or commands
- New capabilities added

**PATCH (0.0.x)** - Bug fixes and small improvements:
- Commits starting with `fix:`
- Commits starting with `refactor:`, `docs:`, `test:`, `chore:`
- Performance improvements
- Documentation updates

**Default:** If no conventional commits, use PATCH.

## Instructions

### 1. Analyze Changes

Get the current version and commits since last tag:

```bash
# Get latest tag (current version)
git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"

# Get commits since last tag
git log $(git describe --tags --abbrev=0 2>/dev/null || git rev-list --max-parents=0 HEAD)..HEAD --oneline

# Get detailed commit messages for analysis
git log $(git describe --tags --abbrev=0 2>/dev/null || git rev-list --max-parents=0 HEAD)..HEAD --format="%h %s%n%b"
```

### 2. Determine Version Bump

Based on commits:
- Check for `BREAKING CHANGE` or `!` → MAJOR
- Check for `feat:` → MINOR
- Otherwise → PATCH

Calculate new version:
```bash
# Example: current v1.2.3
# MAJOR: v2.0.0
# MINOR: v1.3.0
# PATCH: v1.2.4
```

### 3. Build Executables

Use the Makefile to build for all platforms:

```bash
make build-all
```

This builds:
- `bin/day-night-cycle-darwin-amd64` (macOS Intel)
- `bin/day-night-cycle-darwin-arm64` (macOS Apple Silicon)

### 4. Create Git Tag

```bash
git tag -a v1.3.0 -m "Release v1.3.0"
git push origin v1.3.0
```

### 5. Generate Release Notes

Write minimalist notes following this format:

```markdown
[One sentence summary of the release]

## Changes

- [User-facing change 1]
- [User-facing change 2]
- [User-facing change 3]

## Fixes

- [Bug fix 1]
- [Bug fix 2]
```

**Guidelines:**
- **Be concise:** One line per change, no elaboration unless critical
- **User perspective:** What changed for users, not how it was implemented
- **Skip noise:** Don't mention refactoring, internal cleanup, or dependency updates unless they affect users
- **Group logically:** Features, fixes, breaking changes
- **No emotion:** Avoid "exciting", "amazing", "awesome" - just facts
- **No redundancy:** If the commit message is clear, use it directly

**Examples of good vs bad notes:**

✅ Good:
- Add support for Cursor IDE theme switching
- Fix launchd schedule generation for negative timezones
- Remove deprecated Name field from LocationConfig

❌ Bad:
- We're excited to announce amazing new Cursor IDE support! 🎉
- This release includes a critical bug fix that was causing issues
- Refactored internal plugin architecture for better maintainability

### 6. Create GitHub Release

**CRITICAL: You MUST always write custom release notes. Never use auto-generated notes.**

Write the release notes to a temporary file, then use the Makefile to create the release:

```bash
# Write notes to a file
cat > /tmp/release-notes.md << 'EOF'
[One sentence summary of the release]

## Changes

- [User-facing change 1]
- [User-facing change 2]

## Fixes

- [Bug fix 1]
EOF

# Create release with custom notes
make release RELEASE_NOTES=/tmp/release-notes.md
```

This will:
- Build all platform binaries with `make build-all`
- Upload binaries to GitHub
- Create release with your custom notes from the file
- Use the tag created earlier

**IMPORTANT:** Always provide the `RELEASE_NOTES` parameter. The Makefile has a fallback to auto-generated notes, but you must never use it. Claude-written release notes are required for all releases.

## Checklist

Before completing the release:

- [ ] Analyzed all commits since last tag
- [ ] Determined correct semantic version bump
- [ ] Calculated new version number
- [ ] Built executables with `make build-all`
- [ ] Verified build artifacts exist in `bin/`
- [ ] Created git tag with new version
- [ ] Pushed tag to origin
- [ ] **Generated minimalist release notes (REQUIRED)**
- [ ] **Wrote release notes to /tmp/release-notes.md (REQUIRED)**
- [ ] **Created GitHub release with `make release RELEASE_NOTES=/tmp/release-notes.md` (REQUIRED)**
- [ ] Verified release appears on GitHub with binaries and custom release notes

## Example Workflow

User: "Cut a new release"

Expected workflow:
1. Run `git describe --tags --abbrev=0` to get current version (e.g., v1.2.0)
2. Run `git log v1.2.0..HEAD --format="%h %s%n%b"` to get commits
3. Analyze commits:
   - Found: `feat: add Cursor IDE plugin`
   - Found: `fix: handle timezone edge case`
   - Found: `refactor: simplify main.go`
4. Determine: MINOR bump (feat: found, no breaking changes)
5. Calculate: v1.2.0 → v1.3.0
6. Run `make build-all` to build binaries
7. Create tag: `git tag -a v1.3.0 -m "Release v1.3.0"`
8. Push tag: `git push origin v1.3.0`
9. Write release notes to file:
   ```bash
   cat > /tmp/release-notes.md << 'EOF'
   Add Cursor IDE support and fix timezone handling.

   ## Changes
   - Add support for Cursor IDE theme switching

   ## Fixes
   - Fix timezone handling for negative offsets
   EOF
   ```
10. Run `make release RELEASE_NOTES=/tmp/release-notes.md` to create GitHub release

## Common Scenarios

### First Release (No Tags)

If `git describe --tags` fails, this is the first release:
- Start with `v0.1.0` or `v1.0.0` depending on maturity
- Include all commits in release notes summary
- Focus on what the tool does, not individual commits

### Patch Release (Only Fixes)

If only `fix:`, `docs:`, `chore:`, `refactor:` commits:
- Bump PATCH version
- Keep notes very brief
- Group related fixes

### Major Release (Breaking Changes)

If `BREAKING CHANGE:` found:
- Bump MAJOR version
- Clearly list breaking changes first
- Explain migration path if needed

## Notes

- **CRITICAL:** Always generate and use custom release notes - never rely on auto-generated notes
- Always write release notes to `/tmp/release-notes.md` before running `make release`
- Always pass `RELEASE_NOTES=/tmp/release-notes.md` to the make release command
- Always verify the Makefile is being used correctly
- The `make release` target requires GitHub CLI (`gh`)
- Never release from a dirty working tree (warn user if uncommitted changes)
- Minimize release note length while maintaining clarity
- When in doubt, bump PATCH (smallest change)
- Breaking changes require explicit documentation
