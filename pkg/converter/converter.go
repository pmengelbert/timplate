package converter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/pmengelbert/timplate/pkg/timesheet"
)

const (
	outDirName  = "._latexFiles/"
	styFileName = "._enumitem.sty"
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
		EnumItemString string
	}
)

var escapeRegex = regexp.MustCompile("([&%$#_{}~\\^])")

func DefaultConverter(infile, outfile string) (*Converter, error) {
	c := &Converter{
		Infile:         infile,
		Outfile:        outfile,
		Buffer:         new(bytes.Buffer),
		TemplateString: timesheetTemplate,
		EnumItemString: enumitem,
	}

	err := c.loadInfileText()
	if err != nil {
		return nil, err
	}

	c.parseTemplate()

	err = c.parseYaml()
	if err != nil {
		return nil, err
	}

	c.executeTemplate()

	return c, nil
}

func (c *Converter) loadInfileText() error {
	var err error
	c.InfileText, err = ioutil.ReadFile(c.Infile)
	if err != nil {
		fmt.Println("error reading file")
		os.Exit(1)
	}
	c.InfileText = escapeRegex.ReplaceAll(c.InfileText, []byte("\\$1"))
	return err
}

func (c *Converter) parseTemplate() {
	c.Template = template.Must(template.New("timesheet").Delims("<<", ">>").
		Parse(c.TemplateString))
}

func (c *Converter) parseYaml() error {
	err := yaml.Unmarshal(c.InfileText, &c.Sheet)
	if err != nil {
		return err
	}

	for i, r := range c.Sheet.Records {
		for _, t := range r.Times {
			a := strings.Split(t, "-")
			for j := range a {
				a[j] = strings.TrimSpace(a[j])
			}

			x, err := timesheet.Parse(a[0])
			if err != nil {
				return err
			}

			y, err := timesheet.Parse(a[1])
			if err != nil {
				return err
			}

			diff := y.DifferenceInHours(x)
			if diff < 0 {
				return fmt.Errorf("bad time period: %v, %v", x, y)
			}

			c.Sheet.Records[i].TimeSum += diff
		}
	}

	return nil
}

func (c *Converter) executeTemplate() {
	c.Sheet.CapitalizeDescriptions()
	c.Template.Execute(c.Buffer, c.Sheet)
}

func (c *Converter) SaveOutfile() error {
	err := ioutil.WriteFile(c.Outfile, c.Buffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Converter) CompilePDF() error {
	ioutil.WriteFile(styFileName, []byte(c.EnumItemString), 0644)

	os.Mkdir(outDirName, 0755)
	cmd := exec.Command("pdflatex", "-output-directory=._latexFiles",
		"-halt-on-error", c.Outfile)

	str, err := cmd.Output()
	if err != nil {
		fmt.Printf("pdflatex encountered an error\n")
		c.cleanUpIntermediateFiles()
		return err
	}

	pdfFilename := strings.TrimSuffix(c.Outfile, path.Ext(c.Outfile)) + ".pdf"
	err = os.Rename(outDirName+pdfFilename, pdfFilename)
	if err != nil {
		fmt.Println("error renaming file")
		return err
	}

	fmt.Println(string(str))
	err = c.cleanUpIntermediateFiles()
	if err != nil {
		fmt.Printf("Was unable to delete intermediate files. Delete them manually in the "+
			"%s directory.\n", outDirName)
	}

	os.Remove(c.Outfile)

	if err != nil {
		return err
	}

	return nil
}

func (c *Converter) cleanUpIntermediateFiles() error {
	err := os.RemoveAll(outDirName)
	if err != nil {
		return err
	}

	err = os.Remove(styFileName)
	if err != nil {
		return err
	}

	return nil
}
