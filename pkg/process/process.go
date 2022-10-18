package process

import (
	"os"
	"os/exec"
	"strings"
	"sync"
)

func New(workDir, commandStr string) *Process {
	commands := ParseCommandLine(commandStr)
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Dir = workDir
	if os.Getenv("WISP_TRACE_PROCESS") == "1" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}
	return &Process{
		WorkDir:     workDir,
		CommandLine: commandStr,
		programName: commands[0],
		args:        commands[1:],
		cmd:         cmd,
		done:        make(chan struct{}),
	}
}

type Process struct {
	WorkDir     string
	CommandLine string

	programName string
	args        []string
	cmd         *exec.Cmd

	firstErr error
	done     chan struct{}
	once     sync.Once
}

func (p *Process) Start() error {
	err := p.cmd.Run()
	p.once.Do(func() {
		close(p.done)
	})
	return err
}

func (p *Process) Stop() {
	select {
	case <-p.done:
	default:
		if p.cmd != nil && p.cmd.Process != nil {
			_ = p.cmd.Process.Signal(os.Kill)
			<-p.done
		}
	}
}

func ParseCommandLine(s string) []string {
	commands := make([]string, 0, 2)
	for _, value := range strings.Split(s, " ") {
		if newValue := strings.TrimSpace(value); newValue != "" {
			commands = append(commands, newValue)
		}
	}
	return commands
}
