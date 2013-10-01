package cli

import (
	//"github.com/capitancambio/go-subcommand"
	"github.com/daisy-consortium/pipeline-clientlib-go"
	"os"
	"testing"
)

var SCRIPT pipeline.Script = pipeline.Script{
	Nicename:    "test-to-test",
	Description: "Mocked script",
	Homepage:    "daisy.org",
	Href:        "daisy.org/test",
	Id:          "test",
	Options: []pipeline.Option{
		pipeline.Option{
			Required:  true,
			Sequence:  false,
			Name:      "test-opt",
			Ordered:   false,
			Mediatype: "xml",
			Desc:      "I'm a test option",
			Type:      "anyFileURI",
		},
		pipeline.Option{
			Required:  false,
			Sequence:  false,
			Name:      "another-opt",
			Ordered:   false,
			Mediatype: "xml",
			Desc:      "I'm a test option",
			Type:      "boolean",
		},
	},
	Inputs: []pipeline.Input{
		pipeline.Input{
			Desc:      "input port",
			Mediatype: "application/x-dtbook+xml",
			Name:      "source",
			Sequence:  true,
		},
		pipeline.Input{
			Desc:      "input port not seq",
			Mediatype: "application/x-dtbook+xml",
			Name:      "single",
			Sequence:  false,
		},
	},
}

var in1, in2 = "tmp/dir1/file.xml", "tmp/dir2/file.xml"

func TestCliAddScriptCommand(t *testing.T) {
	link := PipelineLink{pipeline: newPipelineTest(false), config: &Config{Starting: false}}
	cli, err := NewCli("testprog", link)
	if err != nil {
		t.Error("Unexpected error")
	}
	cli.AddScriptCommand("test", "", func(string, ...string) error { return nil })
	if cli.Scripts[0].Name != "test" {
		t.Error("Add script is not adding scripts to the list")
	}
}

func TestCliAddCommand(t *testing.T) {
	link := PipelineLink{pipeline: newPipelineTest(false), config: &Config{Starting: false}}
	cli, err := NewCli("testprog", link)
	if err != nil {
		t.Error("Unexpected error")
	}
	cli.AddCommand("stest", "", func(string, ...string) error { return nil })
	if cli.StaticCommands[0].Name != "stest" {
		t.Error("Add Command is not adding commands to the list")
	}

	if len(cli.Scripts) != 0 {
		t.Error("Scripts is not empty")
	}
}

func TestCliNonRequiredOptions(t *testing.T) {
	link := PipelineLink{pipeline: newPipelineTest(false), config: &Config{Starting: false}}
	cli, err := NewCli("testprog", link)
	if err != nil {
		t.Error("Unexpected error")
	}
	cli.AddScripts([]pipeline.Script{SCRIPT}, link, true)
	//parser.Parse([]string{"test","--i-source","value"})
	err = cli.Run([]string{"test", "-o", os.TempDir(), "--i-source", "./tmp/file", "--i-single", "./tmp/file2", "--x-test-opt", "./myfile.xml"})
	if err != nil {
		t.Error("Non required option threw an error")
	}
}

func TestPrintHelpErrors(t *testing.T) {
	link := PipelineLink{pipeline: newPipelineTest(false), config: &Config{Starting: false}}
	cli, err := NewCli("testprog", link)
	if err != nil {
		t.Error("Unexpected error")
	}
	cli.AddScripts([]pipeline.Script{SCRIPT}, link, true)
	//more than one parameter fail
	err = printHelp(*cli, "one", "two")
	if err == nil {
		t.Error("Expected error (more than one param) is nil")
	}
	err = printHelp(*cli, "one")
	if err == nil {
		t.Error("Expected error (unknown command) is nil")
	}

}

func TestClientNew(t *testing.T) {
	link := PipelineLink{pipeline: newPipelineTest(false), config: &Config{Starting: false}}
	cli, err := NewCli("testprog", link)
	if err != nil {
		t.Error("Unexpected error")
	}
	cli.AddNewClientCommand(link)
	//Bad role
	err = cli.Run([]string{"create", "-i", "paco", "-r", "PLUMBER", "-s", "sshh"})
	if err == nil {
		t.Error("Bad role didn't err")
	}
	err = cli.Run([]string{"create", "-r", "ADMIN", "-s", "sshh"})
	if err == nil {
		t.Error("No id didn't err")
	}
	err = cli.Run([]string{"create", "-r", "ADMIN", "-i", "paco"})
	if err == nil {
		t.Error("No no secret didn't err")
	}
}

func TestClientDelete(t *testing.T) {
	link := PipelineLink{pipeline: newPipelineTest(false), config: &Config{Starting: false}}
	cli, err := NewCli("testprog", link)
	if err != nil {
		t.Error("Unexpected error")
	}
	cli.AddDeleteClientCommand(link)
	//Bad role
	err = cli.Run([]string{"delete"})
	if err == nil {
		t.Error("Bad number of args didn't err")
	}

	err = cli.Run([]string{"delete", "uno", "due"})
	if err == nil {
		t.Error("Bad number of args didn't err")
	}
}
