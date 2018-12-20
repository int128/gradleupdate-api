package repositories

import (
	"context"
	"encoding/base64"
	"github.com/int128/gradleupdate/infrastructure"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
)

type Repository struct{}

func (r *Repository) Get(ctx context.Context, id domain.RepositoryIdentifier) (domain.Repository, error) {
	client := infrastructure.GitHubClient(ctx)
	repository, resp, err := client.Repositories.Get(ctx, id.Owner, id.Name)
	if resp != nil && resp.StatusCode == 404 {
		return domain.Repository{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return domain.Repository{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.Repository{
		RepositoryIdentifier: domain.RepositoryIdentifier{
			Owner: repository.GetOwner().GetLogin(),
			Name:  repository.GetName(),
		},
		Description: repository.GetDescription(),
		AvatarURL:   repository.GetOwner().GetAvatarURL(),
		DefaultBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{
				Owner: repository.GetOwner().GetLogin(),
				Name:  repository.GetName(),
			},
			Name: repository.GetDefaultBranch(),
		},
	}, nil
}

func (r *Repository) GetFile(ctx context.Context, id domain.RepositoryIdentifier, path string) (domain.File, error) {
	client := infrastructure.GitHubClient(ctx)
	fc, _, resp, err := client.Repositories.GetContents(ctx, id.Owner, id.Name, path, nil)
	if resp != nil && resp.StatusCode == 404 {
		return domain.File{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return domain.File{}, errors.Wrapf(err, "GitHub API returned error")
	}
	if fc == nil {
		return domain.File{}, errors.Errorf("Expected file but found directory %s", path)
	}
	var content []byte
	switch fc.GetEncoding() {
	case "base64":
		buf := make([]byte, base64.StdEncoding.DecodedLen(len(*fc.Content)))
		n, err := base64.StdEncoding.Decode(buf, []byte(*fc.Content))
		if err != nil {
			return domain.File{}, errors.Wrapf(err, "Could not decode content")
		}
		content = buf[:n]
	default:
		content = []byte(*fc.Content)
	}
	return domain.File{
		Path:    path,
		Content: content,
	}, nil
}

func (r *Repository) Fork(ctx context.Context, id domain.RepositoryIdentifier) (domain.Repository, error) {
	client := infrastructure.GitHubClient(ctx)
	fork, resp, err := client.Repositories.CreateFork(ctx, id.Owner, id.Name, &github.RepositoryCreateForkOptions{})
	if resp != nil && resp.StatusCode == 404 {
		return domain.Repository{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			// Fork in progress
		} else {
			return domain.Repository{}, errors.Wrapf(err, "GitHub API returned error")
		}
	}
	return domain.Repository{
		RepositoryIdentifier: domain.RepositoryIdentifier{
			Owner: fork.GetOwner().GetLogin(),
			Name:  fork.GetName(),
		},
		Description: fork.GetDescription(),
		AvatarURL:   fork.GetOwner().GetAvatarURL(),
		DefaultBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{
				Owner: fork.GetOwner().GetLogin(),
				Name:  fork.GetName(),
			},
			Name: fork.GetDefaultBranch(),
		},
	}, nil
}
