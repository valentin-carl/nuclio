import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    discardResponseBodies: true,

    scenarios: {
  
      contacts: {
  
        executor: 'per-vu-iterations',
  
        vus: 15,
  
        iterations: 20,
  
        maxDuration: '30s',
  
      },
  
    },
};

const URL = 'http://localhost:8070/api/function_invocations';

export default function () {
  const vanilla = {
    method: 'POST',
    url: URL,
    body: "5000000",
    params: {
      headers: { "X-Nuclio-Function-Name": 'vanilla' },
    },
  };
  const second = {
    method: 'POST',
    url: URL,
    body: "5000000",
    params: {
      headers: { "X-Nuclio-Function-Name": 'second' },
    },
  };


  const res = http.batch([vanilla, second]);
  


  check(res[0], { 'status was 200': (r) => r.status == 200 });
}


