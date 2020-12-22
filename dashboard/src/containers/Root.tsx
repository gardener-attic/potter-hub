import { ConnectedRouter } from "connected-react-router";
import * as React from "react";
import { Provider } from "react-redux";

import Layout from "../components/Layout";
import store, { history } from "../store";
import ConfigLoaderContainer from "./ConfigLoaderContainer";
import HeaderContainer from "./HeaderContainer";
import Routes from "./RoutesContainer";
import TabTitle from "./TabTitleContainer"
import URLParamsProvider from "./URLParamsProviderContainer"

class Root extends React.Component {
  public render() {
    return (
      <Provider store={store}>
        <URLParamsProvider>
          <TabTitle>
            <ConfigLoaderContainer>
              <ConnectedRouter history={history}>
                <Layout headerComponent={HeaderContainer}>
                  <Routes />
                </Layout>
              </ConnectedRouter>
            </ConfigLoaderContainer>
          </TabTitle>
        </URLParamsProvider>
      </Provider>
    );
  }
}

export default Root;
