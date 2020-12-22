import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";
import { getType } from "typesafe-actions";
import actions from ".";
import Config from "../shared/Config";

const mockStore = configureMockStore([thunk]);

let store: any;
const testConfig = "have you tried to turn it off and on again";
const testControllerAppVersion = "0.1"

beforeEach(() => {
  Config.getConfig = jest.fn(() => testConfig);
  Config.getControllerAppVersion = jest.fn(() => testControllerAppVersion)
  store = mockStore();
});

describe("getConfig", () => {
  it("dispatches request config and its returned value", async () => {
    const expectedActions = [
      {
        type: getType(actions.config.requestConfig),
      },
      {
        payload: testConfig,
        type: getType(actions.config.receiveConfig),
      },
      {
        payload: testControllerAppVersion,
        type: getType(actions.config.receiveControllerAppVersion),
      },
    ];

    await store.dispatch(actions.config.getConfig());
    expect(store.getActions()).toEqual(expectedActions);
  });
});
