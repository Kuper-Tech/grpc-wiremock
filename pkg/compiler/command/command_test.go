package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/compiler/compilecontract"
)

func Test_buildArgs(t *testing.T) {
	tests := []struct {
		name string

		extraProtoPaths []string
		contract        compilecontract.Contract
		pluginValues    map[protoCompilerPlugin]pluginValues

		want Command
	}{
		{
			pluginValues: CompileGoPackages("/here/go"),
			contract: compilecontract.Contract{
				HeaderPath:   "/test/inner/dir/grpc/sample.proto",
				ImportsPaths: []string{"/test/inner/dir/grpc/import/sample.proto"},
			},
			want: Command{
				flags: []flag{
					createProtoPath("/test/inner/dir/grpc"),
					createProtoPath("/test/inner/dir/grpc/import"),
					createPluginOut(goPlugin, "/here/go"),
					createPluginOut(goGRPCPlugin, "/here/go"),
				},
				arguments: []string{
					"/test/inner/dir/grpc/sample.proto",
					"/test/inner/dir/grpc/import/sample.proto",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildArgs(tt.contract, tt.extraProtoPaths, tt.pluginValues)
			assert.ElementsMatch(t, tt.want.arguments, got.arguments)
			assert.ElementsMatch(t, tt.want.flags, got.flags)
		})
	}
}

func Test_command_Args(t *testing.T) {
	tests := []struct {
		name string

		flags     []flag
		arguments []string

		want []string
	}{
		{
			flags: []flag{
				createPluginOut(goPlugin, "/here/go"),
				createPluginOut(goGRPCPlugin, "/here/grpc"),
				createProtoPath("/test/inner/dir/grpc"),
				createProtoPath("/test/inner/dir/grpc/import"),
			},
			arguments: []string{
				"/test/inner/dir/grpc/sample.proto",
				"/test/inner/dir/grpc/import/sample.proto",
			},
			want: []string{
				"--go_out=/here/go",
				"--go-grpc_out=/here/grpc",
				"--proto_path=/test/inner/dir/grpc",
				"--proto_path=/test/inner/dir/grpc/import",
				"/test/inner/dir/grpc/sample.proto",
				"/test/inner/dir/grpc/import/sample.proto",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{flags: tt.flags, arguments: tt.arguments}
			got := c.Args()
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func Test_command_String(t *testing.T) {
	tests := []struct {
		name string

		flags     []flag
		arguments []string

		want string
	}{
		{
			flags: []flag{
				createPluginOut(goPlugin, "/here/go"),
				createPluginOut(goGRPCPlugin, "/here/grpc"),
				createProtoPath("/test/inner/dir/grpc"),
				createProtoPath("/test/inner/dir/grpc/import"),
			},
			arguments: []string{
				"/test/inner/dir/grpc/sample.proto",
				"/test/inner/dir/grpc/import/sample.proto",
			},
			want: "\n\t--go_out=/here/go \\\n\t--go-grpc_out=/here/grpc \\\n\t--proto_path=/test/inner/dir/grpc \\\n\t--proto_path=/test/inner/dir/grpc/import \\\n\t/test/inner/dir/grpc/sample.proto \\\n\t/test/inner/dir/grpc/import/sample.proto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{flags: tt.flags, arguments: tt.arguments}
			got := c.String()
			assert.Equal(t, tt.want, got)
		})
	}
}
