package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/compiler/compilecontract"
)

type protoCompilerPlugin string

const (
	goPlugin        = "go"
	goGRPCPlugin    = "go-grpc"
	goOpenAPIPlugin = "openapi"
)

type pluginValues struct {
	out string
}

type flag struct {
	key   string
	value string
}

func createProtoPath(value string) flag {
	const protoPathFlag = "proto_path"
	return flag{
		key:   protoPathFlag,
		value: value,
	}
}

func createPluginOut(plugin protoCompilerPlugin, value string) flag {
	const tplFlag = "%s_out"
	return flag{
		key:   fmt.Sprintf(tplFlag, plugin),
		value: value,
	}
}

func (a flag) String() string {
	return fmt.Sprintf("--%s=%s", a.key, a.value)
}

type Command struct {
	flags     []flag
	arguments []string
}

func BuildArgs(
	contract compilecontract.Contract,
	extraProtoPaths []string,
	pluginOutputs map[protoCompilerPlugin]pluginValues,
) Command {
	var flags []flag

	protoPaths := append(extraProtoPaths, getProtoPaths(contract)...)

	for _, path := range protoPaths {
		flags = append(flags, createProtoPath(path))
	}

	for plugin, values := range pluginOutputs {
		flags = append(flags, createPluginOut(plugin, values.out))
	}

	arguments := contractToArguments(contract)

	return Command{flags: flags, arguments: arguments}
}

func contractToArguments(contract compilecontract.Contract) []string {
	arguments := []string{contract.HeaderPath}
	arguments = append(arguments, contract.ImportsPaths...)

	return arguments
}

func getProtoPaths(contract compilecontract.Contract) (protoPaths []string) {
	headerDir := filepath.Dir(contract.HeaderPath)

	uniq := map[string]struct{}{
		headerDir: {},
	}

	protoPaths = append(protoPaths, headerDir)

	for _, path := range contract.ImportsPaths {
		dir := filepath.Dir(path)

		if _, exists := uniq[dir]; exists {
			continue
		}

		uniq[dir] = struct{}{}
		protoPaths = append(protoPaths, dir)
	}

	return protoPaths
}

func (c *Command) Args() []string {
	return c.toSlice()
}

func (c *Command) String() string {
	return fmt.Sprintf(
		"\n\t%s",
		strings.Join(c.toSlice(), " \\\n\t"),
	)
}

func (c *Command) toSlice() []string {
	var args []string

	for _, fl := range c.flags {
		args = append(args, fl.String())
	}
	args = append(args, c.arguments...)

	return args
}

type CompileCommand map[protoCompilerPlugin]pluginValues

func CompileGoPackages(path string) CompileCommand {
	return map[protoCompilerPlugin]pluginValues{
		goPlugin:     {out: path},
		goGRPCPlugin: {out: path},
	}
}

func CompileOpenAPI(path string) CompileCommand {
	return map[protoCompilerPlugin]pluginValues{
		goOpenAPIPlugin: {out: path},
	}
}
