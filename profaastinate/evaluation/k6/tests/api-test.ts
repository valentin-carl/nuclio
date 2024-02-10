// import necessary module
import http from "k6/http";
import { Endpoint, FunctionName, HeaderKey } from "../enums/public-api";
import { check } from "k6";

function buildHeader(functionName: string, deadline: number = 1000, namespace = 'nuclio')  {
  var functionHeader = new Headers();
  functionHeader.append(HeaderKey.FUNCTION_NAME, functionName);
  functionHeader.append(HeaderKey.FUNCTION_NAMESPACE, namespace);
  functionHeader.append(HeaderKey.ASYNC_DEADLINE, `${deadline}`);
  return convertHeadersToPlainObject(functionHeader);
}

function convertHeadersToPlainObject(headers: Headers): { [name: string]: string } {
  let plainHeaders: { [name: string]: string } = {};
  headers.forEach((value, key) => {
      plainHeaders[key] = value;
  });
  return plainHeaders;
}

export const options = {
  scenarios: {
      constant_request_rate: {
          executor: 'constant-arrival-rate',
          rate: 1,
          timeUnit: '1s', // 1000 iterations per second, i.e. 1000 RPS
          duration: '30s',
          preAllocatedVUs: 1, // the size of the VU (virtual user) pool for the executor
          maxVUs: 100, // the maximum number of VUs the executor is allowed to scale to
      },
  }
};

function callHello1() {
  const plainFunctionHeader = buildHeader(FunctionName.HELLO_1, 10000, 'default');
  const res = http.get(Endpoint.INVOCATION_URL, { headers: plainFunctionHeader });
  check(res, { 'status was 200': (r) => r.status == 200 });
}

export default function () {
  callHello1();
}
