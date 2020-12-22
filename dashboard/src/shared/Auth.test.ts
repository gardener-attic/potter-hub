import Axios from "axios";
import * as jwt from "jsonwebtoken";
import { Auth } from "./Auth";

describe("Auth", () => {
  beforeEach(() => {
    Axios.get = jest.fn();
  });
  afterEach(() => {
    jest.resetAllMocks();
  });
  it("should get an URL with the given token", async () => {
    const mock = jest.fn();
    Axios.get = mock;
    await Auth.validateToken("foo");
    expect(mock.mock.calls[0]).toEqual([
      "api/ui-backend/reverse-proxy/v1/",
      { headers: { Authorization: "Bearer foo" } },
    ]);
  });

  describe("when there is an error", () => {
    [
      {
        name: "should throw an invalid token error for 401 responses",
        response: { status: 401, data: "ignored anyway" },
        expectedError: new Error("invalid token"),
      },
      {
        name: "should throw a standard error for a 404 response",
        response: { status: 404, data: "Not found" },
        expectedError: new Error("404: Not found"),
      },
      {
        name: "should throw a standard error for a 500 response",
        response: { status: 500, data: "Server exception" },
        expectedError: new Error("500: Server exception"),
      },
      {
        name: "should throw a standard error for a 403 response",
        response: { status: 403, data: "Not Allowed" },
        expectedError: new Error("token not authorised for cluster"),
      },
    ].forEach(testCase => {
      it(testCase.name, async () => {
        const mock = jest.fn(() => {
          return Promise.reject({ response: testCase.response });
        });
        Axios.get = mock;
        // TODO(absoludity): tried using `expect(fn()).rejects.toThrow()` but it seems we need
        // to upgrade jest for `toThrow()` to work with async.
        let err = null;
        try {
          await Auth.validateToken("foo");
        } catch (e) {
          err = e;
        } finally {
          expect(err).toEqual(testCase.expectedError);
        }
      });
    });
  });

  describe("fetchOIDCToken", () => {
    it("should fetch a token", async () => {
      Axios.head = jest.fn(path => {
        expect(path).toEqual("");
        return { headers: { authorization: "Bearer foo" } };
      });
      const token = await Auth.fetchOIDCToken();
      expect(token).toEqual("foo");
    });
  });
  it("should not return a token if the info is not present", async () => {
    Axios.head = jest.fn(() => {
      return {};
    });
    const token = await Auth.fetchOIDCToken();
    expect(token).toEqual(null);
  });
  it("should not return a token if the call fails", async () => {
    Axios.head = jest.fn(() => {
      throw new Error();
    });
    const token = await Auth.fetchOIDCToken();
    expect(token).toEqual(null);
  });

  describe("defaultNamespaceFromToken", () => {
    const customNamespace = "kubeapps-user";

    it("should return the k8s namespace from the jwt token", () => {
      const token = jwt.sign(
        {
          iss: "kubernetes/serviceaccount",
          "kubernetes.io/serviceaccount/namespace": customNamespace,
        },
        "secret",
      );

      const defaultNamespace = Auth.defaultNamespaceFromToken(token);

      expect(defaultNamespace).toEqual(customNamespace);
    });

    it("should return default if the namespace is not present", () => {
      const token = jwt.sign(
        {
          iss: "kubernetes/serviceaccount",
        },
        "secret",
      );

      const defaultNamespace = Auth.defaultNamespaceFromToken(token);

      expect(defaultNamespace).toEqual("_all");
    });

    it("should return default if the token cannot be decoded", () => {
      const token = "not a jwt token";

      const defaultNamespace = Auth.defaultNamespaceFromToken(token);

      expect(defaultNamespace).toEqual("_all");
    });
  });
});
