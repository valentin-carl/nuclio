import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    discardResponseBodies: true,

    scenarios: {
  
      contacts: {
  
        executor: 'per-vu-iterations',
  
        vus: 20,
  
        iterations: 20,
  
        maxDuration: '30s',
  
      },
  
    },
};

  const req3 = {
    method: 'GET',
    url: 'http://localhost:8070/api/function_invocationst',
    params: {
        headers: {
            'X-Nuclio-Function-Name': 'hello-1',
            'X-Nuclio-Function-Namespace': 'default'
          }
    },
  };

let count = 0;

export default function () {
  const headers = {
    'X-Nuclio-Function-Name': 'hello-1',
    'X-Nuclio-Function-Namespace': 'default'
  };
  const res = http.get('http://localhost:8070/api/function_invocations', { headers: headers });

  check(res, { 'status was 200': (r) => r.status == 200 });
}


