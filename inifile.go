// Package inifile implements parsing a simple ini-file
package inifile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// index is section/name pair
type index struct {
	section string
	name    string
}

// IniFile stores the parsed values from a ini-file
type IniFile struct {
	sections []string
	data     map[index]string
}

// Get returns the value of variable from the specified section (use "" for an unnamed section)
func (i *IniFile) Get(section string, name string) string {
	return i.data[index{section: section, name: name}]
}

// Sections returns a list of partitions
func (i *IniFile) Sections() []string {
	return i.sections
}

// envFilename reads an environment variable that matches the static filename
func envFilename(filename string) string {
	if filepath.Base(filename) != filename {
		return filename
	}
	if name := os.Getenv(strings.ToUpper(strings.Replace(filename, ".", "_", -1))); name != "" {
		return name
	}
	return filename
}

// Read parses the specified ini-file
func Read(filename string) (*IniFile, error) {
	data, err := ioutil.ReadFile(envFilename(filename))
	if err != nil {
		return nil, err
	}

	output := &IniFile{data: map[index]string{}}
	section := ""

	for no, s := range strings.Split(string(data), "\n") {
		// remove comment
		i := strings.IndexAny(s, "#;")
		if i >= 0 {
			s = s[0:i]
		}

		// ignore blank characters and line feeds
		s = strings.Trim(s, " \t\n\r")
		if len(s) < 1 {
			continue
		}

		if s[0] == '[' {
			// section name
			section = strings.Trim(s, "[] \t")
			output.sections = append(output.sections, section)

		} else {
			i = strings.IndexByte(s, '=')
			if i > 0 {
				// name = value
				name := strings.Trim(s[0:i], " \t")
				value := strings.Trim(s[i+1:], " '\"\t")
				output.data[index{section: section, name: name}] = value
			} else {
				return nil, fmt.Errorf("Unknown value in file %s in line %d", filename, no+1)
			}
		}
	}
	return output, nil
}
