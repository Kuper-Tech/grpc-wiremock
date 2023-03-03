package {{ .PackageHeader }}

import (
	"google.golang.org/grpc"

	"{{ .GoPackage }}"
)

type Service struct {
	{{ .Package }}.Unimplemented{{ .Service }}Server
}

func NewService() *Service {
	return &Service{}
}

func (p *Service) RegisterGRPC(server *grpc.Server) {
	{{ .Package }}.Register{{ .Service }}Server(server, p)
}
