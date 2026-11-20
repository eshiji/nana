// Package dockerfile reads facts from Dockerfiles using Docker's own parser.
package dockerfile

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Arg is one ARG declaration. Stage is "" for global args (before the first FROM),
// otherwise the stage alias, or its index when the stage has no alias, as Docker
// refers to stages in COPY --from.
type Arg struct {
	Name       string `json:"name"`
	Default    string `json:"default"`
	HasDefault bool   `json:"has_default"`
	Stage      string `json:"stage"`
	Line       int    `json:"line"`
}

var errArgWithoutName = errors.New("ARG requires at least one name")

// ParseArgs returns every ARG in declaration order. Defaults are returned as written:
// surrounding quotes are removed, references such as ${OTHER} are not expanded.
func ParseArgs(r io.Reader) ([]Arg, error) {
	result, err := parser.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parse Dockerfile: %w", err)
	}
	args := []Arg{}
	stage, stageIndex := "", 0
	for _, node := range result.AST.Children {
		switch strings.ToUpper(node.Value) {
		case "FROM":
			stage = stageName(node, stageIndex)
			stageIndex++
		case "ARG":
			declared, err := argsFrom(node, stage)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", node.StartLine, err)
			}
			args = append(args, declared...)
		}
	}
	return args, nil
}

func stageName(from *parser.Node, index int) string {
	words := []string{}
	for word := from.Next; word != nil; word = word.Next {
		words = append(words, word.Value)
	}
	if len(words) == 3 && strings.EqualFold(words[1], "AS") {
		return words[2]
	}
	return strconv.Itoa(index)
}

func argsFrom(node *parser.Node, stage string) ([]Arg, error) {
	if node.Next == nil {
		return nil, errArgWithoutName
	}
	args := []Arg{}
	for word := node.Next; word != nil; word = word.Next {
		name, value, hasDefault := strings.Cut(word.Value, "=")
		args = append(args, Arg{
			Name:       name,
			Default:    unquote(value),
			HasDefault: hasDefault,
			Stage:      stage,
			Line:       node.StartLine,
		})
	}
	return args, nil
}

func unquote(value string) string {
	for _, quote := range []string{`"`, `'`} {
		if len(value) >= 2 && strings.HasPrefix(value, quote) && strings.HasSuffix(value, quote) {
			return value[1 : len(value)-1]
		}
	}
	return value
}
