---
id: decision:config-file-path-resolution
type: decision
title: Config File Path Resolution
---
Resolve the first readable TOML path from explicit, extra, user, and system candidates; never merge files.

```yaml
status: accepted
lookup_inputs:
  - vendor_name: required API argument when search reaches configdir vendor segment
  - tool_name: required application name when search reaches configdir
  - file_name: config file basename e.g. config.toml
  - optional --config-path from CLI when present
  - optional ExtraConfigReadPaths file paths in caller order
api_sketch:
  - 'ResolveWithExtras(vendor, tool, fileName, explicitPath, extraReadPaths) (path string, ok bool, err error)'
  - vendor and tool are passed by the app through configbind public API
exclusive_file:
  - at most one file path is chosen for TOML load
  - stop at the first readable candidate
  - extra, user, and system files are never content-merged
priority_high_to_low:
  - id: explicit_config_path
    source: LoadOptions.ExplicitConfigPath or cliparser flag --config-path
    meaning: explicit filesystem path to the config file
    wins_over: all remaining candidates
  - id: extra_config_read_paths
    source: LoadOptions.ExtraConfigReadPaths
    meaning: optional direct file paths checked in slice order
    missing: skip without error
  - id: user_config_dir
    source: system:configdir user-level folder for vendor+tool
    meaning: first existing file_name under user config dir
  - id: system_config_dir
    source: system:configdir system-level folder for vendor+tool
    meaning: used only when user path does not contain file_name
search_helper: system:configdir
rules:
  - if --config-path is set, use that path only; skip directory search
  - if --config-path is set but missing or unreadable, return error; no fallback to configdir
  - after no explicit path, return the first readable ExtraConfigReadPaths entry
  - missing or unreadable extra entries are ignored
  - when no extra entry exists, QueryFolderContainsFile-style search user then system
  - if no file found without --config-path, TOML layer is absent; defaults/env/CLI still apply
  - --config-path is a process-level path flag, not a Bind field under a prefix table
  - vendor_name and tool_name are always supplied by the application API caller
related:
  - requirement:config-file-discovery
  - system:configdir
  - decision:cli-flag-naming
  - concept:cli-option-codegen
  - requirement:layered-config-load
  - flow:config-load
  - system:configbind
```
