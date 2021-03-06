package pkg

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"k8s.io/helm/pkg/engine"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Assets struct {
	dryRun                 bool
	rootDir                string
	valuesOverrideFile     string
	valuesOverrideFileMode os.FileMode
	valuesOverride         map[string]interface{}
}

func NewAssets(rootDir string, valuesFile string, dryRun bool) (Assets, error) {
	values, err := vals(valuesFile)
	if err != nil {
		return Assets{}, err
	}

	info, err := os.Stat(rootDir)
	if err != nil {
		return Assets{}, fmt.Errorf("The assets directory must be present")
	}
	return Assets{dryRun, rootDir, valuesFile, info.Mode(), values}, nil
}

func (a Assets) Render() error {
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
		return err
	}

	return a.updateValuesFile(assetMap)
}

func (a Assets) updateValuesFile(assetMap map[string]string) error {
	cfgMapData := make(map[string]string)
	for cfgMapKey, assetContent := range assetMap {
		cfgMapData[cfgMapKey] = assetContent
	}
	a.valuesOverride["data"] = cfgMapData
	updatedValues, err := yaml.Marshal(a.valuesOverride)
	if err != nil {
		return nil
	}
	fmt.Println(string(updatedValues))
	if a.dryRun {
		fmt.Printf("Cowardly refuse to overwrite %s. Please use --i-am-sure\n", a.valuesOverrideFile)
	} else {
		ioutil.WriteFile(a.valuesOverrideFile, updatedValues, a.valuesOverrideFileMode)
	}
	return nil
}

func renderSingle(path string, vals map[string]interface{}) (string, error) {
	fmt.Printf("Rendering %s...\n", path)
	t := template.New(path).Funcs(engine.FuncMap()).Delims("<<", ">>")
	fileContent, _ := ioutil.ReadFile(path)
	t.Parse(string(fileContent))
	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]interface{}{"Values": vals}); err != nil {
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
