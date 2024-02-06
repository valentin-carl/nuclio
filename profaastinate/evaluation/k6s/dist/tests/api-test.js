"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.options = void 0;
// import necessary module
const http_1 = __importDefault(require("k6/http"));
const public_api_1 = require("../enums/public-api");
const k6_1 = require("k6");
function buildHeader(functionName, deadline = 1000, namespace = 'nuclio') {
    var functionHeader = new Headers();
    functionHeader.append(public_api_1.HeaderKey.FUNCTION_NAME, functionName);
    functionHeader.append(public_api_1.HeaderKey.FUNCTION_NAMESPACE, namespace);
    functionHeader.append(public_api_1.HeaderKey.ASYNC_DEADLINE, `${deadline}`);
    return convertHeadersToPlainObject(functionHeader);
}
function convertHeadersToPlainObject(headers) {
    let plainHeaders = {};
    headers.forEach((value, key) => {
        plainHeaders[key] = value;
    });
    return plainHeaders;
}
exports.options = {
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
    const plainFunctionHeader = buildHeader(public_api_1.FunctionName.HELLO_1, 10000, 'default');
    const res = http_1.default.get(public_api_1.Endpoint.INVOCATION_URL, { headers: plainFunctionHeader });
    (0, k6_1.check)(res, { 'status was 200': (r) => r.status == 200 });
}
function default_1() {
    callHello1();
}
exports.default = default_1;
