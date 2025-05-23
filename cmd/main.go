package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

func parse(file string, values map[string]any) error {
	t := template.New(file)
	t = t.Funcs(sprig.TxtFuncMap())

	t, err := t.ParseFiles(file)
	if err != nil {
		return err
	}

	t.Option("missingkey=error")
	// t.Option("missingkey=zero")

	if err := t.ExecuteTemplate(os.Stdout, filepath.Base(file), values); err != nil {
		return err
	}

	return nil
}

func main() {
	valuesFile := flag.String("values", "", "--values <file>")
	flag.Parse()

	var values map[string]any
	if valuesFile != nil && *valuesFile != "" {
		b, err := os.ReadFile(*valuesFile)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(b, &values); err != nil {
			panic(err)
		}
	}

	args := flag.CommandLine.Args()
	if len(args) != 1 {
		panic(fmt.Errorf("must have 1 non-flag argument"))
	}

	if err := parse(args[0], values); err != nil {
		panic(err)
	}
}
