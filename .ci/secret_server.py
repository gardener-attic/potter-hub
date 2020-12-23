import util
import kube.ctx
from pprint import pprint

class SecretServer:
    def __init__(self, cfg_set_name: str = 'hub'):
        self.cfg_set = util.ctx().cfg_factory().cfg_set(cfg_set_name)

    def get_landscape_config(self, cfg_name: str):
        return self.cfg_set.hub(cfg_name)

    def get_repository_config(self):
        return self.cfg_set.hub_repositories("repositories")

    def get_kubeconfig(self, cfg_name: str):
        return self.cfg_set.kubernetes(cfg_name).kubeconfig()

    def get_container_registry_creds(self, cdf_name: str):
        return self.cfg_set.container_registry(cdf_name).credentials()

    def get_kube_client(self, kubecfg_name: str):
        return kube.ctx.Ctx(self.cfg_set.kubernetes(kubecfg_name).kubeconfig())