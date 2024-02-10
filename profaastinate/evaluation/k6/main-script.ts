import http from "k6/http";
import {check} from "k6";

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

export enum Duration {
    THIRTY_SECONDS = "30s",
}

/**
 * Enum for Nuclio header keys
 */
export enum HeaderKey {
    FUNCTION_NAME = "x-nuclio-function-name",
    FUNCTION_NAMESPACE = "x-nuclio-function-namespace",
    ASYNC_DEADLINE = "x-profaastinate-process-deadline",
}

export function buildHeader(functionName: string, deadline?: number, namespace?: string): { [name: string]: string } {
    let headers: { [name: string]: string } = {};
    headers[HeaderKey.FUNCTION_NAME] = functionName;
    if (deadline) {
        headers[HeaderKey.ASYNC_DEADLINE] = `${deadline}`;
    }
    if (namespace) {
        headers[HeaderKey.FUNCTION_NAMESPACE] = namespace;
    }
    return headers;
}

export const options = {
    scenarios: {
        ramping_arrival_rate: {
            executor: 'ramping-arrival-rate',
            startRate: 250,
            timeUnit: Duration.THIRTY_SECONDS, // 50 iterations per minute
            preAllocatedVUs: 50, // pre-allocate 50 VUs
            maxVUs: 100,
            stages: [
                {target: 500, duration: Duration.THIRTY_SECONDS}, // ramp to 100 iterations per minute
                {target: 1000, duration: Duration.THIRTY_SECONDS}, // stay at 100 for 5 minutes
                {target: 100, duration: Duration.THIRTY_SECONDS},   // ramp down to 10 iterations per minute
            ],
        },
    },
};

function callAsyncHello1() {
    const plainFunctionHeader = buildHeader(
        FunctionName.HELLO_1,
        30000
    );
    const res = http.get(Endpoint.INVOCATION_URL, {
        headers: plainFunctionHeader,
    });
    check(res, {"status was 204": (r) => r.status == StatusCode.OK});
}

function callSyncHello1() {
    const plainFunctionHeader = buildHeader(
        FunctionName.HELLO_1
    );
    const res = http.get(Endpoint.INVOCATION_URL, {
        headers: plainFunctionHeader,
    });

    if (res.status !== StatusCode.OK) {
        http.get(Endpoint.EVALUATION_URL, {headers: plainFunctionHeader});
    }
    check(res, {"status was 200": (r) => r.status == StatusCode.OK});
}

export default function () {
    callAsyncHello1();
}
