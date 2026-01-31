## 1. Update Installation Scripts

- [x] 1.1 Update `REPO` variable in `scripts/install.sh` from `sunzhaoyang/aiq` to `sunetic/aiq`
- [x] 1.2 Update `REPO` variable in `scripts/install.bat` from `sunzhaoyang/aiq` to `sunetic/aiq`
- [x] 1.3 Verify installation script syntax correctness (bash and batch scripts)

## 2. Update Documentation

- [x] 2.1 Update installation command URL in `README.md` line 40: `https://raw.githubusercontent.com/sunzhaoyang/aiq/main/scripts/install.sh` → `https://raw.githubusercontent.com/sunetic/aiq/main/scripts/install.sh`
- [x] 2.2 Update Windows installation command URL in `README.md` line 47: `https://raw.githubusercontent.com/sunzhaoyang/aiq/main/scripts/install.bat` → `https://raw.githubusercontent.com/sunetic/aiq/main/scripts/install.bat`
- [x] 2.3 Update git clone command in `README.md` line 82: `https://github.com/sunzhaoyang/aiq.git` → `https://github.com/sunetic/aiq.git`
- [x] 2.4 Update issue link in `README.md` line 186: `https://github.com/aiq/aiq/issues` → `https://github.com/sunetic/aiq/issues` (if exists)
- [x] 2.5 Check and update relevant URLs in `README_CN.md` (if `sunzhaoyang/aiq` exists or links need updating)
- [x] 2.6 Verify all link formats are correct, ensure URLs are accessible

## 3. Comprehensive Check

- [x] 3.1 Use grep to search all `sunzhaoyang` references, confirm no files missed
- [x] 3.2 Use grep to search all GitHub URL references (`raw.githubusercontent.com`, `github.com/.*/releases`, `github.com/.*/blob`, `github.com/.*/tree`), confirm all updated
- [x] 3.3 Confirm archived documents (`openspec/changes/archive/`) remain unchanged
- [x] 3.4 Confirm module path `github.com/aiq/aiq` in `go.mod` remains unchanged (conforms to design decision)

## 4. Testing and Verification

- [x] 4.1 Verify `scripts/install.sh` syntax is correct (using `bash -n scripts/install.sh`)
- [x] 4.2 Verify `scripts/install.bat` syntax is correct (test in Windows environment, or use syntax checking tools)
- [x] 4.3 Verify link formats in documentation are correct (check Markdown link syntax)
- [x] 4.4 Confirm code builds normally (run `go build` to verify module path change doesn't affect build)
- [x] 4.5 If new repository already has releases, verify installation scripts can correctly download from new repository (optional, depends on new repository status)
