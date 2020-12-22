import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";
import { getType } from "typesafe-actions";
import Namespace from "../shared/Namespace";
import { errorNamespaces, fetchNamespaces, receiveNamespaces, requestNamespaces, setNamespace } from "./namespace";

const mockStore = configureMockStore([thunk]);

let store: any;

beforeEach(() => {
  store = mockStore();
});

// Regular action creators
// Regular action creators
interface ITestCase {
  name: string;
  action: (...args: any[]) => any;
  args?: any;
  payload?: any;
}

const actionTestCases: ITestCase[] = [
  { name: "setNamespace", action: setNamespace, args: "jack", payload: "jack" },
  {
    name: "receiveNamespces",
    action: receiveNamespaces,
    args: ["jack", "danny"],
    payload: ["jack", "danny"],
  },
];

actionTestCases.forEach(tc => {
  describe(tc.name, () => {
    it("has expected structure", () => {
      expect(tc.action.call(null, tc.args)).toEqual({
        type: getType(tc.action),
        payload: tc.payload,
      });
    });
  });
});

// Async action creators
describe("fetchNamespaces", () => {
  it("dispatches the list of namespace names if no error", async () => {
    Namespace.list = jest.fn().mockImplementationOnce(() => {
      return {
        items: [{ metadata: { name: "overlook-hotel" } }, { metadata: { name: "room-217" } }],
      };
    });
    const expectedActions = [
      {
        type: getType(requestNamespaces),
      },
      {
        type: getType(receiveNamespaces),
        payload: ["overlook-hotel", "room-217"],
      },
    ];

    await store.dispatch(fetchNamespaces());
    expect(store.getActions()).toEqual(expectedActions);
  });

  it("dispatches errorNamespaces if error listing namespaces", async () => {
    const err = new Error("Bang!");
    Namespace.list = jest.fn().mockImplementationOnce(() => Promise.reject(err));
    const expectedActions = [
      {
        type: getType(requestNamespaces),
      },
      {
        type: getType(errorNamespaces),
        payload: { err, op: "list" },
      },
      {
        type: getType(setNamespace),
        payload: "",
      },
    ];

    await store.dispatch(fetchNamespaces());

    expect(store.getActions()).toEqual(expectedActions);
  });
});
