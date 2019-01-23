package gateways

import (
	"context"
	"net/http"
	"time"

	"github.com/int128/gradleupdate/domain"
)

//go:generate mockgen -destination mock_gateways/badge_last_access.go -package mock_gateways github.com/int128/gradleupdate/gateways/interfaces BadgeLastAccessRepository

type BadgeLastAccessRepository interface {
	Save(ctx context.Context, a domain.BadgeLastAccess) error
	FindBySince(ctx context.Context, since time.Time) ([]domain.BadgeLastAccess, error)
}

//go:generate mockgen -destination mock_gateways/repo_last_scan.go -package mock_gateways github.com/int128/gradleupdate/gateways/interfaces RepositoryLastScanRepository

type RepositoryLastScanRepository interface {
	Save(ctx context.Context, a domain.RepositoryLastScan) error
}

type PushBranchRequest struct {
	BaseBranch    domain.Branch
	HeadBranch    domain.BranchID
	CommitMessage string
	CommitFiles   []domain.File
}

type GitService interface {
	CreateBranch(ctx context.Context, req PushBranchRequest) (*domain.Branch, error)
	UpdateForceBranch(ctx context.Context, req PushBranchRequest) (*domain.Branch, error)
}

//go:generate mockgen -destination mock_gateways/gradle.go -package mock_gateways github.com/int128/gradleupdate/gateways/interfaces GradleService

type GradleService interface {
	GetCurrentVersion(ctx context.Context) (domain.GradleVersion, error)
}

//go:generate mockgen -destination mock_gateways/pull_request.go -package mock_gateways github.com/int128/gradleupdate/gateways/interfaces PullRequestRepository

type PullRequestRepository interface {
	Create(ctx context.Context, pull domain.PullRequest) (*domain.PullRequest, error)
	FindByBranch(ctx context.Context, baseBranch, headBranch domain.BranchID) (*domain.PullRequest, error)
}

//go:generate mockgen -destination mock_gateways/repo.go -package mock_gateways github.com/int128/gradleupdate/gateways/interfaces RepositoryRepository

type RepositoryRepository interface {
	Get(context.Context, domain.RepositoryID) (*domain.Repository, error)
	GetFileContent(context.Context, domain.RepositoryID, string) (domain.FileContent, error)
	GetReadme(ctx context.Context, id domain.RepositoryID) (domain.FileContent, error)
	Fork(context.Context, domain.RepositoryID) (*domain.Repository, error)
	GetBranch(ctx context.Context, branch domain.BranchID) (*domain.Branch, error)
}

type ResponseCacheRepository interface {
	Find(ctx context.Context, req *http.Request) (*http.Response, error)
	Save(ctx context.Context, req *http.Request, resp *http.Response) error
	Remove(ctx context.Context, req *http.Request) error
}

//go:generate mockgen -destination mock_gateways/errors.go -package mock_gateways github.com/int128/gradleupdate/gateways/interfaces RepositoryError

type RepositoryError interface {
	error
	NoSuchEntity() bool
}
