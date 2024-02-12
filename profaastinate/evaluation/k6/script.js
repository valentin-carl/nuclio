import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    discardResponseBodies: true,

    scenarios: {
  
      contacts: {
  
        executor: 'per-vu-iterations',
  
        vus: 4,
  
        iterations: 100,
  
        maxDuration: '15m',
  
      },
  
    },
};


const URL = 'http://localhost:8070/api/function_invocations';

let count = 0

 export default function () {
  const vanilla = {
    method: 'POST',
    url: URL,
    params: {
      headers: { "X-Nuclio-Function-Name": 'vanilla', "X-Profaastinate-Process-Deadline": (count % 12 == 0 ) ? "5000" : "4000000" ,  "MAX": "5000000"},
    },
  };

    const second = {
      method: 'POST',
      url: URL,
      params: {
        headers: { "X-Nuclio-Function-Name": 'second', "X-Profaastinate-Process-Deadline": ( count % 6 == 0 ) ? "5000" : "4000000", "MAX": "4000000"},
      }
    }

    const vanilla_2 = {
      method: 'POST',
      url: URL,
      params: {
        headers: { "X-Nuclio-Function-Name": 'vanilla', "X-Profaastinate-Process-Deadline": ( count % 12 == 0 ) ? "5000" : "4000000",  "MAX": "5000000"},
      },
    };
  
      const second_2 = {
        method: 'POST',
        url: URL,
        params: {
          headers: { "X-Nuclio-Function-Name": 'second', "X-Profaastinate-Process-Deadline":(count % 24 == 0 ) ? "5000" : "4000000", "MAX": "4000000"},
        }
      }


  const res = http.batch([vanilla, second, vanilla_2, second_2]);
  
  if (count < 7)
    count++
  else {
    count++
    sleep(5)
  }

  check(res[0], { 'status was 200': (r) => r.status == 200 });
}


