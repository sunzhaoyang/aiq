## Context

Repository has migrated from `sunzhaoyang/aiq` to `sunetic/aiq`. Current codebase contains multiple references to old repository path, including:
- Repository variables in installation scripts (`scripts/install.sh` and `scripts/install.bat`)
- Installation commands and clone commands in README documentation
- Other documentation references that may exist

**Current State**:
- `scripts/install.sh`: `REPO="sunzhaoyang/aiq"`
- `scripts/install.bat`: `set "REPO=sunzhaoyang/aiq"`
- `README.md`: 3 `sunzhaoyang/aiq` URLs
- `go.mod`: `module github.com/aiq/aiq` (module path, needs evaluation)

**Constraints**:
- Must maintain backward compatibility (if users already installed, should not affect existing installations)
- Archived historical documents (`openspec/changes/archive/`) don't need updating, maintain historical record integrity

## Goals / Non-Goals

**Goals:**
- Update repository path references in all active files from `sunzhaoyang/aiq` to `sunetic/aiq`
- Ensure installation scripts can correctly download binaries from new repository
- Ensure links in documentation point to new repository
- Maintain codebase consistency and accuracy

**Non-Goals:**
- Do not update archived historical documents (maintain historical records)
- Do not modify import paths in Go code (`github.com/aiq/aiq` is module path, can differ from repository path)
- Do not make functional changes, only update path references

## Decisions

### Decision 1: Update Repository Variables in Installation Scripts
**Choice**: Directly update `REPO` variable from `sunzhaoyang/aiq` to `sunetic/aiq`

**Rationale**: 
- Installation scripts depend on GitHub API and Releases downloads, must point to correct repository
- This is critical path for user installation, must be accurate

**Alternatives Considered**:
- Use environment variables: Increases complexity, users need additional configuration
- Auto-detection: Unreliable, may detect incorrectly

### Decision 2: Update URLs in README Documentation
**Choice**: Update all `sunzhaoyang/aiq` URL references to `sunetic/aiq`

**Rationale**:
- Documentation is user's first point of contact, must be accurate
- Includes installation commands, clone commands, and issue links

**Alternatives Considered**:
- Use redirects: GitHub may not support cross-user redirects
- Keep old links: Would cause user confusion

### Decision 3: Keep Go Module Path Unchanged
**Choice**: `module github.com/aiq/aiq` in `go.mod` remains unchanged

**Rationale**:
- Go module path is internal code identifier, can differ from actual GitHub repository path
- Changing module path would require updating all import statements, large impact scope
- If unification needed in future, can be handled as independent migration task

**Alternatives Considered**:
- Update to `github.com/sunetic/aiq`: Would require updating all import paths, high risk, low benefit

### Decision 4: Do Not Update Archived Documents
**Choice**: Historical documents under `openspec/changes/archive/` directory remain unchanged

**Rationale**:
- Archived documents are historical records, should remain as-is to reflect state at that time
- Updating historical documents would destroy historical accuracy

## Risks / Trade-offs

**Risk 1: Users Using Old Installation Commands**
- **Risk**: Users may copy old installation commands from old documentation or cache
- **Mitigation**: README already updated, search engines will gradually index new content; GitHub may provide redirects (if configured)

**Risk 2: Go Module Path Inconsistency with Repository Path**
- **Risk**: Developers may be confused why module path and repository path differ
- **Mitigation**: Explain in README that this is normal, Go module paths can be independent of repository paths; can migrate separately in future if needed

**Risk 3: Missing Some References**
- **Risk**: May have hidden configuration files or references in comments not discovered
- **Mitigation**: Use grep for comprehensive search, check all possible locations; confirm again during code review

**Trade-offs**:
- **Simplicity vs Consistency**: Choose to keep Go module path unchanged, sacrifice consistency but avoid large-scale refactoring risk

## Migration Plan

### Step 1: Update Installation Scripts
1. Update `REPO` variable in `scripts/install.sh`
2. Update `REPO` variable in `scripts/install.bat`
3. Verify script syntax correctness

### Step 2: Update Documentation
1. Update all `sunzhaoyang/aiq` URLs in `README.md`
2. Check and update relevant URLs in `README_CN.md` (if exist)
3. Verify all link formats are correct

### Step 3: Comprehensive Check
1. Use grep to search all possible repository path references
2. Confirm no files missed
3. Verify decision not to update archived documents

### Step 4: Testing and Verification
1. Verify installation scripts can download from new repository (if new repository already has releases)
2. Verify links in documentation are accessible
3. Confirm code builds normally (module path unchanged)

### Rollback Strategy
- If issues found, can quickly rollback file changes
- Since only path updates, no functional changes involved, low rollback risk

## Open Questions

1. **GitHub Redirect**: Has old repository `sunzhaoyang/aiq` configured redirect to new repository? This affects users using old links.
2. **Go Module Path**: Will we need to unify module path and repository path in future? If needed, should be handled as independent task.
