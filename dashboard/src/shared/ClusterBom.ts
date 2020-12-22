import { axiosWithAuth } from "./AxiosInstance";
import { CLUSTERBOM_URL } from "./Kube";
import { IClusterBom, IClusterBomList } from "./types";

export class ClusterBom {
  public static async list() {
    const { data } = await axiosWithAuth.get<IClusterBomList>(CLUSTERBOM_URL);
    return data;
  }

  public static async update(clusterBom: IClusterBom) {
    await axiosWithAuth.put(
      CLUSTERBOM_URL + "/" + clusterBom.metadata.name,
      JSON.stringify(clusterBom),
    );
  }

}
