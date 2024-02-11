import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    discardResponseBodies: true,

    scenarios: {
  
      contacts: {
  
        executor: 'per-vu-iterations',
  
        vus: 16,
  
        iterations: 30,
  
        maxDuration: '15min',
  
      },
  
    },
};

const URL = 'http://localhost:8070/api/function_invocations';

export default function () {
  const vanilla = {
    method: 'POST',
    url: URL,
    params: {
      headers: { "X-Nuclio-Function-Name": 'vanilla',  "MAX": "5000000"},
    },
  };
  const second = {
    method: 'POST',
    url: URL,
    params: {
      headers: { "X-Nuclio-Function-Name": 'second',  "MAX": "2500000"},
    },
  };


  const res = http.batch([vanilla, second]);
  


  check(res[0], { 'status was 200': (r) => r.status == 200 });
}


