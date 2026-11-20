package dockerfile

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name       string
		dockerfile string
		want       []Arg
	}{
		{
			name:       "arg with a default",
			dockerfile: "FROM scratch\nARG EXAMPLE_VERSION=1.2.3\n",
			want:       []Arg{{Name: "EXAMPLE_VERSION", Default: "1.2.3", HasDefault: true, Stage: "0", Line: 2}},
		},
		{
			name:       "arg without a default",
			dockerfile: "FROM scratch\nARG EXAMPLE_TOKEN\n",
			want:       []Arg{{Name: "EXAMPLE_TOKEN", Stage: "0", Line: 2}},
		},
		{
			name:       "several names on one line",
			dockerfile: "FROM scratch\nARG TARGETOS TARGETARCH\n",
			want: []Arg{
				{Name: "TARGETOS", Stage: "0", Line: 2},
				{Name: "TARGETARCH", Stage: "0", Line: 2},
			},
		},
		{
			name:       "quotes around the default are removed",
			dockerfile: "FROM scratch\nARG EXAMPLE_VERSION=\"1.2.3\"\n",
			want:       []Arg{{Name: "EXAMPLE_VERSION", Default: "1.2.3", HasDefault: true, Stage: "0", Line: 2}},
		},
		{
			name:       "global args and named stages",
			dockerfile: "ARG BASE_TAG=1.0\nFROM example.com/base:${BASE_TAG} AS build\nARG BUILD_MODE=release\n",
			want: []Arg{
				{Name: "BASE_TAG", Default: "1.0", HasDefault: true, Stage: "", Line: 1},
				{Name: "BUILD_MODE", Default: "release", HasDefault: true, Stage: "build", Line: 3},
			},
		},
		{
			name:       "comments and line continuations",
			dockerfile: "FROM scratch\n# ARG IGNORED=1\nARG FIRST=1 \\\n    SECOND=2\n",
			want: []Arg{
				{Name: "FIRST", Default: "1", HasDefault: true, Stage: "0", Line: 3},
				{Name: "SECOND", Default: "2", HasDefault: true, Stage: "0", Line: 3},
			},
		},
		{
			name:       "a default referencing another arg stays as written",
			dockerfile: "FROM scratch\nARG EXAMPLE_VERSION=1.2.3\nARG EXAMPLE_URL=https://example.com/v${EXAMPLE_VERSION}\n",
			want: []Arg{
				{Name: "EXAMPLE_VERSION", Default: "1.2.3", HasDefault: true, Stage: "0", Line: 2},
				{Name: "EXAMPLE_URL", Default: "https://example.com/v${EXAMPLE_VERSION}", HasDefault: true, Stage: "0", Line: 3},
			},
		},
		{
			name:       "declaration order across stages",
			dockerfile: "FROM scratch\nARG ZULU=1\nFROM scratch\nARG ALPHA=2\n",
			want: []Arg{
				{Name: "ZULU", Default: "1", HasDefault: true, Stage: "0", Line: 2},
				{Name: "ALPHA", Default: "2", HasDefault: true, Stage: "1", Line: 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(strings.NewReader(tt.dockerfile))
			if err != nil {
				t.Fatalf("ParseArgs() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArgs() =\n%+v\nwant\n%+v", got, tt.want)
			}
		})
	}
}

func TestParseArgsRejectsAnInvalidDockerfile(t *testing.T) {
	if _, err := ParseArgs(strings.NewReader("ARG\n")); err == nil {
		t.Fatal("ParseArgs() error = nil, want an error for ARG without a name")
	}
}
