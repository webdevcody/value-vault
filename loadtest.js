import http from "k6/http";
import { check } from "k6";

export const options = {
  stages: [{ target: 1000, duration: "10s" }],
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

function writeKey() {
  const key = randomString(6);
  const res = http.post(
    `http://localhost:8080/keys/${key}`,
    JSON.stringify({
      json: {
        key: "test",
        value: "test",
      },
    })
  );
  // Validate response status
  check(res, { "status was 200": (r) => r.status == 201 });
}

function readValue() {
  const key = randomString(6);
  const res = http.get(`http://localhost:7777/keys/${key}`);
  check(res, { "status was 200": (r) => r.status == 200 });
}

// Simulated user behavior
export default function () {
  writeKey();
  // readValue();
  // getTest();
}
