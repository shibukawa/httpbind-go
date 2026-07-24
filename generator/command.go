package generator

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// CommandIO contains process state injected into a command execution.
type CommandIO struct {
	Stdin            io.Reader
	Stdout           io.Writer
	Stderr           io.Writer
	WorkingDirectory string
	Environment      []string
}

// Command is one independently testable subcommand.
type Command struct {
	Name    string
	Summary string
	Run     func(context.Context, []string, CommandIO) int
}

// GenerateCommand creates the tinybind generate subcommand.
func GenerateCommand(options Options) Command {
	return Command{
		Name:    "generate",
		Summary: "generate package-local tinybind artifacts",
		Run: func(ctx context.Context, args []string, streams CommandIO) int {
			return runGenerate(ctx, args, streams, options)
		},
	}
}

// CommandSet is an immutable command dispatcher.
type CommandSet struct {
	commands map[string]Command
	names    []string
}

// NewCommandSet validates commands and constructs an immutable dispatcher.
func NewCommandSet(commands ...Command) (CommandSet, error) {
	set := CommandSet{commands: make(map[string]Command, len(commands))}
	for _, command := range commands {
		if command.Name == "" || strings.ContainsAny(command.Name, " \t\r\n") || command.Run == nil {
			return CommandSet{}, fmt.Errorf("generator: invalid command %q", command.Name)
		}
		if _, exists := set.commands[command.Name]; exists {
			return CommandSet{}, fmt.Errorf("generator: duplicate command %q", command.Name)
		}
		set.commands[command.Name] = command
		set.names = append(set.names, command.Name)
	}
	if len(set.names) == 0 {
		return CommandSet{}, fmt.Errorf("generator: command set is empty")
	}
	sort.Strings(set.names)
	return set, nil
}

// MustCommandSet is NewCommandSet for process setup code and panics on invalid commands.
func MustCommandSet(commands ...Command) CommandSet {
	set, err := NewCommandSet(commands...)
	if err != nil {
		panic(err)
	}
	return set
}

// Run dispatches one command without reading process globals or terminating the process.
func (set CommandSet) Run(ctx context.Context, args []string, streams CommandIO) int {
	if ctx == nil {
		ctx = context.Background()
	}
	if streams.Stdout == nil {
		streams.Stdout = io.Discard
	}
	if streams.Stderr == nil {
		streams.Stderr = io.Discard
	}
	if len(args) == 0 {
		set.writeHelp(streams.Stdout)
		return 0
	}
	if args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		if len(args) == 1 {
			set.writeHelp(streams.Stdout)
			return 0
		}
		command, ok := set.commands[args[1]]
		if !ok {
			fmt.Fprintf(streams.Stderr, "unknown command %q\n", args[1])
			return 2
		}
		return command.Run(ctx, []string{"-h"}, streams)
	}
	command, ok := set.commands[args[0]]
	if !ok {
		fmt.Fprintf(streams.Stderr, "unknown command %q\n", args[0])
		set.writeHelp(streams.Stderr)
		return 2
	}
	return command.Run(ctx, args[1:], streams)
}

func (set CommandSet) writeHelp(writer io.Writer) {
	fmt.Fprintln(writer, "commands:")
	for _, name := range set.names {
		command := set.commands[name]
		fmt.Fprintf(writer, "  %-12s %s\n", command.Name, command.Summary)
	}
}

// Main owns only the outer process boundary.
func Main(set CommandSet) {
	workingDirectory, _ := os.Getwd()
	os.Exit(set.Run(context.Background(), os.Args[1:], CommandIO{
		Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr,
		WorkingDirectory: workingDirectory, Environment: os.Environ(),
	}))
}
