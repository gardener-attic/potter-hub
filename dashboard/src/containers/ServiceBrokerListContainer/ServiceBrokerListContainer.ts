import { connect } from "react-redux";
import { Action } from "redux";
import { ThunkDispatch } from "redux-thunk";

import actions from "../../actions";
import ServiceBrokerList from "../../components/Config/ServiceBrokerList";
import { IServiceBroker } from "../../shared/ServiceCatalog";
import { IStoreState } from "../../shared/types";

function mapStateToProps({ serviceCatalog: catalog }: IStoreState) {
  return {
    brokers: catalog.brokers,
    errors: catalog.errors,
    isInstalled: catalog.isServiceCatalogInstalled,
  };
}

function mapDispatchToProps(dispatch: ThunkDispatch<IStoreState, null, Action>) {
  return {
    checkCatalogInstalled: async () => {
      dispatch(actions.serviceCatalog.checkCatalogInstalled());
    },
    getBrokers: () => dispatch(actions.serviceCatalog.getBrokers()),
    sync: (broker: IServiceBroker) => dispatch(actions.serviceCatalog.sync(broker)),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceBrokerList);
