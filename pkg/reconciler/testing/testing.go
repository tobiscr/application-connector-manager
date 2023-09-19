package testing

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	"github.com/kyma-project/application-connector-manager/pkg/yaml"
	"golang.org/x/exp/slices"
)

var (
	NamesFnApply        = "github.com/kyma-project/application-connector-manager/pkg/reconciler.sFnApply"
	NamesFnUpdateStatus = "github.com/kyma-project/application-connector-manager/pkg/reconciler.sFnUpdate.stopWithErrorAndNoRequeue.sFnUpdateStatus.func2"
	// test data files
	TdUpdateAcmValid = "acm-valid.yaml"
)

type StateTest string

var (
	SfnUpdate StateTest = "update"
)

func LoadTestData(st StateTest) (map[string][]unstructured.Unstructured, error) {
	fullDirPath, err := filepath.Abs(path.Join("testdata", string(st)))
	if err != nil {
		return nil, fmt.Errorf("unable to determine absolute path: %w", err)
	}
	// open directory containing test data
	dirFile, err := os.Open(fullDirPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open test data directory: %w", err)
	}
	// list test data files
	dirEntries, err := dirFile.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("unable to read test data directory: %w", err)
	}
	// filter out non yaml files
	dirEntries = slices.DeleteFunc(dirEntries, func(e os.DirEntry) bool {
		info, err := e.Info()
		if err != nil {
			return false
		}

		isYAML := strings.HasSuffix(info.Name(), ".yaml")
		return info.IsDir() || !isYAML
	})
	// prepare results
	result := map[string][]unstructured.Unstructured{}
	for _, e := range dirEntries {
		fullElementPath := path.Join(fullDirPath, e.Name())
		file, err := os.Open(fullElementPath)
		if err != nil {
			return nil, fmt.Errorf("unable to open test data file: %w", err)
		}
		// extract data from data file
		data, err := yaml.LoadData(file)
		if err != nil {
			return nil, fmt.Errorf("unable to extrat content from test data file: %w", err)
		}
		result[e.Name()] = data
	}
	return result, nil
}
