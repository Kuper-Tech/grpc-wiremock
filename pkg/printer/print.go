package printer

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
)

func Print(descriptors []*desc.FileDescriptor, outputPath string) error {
	printer := protoprint.Printer{}

	if err := printer.PrintProtosToFileSystem(descriptors, outputPath); err != nil {
		return fmt.Errorf("print: %w", err)
	}

	return nil
}
