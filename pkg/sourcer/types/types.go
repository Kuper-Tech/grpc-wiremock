package types

import (
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

// SourceLoadType represents how exactly contracts must be loaded.
type SourceLoadType int

const (
	// UnknownLoadType represents default value of SourceLoadType.
	UnknownLoadType SourceLoadType = iota

	// SourceDirType represents that provided input is directory with the contracts.
	SourceDirType SourceLoadType = iota

	// SourceSingleType represents that provided input is single contract file.
	SourceSingleType
)

type SourceFileType string

const (
	// UnknownFileType represents unknown contracts type.
	UnknownFileType = "unknown"

	// ProtoType represents contracts described in .proto files.
	ProtoType SourceFileType = "proto"

	// OpenAPIType represents contracts described in .yaml files.
	OpenAPIType SourceFileType = "yaml"
)

func (s SourceFileType) Is(that string) bool {
	return string(s) == that
}

func SourceFileTypeFromPath(path string) SourceFileType {
	extension := fsutils.GetFileExt(path)

	switch {
	case ProtoType.Is(extension):
		return ProtoType
	case OpenAPIType.Is(extension):
		return OpenAPIType
	}

	return ""
}

type ProtoMethodType int

const (
	UnaryType ProtoMethodType = iota
	ClientSideStreamingType
	ServerSideStreamingType
	BidirectionalStreamingType
)

func MethodType(isClientSideStream, isServerSideStream bool) ProtoMethodType {
	if isServerSideStream {
		if isClientSideStream {
			return BidirectionalStreamingType
		}
		return ServerSideStreamingType
	}

	if isClientSideStream {
		return ClientSideStreamingType
	}

	return UnaryType
}

func (t ProtoMethodType) String() string {
	switch t {
	case UnaryType:
		return "unary"
	case ClientSideStreamingType:
		return "client side streaming"
	case ServerSideStreamingType:
		return "server side streaming"
	case BidirectionalStreamingType:
		return "bidirectional streaming"
	}

	return ""
}
