package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
	"github.com/pkg/errors"
)

func TestSendUpdate_Do(t *testing.T) {
	ctx := context.Background()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}
	fixedTime := &gatewaysTestDoubles.FixedTime{
		NowValue: time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC),
	}
	readmeContent := git.FileContent("![Gradle Status](https://gradleupdate.appspot.com/owner/repo/status.svg)")

	t.Run("SuccessfullyUpdated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gatewaysTestDoubles.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(readmeContent, nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendPullRequest := usecasesTestDoubles.NewMockSendPullRequest(ctrl)
		sendPullRequest.EXPECT().Do(ctx, usecases.SendPullRequestRequest{
			Base:           repositoryID,
			HeadBranchName: "gradle-5.0-owner",
			CommitMessage:  "Gradle 5.0",
			CommitFiles: []git.File{{
				Path:    gradle.WrapperPropertiesPath,
				Content: testdata.GradleWrapperProperties50,
			}},
			Title: "Gradle 5.0",
			Body: `Gradle 5.0 is available.

This is sent by @gradleupdate. See https://gradleupdate.appspot.com/owner/repo/status for more.`,
		}).Return(nil)

		u := SendUpdate{
			RepositoryRepository:    repositoryRepository,
			GradleReleaseRepository: gradleService,
			SendPullRequest:         sendPullRequest,
			Time:                    fixedTime,
		}
		err := u.Do(ctx, repositoryID)
		if err != nil {
			t.Fatalf("error while Do: %+v", err)
		}
	})

	t.Run("AlreadyHasLatestGradle", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gatewaysTestDoubles.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(readmeContent, nil)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "4.10.2"}, nil)

		sendPullRequest := usecasesTestDoubles.NewMockSendPullRequest(ctrl)
		u := SendUpdate{
			RepositoryRepository:    repositoryRepository,
			GradleReleaseRepository: gradleService,
			SendPullRequest:         sendPullRequest,
			Time:                    fixedTime,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.AlreadyHasLatestGradle {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.AlreadyHasLatestGradle, preconditionViolation)
		}
	})

	t.Run("NoGradleVersion/NoGradleWrapperProperties", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gatewaysTestDoubles.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(readmeContent, nil).MaxTimes(1)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(nil, &gatewaysTestDoubles.NoSuchEntityError{})

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil).MaxTimes(1)

		sendPullRequest := usecasesTestDoubles.NewMockSendPullRequest(ctrl)
		u := SendUpdate{
			RepositoryRepository:    repositoryRepository,
			GradleReleaseRepository: gradleService,
			SendPullRequest:         sendPullRequest,
			Time:                    fixedTime,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.NoGradleWrapperProperties {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoGradleWrapperProperties, preconditionViolation)
		}
	})

	t.Run("NoGradleVersion/NoGradleVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gatewaysTestDoubles.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(readmeContent, nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(git.FileContent("INVALID"), nil)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendPullRequest := usecasesTestDoubles.NewMockSendPullRequest(ctrl)
		u := SendUpdate{
			RepositoryRepository:    repositoryRepository,
			GradleReleaseRepository: gradleService,
			SendPullRequest:         sendPullRequest,
			Time:                    fixedTime,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.NoGradleVersion {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoGradleVersion, preconditionViolation)
		}
	})

	t.Run("NoReadmeBadge/NoReadme", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gatewaysTestDoubles.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(nil, &gatewaysTestDoubles.NoSuchEntityError{})
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil).MaxTimes(1)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil).MaxTimes(1)

		sendPullRequest := usecasesTestDoubles.NewMockSendPullRequest(ctrl)
		u := SendUpdate{
			RepositoryRepository:    repositoryRepository,
			GradleReleaseRepository: gradleService,
			SendPullRequest:         sendPullRequest,
			Time:                    fixedTime,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.NoReadme {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoReadme, preconditionViolation)
		}
	})

	t.Run("NoReadmeBadge/NoReadmeBadge", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gatewaysTestDoubles.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(git.FileContent("INVALID"), nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendPullRequest := usecasesTestDoubles.NewMockSendPullRequest(ctrl)
		u := SendUpdate{
			RepositoryRepository:    repositoryRepository,
			GradleReleaseRepository: gradleService,
			SendPullRequest:         sendPullRequest,
			Time:                    fixedTime,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.NoReadmeBadge {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoReadmeBadge, preconditionViolation)
		}
	})
}
