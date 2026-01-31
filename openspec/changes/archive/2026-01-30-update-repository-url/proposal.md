## Why

Repository path has migrated from `sunzhaoyang/aiq` to `sunetic/aiq`. Need to update all references to old repository path in the project to ensure installation scripts, documentation, and code repository references point to the new correct path.

## What Changes

- Update repository path references in installation scripts (`scripts/install.sh` and `scripts/install.bat`)
- Update repository URLs in README documentation (`README.md` and `README_CN.md`)
- Check and update repository path references that may exist in code
- **Note**: Module path in `go.mod` should usually remain unchanged unless Go module path migration is actually needed

## Capabilities

### New Capabilities
<!-- No new capabilities -->

### Modified Capabilities
- `installation-script`: Installation scripts need to update GitHub repository path references from `sunzhaoyang/aiq` to `sunetic/aiq`

## Impact

**Affected Files**:
- `scripts/install.sh` - Unix/Linux/macOS installation script
- `scripts/install.bat` - Windows installation script
- `README.md` - English documentation
- `README_CN.md` - Chinese documentation
- Other documentation or comments that may contain repository references

**Impact Scope**:
- User installation experience: Ensure installation scripts can download binaries from correct repository
- Documentation accuracy: Ensure links and examples in documentation point to correct repository
- Developer experience: Ensure developers can find correct repository address

**Notes**:
- Module path `github.com/aiq/aiq` in `go.mod` may need evaluation for whether update is needed
- Need to check if other configuration files or documentation reference old path
