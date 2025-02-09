package test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/defsec/pkg/scanners/helm/parser"
)

func Test_helm_parser(t *testing.T) {

	tests := []struct {
		testName  string
		chartName string
	}{
		{
			testName:  "Parsing directory 'testchart'",
			chartName: "testchart",
		},
	}

	for _, test := range tests {
		chartName := test.chartName

		t.Logf("Running test: %s", test.testName)

		helmParser := parser.New(chartName)
		err := helmParser.ParseFS(context.TODO(), os.DirFS(filepath.Join("testdata", chartName)), ".")
		require.NoError(t, err)
		manifests, err := helmParser.RenderedChartFiles()
		require.NoError(t, err)

		assert.Len(t, manifests, 3)

		for _, manifest := range manifests {
			expectedPath := filepath.Join("testdata", "expected", manifest.TemplateFilePath)

			expectedContent, err := os.ReadFile(expectedPath)
			require.NoError(t, err)

			assert.Equal(t, strings.ReplaceAll(string(expectedContent), "\r\n", "\n"), strings.ReplaceAll(manifest.ManifestContent, "\r\n", "\n"))
		}
	}
}

func Test_helm_tarball_parser(t *testing.T) {

	tests := []struct {
		testName    string
		chartName   string
		archiveFile string
	}{
		{
			testName:    "standard tarball",
			chartName:   "mysql",
			archiveFile: "mysql-8.8.26.tar",
		},
		{
			testName:    "gzip tarball with tar.gz extension",
			chartName:   "mysql",
			archiveFile: "mysql-8.8.26.tar.gz",
		},
		{
			testName:    "gzip tarball with tgz extension",
			chartName:   "mysql",
			archiveFile: "mysql-8.8.26.tgz",
		},
	}

	for _, test := range tests {

		t.Logf("Running test: %s", test.testName)

		testPath := filepath.Join("testdata", test.archiveFile)

		testTemp := t.TempDir()
		testFileName := filepath.Join(testTemp, test.archiveFile)
		require.NoError(t, copyArchive(testPath, testFileName))

		testFs := os.DirFS(testTemp)

		helmParser := parser.New(test.chartName)
		err := helmParser.ParseFS(context.TODO(), testFs, ".")
		require.NoError(t, err)

		manifests, err := helmParser.RenderedChartFiles()
		require.NoError(t, err)

		assert.Len(t, manifests, 6)

		oneOf := []string{
			"configmap.yaml",
			"statefulset.yaml",
			"svc-headless.yaml",
			"svc.yaml",
			"secrets.yaml",
			"serviceaccount.yaml",
		}

		for _, manifest := range manifests {
			filename := filepath.Base(manifest.TemplateFilePath)
			assert.Contains(t, oneOf, filename)

			if strings.HasSuffix(manifest.TemplateFilePath, "secrets.yaml") {
				continue
			}
			expectedPath := filepath.Join("testdata", "expected", manifest.TemplateFilePath)

			expectedContent, err := os.ReadFile(expectedPath)
			require.NoError(t, err)

			assert.Equal(t, strings.ReplaceAll(string(expectedContent), "\r\n", "\n"), strings.ReplaceAll(manifest.ManifestContent, "\r\n", "\n"))
		}
	}
}
