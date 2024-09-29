package generator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/EvgenyOvsov/protoc-gateway-impl/pkg"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type Generator interface {
	Handle(generator *protogen.Plugin) error
	SetParamsFunc(name, value string) error
}

type Option struct {
	module string
}

type ImplGenerator struct {
	wg sync.WaitGroup
	o  Option
}

type serviceGenerator struct {
	interfaceName      string
	implementationName string

	prefix    string
	options   *Option
	protoGen  *protogen.Plugin
	service   *protogen.Service
	protoFile *protogen.File
}

func (g *ImplGenerator) SetParamsFunc(name, value string) error {
	switch strings.ToLower(name) {
	case "out":
		g.o.module = value
		return nil
	default:
		return errors.New("unknown param " + name)
	}
}

func (g *ImplGenerator) Validate() error {
	if g.o.module == "" {
		return errors.New("out param is required")
	}
	return nil
}

func NewGenerator() Generator {
	return &ImplGenerator{
		wg: sync.WaitGroup{},
	}
}

func (g *ImplGenerator) Handle(generator *protogen.Plugin) error {
	if err := g.Validate(); err != nil {
		return err
	}

	for _, f := range generator.Files {
		if f.Generate {
			g.wg.Add(1)
			go func(f *protogen.File) {
				defer g.wg.Done()
				if err := g.generate(f, generator); err != nil {
					log.Println("ERROR", err.Error())
				}
			}(f)
		}
	}
	g.wg.Wait()
	return nil
}

func (g *ImplGenerator) generate(f *protogen.File, generator *protogen.Plugin) error {
	for _, service := range f.Services {
		serviceGenerator := &serviceGenerator{
			interfaceName:      service.GoName,
			implementationName: string(unicode.ToLower(rune(service.GoName[0]))) + service.GoName[1:],
			prefix:             g.o.module + "/" + pkg.ToSnakeCase(service.GoName),
			options:            &g.o,
			protoGen:           generator,
			service:            service,
			protoFile:          f,
		}
		if err := serviceGenerator.generateService(); err != nil {
			log.Println("ERROR", err.Error())
		}
	}
	return nil
}

func (g *serviceGenerator) generateService() error {
	if err := g.generateInterface(); err != nil {
		return err
	}
	if err := g.generateImplementation(g.service, g.protoGen, g.prefix, g.protoFile); err != nil {
		return err
	}
	for _, method := range g.service.Methods {
		if err := g.generateMethod(method, g.protoGen, g.prefix, g.protoFile); err != nil {
			return err
		}
	}
	return nil
}

func (g *serviceGenerator) generateInterface() error {
	file := g.protoGen.NewGeneratedFile(g.prefix+"/service.auto.go", g.protoFile.GoImportPath)

	file.P(pkg.HeaderCantModify)
	file.P(`package `, pkg.ToSnakeCase(g.service.GoName))
	file.P(`import (`)
	file.P(`	"context"`)
	file.P(`	"net/http"`)
	file.P(`	gw `, g.protoFile.GoImportPath)
	file.P(`	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"`)
	file.P(`)`)
	file.P(`type `, g.service.GoName, ` interface {`)
	file.P(`	Apply(ctx context.Context, httpMux *http.ServeMux) error`)
	for _, method := range g.service.Methods {
		file.P(`	`, method.GoName, `(ctx context.Context, req *gw.`, method.Input.GoIdent, `) (*gw.`, method.Output.GoIdent, `, error)`)
	}
	file.P(`}`)
	file.P(`func (s *`, g.implementationName, `) Apply(ctx context.Context, httpMux *http.ServeMux) error {`)
	file.P(`	grpcMux := runtime.NewServeMux()`)
	file.P(`	if err := gw.Register`, g.service.GoName, `HandlerServer(ctx, grpcMux, s); err != nil {`)
	file.P(`		return err`)
	file.P(`	}`)

	for _, method := range g.service.Methods {
		opts := method.Desc.Options().(*descriptor.MethodOptions)
		if proto.HasExtension(opts, annotations.E_Http) {
			ext := proto.GetExtension(opts, annotations.E_Http).(*annotations.HttpRule)
			url := pkg.FindGoogleAnnotationURL(ext)
			if url == "" {
				return fmt.Errorf("no url for RPC %s", method.GoName)
			}
			file.P(`	httpMux.Handle("`, url, `", grpcMux)`)
		}
	}

	file.P(`	return nil`)
	file.P(`}`)
	return nil
}

func (g *serviceGenerator) generateImplementation(srv *protogen.Service, generator *protogen.Plugin, prefix string, f *protogen.File) error {
	filepath, err := filepath.Abs(prefix + "/service.go")
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file := generator.NewGeneratedFile(prefix+"/service.go", f.GoImportPath)

		file.P(pkg.HeaderCanModify)
		file.P(`package `, pkg.ToSnakeCase(srv.GoName))
		file.P(`import (`)
		file.P(`	gw `, g.protoFile.GoImportPath)
		file.P(`)`)
		file.P(`type `, g.implementationName, ` struct {`)
		file.P(`	gw.UnimplementedBasicServiceServer`)
		file.P(`}`)
		file.P(`func New`, srv.GoName, `() `, srv.GoName, ` {`)
		file.P(`	return &`, string(unicode.ToLower(rune(srv.GoName[0]))), srv.GoName[1:], `{}`)
		file.P(`}`)
	}
	return nil
}

func (g *serviceGenerator) generateMethod(m *protogen.Method, generator *protogen.Plugin, prefix string, f *protogen.File) error {
	methodName := pkg.ToSnakeCase(m.GoName)
	packageName := pkg.ToSnakeCase(m.Parent.GoName)
	filepath, err := filepath.Abs(prefix + "/" + methodName + ".go")
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file := generator.NewGeneratedFile(prefix+"/"+methodName+".go", f.GoImportPath)

		file.P(pkg.HeaderCanModify)
		file.P(`package `, packageName)
		file.P(`import (`)
		file.P(`	"context"`)
		file.P(`	gw `, g.protoFile.GoImportPath, ``)
		file.P(`	codes "google.golang.org/grpc/codes"`)
		file.P(`	status "google.golang.org/grpc/status"`)
		file.P(`)`)
		file.P(fmt.Sprintf(`func (s *%s) %s(ctx context.Context, req *gw.%s) (*gw.%s, error) {`,
			g.implementationName, m.GoName, m.Input.GoIdent.GoName, m.Output.GoIdent.GoName))
		file.P(`	return nil, status.Errorf(codes.Unimplemented, "method not implemented")`)
		file.P(`}`)
	}
	return nil
}
