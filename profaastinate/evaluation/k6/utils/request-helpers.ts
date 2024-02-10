import { HeaderKey } from '../enums/public-api';

export function buildHeader(functionName: string, deadline: number = 1000, namespace = 'nuclio'): { [name: string]: string }  {
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