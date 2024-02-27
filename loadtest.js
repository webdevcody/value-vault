import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  stages: [
    { target: 200, duration: "20s" },
    { target: 200, duration: "20s" },
  ],
  noConnectionReuse: true,
  dns: {
    ttl: "0",
    select: "random",
  },
};

const BASE_URL = "http://localhost";

function randomString(length) {
  var result = "";
  var characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}

function writeKey(key, value, traceId) {
  const res = http.post(
    // `http://localhost:${getRandomNode()}/keys/${key}`,
    `${BASE_URL}/keys/${key}`,
    JSON.stringify(value),
    {
      headers: {
        "X-Trace-Id": traceId,
      },
    }
  );
  // Validate response status
  check(res, { "status was 201": (r) => r.status == 201 });
}

function readValue(key, expectedValue, traceId) {
  const res = http.get(`${BASE_URL}/keys/${key}`, {
    headers: {
      "X-Trace-Id": traceId,
    },
  });
  // const res = http.get(`http://localhost:${getRandomNode()}/keys/${key}`);
  check(res, {
    "Response body contains expected string": (res) => {
      if (!res.body.includes(expectedValue)) {
        console.log(`expected ${expectedValue} but got ${res.body}`);
        console.log(key);
      }
      return res.body.includes(expectedValue);
    },
  });
}

function getStatus(key, expectedValue, traceId) {
  const res = http.get(`${BASE_URL}/status`, {
    headers: {
      "X-Trace-Id": traceId,
    },
  });
  check(res, { "status was 200": (r) => r.status == 200 });
}

export default function () {
  const randomValue = randomString(32);
  const randomTrace = randomString(32);
  const key = randomString(6);
  writeKey(key, randomValue, randomTrace);
  sleep(5);
  readValue(key, randomValue, randomTrace);
}

// export default function () {
//   getStatus();
// }
