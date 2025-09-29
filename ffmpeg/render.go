package ffmpeg

import (
	"context"
	"fmt"

	"github.com/aquaticcalf/rangastalam/ast"
)

// Renderer provides a high-level interface for rendering projects
type Renderer struct {
	executor   *Executor
	translator *Translator
}

// NewRenderer creates a new renderer
func NewRenderer() *Renderer {
	return &Renderer{
		executor: NewExecutor(),
	}
}

// RenderOptions provides configuration for rendering
type RenderOptions struct {
	OutputPath string
	DryRun     bool
	Verbose    bool
	Quality    string // "high", "medium", "low"
	Format     string // "mp4", "avi", "mov", etc.
}

// DefaultRenderOptions returns sensible defaults
func DefaultRenderOptions() *RenderOptions {
	return &RenderOptions{
		OutputPath: "output.mp4",
		DryRun:     false,
		Verbose:    false,
		Quality:    "medium",
		Format:     "mp4",
	}
}

// Render converts a project to video using ffmpeg
func (r *Renderer) Render(project *ast.Project, options *RenderOptions) error {
	return r.RenderWithContext(context.Background(), project, options)
}

// RenderWithContext renders with cancellation support
func (r *Renderer) RenderWithContext(ctx context.Context, project *ast.Project, options *RenderOptions) error {
	if project == nil {
		return fmt.Errorf("project cannot be nil")
	}
	
	if options == nil {
		options = DefaultRenderOptions()
	}
	
	// Set up executor options
	r.executor.SetDryRun(options.DryRun)
	r.executor.SetVerbose(options.Verbose)
	
	// Check ffmpeg availability
	if err := r.executor.CheckFFmpeg(); err != nil {
		return fmt.Errorf("ffmpeg check failed: %w", err)
	}
	
	// Create translator
	r.translator = NewTranslator(project)
	
	// Translate project to ffmpeg command
	cmd, err := r.translator.Translate()
	if err != nil {
		return fmt.Errorf("translation failed: %w", err)
	}
	
	// Update output path and codec options based on settings
	r.configureOutput(cmd, options)
	
	// Execute the command
	if err := r.executor.ExecuteWithContext(ctx, cmd); err != nil {
		return fmt.Errorf("rendering failed: %w", err)
	}
	
	return nil
}

// configureOutput sets up output parameters based on render options
func (r *Renderer) configureOutput(cmd *Command, options *RenderOptions) {
	// Clear existing outputs
	cmd.Outputs = nil
	
	// Build codec options based on quality
	var codecOpts []string
	
	switch options.Quality {
	case "high":
		codecOpts = []string{"-c:v", "libx264", "-preset", "slow", "-crf", "18", "-c:a", "aac", "-b:a", "320k"}
	case "low":
		codecOpts = []string{"-c:v", "libx264", "-preset", "ultrafast", "-crf", "28", "-c:a", "aac", "-b:a", "128k"}
	default: // medium
		codecOpts = []string{"-c:v", "libx264", "-preset", "medium", "-crf", "23", "-c:a", "aac", "-b:a", "192k"}
	}
	
	// Add format-specific options
	switch options.Format {
	case "avi":
		codecOpts = append(codecOpts, "-f", "avi")
	case "mov":
		codecOpts = append(codecOpts, "-f", "mov")
	default: // mp4
		codecOpts = append(codecOpts, "-f", "mp4")
	}
	
	cmd.AddOutput(options.OutputPath, codecOpts...)
}

// GetCommand returns the ffmpeg command that would be executed (useful for debugging)
func (r *Renderer) GetCommand(project *ast.Project) (*Command, error) {
	if r.translator == nil {
		r.translator = NewTranslator(project)
	}
	
	return r.translator.Translate()
}

// GetCommandString returns the ffmpeg command as a string
func (r *Renderer) GetCommandString(project *ast.Project, options *RenderOptions) (string, error) {
	cmd, err := r.GetCommand(project)
	if err != nil {
		return "", err
	}
	
	if options == nil {
		options = DefaultRenderOptions()
	}
	
	r.configureOutput(cmd, options)
	
	return cmd.String(), nil
}