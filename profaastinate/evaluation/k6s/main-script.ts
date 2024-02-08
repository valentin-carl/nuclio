import http from "k6/http";
import {check, sleep} from "k6";
import {StatusCodes} from "http-status-codes";

/**
 * Enum for the endpoints used in the application
 */
export enum Endpoint {
  INVOCATION_URL = "http://localhost:8070/api/function_invocations",
  EVALUATION_URL = "http://localhost:8888/evaluation/headers",
}

export enum FunctionName {
  HELLO_1 = "hello-1",
}

export enum StatusCode {
    OK = 204
}

/**
 * Enum for Nuclio header keys
 */
export enum HeaderKey {
  FUNCTION_NAME = "x-nuclio-function-name",
  FUNCTION_NAMESPACE = "x-nuclio-function-namespace",
  ASYNC_DEADLINE = "x-profaastinate-process-deadline",
}

export function buildHeader(functionName: string, deadline?: number, namespace = 'default'): { [name: string]: string } {
    let headers: { [name: string]: string } = {};
    headers[HeaderKey.FUNCTION_NAME] = functionName;
    headers[HeaderKey.FUNCTION_NAMESPACE] = namespace;
    if (deadline) {
        headers[HeaderKey.ASYNC_DEADLINE] = `${deadline}`;
    }
    return headers;
}

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 10,
      timeUnit: "1s", // 1000 iterations per second, i.e. 1000 RPS
      duration: "20s",
      preAllocatedVUs: 1, // the size of the VU (virtual user) pool for the executor
      maxVUs: 10, // the maximum number of VUs the executor is allowed to scale to
    },
  },
};

function callAsyncHello1() {
  const plainFunctionHeader = buildHeader(
    FunctionName.HELLO_1,
    10000
  );
  const res = http.get(Endpoint.INVOCATION_URL, {
    headers: plainFunctionHeader,
  });
  check(res, { "status was 204": (r) => r.status == StatusCode.OK });
}

function callSyncHello1() {
    const plainFunctionHeader = buildHeader(
        FunctionName.HELLO_1
    );
    const res = http.get(Endpoint.INVOCATION_URL, {
        headers: plainFunctionHeader,
    });

    if (res.status !== StatusCode.OK) {
        http.get(Endpoint.EVALUATION_URL, { headers: plainFunctionHeader });
    }
    check(res, { "status was 200": (r) => r.status == StatusCode.OK });
}

export default function () {
    callSyncHello1();
}
