package ffmpeg

import (
	"fmt"
	"strings"
)

// Command represents an ffmpeg command with its arguments
type Command struct {
	Inputs  []Input
	Filters []Filter
	Outputs []Output
}

// Input represents an input source for ffmpeg
type Input struct {
	Path    string
	Options []string
}

// Filter represents a filter or filter chain
type Filter struct {
	Name   string
	Params map[string]string
	Inputs []string
	Output string
}

// Output represents an output destination
type Output struct {
	Path    string
	Options []string
}

// String generates the complete ffmpeg command as a string
func (c *Command) String() string {
	var parts []string
	
	// Add inputs
	for _, input := range c.Inputs {
		if len(input.Options) > 0 {
			parts = append(parts, strings.Join(input.Options, " "))
		}
		parts = append(parts, "-i", input.Path)
	}
	
	// Add filters if any
	if len(c.Filters) > 0 {
		filterComplex := c.buildFilterComplex()
		if filterComplex != "" {
			parts = append(parts, "-filter_complex", fmt.Sprintf("'%s'", filterComplex))
		}
	}
	
	// Add outputs
	for _, output := range c.Outputs {
		if len(output.Options) > 0 {
			parts = append(parts, strings.Join(output.Options, " "))
		}
		parts = append(parts, output.Path)
	}
	
	return "ffmpeg " + strings.Join(parts, " ")
}

// buildFilterComplex creates the filter_complex string from filters
func (c *Command) buildFilterComplex() string {
	if len(c.Filters) == 0 {
		return ""
	}
	
	var filterStrings []string
	for _, filter := range c.Filters {
		filterStr := c.buildSingleFilter(filter)
		if filterStr != "" {
			filterStrings = append(filterStrings, filterStr)
		}
	}
	
	return strings.Join(filterStrings, "; ")
}

// buildSingleFilter creates a single filter string
func (c *Command) buildSingleFilter(filter Filter) string {
	var parts []string
	
	// Add inputs
	if len(filter.Inputs) > 0 {
		parts = append(parts, strings.Join(filter.Inputs, ""))
	}
	
	// Add filter name
	parts = append(parts, filter.Name)
	
	// Add parameters
	if len(filter.Params) > 0 {
		var params []string
		for k, v := range filter.Params {
			if k == "" {
				// Special case: empty key means the value is the entire parameter string
				params = append(params, v)
			} else if v == "" {
				params = append(params, k)
			} else {
				params = append(params, fmt.Sprintf("%s=%s", k, v))
			}
		}
		if len(params) > 0 {
			parts = append(parts, fmt.Sprintf("=%s", strings.Join(params, ":")))
		}
	}
	
	// Add output label
	if filter.Output != "" {
		parts = append(parts, filter.Output)
	}
	
	return strings.Join(parts, "")
}

// AddInput adds an input source to the command
func (c *Command) AddInput(path string, options ...string) {
	c.Inputs = append(c.Inputs, Input{
		Path:    path,
		Options: options,
	})
}

// AddFilter adds a filter to the command
func (c *Command) AddFilter(name string, inputs []string, output string, params map[string]string) {
	c.Filters = append(c.Filters, Filter{
		Name:   name,
		Params: params,
		Inputs: inputs,
		Output: output,
	})
}

// AddOutput adds an output destination to the command
func (c *Command) AddOutput(path string, options ...string) {
	c.Outputs = append(c.Outputs, Output{
		Path:    path,
		Options: options,
	})
}