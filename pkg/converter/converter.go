package converter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/pmengelbert/timplate/pkg/timesheet"
)

type (
	Converter struct {
		Infile         string
		InfileText     []byte
		Outfile        string
		Buffer         *bytes.Buffer
		Sheet          timesheet.Sheet
		Template       *template.Template
		TemplateString string
	}
)

func DefaultConverter(infile, outfile string) (*Converter, error) {
	c := &Converter{
		Infile:         infile,
		Outfile:        outfile,
		Buffer:         new(bytes.Buffer),
		TemplateString: timesheetTemplate,
	}

	err := c.loadInfileText()
	if err != nil {
		return nil, err
	}

	c.parseTemplate()
	err = c.executeTemplate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Converter) loadInfileText() error {
	var err error
	c.InfileText, err = ioutil.ReadFile(c.Infile)
	if err != nil {
		fmt.Println("error reading file")
		os.Exit(1)
	}
	return err
}

func (c *Converter) parseTemplate() {
	c.Template = template.Must(template.New("timesheet").Delims("<<", ">>").
		Parse(c.TemplateString))
}

func (c *Converter) executeTemplate() error {
	err := yaml.Unmarshal(c.InfileText, &c.Sheet)
	if err != nil {
		return err
	}

	c.Sheet.CapitalizeDescriptions()
	c.Template.Execute(c.Buffer, c.Sheet)
	return nil
}

func (c *Converter) SaveOutfile() error {
	err := ioutil.WriteFile(c.Outfile, c.Buffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Converter) CompilePDF() error {
	cmd := exec.Command("pdflatex", c.Outfile)
	str, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(str))
	err = c.cleanUpIntermediateFiles()
	if err != nil {
		return err
	}

	return nil
}

func (c *Converter) cleanUpIntermediateFiles() error {
	baseName := strings.Split(c.Outfile, ".")[0]
	for _, s := range []string{".aux", ".log", ".tex"} {
		os.Remove(baseName + s)
	}
	os.Remove(c.Outfile)

	return nil
}
