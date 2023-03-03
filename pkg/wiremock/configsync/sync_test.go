package configsync

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

var (
	osfs       = afero.NewOsFs()
	projectDir = path.Join(fsutils.CurrentDir(), "../../..")
)

func TestSyncWiremockConfig(t *testing.T) {
	tests := []struct {
		name string

		fs           afero.Fs
		wiremockPath string
		templatePath string

		source config.Wiremock
		want   config.Wiremock
	}{
		{
			fs:           osfs,
			wiremockPath: filepath.Join(projectDir, "tests", "two-services"),
			templatePath: filepath.Join(projectDir, "static/tests/data/wiremock/with-domains"),
			source: config.Wiremock{
				Services: []config.Service{
					{Name: "awesome", Port: 8000, RootDir: "/home/mock/awesome"},
					{Name: "push-sender", Port: 8001, RootDir: "/home/mock/push-sender"},
				},
			},
			want: config.Wiremock{
				Services: []config.Service{
					{Name: "awesome", Port: 8000, RootDir: "/home/mock/awesome"},
					{Name: "push-sender", Port: 8001, RootDir: "/home/mock/push-sender"},
				},
			},
		},
		{
			wiremockPath: filepath.Join(projectDir, "tests", "two-services"),
			templatePath: filepath.Join(projectDir, "static/tests/data/wiremock/with-domains"),

			fs: osfs,
			source: config.Wiremock{
				Services: []config.Service{
					{Name: "domain1", Port: 8001, RootDir: "/home/mock/domain1"},
					{Name: "domain2", Port: 8002, RootDir: "/home/mock/domain2"},
				},
			},
			want: config.Wiremock{
				Services: []config.Service{
					{Name: "awesome", Port: 8003, RootDir: filepath.Join(projectDir, "tests/two-services/awesome")},
					{Name: "push-sender", Port: 8004, RootDir: filepath.Join(projectDir, "tests/two-services/push-sender")},
				},
			},
		},
		{
			wiremockPath: filepath.Join(projectDir, "tests", "two-services"),
			templatePath: filepath.Join(projectDir, "static/tests/data/wiremock/with-domains"),

			fs: osfs,
			source: config.Wiremock{
				Services: []config.Service{
					{Name: "awesome", Port: 8000, RootDir: "/home/mock/awesome"},
					{Name: "domain2", Port: 8002, RootDir: "/home/mock/domain2"},
				},
			},
			want: config.Wiremock{
				Services: []config.Service{
					{Name: "awesome", Port: 8000, RootDir: "/home/mock/awesome"},
					{Name: "push-sender", Port: 8003, RootDir: filepath.Join(projectDir, "tests/two-services/push-sender")},
				},
			},
		},
		{
			wiremockPath: filepath.Join(projectDir, "tests", "two-services"),
			templatePath: filepath.Join(projectDir, "static/tests/data/wiremock/with-domains"),

			fs:     osfs,
			source: config.Wiremock{},
			want: config.Wiremock{
				Services: []config.Service{
					{Name: "awesome", Port: 8000, RootDir: filepath.Join(projectDir, "tests/two-services/awesome")},
					{Name: "push-sender", Port: 8001, RootDir: filepath.Join(projectDir, "tests/two-services/push-sender")},
				},
			},
		},
	}

	for _, tt := range tests {
		err := fsutils.CopyDir(osfs, osfs, tt.templatePath, tt.wiremockPath, false)
		require.NoError(t, err)

		t.Run(tt.name, func(t *testing.T) {
			got, err := SyncWiremockConfig(tt.fs, tt.source, tt.wiremockPath)
			require.NoError(t, err)

			require.ElementsMatch(t, got.Services, tt.want.Services)
		})
	}
}

func TestSyncWiremockConfig_Dynamic(t *testing.T) {
	testCase := struct {
		name string

		fs afero.Fs

		wiremockPath string
		templatePath string

		sourceConfig     config.Wiremock
		wantBeforeDelete config.Wiremock
		wantAfterDelete  config.Wiremock
	}{
		fs:           osfs,
		wiremockPath: filepath.Join(projectDir, "tests", "two-services"),
		templatePath: filepath.Join(projectDir, "static/tests/data/wiremock/with-domains"),

		sourceConfig: config.Wiremock{
			Services: []config.Service{
				{Name: "awesome", Port: 8000, RootDir: "/home/mock/awesome"},
				{Name: "domain2", Port: 8002, RootDir: "/home/mock/domain2"},
			},
		},
		wantBeforeDelete: config.Wiremock{
			Services: []config.Service{
				{Name: "awesome", Port: 8000, RootDir: "/home/mock/awesome"},
				{Name: "push-sender", Port: 8003, RootDir: filepath.Join(projectDir, "tests/two-services/push-sender")},
			},
		},
		wantAfterDelete: config.Wiremock{
			Services: []config.Service{
				{Name: "push-sender", Port: 8003, RootDir: filepath.Join(projectDir, "tests/two-services/push-sender")},
			},
		},
	}

	err := fsutils.CopyDir(testCase.fs, testCase.fs, testCase.templatePath, testCase.wiremockPath, false)
	require.NoError(t, err)

	t.Run(testCase.name, func(t *testing.T) {
		got, err := SyncWiremockConfig(testCase.fs, testCase.sourceConfig, testCase.wiremockPath)
		require.NoError(t, err)

		require.ElementsMatch(t, got.Services, testCase.wantBeforeDelete.Services)
	})

	err = testCase.fs.RemoveAll(filepath.Join(testCase.wiremockPath, "awesome"))
	require.NoError(t, err)

	t.Run(testCase.name, func(t *testing.T) {
		got, err := SyncWiremockConfig(testCase.fs, testCase.wantBeforeDelete, testCase.wiremockPath)
		require.NoError(t, err)

		require.ElementsMatch(t, got.Services, testCase.wantAfterDelete.Services)
	})
}
