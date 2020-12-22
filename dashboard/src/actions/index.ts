import { push } from "connected-react-router";

import * as apps from "./apps";
import * as auth from "./auth";
import * as charts from "./charts";
import * as clusterBom from "./clusterbom";
import * as config from "./config";
import * as kube from "./kube";
import * as namespace from "./namespace";
import * as repos from "./repos";
import * as serviceCatalog from "./serviceCatalog";

export default {
  apps,
  auth,
  clusterBom,
  serviceCatalog,
  charts,
  config,
  kube,
  namespace,
  repos,
  shared: {
    pushSearchFilter: (f: string) => push(`?q=${f}`),
  },
};
