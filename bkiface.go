package buildkite

import (
	"context"
	"io"
)

//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -destination=bkmock/services.go -source=bkiface.go -package=bkmock

type AccessTokens interface {
	Get(context.Context) (AccessToken, *Response, error)
	Revoke(context.Context) (*Response, error)
}

type Agents interface {
	Create(context.Context, string, Agent) (Agent, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	Get(context.Context, string, string) (Agent, *Response, error)
	List(context.Context, string, *AgentListOptions) ([]Agent, *Response, error)
	Stop(context.Context, string, string, bool) (*Response, error)
}

type Annotations interface {
	Create(context.Context, string, string, string, AnnotationCreate) (Annotation, *Response, error)
	ListByBuild(context.Context, string, string, string, *AnnotationListOptions) ([]Annotation, *Response, error)
}

type Artifacts interface {
	DownloadArtifactByURL(context.Context, string, io.Writer) (*Response, error)
	ListByBuild(context.Context, string, string, string, *ArtifactListOptions) ([]Artifact, *Response, error)
	ListByJob(context.Context, string, string, string, string, *ArtifactListOptions) ([]Artifact, *Response, error)
}

type Builds interface {
	Cancel(context.Context, string, string, string) (Build, error)
	Create(context.Context, string, string, CreateBuild) (Build, *Response, error)
	Get(context.Context, string, string, string, *BuildsListOptions) (Build, *Response, error)
	List(context.Context, *BuildsListOptions) ([]Build, *Response, error)
	ListByOrg(context.Context, string, *BuildsListOptions) ([]Build, *Response, error)
	ListByPipeline(context.Context, string, string, *BuildsListOptions) ([]Build, *Response, error)
	Rebuild(context.Context, string, string, string) (Build, error)
}

type ClusterQueues interface {
	Create(context.Context, string, string, ClusterQueueCreate) (ClusterQueue, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
	Get(context.Context, string, string, string) (ClusterQueue, *Response, error)
	List(context.Context, string, string, *ClusterQueuesListOptions) ([]ClusterQueue, *Response, error)
	Pause(context.Context, string, string, string, ClusterQueuePause) (ClusterQueue, *Response, error)
	Resume(context.Context, string, string, string) (*Response, error)
	Update(context.Context, string, string, string, ClusterQueueUpdate) (ClusterQueue, *Response, error)
}

type ClusterTokens interface {
	Create(context.Context, string, string, ClusterTokenCreateUpdate) (ClusterToken, *Response, error)
	Delete(context.Context, string, string, string) (*Response, error)
	Get(context.Context, string, string, string) (ClusterToken, *Response, error)
	List(context.Context, string, string, *ClusterTokensListOptions) ([]ClusterToken, *Response, error)
	Update(context.Context, string, string, string, ClusterTokenCreateUpdate) (ClusterToken, *Response, error)
}

type Clusters interface {
	Create(context.Context, string, ClusterCreate) (Cluster, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	Get(context.Context, string, string) (Cluster, *Response, error)
	List(context.Context, string, *ClustersListOptions) ([]Cluster, *Response, error)
	Update(context.Context, string, string, ClusterUpdate) (Cluster, *Response, error)
}

type FlakyTests interface {
	List(context.Context, string, string, *FlakyTestsListOptions) ([]FlakyTest, *Response, error)
}

type Jobs interface {
	GetJobEnvironmentVariables(context.Context, string, string, string, string) (JobEnvs, *Response, error)
	GetJobLog(context.Context, string, string, string, string) (JobLog, *Response, error)
	RetryJob(context.Context, string, string, string, string) (Job, *Response, error)
	UnblockJob(context.Context, string, string, string, string, *JobUnblockOptions) (Job, *Response, error)
}

type Organizations interface {
	Get(context.Context, string) (Organization, *Response, error)
	List(context.Context, *OrganizationListOptions) ([]Organization, *Response, error)
}

type PackageRegistries interface {
	Create(context.Context, string, CreatePackageRegistryInput) (PackageRegistry, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	Get(context.Context, string, string) (PackageRegistry, *Response, error)
	List(context.Context, string) ([]PackageRegistry, *Response, error)
	Update(context.Context, string, string, UpdatePackageRegistryInput) (PackageRegistry, *Response, error)
}

type Packages interface {
	Create(context.Context, string, string, CreatePackageInput) (Package, *Response, error)
	Get(context.Context, string, string, string) (Package, *Response, error)
	RequestPresignedUpload(context.Context, string, string) (*PackagePresignedUpload, *Response, error)
}

type PipelineTemplates interface {
	Create(context.Context, string, PipelineTemplateCreateUpdate) (PipelineTemplate, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	Get(context.Context, string, string) (PipelineTemplate, *Response, error)
	List(context.Context, string, *PipelineTemplateListOptions) ([]PipelineTemplate, *Response, error)
	Update(context.Context, string, string, PipelineTemplateCreateUpdate) (PipelineTemplate, *Response, error)
}

type Pipelines interface {
	AddWebhook(context.Context, string, string) (*Response, error)
	Archive(context.Context, string, string) (*Response, error)
	Create(context.Context, string, CreatePipeline) (Pipeline, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	Get(context.Context, string, string) (Pipeline, *Response, error)
	List(context.Context, string, *PipelineListOptions) ([]Pipeline, *Response, error)
	Unarchive(context.Context, string, string) (*Response, error)
	Update(context.Context, string, string, UpdatePipeline) (Pipeline, *Response, error)
}

type Users interface {
	CurrentUser(context.Context) (User, *Response, error)
}

type Teams interface {
	List(context.Context, string, *TeamsListOptions) ([]Team, *Response, error)
}

type TestRuns interface {
	Get(context.Context, string, string, string) (TestRun, *Response, error)
	List(context.Context, string, string, *TestRunsListOptions) ([]TestRun, *Response, error)
}

type TestSuites interface {
	Create(context.Context, string, TestSuiteCreate) (TestSuite, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	Get(context.Context, string, string) (TestSuite, *Response, error)
	List(context.Context, string, *TestSuiteListOptions) ([]TestSuite, *Response, error)
	Update(context.Context, string, string, TestSuite) (TestSuite, *Response, error)
}

type Tests interface {
	Get(context.Context, string, string, string) (Test, *Response, error)
}
