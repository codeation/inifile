// Package inifile implements parsing a simple ini-file.
package inifile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var errorValue = errors.New("unknown variable value")

const sides = 2

// index is section/name pair.
type index struct {
	section string
	name    string
}

// IniFile stores the parsed values from a ini-file.
type IniFile struct {
	sections       []string
	data           map[index]string
	commandEnabled bool
}

// command makes substitution if external command or filename is encapsulated.
func command(s string) string {
	if strings.HasPrefix(s, "$(<") && strings.HasSuffix(s, ")") {
		// make filename contents substitution in line "$(<filename)"
		filename := strings.TrimSuffix(strings.TrimPrefix(s, "$(<"), ")")
		if data, err := ioutil.ReadFile(filename); err == nil {
			return string(data)
		}
	} else if strings.HasPrefix(s, "$(") && strings.HasSuffix(s, ")") {
		// encapsulate command output for line "$(command line)"
		commandLine := strings.TrimSuffix(strings.TrimPrefix(s, "$("), ")")
		commands := strings.Split(commandLine, " ")
		if data, err := exec.Command(commands[0], commands[1:]...).Output(); err == nil {
			return string(data)
		}
	}
	// without substitution
	return s
}

// Command enables or disables external commands encapsulation.
// Note that executing external commands can lead to application vulnerabilities.
func (ini *IniFile) Command(enabled bool) {
	ini.commandEnabled = enabled
}

// Get returns the value of variable from the specified section
// (use "" for an unnamed section).
// Get makes substitution if commands encapsulation is enabled.
// "$(<filename)" format for the file contents.
// "$(command line)" format to encapsulate command output.
func (ini *IniFile) Get(section string, name string) string {
	s := ini.data[index{section: section, name: name}]
	if ini.commandEnabled {
		s = command(s)
	}

	return s
}

// Sections returns a list of partitions.
func (ini *IniFile) Sections() []string {
	return ini.sections
}

// envFilename reads an environment variable that matches the base file name.
func envFilename(filename string) string {
	// The environment variable is checked only if
	// the filename does not contain the directory name
	if filepath.Base(filename) != filename {
		return filename
	}
	// The name of the environment variable is the file name in uppercase,
	// the dot is replaced by an underscore
	envName := strings.ToUpper(strings.Replace(filename, ".", "_", -1))
	if name := os.Getenv(envName); name != "" {
		return name
	}
	// Environment variable not found
	return filename
}

// Read parses the specified ini-file.
// You can specify a full path to INI file via the environment variable.
// The environment variable is checked only if
// the filename does not contain the directory name.
// The name of the environment variable is the file name in uppercase,
// the dot is replaced by an underscore.
func Read(filename string) (*IniFile, error) {
	data, err := ioutil.ReadFile(envFilename(filename))
	if err != nil {
		return nil, fmt.Errorf("readFile: %w", err)
	}

	ini := &IniFile{
		data: map[index]string{},
	}

	var sectionName string

	for no, s := range strings.Split(string(data), "\n") {
		// remove comments and trim spaces or line feeds
		if pos := strings.IndexAny(s, "#;"); pos >= 0 {
			s = s[0:pos]
		}

		s = strings.TrimSpace(s)

		// ignore empty line
		if len(s) == 0 {
			continue
		}

		// new section name
		if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
			sectionName = strings.TrimSuffix(strings.TrimPrefix(s, "["), "]")
			ini.sections = append(ini.sections, sectionName)
			continue
		}

		// parsing the name=value pair
		ss := strings.SplitN(s, "=", sides)
		if len(ss) != sides {
			return nil, fmt.Errorf("%w (%s, %d)", errorValue, filename, no+1)
		}

		ini.data[index{section: sectionName, name: strings.TrimSpace(ss[0])}] = strings.TrimSpace(ss[1])
	}
	return ini, nil
}
