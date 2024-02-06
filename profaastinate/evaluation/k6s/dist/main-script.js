"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.options = exports.buildHeader = exports.HeaderKey = exports.FunctionName = exports.Endpoint = void 0;
const http_1 = __importDefault(require("k6/http"));
const k6_1 = require("k6");
/**
 * Enum for the endpoints used in the application
 */
var Endpoint;
(function (Endpoint) {
    Endpoint["INVOCATION_URL"] = "http://localhost:8070/api/function_invocations";
})(Endpoint || (exports.Endpoint = Endpoint = {}));
var FunctionName;
(function (FunctionName) {
    FunctionName["HELLO_1"] = "hello-1";
})(FunctionName || (exports.FunctionName = FunctionName = {}));
/**
 * Enum for Nuclio header keys
 */
var HeaderKey;
(function (HeaderKey) {
    HeaderKey["FUNCTION_NAME"] = "x-nuclio-function-name";
    HeaderKey["FUNCTION_NAMESPACE"] = "x-nuclio-function-namespace";
    HeaderKey["ASYNC_DEADLINE"] = "x-nuclio-async-deadline";
})(HeaderKey || (exports.HeaderKey = HeaderKey = {}));
function buildHeader(functionName, deadline = 1000, namespace = 'nuclio') {
    let headers = {};
    headers[HeaderKey.FUNCTION_NAME] = functionName;
    headers[HeaderKey.FUNCTION_NAMESPACE] = namespace;
    headers[HeaderKey.ASYNC_DEADLINE] = `${deadline}`;
    return headers;
}
exports.buildHeader = buildHeader;
exports.options = {
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
    const plainFunctionHeader = buildHeader(FunctionName.HELLO_1, 10000, "default");
    const res = http_1.default.get(Endpoint.INVOCATION_URL, {
        headers: plainFunctionHeader,
    });
    (0, k6_1.check)(res, { "status was 200": (r) => r.status == 200 });
}
function default_1() {
    callHello1();
}
exports.default = default_1;
