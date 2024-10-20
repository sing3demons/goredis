import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 10,
    duration: '30s',
};

export default function () {
    const url = `${__ENV.BASE_URL || 'https://host.docker.internal:8000'}/api`;
    
    // Make the request, ignoring SSL certificate validation
    const res = http.get(url, { tags: { my_custom_tag: 'test' }, insecure: true });

    // Check the response status
    check(res, {
        'is status 200': (r) => r.status === 200,
        'response body is not empty': (r) => r.body.length > 0,
    });

    sleep(1); // Sleep for 1 second between requests
}
