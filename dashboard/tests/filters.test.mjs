import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";
import ts from "typescript";

const source = await readFile(new URL("../src/lib/filters.ts", import.meta.url), "utf8");
const { outputText } = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.ESNext } });
const { emptyFilters, filterRecords, exportRecords } = await import(`data:text/javascript;base64,${Buffer.from(outputText).toString("base64")}`);
const record = (lineNumber, sourceFile, timestamp, status, line = "GET /home", pattern = "request") => ({
  lineNumber, sourceFile, line, pattern,
  params: { ...(timestamp === null ? {} : { timestamp: { value: timestamp, type: "time" } }), ...(status === null ? {} : { status: { value: status, type: "int" } }) },
});
const records = [
  record(0, "a.log", "2024-01-01T12:00:00Z", 200),
  record(1, "b.log", "2024-01-01T13:00:00Z", 500),
  record(2, "a.log", "2024-01-01T14:00:00Z", 500),
  record(3, undefined, null, null, "broken GET", ""),
  record(4, "a.log", "invalid", 500),
];
const select = (overrides, input = records) => filterRecords(input, { ...emptyFilters(), ...overrides });

test("empty filters retain records and source metadata", () => {
  assert.deepEqual(select({}), records);
});

test("time bounds are inclusive and exclude missing or invalid timestamps", () => {
  assert.deepEqual(select({ timestampToken: "timestamp", start: "2024-01-01T12:00:00Z", end: "2024-01-01T13:00:00Z" }), records.slice(0, 2));
  assert.deepEqual(select({ timestampToken: "timestamp", end: "2024-01-01T12:00:00Z" }), [records[0]]);
  assert.deepEqual(select({ timestampToken: "timestamp", start: "2024-01-01T14:00:00Z" }), [records[2]]);
  assert.deepEqual(select({ timestampToken: "timestamp", start: "2024-02-01", end: "2024-01-01" }), []);
});

test("token filters require presence and compare exact values including zero and empty strings", () => {
  assert.deepEqual(select({ token: "status", value: "500" }), [records[1], records[2], records[4]]);
  assert.deepEqual(select({ token: "status" }), [records[0], records[1], records[2], records[4]]);
  assert.deepEqual(select({ token: "status", value: "50" }), []);
  const edgeRecords = [record(0, "a", null, 0), record(1, "a", null, "")];
  assert.deepEqual(select({ token: "status", value: "0" }, edgeRecords), [edgeRecords[0]]);
  assert.deepEqual(select({ token: "status", value: "" }, edgeRecords), [edgeRecords[1]]);
});

test("source, token, time, and case-insensitive raw text filters intersect", () => {
  assert.deepEqual(select({ source: "a.log", token: "status", value: "500", timestampToken: "timestamp", start: "2024-01-01", end: "2024-01-02", query: "get /HOME" }), [records[2]]);
  assert.deepEqual(select({ source: "", failuresOnly: true, query: "get" }), [records[3]]);
  assert.deepEqual(select({ source: "missing.log" }), []);
});

test("export contains all matching records, preserves metadata and multiline text, and follows refreshes", () => {
  const filters = { ...emptyFilters(), source: "a.log" };
  const updated = [...records, record(5, "a.log", null, 200, 'line one\n"line two"')];
  const matching = filterRecords(updated, filters);
  assert.deepEqual(JSON.parse(exportRecords(matching)), [records[0], records[2], records[4], updated[5]]);
  const many = Array.from({ length: 125 }, (_, i) => record(i, "a.log", null, 200));
  assert.equal(JSON.parse(exportRecords(filterRecords(many, filters))).length, 125);
  assert.deepEqual(JSON.parse(exportRecords([])), []);
});
