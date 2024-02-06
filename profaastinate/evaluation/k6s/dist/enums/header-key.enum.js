"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HeaderKey = void 0;
/**
 * Enum for Nuclio header keys
 */
var HeaderKey;
(function (HeaderKey) {
    HeaderKey["FUNCTION_NAME"] = "x-nuclio-function-name";
    HeaderKey["FUNCTION_NAMESPACE"] = "x-nuclio-function-namespace";
    HeaderKey["ASYNC_DEADLINE"] = "x-nuclio-async-deadline";
})(HeaderKey || (exports.HeaderKey = HeaderKey = {}));
