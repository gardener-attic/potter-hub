import axios from "axios";

import { UI_BACKEND_ROOT_URL } from "./Kube";

export interface IConfig {
  namespace: string;
  appVersion: string;
  error?: Error;
  featuredChartIds?: string[];
  defaultRepo?: string;
  generalRepoInfo?: string;
  staticCatalogInfo?: string;
  appName: string;
  footer?: IFooterConfig;
  header?: IHeaderConfig;
}

export interface IHeaderConfig {
  helpMenu?: IHelpMenuConfig;
}

export interface IHelpMenuConfig {
  links: ILink[];
}

export interface IFooterConfig {
  sectionOne?: IFooterSectionConfig;
  sectionTwo?: IFooterSectionConfig;
}

export interface IFooterSectionConfig {
  title: string;
  links: ILink[];
}

export interface ILink {
  title: string;
  href: string;
}

export default class Config {
  public static async getConfig() {
    const url = Config.ConfigAPIEndpoint;
    const { data } = await axios.get<IConfig>(url);

    // Development environment config overrides
    // TODO(miguel) Rename env variable to KUBEAPPS_NAMESPACE once/if we eject create-react-app
    // Currently we are using REACT_APP_* because it's the only way to inject env variables in a sealed setup.
    // Please note that this env variable gets mapped in the run command in the package.json file
    if (process.env.NODE_ENV !== "production") {
      data.appName = "Potter"
      if (process.env.REACT_APP_KUBEAPPS_NS) {
        data.namespace = process.env.REACT_APP_KUBEAPPS_NS;
      }
      data.featuredChartIds = [
        "bitnami/mongodb",
        "bitnami/apache",
      ];
      data.generalRepoInfo =
        "This is a test message for 'All repositories'.<br/>This should appear on a new line.";
      data.staticCatalogInfo = "Test info message for the static catalog.<br/>This should appear on a new line.";
      data.footer = {
        sectionOne: {
          title: "Section 1",
          links: [
            {
              title: "link-1-title",
              href: ""
            }
          ],
        },
        sectionTwo: {
          title: "Section 2",
          links: [
            {
              title: "link-2-title",
              href: ""
            }
          ],
        }
      }
      data.header = {
        helpMenu: {
          links: [
            {
              title: "link-3-title",
              href: ""
            },
            {
              title: "link-4-title",
              href: ""
            }
          ]
        }
      }
    }

    return data;
  }

  public static async getControllerAppVersion() {
    const url = Config.ControllerAppVersionAPIEndpoint;
    const { data } = await axios.get<string>(url);
    return data
  }

  private static ConfigAPIEndpoint: string = "config.json";
  private static ControllerAppVersionAPIEndpoint: string = UI_BACKEND_ROOT_URL + "/controller-version"
}
