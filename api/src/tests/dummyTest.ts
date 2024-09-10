import assert from "node:assert";
import { test } from "node:test";

test("synchronous passing test", () => {
    assert.equal(1, 1);
});

test("asynchronous passing test", async () => {
    assert.equal(1, 1);
});
