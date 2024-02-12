import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    discardResponseBodies: true,

    scenarios: {
  
      contacts: {
  
        executor: 'per-vu-iterations',
  
        vus: 1,
  
        iterations: 800,
  
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
      headers: { "X-Nuclio-Function-Name": 'vanilla', "X-Profaastinate-Process-Deadline": count > 400 && Math.random() >= 0.3 ? "5000" : "4000000" ,  "MAX": "5000000"},
    },
  };

    const second = {
      method: 'POST',
      url: URL,
      params: {
        headers: { "X-Nuclio-Function-Name": 'second', "X-Profaastinate-Process-Deadline": count > 400 && Math.random() >= 0.4 ? "5000" : "4000000", "MAX": "4000000"},
      }
    }


  const res = http.batch([vanilla, second]);
  
  if (count < 400 || count % 5 != 0 )
    count++
  else {
    count++
    sleep(9)
  }

  check(res[0], { 'status was 200': (r) => r.status == 200 });
}


