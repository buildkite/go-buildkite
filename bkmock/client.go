package bkmock

import "github.com/buildkite/go-buildkite/v4"

func NewClient() buildkite.Client {
	return buildkite.Client{
		AccessTokens:             &MockAccessTokens{},
		Agents:                   &MockAgents{},
		Annotations:              &MockAnnotations{},
		Artifacts:                &MockArtifacts{},
		Builds:                   &MockBuilds{},
		ClusterQueues:            &MockClusterQueues{},
		ClusterTokens:            &MockClusterTokens{},
		Clusters:                 &MockClusters{},
		FlakyTests:               &MockFlakyTests{},
		Jobs:                     &MockJobs{},
		Organizations:            &MockOrganizations{},
		PackageRegistriesService: &MockPackageRegistries{},
		PackagesService:          &MockPackages{},
		PipelineTemplates:        &MockPipelineTemplates{},
		Pipelines:                &MockPipelines{},
		User:                     &MockUsers{},
		Teams:                    &MockTeams{},
		TestRuns:                 &MockTestRuns{},
		TestSuites:               &MockTestSuites{},
		Tests:                    &MockTests{},
	}
}
