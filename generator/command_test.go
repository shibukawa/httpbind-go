package generator_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/shibukawa/tinybind-go/generator"
)

func TestCommandSetComposesFrameworkLifecycleCommands(t *testing.T) {
	var calls []string
	var seenIO generator.CommandIO
	command := func(name string) generator.Command {
		return generator.Command{
			Name: name, Summary: name + " summary",
			Run: func(ctx context.Context, args []string, streams generator.CommandIO) int {
				if err := ctx.Err(); err != nil {
					return 1
				}
				seenIO = streams
				calls = append(calls, name+":"+strings.Join(args, ","))
				return 0
			},
		}
	}
	set, err := generator.NewCommandSet(
		generator.GenerateCommand(generator.DefaultOptions()),
		command("init"), command("build"), command("watch"),
	)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	streams := generator.CommandIO{Stdout: &stdout, Stderr: &stderr, WorkingDirectory: t.TempDir(), Environment: []string{"MODE=test"}}
	if code := set.Run(context.Background(), []string{"build", "./service"}, streams); code != 0 {
		t.Fatalf("build exit=%d stderr=%s", code, stderr.String())
	}
	if len(calls) != 1 || calls[0] != "build:./service" {
		t.Fatalf("calls=%v", calls)
	}
	if seenIO.WorkingDirectory != streams.WorkingDirectory ||
		len(seenIO.Environment) != 1 || seenIO.Environment[0] != "MODE=test" ||
		seenIO.Stdout != &stdout || seenIO.Stderr != &stderr {
		t.Fatalf("injected CommandIO=%+v", seenIO)
	}
	if code := set.Run(context.Background(), []string{"--help"}, streams); code != 0 {
		t.Fatalf("help exit=%d", code)
	}
	for _, name := range []string{"build", "generate", "init", "watch"} {
		if !strings.Contains(stdout.String(), name) {
			t.Fatalf("help missing %q: %s", name, stdout.String())
		}
	}
}

func TestCommandSetRejectsDuplicateNames(t *testing.T) {
	command := generator.Command{Name: "generate", Run: func(context.Context, []string, generator.CommandIO) int { return 0 }}
	if _, err := generator.NewCommandSet(command, command); err == nil {
		t.Fatal("expected duplicate command error")
	}
}

func TestGenerateCommandHelpSucceeds(t *testing.T) {
	set := generator.MustCommandSet(generator.GenerateCommand(generator.DefaultOptions()))
	var stdout, stderr bytes.Buffer
	if code := set.Run(context.Background(), []string{"generate", "--help"}, generator.CommandIO{
		Stdout: &stdout, Stderr: &stderr,
	}); code != 0 {
		t.Fatalf("help exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "Usage of generate") {
		t.Fatalf("help output=%q", stderr.String())
	}
}

func TestGeneratePackageHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := generator.New(generator.DefaultOptions()).GeneratePackage(ctx, generator.GenerateRequest{Dir: t.TempDir()})
	if err == nil || err != context.Canceled {
		t.Fatalf("error=%v", err)
	}
}
