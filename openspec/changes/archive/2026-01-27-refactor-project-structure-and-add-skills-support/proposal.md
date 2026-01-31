## Why

Currently, files in the `~/.aiqconfig/` directory are directly placed in the root directory (`config.yaml`, `sources.yaml`, `session_*.json`). As features expand (Skills, tools, etc. will be added soon), the directory will become cluttered. We need to plan the user configuration directory structure at the architecture level, establish clear hierarchy and naming conventions to avoid future refactoring costs and user data migration issues.

Additionally, the project needs to support Claude Skills to enhance AI Agent capabilities, allowing users to extend functionality through custom Skills. Skills need to be stored in the user configuration directory and should be included in the directory organization planning.

## What Changes

1. **User Configuration Directory Structure Refactoring**: Reorganize the `~/.aiqconfig/` directory structure to establish clear hierarchy:
   - Categorize configuration files into subdirectories by type (e.g., `config/`, `sessions/`, `skills/`, `tools/`, etc.)
   - Establish unified naming conventions and file organization rules
   - Define directory organization standards for future expansion
   - **Note**: The project is in early stage, backward compatibility is not considered, users need to manually migrate existing files

2. **Claude Skills Support**:
   - Implement Skills loading mechanism, support loading user-defined Skills from `~/.aiqconfig/skills/` directory
   - Parse SKILL.md format files (YAML frontmatter + Markdown content)
   - **Progressive Loading**: Dynamically load relevant Skills based on user queries and context needs, rather than loading all Skills at once
   - **Prompt Management Mechanism**:
     - Monitor prompt length and token usage
     - Implement prompt compression strategy (compress or evict low-priority content when approaching token limits)
     - Establish Skills priority and eviction mechanism to avoid performance impact from overly long prompts
   - **Built-in Tools**: Provide basic toolset to support common operations in Skills:
     - HTTP request tool (support URL operations in Skills)
     - Command execution tool (support command invocation in Skills)
     - File operation tool (read, write, list files)
     - Other basic tools (extend based on Skills requirements)
   - Implement functionality similar to Claude Agent SDK (Go version)

3. **Configuration Directory Organization Standards**: Establish organization rules and best practices documentation for `~/.aiqconfig/` directory to ensure structure remains clear during future expansion

## Capabilities

### New Capabilities

- `user-config-directory-organization`: Refactor user configuration directory (`~/.aiqconfig/`) organization structure, establish clear hierarchy and naming conventions
- `claude-skills-support`: Support Claude Skills format, implement Skills loading, parsing and integration mechanisms

### Modified Capabilities

- `sql-interactive-mode`: Need to integrate Skills content into prompts to enhance AI Agent capabilities
- `configuration-management`: May need to add Skills directory configuration items

## Impact

- **User Configuration Directory Structure**: Refactor `~/.aiqconfig/` directory structure
- **Configuration File Paths**: All code accessing configuration files needs to update paths (`config.yaml` → `config/config.yaml`, `sources.yaml` → `config/sources.yaml`, `session_*.json` → `sessions/session_*.json`)
- **Configuration System**: Need to add Skills-related configuration (Skills directory path, etc.)
- **LLM Integration**: Need to integrate Skills content into prompt building logic, implement progressive loading and prompt management
- **Prompt Management**: Need to implement token monitoring, prompt compression, and content eviction mechanism
- **Tool System**: Need to extend tools registration mechanism, provide built-in toolset to support operations in Skills
- **Skills Matching**: Need to implement Skills matching algorithm with user queries to determine when to load which Skills
- **Dependency Management**: May need to add YAML parsing library (if not already used) for parsing Skills frontmatter
- **Documentation**: Need to update configuration directory structure documentation and Skills usage guide
