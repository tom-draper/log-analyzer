import assert from "node:assert/strict";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { fileURLToPath, pathToFileURL } from "node:url";
import test from "node:test";
import { compile } from "svelte/compiler";
import { render } from "svelte/server";

test("dashboard summary renders the time range of a large log", async () => {
  const filename = fileURLToPath(new URL("../src/lib/DashboardMetrics.svelte", import.meta.url));
  const source = await readFile(filename, "utf8");
  const { js } = compile(source, { filename, generate: "server" });
  const directory = await mkdtemp(new URL(".metrics-", import.meta.url));
  try {
    const compiledFile = `${directory}/DashboardMetrics.mjs`;
    await writeFile(compiledFile, js.code);
    const { default: DashboardMetrics } = await import(pathToFileURL(compiledFile).href);
    const first = "2024-01-01T12:00:00Z";
    const last = "2024-02-01T12:00:00Z";
    const extraction = Array.from({ length: 200_000 }, (_, i) => ({
      params: { timestamp: { value: i === 0 ? first : last, type: "time" } },
      pattern: "timestamp",
      lineNumber: i,
      line: i === 0 ? first : last,
    }));
    const { body } = render(DashboardMetrics, {
      props: { data: { extraction, config: { tokens: ["timestamp"], patterns: ["timestamp"] }, locations: {} } },
    });
    assert.ok(body.includes(`${new Date(first).toLocaleDateString()} – ${new Date(last).toLocaleDateString()}`));
    assert.ok(body.includes("100.0%"));
  } finally {
    await rm(directory, { recursive: true, force: true });
  }
});
