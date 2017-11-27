package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	helmdeployLinter "github.com/remyLemeunier/k8s-deploy/linter"
	"github.com/remyLemeunier/k8s-deploy/releases"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	helmlint "k8s.io/helm/pkg/lint"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "",
	Long:  ``,
	RunE:  lint,
}

var displayTempate bool

func init() {
	RootCmd.AddCommand(lintCmd)
	lintCmd.PersistentFlags().BoolVar(&displayTempate, "display-template", false, "Display the rendered template")

}

func lint(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("First args must be the release path.")
	}

	manifest, err := releases.ReadManifestFromFile(path.Join(args[0], "release.yaml"))
	if err != nil {
		return err
	}

	chartPath := path.Join(args[0], manifest.Chart)

	values, err := loadValues(args[0], manifest)
	if err != nil {
		return err
	}
	config := &chart.Config{Raw: string(values), Values: map[string]*chart.Value{}}
	chart, err := chartutil.Load(chartPath)
	if err != nil {
		return err
	}

	options := chartutil.ReleaseOptions{}
	vals, err := chartutil.ToRenderValues(chart, config, options)
	if err != nil {
		return err
	}

	tmpChartPath, err := createTemporaryChartConfig(chart, vals, manifest, chartPath)
	if err != nil {
		return err
	}
	//cleanup
	//defer os.RemoveAll(tmpChartPath)
	fmt.Println("Linter from helm-lint")
	linter := helmlint.All(tmpChartPath)
	for _, message := range linter.Messages {
		fmt.Println(message.Error())
	}

	fmt.Println("-----------------------")
	fmt.Println("Linter from helm-deploy")
	deployLinter := helmdeployLinter.All(tmpChartPath)
	for _, message := range deployLinter.Messages {
		fmt.Println(message.Error())
	}

	fmt.Println("-----------------------")

	if displayTempate {
		// Display templates
		templatesPath := tmpChartPath + "/templates"
		files, err := ioutil.ReadDir(templatesPath)
		if err != nil {
			return err
		}

		for _, f := range files {
			data, err := ioutil.ReadFile(templatesPath + "/" + f.Name())
			if err != nil {
				return err
			}
			fmt.Println(f.Name())
			fmt.Println(string(data))
			fmt.Println("--------------------------------------")
		}
	}

	return nil
}

func createTemporaryChartConfig(chart *chart.Chart, vals map[string]interface{}, manifest *releases.Manifest, chartPath string) (string, error) {
	renderer := engine.New()
	out, err := renderer.Render(chart, vals)
	if err != nil {
		return "", err
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	baseDir := path.Base(manifest.Chart)
	os.Mkdir(tmpDir+"/"+baseDir, 0777)
	if err != nil {
		return "", err
	}

	dir := tmpDir + "/" + baseDir
	// Recreate the complete chart dir with rendered template.
	err = cp(dir+"/Chart.yaml", chartPath+"/Chart.yaml")
	if err != nil {
		return "", err
	}

	err = cp(dir+"/values.yaml", chartPath+"/values.yaml")
	if err != nil {
		return "", err
	}

	// Recreate templates
	os.Mkdir(dir+"/templates", 0777)
	if err != nil {
		return "", err
	}
	for key, value := range out {
		_, fileName := path.Split(key)
		tmpFileName := dir + "/templates/" + fileName
		err := ioutil.WriteFile(tmpFileName, []byte(value), 0777)
		if err != nil {
			return "", err
		}
	}

	return dir, nil
}

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func loadValues(releasePath string, manifest *releases.Manifest) ([]byte, error) {
	defaultValues, err := ioutil.ReadFile(path.Join(releasePath, manifest.Chart, "values.yaml"))
	if err != nil {
		return nil, err
	}

	filesPath := make([]string, len(manifest.ValueFiles))
	for k, v := range manifest.ValueFiles {
		filesPath[k] = path.Join(releasePath, v)
	}

	values, err := releases.LoadValues(filesPath, []string{})
	if err != nil {
		return nil, err
	}

	var defaultValuesMap map[string]interface{}
	err = yaml.Unmarshal(defaultValues, &defaultValuesMap)
	if err != nil {
		return nil, err
	}

	var valuesMap map[string]interface{}
	err = yaml.Unmarshal(values, &valuesMap)
	if err != nil {
		return nil, err
	}

	valuesMerged := mergeValues(defaultValuesMap, valuesMap)

	v, err := yaml.Marshal(valuesMerged)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Borrowed from helm
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
