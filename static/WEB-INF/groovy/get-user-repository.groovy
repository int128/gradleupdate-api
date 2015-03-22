import groovy.json.JsonBuilder
import infrastructure.GitHub
import service.GitHubRepositoryService
import util.CrossOriginPolicy

import static com.google.appengine.api.utils.SystemProperty.Environment.Value.Development

CrossOriginPolicy.allowOrigin(response, headers)

assert params.fullName
assert app.env.name == Development || headers.Authorization

final gitHub = GitHub.authorizationOrDefault(headers.Authorization)
final service = new GitHubRepositoryService(gitHub)
final entity = service.query(params.fullName)

if (entity) {
    response.contentType = 'application/json'
    println new JsonBuilder(entity)
} else {
    response.sendError(404, 'No Admin Permission')
}
