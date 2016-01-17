// ThinkGo 部署工具，支持新建项目，支持热编译并运行程序。
package deploy

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
)

const version = "1.3.0"

type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) int

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short template.HTML

	// Long is the long message shown in the 'go help <this-command>' output.
	Long template.HTML

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(string(c.Long)))
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var commands = []*Command{
	cmdNew,
	cmdRun,
}

func Deploy() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			os.Exit(cmd.Run(cmd, args))
			return
		}
	}

	fmt.Fprintf(os.Stderr, "thinkgo: unknown subcommand %q\nRun 'thinkgo help' for usage.\n", args[0])
	os.Exit(2)
}

var usageTemplate = `ThinkGo is a tool for managing thinkgo framework.

Usage:

	thinkgo command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "thinkgo help [command]" for more information about a command.

Additional help topics:
{{range .}}{{if not .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "thinkgo help [topic]" for more information about that topic.

`

var helpTemplate = `{{if .Runnable}}usage: thinkgo {{.UsageLine}}

{{end}}{{.Long | trim}}
`

func usage() {
	tmpl(os.Stdout, usageTemplate, commands)
	os.Exit(2)
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
		return template.HTML(strings.TrimSpace(string(s)))
	}})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func help(args []string) {
	if len(args) == 0 {
		usage()
		// not exit 2: succeeded at 'go help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stdout, "usage: thinkgo help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'thinkgo help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'go help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stdout, "Unknown help topic %#q.  Run 'thinkgo help'.\n", arg)
	os.Exit(2) // failed at 'thinkgo help cmd'
}
