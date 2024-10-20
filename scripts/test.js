import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 10,
    duration: '30s',
    insecureSkipTLSVerify: true,  // Globally disable TLS verification
};

export default function () {
    // const url = 'https://localhost:8000/api';
    const url = `${__ENV.BASE_URL || 'https://localhost:8000'}/api`;

    const res = http.get(url);

    check(res, {
        'is status 200': (r) => r.status === 200,
        'response body contains message': (r) => r.body.includes('"message":"Hello, World!"'),
    });

    console.log(`Response body: ${res.body}`);

    sleep(1);
}
