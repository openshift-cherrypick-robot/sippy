// nolint
package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apitype "github.com/openshift/sippy/pkg/apis/api"
)

func fakeComponentAndCapabilityGetter(test apitype.ComponentTestIdentification, stats apitype.ComponentTestStatus) (string, []string) {
	name := stats.TestName
	known := map[string]struct {
		component    string
		capabilities []string
	}{
		"test 1": {
			component:    "component 1",
			capabilities: []string{"cap1"},
		},
		"test 2": {
			component:    "component 2",
			capabilities: []string{"cap21", "cap22"},
		},
		"test 3": {
			component:    "component 1",
			capabilities: []string{"cap1"},
		},
	}
	if comCap, ok := known[name]; ok {
		return comCap.component, comCap.capabilities
	}
	return "other", []string{"other"}
}

var (
	defaultAdvancedOption = apitype.ComponentReportRequestAdvancedOptions{
		Confidence:     95,
		PityFactor:     5,
		MinimumFailure: 3,
	}
	defaultComponentReportGenerator = componentReportGenerator{
		gcsBucket:                             "test-platform-results",
		ComponentReportRequestVariantOptions:  apitype.ComponentReportRequestVariantOptions{GroupBy: "cloud,arch,network"},
		ComponentReportRequestAdvancedOptions: defaultAdvancedOption,
	}
	groupByVariantComponentReportGenerator = componentReportGenerator{
		gcsBucket:                             "test-platform-results",
		ComponentReportRequestVariantOptions:  apitype.ComponentReportRequestVariantOptions{GroupBy: "cloud,arch,network,variants"},
		ComponentReportRequestAdvancedOptions: defaultAdvancedOption,
	}
	componentPageGenerator = componentReportGenerator{
		gcsBucket: "test-platform-results",
		ComponentReportRequestTestIdentificationOptions: apitype.ComponentReportRequestTestIdentificationOptions{
			Component: "component 2",
		},
		ComponentReportRequestVariantOptions: apitype.ComponentReportRequestVariantOptions{
			GroupBy: "cloud,arch,network",
		},
		ComponentReportRequestAdvancedOptions: defaultAdvancedOption,
	}
	capabilityPageGenerator = componentReportGenerator{
		gcsBucket: "test-platform-results",
		ComponentReportRequestTestIdentificationOptions: apitype.ComponentReportRequestTestIdentificationOptions{
			Component:  "component 2",
			Capability: "cap22",
		},
		ComponentReportRequestVariantOptions: apitype.ComponentReportRequestVariantOptions{
			GroupBy: "cloud,arch,network",
		},
		ComponentReportRequestAdvancedOptions: defaultAdvancedOption,
	}
	testPageGenerator = componentReportGenerator{
		gcsBucket: "test-platform-results",
		ComponentReportRequestTestIdentificationOptions: apitype.ComponentReportRequestTestIdentificationOptions{
			Component:  "component 2",
			Capability: "cap22",
			TestID:     "2",
		},
		ComponentReportRequestVariantOptions: apitype.ComponentReportRequestVariantOptions{
			GroupBy: "cloud,arch,network",
		},
		ComponentReportRequestAdvancedOptions: defaultAdvancedOption,
	}
	testDetailsGenerator = componentReportGenerator{
		gcsBucket: "test-platform-results",
		ComponentReportRequestTestIdentificationOptions: apitype.ComponentReportRequestTestIdentificationOptions{
			Component:  "component 1",
			Capability: "cap11",
			TestID:     "1",
		},
		ComponentReportRequestVariantOptions: apitype.ComponentReportRequestVariantOptions{
			GroupBy:  "cloud,arch,network",
			Platform: "aws",
			Arch:     "amd64",
			Network:  "ovn",
			Upgrade:  "upgrade-micro",
			Variant:  "standard",
		},
		ComponentReportRequestAdvancedOptions: defaultAdvancedOption,
	}
)

func TestGenerateComponentReport(t *testing.T) {
	awsAMD64OVNTest := apitype.ComponentTestIdentification{
		TestID:       "1",
		Platform:     "aws",
		Arch:         "amd64",
		Network:      "ovn",
		Upgrade:      "upgrade-micro",
		FlatVariants: "standard",
	}
	awsAMD64SDNTest := apitype.ComponentTestIdentification{
		TestID:       "2",
		Platform:     "aws",
		Arch:         "amd64",
		Network:      "sdn",
		Upgrade:      "upgrade-micro",
		FlatVariants: "standard",
	}
	awsAMD64OVN2Test := apitype.ComponentTestIdentification{
		TestID:   "3",
		Platform: "aws",
		Arch:     "amd64",
		Network:  "ovn",
		Upgrade:  "upgrade-micro",
	}
	awsAMD64OVNVariantsTest := apitype.ComponentTestIdentification{
		TestID:       "1",
		Platform:     "aws",
		Arch:         "amd64",
		Network:      "ovn",
		Upgrade:      "upgrade-micro",
		FlatVariants: "standard,fips",
	}
	awsAMD64OVNBaseTestStats90Percent := apitype.ComponentTestStatus{
		TestName:     "test 1",
		Variants:     []string{"standard"},
		TotalCount:   1000,
		FlakeCount:   10,
		SuccessCount: 900,
	}
	awsAMD64OVNBaseTestStats50Percent := apitype.ComponentTestStatus{
		TestName:     "test 1",
		Variants:     []string{"standard"},
		TotalCount:   1000,
		FlakeCount:   10,
		SuccessCount: 500,
	}
	awsAMD64OVNBaseTestStatsVariants90Percent := apitype.ComponentTestStatus{
		TestName:     "test 1",
		Variants:     []string{"standard", "fips"},
		TotalCount:   1000,
		FlakeCount:   10,
		SuccessCount: 900,
	}
	awsAMD64OVNSampleTestStats90Percent := apitype.ComponentTestStatus{
		TestName:     "test 1",
		Variants:     []string{"standard"},
		TotalCount:   100,
		FlakeCount:   1,
		SuccessCount: 90,
	}
	awsAMD64OVNSampleTestStats85Percent := apitype.ComponentTestStatus{
		TestName:     "test 1",
		Variants:     []string{"standard"},
		TotalCount:   100,
		FlakeCount:   1,
		SuccessCount: 85,
	}
	awsAMD64OVNSampleTestStats50Percent := apitype.ComponentTestStatus{
		TestName:     "test 1",
		Variants:     []string{"standard"},
		TotalCount:   100,
		FlakeCount:   1,
		SuccessCount: 50,
	}
	awsAMD64OVNSampleTestStatsTiny := apitype.ComponentTestStatus{
		TestName:     "test 1",
		Variants:     []string{"standard"},
		TotalCount:   3,
		FlakeCount:   0,
		SuccessCount: 1,
	}
	awsAMD64OVNSampleTestStatsVariants90Percent := apitype.ComponentTestStatus{
		TestName:     "test 1",
		Variants:     []string{"standard", "fips"},
		TotalCount:   100,
		FlakeCount:   1,
		SuccessCount: 90,
	}
	awsAMD64SDNBaseTestStats90Percent := apitype.ComponentTestStatus{
		TestName:     "test 2",
		Variants:     []string{"standard"},
		TotalCount:   1000,
		FlakeCount:   10,
		SuccessCount: 900,
	}
	awsAMD64SDNBaseTestStats50Percent := apitype.ComponentTestStatus{
		TestName:     "test 2",
		Variants:     []string{"standard"},
		TotalCount:   1000,
		FlakeCount:   10,
		SuccessCount: 500,
	}
	awsAMD64SDNSampleTestStats90Percent := apitype.ComponentTestStatus{
		TestName:     "test 2",
		Variants:     []string{"standard"},
		TotalCount:   100,
		FlakeCount:   1,
		SuccessCount: 90,
	}
	awsAMD64OVN2BaseTestStats90Percent := apitype.ComponentTestStatus{
		TestName:     "test 3",
		Variants:     []string{"standard"},
		TotalCount:   1000,
		FlakeCount:   10,
		SuccessCount: 900,
	}
	awsAMD64OVN2SampleTestStats80Percent := apitype.ComponentTestStatus{
		TestName:     "test 3",
		Variants:     []string{"standard"},
		TotalCount:   100,
		FlakeCount:   1,
		SuccessCount: 80,
	}
	columnAWSAMD64OVN := apitype.ComponentReportColumnIdentification{
		Platform: "aws",
		Arch:     "amd64",
		Network:  "ovn",
	}
	columnAWSAMD64OVNVariantFipsStandard := apitype.ComponentReportColumnIdentification{
		Platform: "aws",
		Arch:     "amd64",
		Network:  "ovn",
		Variant:  "standard,fips",
	}
	columnAWSAMD64SDN := apitype.ComponentReportColumnIdentification{
		Platform: "aws",
		Arch:     "amd64",
		Network:  "sdn",
	}
	columnAWSAMD64SDNVariantStandard := apitype.ComponentReportColumnIdentification{
		Platform: "aws",
		Arch:     "amd64",
		Network:  "sdn",
		Variant:  "standard",
	}
	columnAWSAMD64OVNFull := apitype.ComponentReportColumnIdentification{
		Platform: "aws",
		Arch:     "amd64",
		Network:  "ovn",
		Upgrade:  "upgrade-micro",
		Variant:  "standard",
	}
	columnAWSAMD64SDNFull := apitype.ComponentReportColumnIdentification{
		Platform: "aws",
		Arch:     "amd64",
		Network:  "sdn",
		Upgrade:  "upgrade-micro",
		Variant:  "standard",
	}
	rowComponent1 := apitype.ComponentReportRowIdentification{
		Component: "component 1",
	}
	rowComponent2 := apitype.ComponentReportRowIdentification{
		Component: "component 2",
	}
	rowComponent2Cap21 := apitype.ComponentReportRowIdentification{
		Component:  "component 2",
		Capability: "cap21",
	}
	rowComponent2Cap22 := apitype.ComponentReportRowIdentification{
		Component:  "component 2",
		Capability: "cap22",
	}
	rowComponent2Cap22Test2 := apitype.ComponentReportRowIdentification{
		Component:  "component 2",
		Capability: "cap22",
		TestName:   "test 2",
		TestID:     "2",
	}

	tests := []struct {
		name           string
		generator      componentReportGenerator
		baseStatus     map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus
		sampleStatus   map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus
		expectedReport apitype.ComponentReport
	}{
		{
			name:      "top page test no significant and missing data",
			generator: defaultComponentReportGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats90Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStats85Percent,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: apitype.ComponentReportRowIdentification{
							Component: "component 1",
						},
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.NotSignificant,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.MissingBasisAndSample,
							},
						},
					},
					{
						ComponentReportRowIdentification: apitype.ComponentReportRowIdentification{
							Component: "component 2",
						},
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.NotSignificant,
							},
						},
					},
				},
			},
		},
		{
			name:      "top page test with both improvement and regression",
			generator: defaultComponentReportGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest:  awsAMD64OVNBaseTestStats90Percent,
				awsAMD64OVN2Test: awsAMD64OVN2BaseTestStats90Percent,
				awsAMD64SDNTest:  awsAMD64SDNBaseTestStats50Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest:  awsAMD64OVNSampleTestStats50Percent,
				awsAMD64OVN2Test: awsAMD64OVN2SampleTestStats80Percent,
				awsAMD64SDNTest:  awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent1,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.ExtremeRegression,
								RegressedTests: []apitype.ComponentReportTestSummary{
									{
										ComponentReportTestIdentification: apitype.ComponentReportTestIdentification{
											ComponentReportRowIdentification: apitype.ComponentReportRowIdentification{
												TestName: awsAMD64OVNBaseTestStats90Percent.TestName,
												TestID:   awsAMD64OVNTest.TestID,
											},
											ComponentReportColumnIdentification: apitype.ComponentReportColumnIdentification{
												Platform: columnAWSAMD64OVN.Platform,
												Arch:     columnAWSAMD64OVN.Arch,
												Network:  columnAWSAMD64OVN.Network,
												Upgrade:  awsAMD64OVNTest.Upgrade,
												Variant:  awsAMD64OVNTest.FlatVariants,
											},
										},
										Status: apitype.ExtremeRegression,
									},
									{
										ComponentReportTestIdentification: apitype.ComponentReportTestIdentification{
											ComponentReportRowIdentification: apitype.ComponentReportRowIdentification{
												TestName: awsAMD64OVN2BaseTestStats90Percent.TestName,
												TestID:   awsAMD64OVN2Test.TestID,
											},
											ComponentReportColumnIdentification: apitype.ComponentReportColumnIdentification{
												Platform: columnAWSAMD64OVN.Platform,
												Arch:     columnAWSAMD64OVN.Arch,
												Network:  columnAWSAMD64OVN.Network,
												Upgrade:  awsAMD64OVN2Test.Upgrade,
												Variant:  awsAMD64OVN2Test.FlatVariants,
											},
										},
										Status: apitype.SignificantRegression,
									},
								},
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.MissingBasisAndSample,
							},
						},
					},
					{
						ComponentReportRowIdentification: rowComponent2,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.SignificantImprovement,
							},
						},
					},
				},
			},
		},
		{
			name:      "component page test no significant and missing data",
			generator: componentPageGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats90Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent2Cap21,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.NotSignificant,
							},
						},
					},
					{
						ComponentReportRowIdentification: rowComponent2Cap22,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.NotSignificant,
							},
						},
					},
				},
			},
		},
		{
			name:      "component page test with both improvement and regression",
			generator: componentPageGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats50Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats50Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent2Cap21,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.SignificantImprovement,
							},
						},
					},
					{
						ComponentReportRowIdentification: rowComponent2Cap22,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.SignificantImprovement,
							},
						},
					},
				},
			},
		},
		{
			name:      "capability page test no significant and missing data",
			generator: capabilityPageGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats90Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent2Cap22Test2,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.NotSignificant,
							},
						},
					},
				},
			},
		},
		{
			name:      "capability page test with both improvement and regression",
			generator: capabilityPageGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats50Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStats50Percent,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent2Cap22Test2,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.SignificantImprovement,
							},
						},
					},
				},
			},
		},
		{
			name:      "test page test no significant and missing data",
			generator: testPageGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats90Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent2Cap22Test2,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVNFull,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDNFull,
								Status:                              apitype.NotSignificant,
							},
						},
					},
				},
			},
		},
		{
			name:      "test page test with both improvement and regression",
			generator: testPageGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats50Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStats50Percent,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent2Cap22Test2,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVNFull,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDNFull,
								Status:                              apitype.SignificantImprovement,
							},
						},
					},
				},
			},
		},
		{
			name: "top page test confidence 90 result in regression",
			generator: componentReportGenerator{
				ComponentReportRequestVariantOptions: apitype.ComponentReportRequestVariantOptions{GroupBy: "cloud,arch,network"},
				ComponentReportRequestAdvancedOptions: apitype.ComponentReportRequestAdvancedOptions{
					Confidence:     90,
					PityFactor:     5,
					MinimumFailure: 3,
				},
			},
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats90Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStats85Percent,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent1,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.SignificantRegression,
								RegressedTests: []apitype.ComponentReportTestSummary{
									{
										ComponentReportTestIdentification: apitype.ComponentReportTestIdentification{
											ComponentReportRowIdentification: apitype.ComponentReportRowIdentification{
												TestName: awsAMD64OVNBaseTestStats90Percent.TestName,
												TestID:   awsAMD64OVNTest.TestID,
											},
											ComponentReportColumnIdentification: apitype.ComponentReportColumnIdentification{
												Platform: columnAWSAMD64OVN.Platform,
												Arch:     columnAWSAMD64OVN.Arch,
												Network:  columnAWSAMD64OVN.Network,
												Upgrade:  awsAMD64OVNTest.Upgrade,
												Variant:  awsAMD64OVNBaseTestStats90Percent.Variants[0],
											},
										},
										Status: apitype.SignificantRegression,
									},
								},
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.MissingBasisAndSample,
							},
						},
					},
					{
						ComponentReportRowIdentification: rowComponent2,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.NotSignificant,
							},
						},
					},
				},
			},
		},
		{
			name: "top page test confidence 90 pity 10 result in no regression",
			generator: componentReportGenerator{
				ComponentReportRequestVariantOptions: apitype.ComponentReportRequestVariantOptions{GroupBy: "cloud,arch,network"},
				ComponentReportRequestAdvancedOptions: apitype.ComponentReportRequestAdvancedOptions{
					Confidence:     90,
					PityFactor:     10,
					MinimumFailure: 3,
				},
			},
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats90Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStats85Percent,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent1,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.NotSignificant,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.MissingBasisAndSample,
							},
						},
					},
					{
						ComponentReportRowIdentification: rowComponent2,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.NotSignificant,
							},
						},
					},
				},
			},
		},
		{
			name:      "top page test minimum failure no regression",
			generator: defaultComponentReportGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNBaseTestStats90Percent,
				awsAMD64SDNTest: awsAMD64SDNBaseTestStats90Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNTest: awsAMD64OVNSampleTestStatsTiny,
				awsAMD64SDNTest: awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: rowComponent1,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.NotSignificant,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.MissingBasisAndSample,
							},
						},
					},
					{
						ComponentReportRowIdentification: rowComponent2,
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVN,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDN,
								Status:                              apitype.NotSignificant,
							},
						},
					},
				},
			},
		},
		{
			name:      "top page test group by variant",
			generator: groupByVariantComponentReportGenerator,
			baseStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNVariantsTest: awsAMD64OVNBaseTestStatsVariants90Percent,
				awsAMD64SDNTest:         awsAMD64SDNBaseTestStats90Percent,
			},
			sampleStatus: map[apitype.ComponentTestIdentification]apitype.ComponentTestStatus{
				awsAMD64OVNVariantsTest: awsAMD64OVNSampleTestStatsVariants90Percent,
				awsAMD64SDNTest:         awsAMD64SDNSampleTestStats90Percent,
			},
			expectedReport: apitype.ComponentReport{
				Rows: []apitype.ComponentReportRow{
					{
						ComponentReportRowIdentification: apitype.ComponentReportRowIdentification{
							Component: "component 1",
						},
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVNVariantFipsStandard,
								Status:                              apitype.NotSignificant,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDNVariantStandard,
								Status:                              apitype.MissingBasisAndSample,
							},
						},
					},
					{
						ComponentReportRowIdentification: apitype.ComponentReportRowIdentification{
							Component: "component 2",
						},
						Columns: []apitype.ComponentReportColumn{
							{
								ComponentReportColumnIdentification: columnAWSAMD64OVNVariantFipsStandard,
								Status:                              apitype.MissingBasisAndSample,
							},
							{
								ComponentReportColumnIdentification: columnAWSAMD64SDNVariantStandard,
								Status:                              apitype.NotSignificant,
							},
						},
					},
				},
			},
		},
	}
	componentAndCapabilityGetter = fakeComponentAndCapabilityGetter
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			report := tc.generator.generateComponentTestReport(tc.baseStatus, tc.sampleStatus, []apitype.TestRegression{})
			assert.Equal(t, tc.expectedReport, report, "expected report %+v, got %+v", tc.expectedReport, report)
		})
	}
}

func TestGenerateComponentTestDetailsReport(t *testing.T) {
	prowJob1 := "ProwJob1"
	prowJob2 := "ProwJob2"
	type testStats struct {
		Success int
		Failure int
		Flake   int
	}
	type requiredJobStats struct {
		job string
		testStats
	}
	baseHighSuccessStats := testStats{
		Success: 1000,
		Failure: 100,
		Flake:   50,
	}
	baseLowSuccessStats := testStats{
		Success: 500,
		Failure: 600,
		Flake:   50,
	}
	sampleHighSuccessStats := testStats{
		Success: 100,
		Failure: 9,
		Flake:   4,
	}
	sampleLowSuccessStats := testStats{
		Success: 50,
		Failure: 59,
		Flake:   4,
	}
	testDetailsRowIdentification := apitype.ComponentReportRowIdentification{
		TestID:     testDetailsGenerator.TestID,
		Component:  testDetailsGenerator.Component,
		Capability: testDetailsGenerator.Capability,
	}
	testDetailsColumnIdentification := apitype.ComponentReportColumnIdentification{
		Network:  testDetailsGenerator.Network,
		Upgrade:  testDetailsGenerator.Upgrade,
		Arch:     testDetailsGenerator.Arch,
		Platform: testDetailsGenerator.Platform,
		Variant:  testDetailsGenerator.Variant,
	}
	sampleReleaseStatsTwoHigh := apitype.ComponentReportTestDetailsReleaseStats{
		Release: testDetailsGenerator.SampleRelease.Release,
		ComponentReportTestDetailsTestStats: apitype.ComponentReportTestDetailsTestStats{
			SuccessRate:  0.9203539823008849,
			SuccessCount: 200,
			FailureCount: 18,
			FlakeCount:   8,
		},
	}
	baseReleaseStatsTwoHigh := apitype.ComponentReportTestDetailsReleaseStats{
		Release: testDetailsGenerator.BaseRelease.Release,
		ComponentReportTestDetailsTestStats: apitype.ComponentReportTestDetailsTestStats{
			SuccessRate:  0.9130434782608695,
			SuccessCount: 2000,
			FailureCount: 200,
			FlakeCount:   100,
		},
	}
	sampleTestStatsHigh := apitype.ComponentReportTestDetailsTestStats{
		SuccessRate:  0.9203539823008849,
		SuccessCount: 100,
		FailureCount: 9,
		FlakeCount:   4,
	}
	baseTestStatsHigh := apitype.ComponentReportTestDetailsTestStats{
		SuccessRate:  0.9130434782608695,
		SuccessCount: 1000,
		FailureCount: 100,
		FlakeCount:   50,
	}
	sampleTestStatsLow := apitype.ComponentReportTestDetailsTestStats{
		SuccessRate:  0.4778761061946903,
		SuccessCount: 50,
		FailureCount: 59,
		FlakeCount:   4,
	}
	baseTestStatsLow := apitype.ComponentReportTestDetailsTestStats{
		SuccessRate:  0.4782608695652174,
		SuccessCount: 500,
		FailureCount: 600,
		FlakeCount:   50,
	}
	sampleReleaseStatsOneHigh := apitype.ComponentReportTestDetailsReleaseStats{
		Release: testDetailsGenerator.SampleRelease.Release,
		ComponentReportTestDetailsTestStats: apitype.ComponentReportTestDetailsTestStats{
			SuccessRate:  0.9203539823008849,
			SuccessCount: 100,
			FailureCount: 9,
			FlakeCount:   4,
		},
	}
	baseReleaseStatsOneHigh := apitype.ComponentReportTestDetailsReleaseStats{
		Release: testDetailsGenerator.BaseRelease.Release,
		ComponentReportTestDetailsTestStats: apitype.ComponentReportTestDetailsTestStats{
			SuccessRate:  0.9130434782608695,
			SuccessCount: 1000,
			FailureCount: 100,
			FlakeCount:   50,
		},
	}
	sampleReleaseStatsOneLow := apitype.ComponentReportTestDetailsReleaseStats{
		Release: testDetailsGenerator.SampleRelease.Release,
		ComponentReportTestDetailsTestStats: apitype.ComponentReportTestDetailsTestStats{
			SuccessRate:  0.4778761061946903,
			SuccessCount: 50,
			FailureCount: 59,
			FlakeCount:   4,
		},
	}
	baseReleaseStatsOneLow := apitype.ComponentReportTestDetailsReleaseStats{
		Release: testDetailsGenerator.BaseRelease.Release,
		ComponentReportTestDetailsTestStats: apitype.ComponentReportTestDetailsTestStats{
			SuccessRate:  0.4782608695652174,
			SuccessCount: 500,
			FailureCount: 600,
			FlakeCount:   50,
		},
	}
	tests := []struct {
		name                    string
		generator               componentReportGenerator
		baseRequiredJobStats    []requiredJobStats
		sampleRequiredJobStats  []requiredJobStats
		expectedReport          apitype.ComponentReportTestDetails
		expectedSampleJobRunLen map[string]int
		expectedBaseJobRunLen   map[string]int
	}{
		{
			name:      "one job with high pass base and sample",
			generator: testDetailsGenerator,
			baseRequiredJobStats: []requiredJobStats{
				{
					job:       prowJob1,
					testStats: baseHighSuccessStats,
				},
			},
			sampleRequiredJobStats: []requiredJobStats{
				{
					job:       prowJob1,
					testStats: sampleHighSuccessStats,
				},
			},
			expectedReport: apitype.ComponentReportTestDetails{
				ComponentReportTestIdentification: apitype.ComponentReportTestIdentification{
					ComponentReportRowIdentification:    testDetailsRowIdentification,
					ComponentReportColumnIdentification: testDetailsColumnIdentification,
				},
				SampleStats:  sampleReleaseStatsOneHigh,
				BaseStats:    baseReleaseStatsOneHigh,
				FisherExact:  0.4807457902463764,
				ReportStatus: apitype.NotSignificant,
				JobStats: []apitype.ComponentReportTestDetailsJobStats{
					{
						JobName:     prowJob1,
						SampleStats: sampleTestStatsHigh,
						BaseStats:   baseTestStatsHigh,
						Significant: false,
					},
				},
			},
			expectedSampleJobRunLen: map[string]int{
				prowJob1: 113,
			},
			expectedBaseJobRunLen: map[string]int{
				prowJob1: 1150,
			},
		},
		{
			name:      "one job with high base and low sample pass rate",
			generator: testDetailsGenerator,
			baseRequiredJobStats: []requiredJobStats{
				{
					job:       prowJob1,
					testStats: baseHighSuccessStats,
				},
			},
			sampleRequiredJobStats: []requiredJobStats{
				{
					job:       prowJob1,
					testStats: sampleLowSuccessStats,
				},
			},
			expectedReport: apitype.ComponentReportTestDetails{
				ComponentReportTestIdentification: apitype.ComponentReportTestIdentification{
					ComponentReportRowIdentification:    testDetailsRowIdentification,
					ComponentReportColumnIdentification: testDetailsColumnIdentification,
				},
				SampleStats:  sampleReleaseStatsOneLow,
				BaseStats:    baseReleaseStatsOneHigh,
				FisherExact:  8.209711662216515e-28,
				ReportStatus: apitype.ExtremeRegression,
				JobStats: []apitype.ComponentReportTestDetailsJobStats{
					{
						JobName:     prowJob1,
						SampleStats: sampleTestStatsLow,
						BaseStats:   baseTestStatsHigh,
						Significant: false,
					},
				},
			},
			expectedSampleJobRunLen: map[string]int{
				prowJob1: 113,
			},
			expectedBaseJobRunLen: map[string]int{
				prowJob1: 1150,
			},
		},
		{
			name:      "one job with low base and high sample pass rate",
			generator: testDetailsGenerator,
			baseRequiredJobStats: []requiredJobStats{
				{
					job:       prowJob1,
					testStats: baseLowSuccessStats,
				},
			},
			sampleRequiredJobStats: []requiredJobStats{
				{
					job:       prowJob1,
					testStats: sampleHighSuccessStats,
				},
			},
			expectedReport: apitype.ComponentReportTestDetails{
				ComponentReportTestIdentification: apitype.ComponentReportTestIdentification{
					ComponentReportRowIdentification:    testDetailsRowIdentification,
					ComponentReportColumnIdentification: testDetailsColumnIdentification,
				},
				SampleStats:  sampleReleaseStatsOneHigh,
				BaseStats:    baseReleaseStatsOneLow,
				FisherExact:  4.911246201592593e-22,
				ReportStatus: apitype.SignificantImprovement,
				JobStats: []apitype.ComponentReportTestDetailsJobStats{
					{
						JobName:     prowJob1,
						SampleStats: sampleTestStatsHigh,
						BaseStats:   baseTestStatsLow,
						Significant: false,
					},
				},
			},
			expectedSampleJobRunLen: map[string]int{
				prowJob1: 113,
			},
			expectedBaseJobRunLen: map[string]int{
				prowJob1: 1150,
			},
		},
		{
			name:      "two jobs with high pass rate",
			generator: testDetailsGenerator,
			baseRequiredJobStats: []requiredJobStats{
				{
					job:       prowJob1,
					testStats: baseHighSuccessStats,
				},
				{
					job:       prowJob2,
					testStats: baseHighSuccessStats,
				},
			},
			sampleRequiredJobStats: []requiredJobStats{
				{
					job:       prowJob1,
					testStats: sampleHighSuccessStats,
				},
				{
					job:       prowJob2,
					testStats: sampleHighSuccessStats,
				},
			},
			expectedReport: apitype.ComponentReportTestDetails{
				ComponentReportTestIdentification: apitype.ComponentReportTestIdentification{
					ComponentReportRowIdentification:    testDetailsRowIdentification,
					ComponentReportColumnIdentification: testDetailsColumnIdentification,
				},
				SampleStats:  sampleReleaseStatsTwoHigh,
				BaseStats:    baseReleaseStatsTwoHigh,
				FisherExact:  0.4119831376606586,
				ReportStatus: apitype.NotSignificant,
				JobStats: []apitype.ComponentReportTestDetailsJobStats{
					{
						JobName:     prowJob1,
						SampleStats: sampleTestStatsHigh,
						BaseStats:   baseTestStatsHigh,
						Significant: false,
					},
					{
						JobName:     prowJob2,
						SampleStats: sampleTestStatsHigh,
						BaseStats:   baseTestStatsHigh,
						Significant: false,
					},
				},
			},
			expectedSampleJobRunLen: map[string]int{
				prowJob1: 113,
				prowJob2: 113,
			},
			expectedBaseJobRunLen: map[string]int{
				prowJob1: 1150,
				prowJob2: 1150,
			},
		},
	}
	componentAndCapabilityGetter = fakeComponentAndCapabilityGetter
	for _, tc := range tests {
		baseStats := map[string][]apitype.ComponentJobRunTestStatusRow{}
		sampleStats := map[string][]apitype.ComponentJobRunTestStatusRow{}
		for _, testStats := range tc.baseRequiredJobStats {
			for i := 0; i < testStats.Success; i++ {
				baseStats[testStats.job] = append(baseStats[testStats.job], apitype.ComponentJobRunTestStatusRow{
					ProwJob:      testStats.job,
					TotalCount:   1,
					SuccessCount: 1,
				})
			}
			for i := 0; i < testStats.Failure; i++ {
				baseStats[testStats.job] = append(baseStats[testStats.job], apitype.ComponentJobRunTestStatusRow{
					ProwJob:    testStats.job,
					TotalCount: 1,
				})
			}
			for i := 0; i < testStats.Flake; i++ {
				baseStats[testStats.job] = append(baseStats[testStats.job], apitype.ComponentJobRunTestStatusRow{
					ProwJob:    testStats.job,
					TotalCount: 1,
					FlakeCount: 1,
				})
			}
		}
		for _, testStats := range tc.sampleRequiredJobStats {
			for i := 0; i < testStats.Success; i++ {
				sampleStats[testStats.job] = append(sampleStats[testStats.job], apitype.ComponentJobRunTestStatusRow{
					ProwJob:      testStats.job,
					TotalCount:   1,
					SuccessCount: 1,
				})
			}
			for i := 0; i < testStats.Failure; i++ {
				sampleStats[testStats.job] = append(sampleStats[testStats.job], apitype.ComponentJobRunTestStatusRow{
					ProwJob:    testStats.job,
					TotalCount: 1,
				})
			}
			for i := 0; i < testStats.Flake; i++ {
				sampleStats[testStats.job] = append(sampleStats[testStats.job], apitype.ComponentJobRunTestStatusRow{
					ProwJob:    testStats.job,
					TotalCount: 1,
					FlakeCount: 1,
				})
			}
		}

		t.Run(tc.name, func(t *testing.T) {
			report := tc.generator.generateComponentTestDetailsReport(baseStats, sampleStats)
			assert.Equal(t, tc.expectedReport.ComponentReportRowIdentification, report.ComponentReportRowIdentification, "expected report row identification %+v, got %+v", tc.expectedReport.ComponentReportRowIdentification, report.ComponentReportRowIdentification)
			assert.Equal(t, tc.expectedReport.ComponentReportColumnIdentification, report.ComponentReportColumnIdentification, "expected report column identification %+v, got %+v", tc.expectedReport.ComponentReportColumnIdentification, report.ComponentReportColumnIdentification)
			assert.Equal(t, tc.expectedReport.BaseStats, report.BaseStats, "expected report base stats %+v, got %+v", tc.expectedReport.BaseStats, report.BaseStats)
			assert.Equal(t, tc.expectedReport.SampleStats, report.SampleStats, "expected report sample stats %+v, got %+v", tc.expectedReport.SampleStats, report.SampleStats)
			assert.Equal(t, tc.expectedReport.FisherExact, report.FisherExact, "expected fisher exact number %+v, got %+v", tc.expectedReport.FisherExact, report.FisherExact)
			assert.Equal(t, tc.expectedReport.ReportStatus, report.ReportStatus, "expected report status %+v, got %+v", tc.expectedReport.ReportStatus, report.ReportStatus)
			assert.Equal(t, len(tc.expectedReport.JobStats), len(report.JobStats), "expected len of job stats %+v, got %+v", len(tc.expectedReport.JobStats), report.JobStats)
			for i := range tc.expectedReport.JobStats {
				jobName := report.JobStats[i].JobName
				assert.Equal(t, tc.expectedReport.JobStats[i].JobName, jobName, "expected job name %+v, got %+v", tc.expectedReport.JobStats[i].JobName, jobName)
				assert.Equal(t, tc.expectedReport.JobStats[i].Significant, report.JobStats[i].Significant, "expected per job significant %+v, got %+v", tc.expectedReport.JobStats[i].Significant, report.JobStats[i].Significant)
				assert.Equal(t, tc.expectedReport.JobStats[i].BaseStats, report.JobStats[i].BaseStats, "expected per job base stats for %s to be %+v, got %+v", tc.expectedReport.JobStats[i].JobName, tc.expectedReport.JobStats[i].BaseStats, report.JobStats[i].BaseStats)
				assert.Equal(t, tc.expectedReport.JobStats[i].SampleStats, report.JobStats[i].SampleStats, "expected per job sample stats for %s to be %+v, got %+v", tc.expectedReport.JobStats[i].JobName, tc.expectedReport.JobStats[i].SampleStats, report.JobStats[i].SampleStats)
				assert.Equal(t, tc.expectedSampleJobRunLen[jobName], len(report.JobStats[i].SampleJobRunStats), "expected sample job run counts %+v, got %+v", tc.expectedSampleJobRunLen[jobName], len(report.JobStats[i].SampleJobRunStats))
				assert.Equal(t, tc.expectedBaseJobRunLen[jobName], len(report.JobStats[i].BaseJobRunStats), "expected base job run counts %+v, got %+v", tc.expectedBaseJobRunLen[jobName], len(report.JobStats[i].BaseJobRunStats))
			}
			// assert.Equal(t, tc.expectedReport.ReportStatus, report.ReportStatus, "expected report %+v, got %+v", tc.expectedReport, report)
			// output, _ := json.MarshalIndent(report, "", "    ")
			// fmt.Printf("-----report \n%+v\n", string(output))
		})
	}
}

func Test_componentReportGenerator_normalizeProwJobName(t *testing.T) {
	tests := []struct {
		name          string
		sampleRelease string
		baseRelease   string
		jobName       string
		want          string
	}{
		{
			name:        "base release is removed",
			baseRelease: "4.16",
			jobName:     "periodic-ci-openshift-release-master-ci-4.16-e2e-azure-ovn-upgrade",
			want:        "periodic-ci-openshift-release-master-ci-X.X-e2e-azure-ovn-upgrade",
		},
		{
			name:          "sample release is removed",
			sampleRelease: "4.16",
			jobName:       "periodic-ci-openshift-release-master-ci-4.16-e2e-azure-ovn-upgrade",
			want:          "periodic-ci-openshift-release-master-ci-X.X-e2e-azure-ovn-upgrade",
		},
		{
			name:    "frequency is removed",
			jobName: "periodic-ci-openshift-release-master-ci-test-job-f27",
			want:    "periodic-ci-openshift-release-master-ci-test-job-fXX",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &componentReportGenerator{}
			if tt.baseRelease != "" {
				c.BaseRelease = apitype.ComponentReportRequestReleaseOptions{Release: tt.baseRelease}
			}
			if tt.sampleRelease != "" {
				c.SampleRelease = apitype.ComponentReportRequestReleaseOptions{Release: tt.sampleRelease}
			}

			assert.Equalf(t, tt.want, c.normalizeProwJobName(tt.jobName), "normalizeProwJobName(%v)", tt.jobName)
		})
	}
}

func Test_componentReportGenerator_assessComponentStatus(t *testing.T) {
	tests := []struct {
		name                   string
		sampleTotal            int
		sampleSuccess          int
		sampleFlake            int
		baseTotal              int
		baseSuccess            int
		baseFlake              int
		numberOfIgnoredSamples int
		expectedStatus         apitype.ComponentReportStatus
		expectedFischers       float64
	}{
		{
			name:                   "triaged still regular regression",
			sampleTotal:            16,
			sampleSuccess:          13,
			sampleFlake:            0,
			baseTotal:              15,
			baseSuccess:            14,
			baseFlake:              1,
			numberOfIgnoredSamples: 2,
			expectedStatus:         -4,
			expectedFischers:       0.4827586206896551,
		},
		{
			name:                   "triaged regular regression",
			sampleTotal:            15,
			sampleSuccess:          13,
			sampleFlake:            0,
			baseTotal:              15,
			baseSuccess:            14,
			baseFlake:              1,
			numberOfIgnoredSamples: 2,
			expectedStatus:         -2,
			expectedFischers:       1,
		},
		{
			name:                   "regular regression",
			sampleTotal:            15,
			sampleSuccess:          13,
			sampleFlake:            0,
			baseTotal:              15,
			baseSuccess:            14,
			baseFlake:              1,
			numberOfIgnoredSamples: 0,
			expectedStatus:         -4,
			expectedFischers:       0.2413793103448262,
		},
		{
			name:                   "zero success",
			sampleTotal:            15,
			sampleSuccess:          0,
			sampleFlake:            0,
			baseTotal:              15,
			baseSuccess:            14,
			baseFlake:              1,
			numberOfIgnoredSamples: 0,
			expectedStatus:         -5,
			expectedFischers:       6.446725037893782e-09,
		},
		{
			name:                   "triaged, zero success",
			sampleTotal:            15,
			sampleSuccess:          0,
			sampleFlake:            0,
			baseTotal:              15,
			baseSuccess:            14,
			baseFlake:              1,
			numberOfIgnoredSamples: 15,
			expectedStatus:         -3,
			expectedFischers:       0,
		},

		{
			name:                   "triaged extreme, fixed",
			sampleTotal:            15,
			sampleSuccess:          5,
			sampleFlake:            0,
			baseTotal:              15,
			baseSuccess:            14,
			baseFlake:              1,
			numberOfIgnoredSamples: 10,
			expectedStatus:         -3,
			expectedFischers:       1,
		},

		{
			name:                   "triaged, still extreme",
			sampleTotal:            15,
			sampleSuccess:          5,
			sampleFlake:            0,
			baseTotal:              15,
			baseSuccess:            14,
			baseFlake:              1,
			numberOfIgnoredSamples: 9,
			expectedStatus:         -5,
			expectedFischers:       0.285714285714284,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &componentReportGenerator{}

			status, fischers := c.assessComponentStatus(tt.sampleTotal, tt.sampleSuccess, tt.sampleFlake, tt.baseTotal, tt.baseSuccess, tt.baseFlake, nil, tt.numberOfIgnoredSamples)
			assert.Equalf(t, tt.expectedStatus, status, "assessComponentStatus expected status not equal")
			assert.Equalf(t, tt.expectedFischers, fischers, "assessComponentStatus expected fischers value not equal")
		})
	}
}
