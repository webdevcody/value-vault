import http from "k6/http";
import { check, sleep } from "k6";

let written = 0;
export const options = {
  stages: [{ target: 2, duration: "60s" }],
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

const nodes = ["8080", "8081"];

function getRandomNode() {
  return nodes[Math.floor(Math.random() * nodes.length)];
}

function writeKey(key, value) {
  const res = http.post(
    `http://localhost:${getRandomNode()}/keys/${key}`,
    JSON.stringify(value)
  );
  written++;
  // Validate response status
  check(res, { "status was 200": (r) => r.status == 201 });
}

function readValue(key, expectedValue) {
  const res = http.get(`http://localhost:${getRandomNode()}/keys/${key}`);
  check(res, {
    "Response body contains expected string": (res) => {
      return res.body.includes(expectedValue);
    },
  });
}

// Simulated user behavior
export default function () {
  const randomValue = randomString(32);
  const key = randomString(6);
  writeKey(key, randomValue);
  sleep(1);
  readValue(key, randomValue);
  console.log("written", written);
}
