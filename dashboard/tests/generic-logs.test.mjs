import assert from "node:assert/strict";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { fileURLToPath, pathToFileURL } from "node:url";
import test from "node:test";
import { compile } from "svelte/compiler";
import { render } from "svelte/server";
import ts from "typescript";

const data = (params) => ({
  extraction: params.map((params, i) => ({ params, line: `event ${i}`, lineNumber: i, pattern: "event" })),
  config: { tokens: [], patterns: ["event"] }, locations: {},
});
const number = value => ({ value, type: "int" });
const string = value => ({ value, type: "str" });
const time = value => ({ value, type: "time" });

async function withComponent(name, run) {
  const filename = fileURLToPath(new URL(`../src/lib/${name}.svelte`, import.meta.url));
  const source = (await readFile(filename, "utf8")).replace('"./timeSlots"', '"./timeSlots.cjs"');
  const { js } = compile(source, { filename, generate: "server" });
  const directory = await mkdtemp(new URL(".generic-", import.meta.url));
  try {
    const timeSlots = await readFile(new URL("../src/lib/timeSlots.ts", import.meta.url), "utf8");
    await writeFile(`${directory}/timeSlots.cjs`, ts.transpileModule(timeSlots, { compilerOptions: { module: ts.ModuleKind.CommonJS, esModuleInterop: true } }).outputText);
    const compiledFile = `${directory}/${name}.mjs`;
    await writeFile(compiledFile, js.code);
    const { default: Component } = await import(pathToFileURL(compiledFile).href);
    await run(props => render(Component, { props }).body.replace(/<!--[\s\S]*?-->/g, ""));
  } finally {
    await rm(directory, { recursive: true, force: true });
  }
}

test("non-HTTP logs do not render request health for missing or ambiguous codes", async () => {
  await withComponent("RequestFunnel", render => {
    for (const params of [[], [{ level: string("INFO") }], [{ code: number(500) }], [{ status: number(200) }], [{ status: string("failed") }], [{ status: number(200.5), method: string("GET") }]]) {
      assert.ok(!render({ data: data(params) }).includes("Request health"));
    }
  });
});

test("request health requires HTTP evidence and counts only valid status codes", async () => {
  await withComponent("RequestFunnel", render => {
    const body = render({ data: data([
      { status: number(200), method: string("GET") },
      { status: number(404), method: string("POST") },
      { status: number(503), method: string("GET") },
      { status: number(0) }, { status: number(500.5) }, {},
    ]) });
    assert.ok(body.includes("Request health"));
    assert.match(body, /<strong\b[^>]*>3<\/strong>/);
    assert.equal((body.match(/33.3%/g) ?? []).length, 2);
    assert.ok(render({ data: data([{ http_status: number(204) }]) }).includes("Request health"));
  });
});

test("time trends skip records with missing timestamps and support logs without time", async () => {
  await withComponent("OverTimeGraph", render => {
    const input = data([
      { timestamp: time("2024-01-01T12:00:00Z"), level: string("INFO") },
      { timestamp: time("2024-01-01T13:00:00Z"), level: string("INFO") },
      { level: string("ERROR") },
      { timestamp: time("invalid"), level: string("ERROR") },
    ]);
    const body = render({ data: input, token: "level", timestampToken: "timestamp" });
    assert.ok(body.includes("INFO"));
    assert.ok(!body.includes("ERROR"));
    assert.ok(!render({ data: input, token: "level", timestampToken: null }).includes("freq-graph"));
    assert.ok(!render({ data: data([{ timestamp: time("invalid"), level: string("INFO") }]), token: "level", timestampToken: "timestamp" }).includes("freq-graph"));
  });
});

test("generic summaries and the explorer render text-only and unmatched records", async () => {
  const input = data([{ level: string("INFO"), message: string("started") }, {}]);
  input.extraction[1].pattern = "";
  input.extraction[1].line = "unstructured application output";
  await withComponent("DashboardMetrics", render => {
    const body = render({ data: input });
    assert.ok(body.includes("Time range"));
    assert.ok(body.includes("50.0%"));
    assert.ok(!body.includes("Invalid Date"));
  });
  await withComponent("EventExplorer", render => {
    assert.ok(render({ data: input }).includes("unstructured application output"));
  });
  await withComponent("PatternCoverage", render => {
    const body = render({ data: input });
    assert.ok(body.includes("Unmatched lines"));
    assert.ok(body.includes("50.0%"));
  });
});
