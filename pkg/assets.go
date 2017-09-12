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

func NewAssets(rootDir string, files ValuesOverrideFiles) (Assets, error) {
	values, err := vals(files)
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

func vals(valsFiles ValuesOverrideFiles) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	for _, filePath := range valsFiles {
		currentMap := map[string]interface{}{}
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return map[string]interface{}{}, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return map[string]interface{}{}, fmt.Errorf("failed to parse %s: %s", filePath, err)
		}
		// Merge with the previous map
		base = mergeValues(base, currentMap)
	}

	return base, nil
}

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = nextMap
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
