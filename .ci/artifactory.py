import requests
import ci.util


class ArtifactoryRepoClient:
    def __init__(self, base_url: str, user: str, password: str):
        self.user = user
        self.password = password
        self.routes = ArtifactoryRoutes(base_url=base_url)

    def request(self, method: str, *args, **kwargs):

        res = requests.request(
            method=method,
            auth=(self.user, self.password),
            *args,
            **kwargs,
        )

        if not res.ok:
            print(res.text)
            raise RuntimeError(f'{method} request to url {res.url} '
                               f'failed with {res.status_code=} {res.reason=}'
                               )
        return res

    def upload_artifact(self, artifact_name: str, artifact):
        res = self.request(
            method='PUT',
            url=self.routes.deploy_artifact(artifact_name),
            data=artifact.read(),
        )
        return res


class ArtifactoryRoutes:
    def __init__(self, base_url: str):
        self.base_url = base_url

    def _api_url(self, *args, **kwargs):
        return ci.util.urljoin(self.base_url, *args)

    def deploy_artifact(self, artifact_name: str):
        return self._api_url(artifact_name)
