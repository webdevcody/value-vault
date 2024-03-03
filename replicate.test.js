async function ingest() {
  for (let i = 0; i < 1000; i++) {
    await fetch(`http://localhost/keys/${i}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ value: i }),
    });
  }
}

async function verify() {
  for (let i = 0; i < 1000; i++) {
    const value = await fetch(`http://localhost/keys/${i}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    }).then((r) => r.json());
    if (!value || value.value !== i) {
      console.error(`key of ${i} was not found or had the wrong value`);
    }
  }
}

async function sleep(amount) {
  return new Promise((resolve) => setTimeout(resolve, amount));
}

async function main() {
  // await ingest();
  // console.log("sleeping, run apply now");
  // await sleep(60000);
  await verify();
}

main();
