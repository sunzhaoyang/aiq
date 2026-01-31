## 1. Discovery and Inventory

- [x] 1.1 Run grep to find all files containing Chinese text in `internal/` directory (Found 60+ files)
- [x] 1.2 Run grep to find all files containing Chinese text in `openspec/` directory (Found 100+ files)
- [x] 1.3 Create inventory list of files to translate, categorized by priority
- [x] 1.4 Verify README_CN.md is excluded from inventory

## 2. Code Comments Translation

- [x] 2.1 Translate Chinese comments in `internal/tool/builtin/command_tool_test.go`
- [x] 2.2 Translate Chinese comments in `internal/llm/client_test.go`
- [ ] 2.3 Search for any other Go source files with Chinese comments and translate
- [ ] 2.4 Review translations for technical accuracy
- [ ] 2.5 Verify comment formatting is preserved
- [ ] 2.6 Commit code comment translations

## 3. Test File Comments Translation

- [ ] 3.1 Translate test function comments (e.g., `// Task 5.4: ...`)
- [ ] 3.2 Translate test scenario descriptions in `t.Run()` calls
- [ ] 3.3 Translate any Chinese error messages in test assertions
- [ ] 3.4 Review test translations for clarity and accuracy
- [ ] 3.5 Verify test functionality remains unchanged
- [ ] 3.6 Commit test comment translations

## 4. Active Spec Files Translation

- [ ] 4.1 Identify all spec files in `openspec/specs/` containing Chinese text
- [ ] 4.2 Translate spec file: `code-comment-translation/spec.md` (if contains Chinese)
- [ ] 4.3 Translate spec file: `spec-document-translation/spec.md` (if contains Chinese)
- [ ] 4.4 Translate spec file: `test-comment-translation/spec.md` (if contains Chinese)
- [ ] 4.5 Translate any other active spec files with Chinese text
- [ ] 4.6 Review translations for technical accuracy and markdown formatting
- [ ] 4.7 Commit active spec file translations

## 5. Active Change Documents Translation

- [ ] 5.1 Identify all change documents in `openspec/changes/` (non-archive) containing Chinese text
- [ ] 5.2 Translate proposal documents (proposal.md files)
- [ ] 5.3 Translate design documents (design.md files)
- [ ] 5.4 Translate spec documents in change directories
- [ ] 5.5 Translate task documents (tasks.md files)
- [ ] 5.6 Review translations for consistency and accuracy
- [ ] 5.7 Commit active change document translations

## 6. Archived Change Documents Translation

- [ ] 6.1 Identify all archived change documents in `openspec/changes/archive/` containing Chinese text
- [ ] 6.2 Translate archived proposal documents
- [ ] 6.3 Translate archived design documents
- [ ] 6.4 Translate archived spec documents
- [ ] 6.5 Translate archived task documents
- [ ] 6.6 Translate archived alternative design documents (if any)
- [ ] 6.7 Review translations for historical accuracy
- [ ] 6.8 Commit archived change document translations

## 7. Verification and Quality Assurance

- [ ] 7.1 Run final grep to verify no Chinese text remains (except README_CN.md)
- [ ] 7.2 Review translations for terminology consistency
- [ ] 7.3 Verify all markdown files render correctly
- [ ] 7.4 Verify code comments don't break line length limits
- [ ] 7.5 Check that test files still compile and run
- [ ] 7.6 Update terminology glossary if new terms were added
- [ ] 7.7 Create summary of translation work completed

## 8. Final Review and Commit

- [ ] 8.1 Review all translated files for consistency
- [ ] 8.2 Verify README_CN.md remains unchanged
- [ ] 8.3 Ensure no functional code changes were made
- [ ] 8.4 Final commit with comprehensive translation changes
- [ ] 8.5 Update change status to complete
