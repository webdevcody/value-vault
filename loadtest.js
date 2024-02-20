import http from "k6/http";
import { check, sleep } from "k6";

// Test configuration
export const options = {
  vus: 200,
  duration: "10s",
  // stages: [
  //   { target: 300, duration: '30s' },
  // ],
};

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

// Simulated user behavior
export default function () {
  const key = randomString(6);
  let res = http.post(
    `http://localhost/keys/${key}`,
    JSON.stringify({
      json: {
        key: "test",
        value: "test",
      },
    })
  );
  // Validate response status
  check(res, { "status was 200": (r) => r.status == 201 });
  // sleep(1);
  // res = http.get(
  //   `http://localhost/keys/${key}`,
  //   JSON.stringify({
  //     json: {
  //       key: "test",
  //       value: "test",
  //     },
  //   })
  // );
  // // Validate response status
  // check(res, { "status was 200": (r) => r.status == 200 });
}
