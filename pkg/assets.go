package pkg

import (
	"k8s.io/helm/pkg/engine"
	"bytes"
	"text/template"
	"fmt"
	"io/ioutil"
	"github.com/ghodss/yaml"
	"strings"
	"path/filepath"
	"os"
	//"sync"
)

type Assets struct {
	rootDir string
	valuesOverride map[string]interface{}
}

func NewAssets(rootDir string, valuesFile string) (Assets, error) {
	values, err := vals(valuesFile)
	if err != nil {
		return Assets{}, err
	}

	return Assets{rootDir, values}, nil
}

func (a Assets) Render() (map[string]string, error) {
	assetMap := make(map[string]string)
	err := filepath.Walk(a.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Use the flattened file path in the configmap
		flattenedFilePath := strings.Replace(path, "/", "_", -1)
		assetContent, err := renderSingle(path, a.valuesOverride)
		if err != nil {
			return err
		}
		assetMap[flattenedFilePath] = assetContent
		return nil
	})

	if err != nil {
		return map[string]string{}, err
	}

	return assetMap, nil
}

type ValuesOverrideFiles []string

func (v *ValuesOverrideFiles) String() string {
	return fmt.Sprint(*v)
}

func (v *ValuesOverrideFiles) Type() string {
	return "ValuesOverrideFiles"
}

func (v *ValuesOverrideFiles) Set(value string) error {
	for _, filePath := range strings.Split(value, ",") {
		*v = append(*v, filePath)
	}
	return nil
}

func renderSingle(path string, vals map[string]interface{}) (string, error) {
	fmt.Printf("Rendering %s...\n", path)
	t := template.New(path).Funcs(engine.FuncMap()).Delims("<<", ">>")
	fileContent, _ := ioutil.ReadFile(path)
	t.Parse(string(fileContent))
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, path, map[string]interface{}{"Values": vals}); err != nil {
		return "", fmt.Errorf("render error in %q: %s", path, err)
	}

	return buf.String(), nil
}

func vals(filePath string) (map[string]interface{}, error) {
	currentMap := map[string]interface{}{}
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return map[string]interface{}{}, err
	}

	if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
		return map[string]interface{}{}, fmt.Errorf("failed to parse %s: %s", filePath, err)
	}

	return currentMap, nil
}

