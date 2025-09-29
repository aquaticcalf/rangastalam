package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// Executor handles the execution of ffmpeg commands
type Executor struct {
	FFmpegPath string
	DryRun     bool
	Verbose    bool
	Timeout    time.Duration
}

// NewExecutor creates a new ffmpeg executor
func NewExecutor() *Executor {
	return &Executor{
		FFmpegPath: "ffmpeg", // assumes ffmpeg is in PATH
		DryRun:     false,
		Verbose:    false,
		Timeout:    10 * time.Minute, // default 10 minute timeout
	}
}

// Execute runs the ffmpeg command
func (e *Executor) Execute(cmd *Command) error {
	return e.ExecuteWithContext(context.Background(), cmd)
}

// ExecuteWithContext runs the ffmpeg command with a context for cancellation
func (e *Executor) ExecuteWithContext(ctx context.Context, cmd *Command) error {
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}
	
	// Build the command arguments
	args := e.buildArgs(cmd)
	
	if e.Verbose || e.DryRun {
		fmt.Printf("Executing: %s %s\n", e.FFmpegPath, strings.Join(args, " "))
	}
	
	if e.DryRun {
		return nil // Don't actually execute in dry run mode
	}
	
	// Create context with timeout
	execCtx := ctx
	if e.Timeout > 0 {
		var cancel context.CancelFunc
		execCtx, cancel = context.WithTimeout(ctx, e.Timeout)
		defer cancel()
	}
	
	// Create the command
	execCmd := exec.CommandContext(execCtx, e.FFmpegPath, args...)
	
	// Set up pipes for stdout and stderr
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	
	// Start the command
	if err := execCmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}
	
	// Handle output
	go e.handleOutput("stdout", stdout)
	go e.handleOutput("stderr", stderr)
	
	// Wait for completion
	if err := execCmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg execution failed: %w", err)
	}
	
	return nil
}

// buildArgs converts the Command to command line arguments
func (e *Executor) buildArgs(cmd *Command) []string {
	var args []string
	
	// Always add -y to overwrite output files without asking
	args = append(args, "-y")
	
	// Add inputs
	for _, input := range cmd.Inputs {
		if len(input.Options) > 0 {
			args = append(args, input.Options...)
		}
		args = append(args, "-i", input.Path)
	}
	
	// Add filters if any
	if len(cmd.Filters) > 0 {
		filterComplex := cmd.buildFilterComplex()
		if filterComplex != "" {
			args = append(args, "-filter_complex", filterComplex)
		}
	}
	
	// Add outputs
	for _, output := range cmd.Outputs {
		if len(output.Options) > 0 {
			args = append(args, output.Options...)
		}
		args = append(args, output.Path)
	}
	
	return args
}

// handleOutput reads and optionally prints command output
func (e *Executor) handleOutput(name string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if e.Verbose {
			fmt.Printf("[%s] %s\n", name, line)
		}
	}
}

// CheckFFmpeg verifies that ffmpeg is available and working
func (e *Executor) CheckFFmpeg() error {
	cmd := exec.Command(e.FFmpegPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("ffmpeg not found or not working: %w", err)
	}
	
	if e.Verbose {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 0 {
			fmt.Printf("Found: %s\n", lines[0])
		}
	}
	
	return nil
}

// SetDryRun enables or disables dry run mode
func (e *Executor) SetDryRun(dryRun bool) {
	e.DryRun = dryRun
}

// SetVerbose enables or disables verbose output
func (e *Executor) SetVerbose(verbose bool) {
	e.Verbose = verbose
}

// SetTimeout sets the execution timeout
func (e *Executor) SetTimeout(timeout time.Duration) {
	e.Timeout = timeout
}

// SetFFmpegPath sets the path to the ffmpeg executable
func (e *Executor) SetFFmpegPath(path string) {
	e.FFmpegPath = path
}