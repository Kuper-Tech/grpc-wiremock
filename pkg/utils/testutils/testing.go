package testutils

import (
	"fmt"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
)

type TestStatus int

const (
	SkipStatus TestStatus = iota
	SuccessStatus
	FailStatus
	UnknownStatus
)

func (t TestStatus) String() string {
	switch t {
	case SkipStatus:
		return "skip"
	case FailStatus:
		return "fail"
	case SuccessStatus:
		return "success"
	}

	return "unknown"
}

func ReadTestsStatuses(fs afero.Fs, path string) (TestToCases, error) {
	body, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, fmt.Errorf("file %s read: %w", path, err)
	}

	rawTests := map[string]interface{}{}

	if err = yaml.Unmarshal(body, &rawTests); err != nil {
		return nil, fmt.Errorf("file %s unmarshal: %w", path, err)
	}

	testsToCases := TestToCases{}

	for testName, rawTestCases := range rawTests {
		skip := map[string]struct{}{}
		fail := map[string]struct{}{}
		success := map[string]struct{}{}

		testCases, ok := rawTestCases.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("parse tests statuses")
		}

		for caseStatus, rawCaseNames := range testCases {
			sliceCaseNames, ok := rawCaseNames.([]interface{})
			if !ok {
				continue
			}

			var caseNames []string
			for _, rawCaseName := range sliceCaseNames {
				caseNames = append(caseNames, rawCaseName.(string))
			}

			switch caseStatus {
			case "skip":
				skip = sliceutils.SliceToMap(caseNames)
			case "fail":
				fail = sliceutils.SliceToMap(caseNames)
			case "success":
				success = sliceutils.SliceToMap(caseNames)
			}
		}

		testsToCases[testName] = testCasesMaps{
			Skip:    skip,
			Fail:    fail,
			Success: success,
		}
	}

	return testsToCases, nil
}

type testCasesMaps struct {
	Skip    map[string]struct{}
	Fail    map[string]struct{}
	Success map[string]struct{}
}

func (t testCasesMaps) getByCase(name string) TestStatus {
	if _, exists := t.Fail[name]; exists {
		return FailStatus
	}

	if _, exists := t.Skip[name]; exists {
		return SkipStatus
	}

	if _, exists := t.Success[name]; exists {
		return SuccessStatus
	}

	return UnknownStatus
}

type TestToCases map[string]testCasesMaps

func (t TestToCases) GetByNameForTest(testName, caseName string) TestStatus {
	return t[testName].getByCase(caseName)
}
