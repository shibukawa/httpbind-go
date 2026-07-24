---
id: requirement:config-file-discovery
type: requirement
title: Config File Discovery
---
Discover the first TOML config file from explicit, extra, user, and system candidates.

```yaml
priority: must
intent: locate config.toml-style files without merging multi-location contents
inputs:
  - vendor_name: required via configbind discovery API when search reaches configdir
  - tool_name: required via configbind discovery API when search reaches configdir
  - file_name: config file name under the tool config directory
  - optional --config-path from CLI
  - optional ExtraConfigReadPaths direct file paths
resolution: decision:config-file-path-resolution
behavior:
  - prefer --config-path over any directory search
  - if --config-path is unreadable or missing, fail with error; no silent fallback
  - else check ExtraConfigReadPaths in slice order; skip missing or unreadable entries
  - stop at the first readable extra file without reading later files
  - else search user config dir then system config dir via system:configdir
  - user path wins when both contain the file; only that file is read
  - system path used only when user path lacks the file
  - never merge keys from extra, user, or system files in one load
  - configdir.New(vendor_name, tool_name) receives both API arguments
cli:
  - flag long name: config-path
  - form: '--config-path <path>'
  - registered by cliparser / configbind process flags
  - maps to explicit path, not a Bind overlay config key under a prefix
tinygo:
  - system:configdir must remain TinyGo-buildable on host targets
  - requirement:configbind-tinygo
acceptance:
  - with --config-path=/tmp/app.toml, that path is used even if user/system files exist
  - with --config-path set to a missing path, load returns an error and does not search configdir
  - the first existing ExtraConfigReadPaths entry wins over later extras and configdir
  - missing ExtraConfigReadPaths entries are ignored
  - without --config-path, user config file is chosen over system when both exist
  - without --config-path and only system file, system file is used
  - extra, user, and system files are not merged
  - discovery API accepts vendor name, tool name, and file name
related:
  - decision:config-file-path-resolution
  - system:configdir
  - requirement:layered-config-load
  - requirement:cli-option-codegen
  - flow:config-load
  - system:configbind
```
