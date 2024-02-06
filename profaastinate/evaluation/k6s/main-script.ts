import http from "k6/http";
import { check } from "k6";

/**
 * Enum for the endpoints used in the application
 */
export enum Endpoint {
  INVOCATION_URL = "http://localhost:8070/api/function_invocations",
}

export enum FunctionName {
  HELLO_1 = "hello-1",
}

/**
 * Enum for Nuclio header keys
 */
export enum HeaderKey {
  FUNCTION_NAME = "x-nuclio-function-name",
  FUNCTION_NAMESPACE = "x-nuclio-function-namespace",
  ASYNC_DEADLINE = "x-nuclio-async-deadline",
}

export function buildHeader(functionName: string, deadline: number = 1000, namespace = 'nuclio'): { [name: string]: string } {
    let headers: { [name: string]: string } = {};
    headers[HeaderKey.FUNCTION_NAME] = functionName;
    headers[HeaderKey.FUNCTION_NAMESPACE] = namespace;
    headers[HeaderKey.ASYNC_DEADLINE] = `${deadline}`;
    return headers;
}

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 1,
      timeUnit: "1s", // 1000 iterations per second, i.e. 1000 RPS
      duration: "30s",
      preAllocatedVUs: 1, // the size of the VU (virtual user) pool for the executor
      maxVUs: 100, // the maximum number of VUs the executor is allowed to scale to
    },
  },
};

function callHello1() {
  const plainFunctionHeader = buildHeader(
    FunctionName.HELLO_1,
    10000,
    "default"
  );
  const res = http.get(Endpoint.INVOCATION_URL, {
    headers: plainFunctionHeader,
  });
  check(res, { "status was 200": (r) => r.status == 200 });
}

export default function () {
  callHello1();
}
