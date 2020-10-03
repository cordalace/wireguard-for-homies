package inputdata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func New(configurators ...Configurator) *Config {
	return NewDefaultConfig().WithOptions(configurators...)
}

func (c *Config) WithOptions(configurators ...Configurator) *Config {
	clonedConfig := c.clone()

	for _, configurator := range configurators {
		configurator(clonedConfig)
	}

	return clonedConfig
}

func NewDefaultConfig() *Config {
	return (&Config{}).WithOptions(
		InputSubdirectory(".inputs"),
		InputFileExtension(""),
	)
}

func InputSubdirectory(name string) Configurator {
	return func(c *Config) {
		c.subDirName = name
	}
}

func InputFileExtension(inputFileExtension string) Configurator {
	return func(c *Config) {
		c.inputFileExtension = inputFileExtension
	}
}

type Configurator func(*Config)

type Config struct {
	subDirName         string
	inputFileExtension string
}

func (c *Config) clone() *Config {
	return &Config{
		subDirName:         c.subDirName,
		inputFileExtension: c.inputFileExtension,
	}
}

func (c *Config) LoadT(t *testing.T) []byte {
	// t.Helper()
	safeTestName := strings.ReplaceAll(t.Name(), "/", "-")
	inputFile := filepath.Join(c.subDirName, safeTestName+c.inputFileExtension)

	buf, err := ioutil.ReadFile(inputFile)

	if os.IsNotExist(err) {
		t.Fatalf("input does not exist for test %s: %s", safeTestName, err)
	}

	if err != nil {
		t.Fatalf("error loading input for test %s: %s", safeTestName, err)
	}

	return buf
}
