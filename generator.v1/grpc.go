package generator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/internal"
	"github.com/taufik-rama/maiden/internal/writer"
)

// GRPC ...
type GRPC struct {
	Output     string
	ConfigGRPC config.ServiceGRPCList

	*sync.WaitGroup
}

// GenerateCommand is the handler for cobra command-line
func (g GRPC) GenerateCommand(*cobra.Command, []string) {

	for serviceName, detail := range g.ConfigGRPC {

		internal.Print("Generating GRPC services `%s`", serviceName)

		dir := g.Output + serviceName
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalln(err)
		}

		filename := dir + string(os.PathSeparator) + serviceName + "-grpc.go"
		writer, err := writer.New(filename)
		if err != nil {
			log.Fatalln(err)
		}
		defer writer.Close()

		// Package & import
		writer.Write("// Code generated by maiden. DO NOT EDIT.")
		writer.WriteEmptyLine()
		writer.Write(`package main`)
		writer.WriteEmptyLine()
		writer.Write(`import "fmt"`)
		writer.Write(`import "net"`)
		writer.Write(`import "context"`)
		writer.Write(`import "google.golang.org/grpc"`)
		writer.Write(`import pb "%s"`, detail.Definition)
		writer.WriteEmptyLine()

		// type of the grpc server
		writer.Write(`type dummyServer struct{}`)
		writer.WriteEmptyLine()
		for name, method := range detail.Methods {

			if strings.TrimSpace(method.Request) == "" {
				log.Printf("`request` field is empty for `%s/%s` grpc service", serviceName, name)
				continue
			} else if strings.TrimSpace(method.Response) == "" {
				log.Printf("`response` field is empty for `%s/%s` grpc service", serviceName, name)
				continue
			}

			writer.Write(`func (d dummyServer) %s(ctx context.Context, req *pb.%s) (*pb.%s, error) {`, name, method.Request, method.Response)
			writer.IncrIndentLevel()
			writer.WriteEmptyLine()

			writer.Write(`if req == nil { return &pb.%s{}, nil }`, method.Response)
			writer.Write("fmt.Printf(\"Incoming request to `%s/%s`: %%s\\n\", req.String())", serviceName, name)
			writer.WriteEmptyLine()

			if detail, ok := detail.Conditions[name]; ok {

				for _, value := range detail {

					if err := writeChecks(&writer, value.Request); err != nil {
						log.Println("invalid service request conditions config on", serviceName, name, err)
						continue
					}

					if err := writeReturn(&writer, value.Response, method.Response, ""); err != nil {
						log.Println("invalid service response conditions config on", serviceName, name, err)
						continue
					}

					writer.WriteEmptyLine()
				}
			}

			writer.Write(`return &pb.%s{}, nil`, method.Response)
			writer.DecrIndentLevel()
			writer.Write(`}`)
			writer.WriteEmptyLine()
		}

		writer.Write(`func main() {`)
		writer.IncrIndentLevel()
		writer.Write("server := grpc.NewServer()")
		writer.Write("pb.Register%sServer(server, dummyServer{})", strings.ReplaceAll(strings.Title(serviceName), "-", ""))
		writer.Write(`if listener, err := net.Listen("tcp", ":%d"); err != nil {`, detail.Port)
		writer.IncrIndentLevel()
		writer.Write(`panic(err)`)
		writer.DecrIndentLevel()
		writer.Write(`} else {`)
		writer.IncrIndentLevel()
		writer.Write(`panic(server.Serve(listener))`)
		writer.DecrIndentLevel()
		writer.Write(`}`)
		writer.DecrIndentLevel()
		writer.Write("}")
	}
}

func writeChecks(writer *writer.Writer, f interface{}) error {

	if f == nil {
		return errors.New("request is nil")
	}

	fields, ok := f.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("not a JSON object, but a %s", reflect.TypeOf(f).Kind().String())
	}

	ifs := []string{}
	for field, value := range fields {
		valueType := reflect.TypeOf(value).Kind()
		if valueType == reflect.String {
			ifs = append(ifs, fmt.Sprintf(`req.%s == "%s"`, field, value))
		} else if isNumber(valueType) {

			if val, ok := value.(float64); ok {
				ifs = append(ifs, fmt.Sprintf(`req.%s == %f`, field, val))
			} else {
				ifs = append(ifs, fmt.Sprintf(`req.%s == %d`, field, value.(int)))
			}
		} else {
			log.Printf("unsupported type `%v`", valueType)
			continue
		}
	}

	writer.Write(`if %s {`, strings.Join(ifs, " && "))
	writer.IncrIndentLevel()

	return nil
}

func isNumber(kind reflect.Kind) bool {
	return kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Float32 ||
		kind == reflect.Float64
}

func writeReturn(writer *writer.Writer, f interface{}, response, field string) error {

	fields, ok := f.(map[interface{}]interface{})
	if !ok {
		return errors.New("not a JSON object")
	}

	if field == "" {
		writer.Write(`return &pb.%s{`, response)
	} else if field == "INARRAY" {
		writer.Write(`&pb.%s{`, response)
	} else {
		writer.Write(`%s: &pb.%s{`, field, response)
	}

	writer.IncrIndentLevel()
	for field, value := range fields {

		if _, ok := field.(string); !ok {
			return fmt.Errorf("response key `%v` is not of string type", field)
		}

		valueType := reflect.TypeOf(value).Kind()
		if valueType == reflect.String {
			writer.Write(`%s: "%s",`, field, value.(string))
		} else if isNumber(valueType) {
			if val, ok := value.(float64); ok {
				writer.Write(`%s: %f,`, field, val)
			} else {
				writer.Write(`%s: %d,`, field, value.(int))
			}
		} else if valueType == reflect.Map {

			if err := writeReturnsFields(writer, field.(string), value); err != nil {
				return err
			}

		} else {
			log.Printf("unsupported type `%v`", valueType)
			continue
		}
	}

	// For the struct
	writer.DecrIndentLevel()
	if field == "" {
		writer.Write(`}, nil`)
	} else {
		writer.Write(`},`)
	}

	// For the `if` block
	if field == "" {
		writer.DecrIndentLevel()
		writer.Write(`}`)
	}

	return nil
}

func writeReturnsFields(writer *writer.Writer, field string, f interface{}) error {

	fields, ok := f.(map[interface{}]interface{})
	if !ok {
		return errors.New("not a JSON object")
	}

	typ, ok := fields["type"]
	if !ok {
		return errors.New("`type` is not defined")
	}

	typType := reflect.TypeOf(typ).Kind()
	if typType != reflect.String && typType != reflect.Slice {
		return errors.New("`type` can only be of type string or array")
	}

	val, ok := fields["values"]
	if !ok {
		return errors.New("`values` is not defined")
	}

	valType := reflect.TypeOf(val).Kind()
	if typType == reflect.String && valType != reflect.Map {
		return errors.New("`values` must be a JSON object if the `type` is not a slice")
	} else if typType == reflect.Slice && valType != reflect.Slice {
		return errors.New("`values` must be a JSON Array if the `type` is a slice")
	}

	if typType == reflect.String {
		typ := typ.(string)
		writeReturn(writer, val, typ, field)
	} else if typType == reflect.Slice {
		typ := typ.([]interface{})
		if len(typ) < 1 {
			return errors.New("`type` field must have a string value")
		} else if _, ok := typ[0].(string); !ok {
			return errors.New("`type` field must have a string value")
		}
		writer.Write("%s: []*pb.%s{", field, typ[0])
		writer.IncrIndentLevel()
		for _, t := range val.([]interface{}) {
			writeReturn(writer, t, typ[0].(string), "INARRAY")
		}
		writer.DecrIndentLevel()
		writer.Write("},")
	}

	return nil
}

// RunCommand is the handler for cobra command-line
func (g GRPC) RunCommand(*cobra.Command, []string) {

	for name := range g.ConfigGRPC {
		internal.Print("Building `%s` service", name)
		if err := g.buildGRPC(g.Output, name); err != nil {
			log.Fatalln(err)
		}
	}

	for name := range g.ConfigGRPC {
		g.WaitGroup.Add(1)
		internal.Print("Starting `%s` service", name)
		go g.startGRPC(g.Output, name, g.WaitGroup)
	}
}

func (g GRPC) buildGRPC(output, name string) error {

	dir := output + name

	{
		command := exec.Command("go", "get", ("." + string(os.PathSeparator) + dir + string(os.PathSeparator) + "..."))
		stderr := strings.Builder{}
		command.Stderr = &stderr
		if err := command.Start(); err != nil {
			return fmt.Errorf("error while fetching deps service `%s`, stderr: %s", name, stderr.String())
		}

		if err := command.Wait(); err != nil {
			return fmt.Errorf("error while fetching deps service `%s`, stderr: %s", name, stderr.String())
		}
	}

	{
		// I'm not using `go run` because that counts as 2 separate process
		command := exec.Command("go", "build", "-o", (dir + string(os.PathSeparator) + name), (dir + string(os.PathSeparator) + name + "-grpc.go"))
		stderr := strings.Builder{}
		command.Stderr = &stderr
		if err := command.Start(); err != nil {
			return fmt.Errorf("error while building service `%s`, stderr: %s", name, stderr.String())
		}

		if err := command.Wait(); err != nil {
			return fmt.Errorf("error while building service `%s`, stderr: %s", name, stderr.String())
		}
	}

	return nil
}

func (g GRPC) startGRPC(output, name string, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()

	command := exec.Command("./" + (output + name + string(os.PathSeparator) + name))
	stderr := strings.Builder{}
	command.Stderr = &stderr
	command.Stdout = os.Stdout
	if err := command.Start(); err != nil {
		log.Printf("error while starting service `%s`, stderr: %s", name, stderr.String())
		return
	}

	if err := command.Wait(); err != nil {
		log.Println("wait:", err, stderr.String())
		return
	}

	if !command.ProcessState.Success() {
		log.Printf("error while listening `%s`, stderr: %s", name, stderr.String())
	}
}
