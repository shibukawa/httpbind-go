---
id: requirement:extensible-generator-command
type: requirement
title: Extensible Framework Generator Command
---
Frameworks compose tinybind generation with framework-owned lifecycle subcommands in one testable command binary.

```yaml
priority: must
command_model:
  api: api:generator-main
  built_in: generate command backed by api:generator-execution
  framework_owned_examples:
    init: create a project from framework templates
    build: run generation then framework build
    watch: observe source changes and rerun generation plus build
composition:
  - framework registers commands without forking or copying tinybind flag parsing
  - framework can add, replace, or omit commands before dispatcher construction
  - each command owns its arguments and FlagSet below the subcommand boundary
  - duplicate command names fail during construction
  - root help lists generated and framework-owned commands consistently
execution:
  - inject args, stdin, stdout, stderr, working directory, environment, and context
  - return exit status or error to the single process entry point
  - only the outer Main boundary may call os.Exit
  - build and watch invoke api:generator-execution in-process rather than recursively executing the generate CLI
  - watch respects context cancellation and does not prescribe a filesystem watcher implementation
  - init, compiler invocation, and file watching remain framework responsibilities
testability:
  - dispatch and each command run without global os.Args or process termination
  - generate command tests use in-memory output writers
  - framework command tests can substitute generator execution
acceptance:
  - one framework binary supports generate, init, build, and watch
  - generate accepts the same package and artifact options as direct generator execution
  - build stops before compilation when generation fails
  - watch reports generation or build errors and can continue according to framework policy
  - command-specific help and unknown-command diagnostics are deterministic
related:
  - api:generator-main
  - api:generator-execution
  - data:generator-options
  - requirement:framework-wrapper-discovery
```
