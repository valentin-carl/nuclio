"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.buildHeader = void 0;
const public_api_1 = require("../enums/public-api");
function buildHeader(functionName, deadline = 1000, namespace = 'nuclio') {
    var functionHeader = new Headers();
    functionHeader.append(public_api_1.HeaderKey.FUNCTION_NAME, functionName);
    functionHeader.append(public_api_1.HeaderKey.FUNCTION_NAMESPACE, namespace);
    functionHeader.append(public_api_1.HeaderKey.ASYNC_DEADLINE, `${deadline}`);
    return convertHeadersToPlainObject(functionHeader);
}
exports.buildHeader = buildHeader;
function convertHeadersToPlainObject(headers) {
    let plainHeaders = {};
    headers.forEach((value, key) => {
        plainHeaders[key] = value;
    });
    return plainHeaders;
}
