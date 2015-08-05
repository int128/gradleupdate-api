import gradle.TemplateRepository
import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fromUser = params.from_user
final fromRepo = params.from_repo
final fromBranch = params.from_branch
final intoRepo = params.into_repo
final intoBranch = params.into_branch
final gradleVersion = params.gradle_version
assert fromUser instanceof String
assert fromRepo instanceof String
assert fromBranch instanceof String
assert intoRepo instanceof String
assert intoBranch instanceof String
assert gradleVersion instanceof String

final gitHub = new GitHub()
final templateRepository = new TemplateRepository(gitHub)

log.info("Creating a tree on $fromRepo")
final tree = templateRepository.createTreeWithGradleWrapper(fromRepo)

log.info("Creating a branch $fromBranch on $fromRepo")
gitHub.createBranch(fromRepo, fromBranch, intoBranch, "Gradle $gradleVersion", tree)

log.info("Queue sending a pull request from $fromBranch into $intoRepo:$intoBranch")
defaultQueue.add(
        url: relativePath(request, '5-pull-request.groovy'),
        params: [
                from_user: fromUser,
                from_branch: fromBranch,
                into_repo: intoRepo,
                into_branch: intoBranch,
                gradle_version: gradleVersion,
        ],
        countdownMillis: 1000)