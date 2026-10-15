<script lang="ts">
  import { emptyFilters, exportRecords, type DashboardFilters } from "./filters";

  export let data: Data;
  export let records: Extraction[];
  export let filters: DashboardFilters;
  $: tokens = [...new Set([...data.extraction.flatMap(record => Object.keys(record.params)), ...(filters.token ? [filters.token] : [])])].sort();
  $: sources = [...new Set([...data.extraction.map(record => record.sourceFile ?? ""), ...(filters.source !== null ? [filters.source] : [])])].sort();
  $: values = [...new Set([...data.extraction.flatMap(record => filters.token in record.params ? [String(record.params[filters.token].value)] : []), ...(filters.value !== null ? [filters.value] : [])])].sort();
  $: invalidRange = filters.start && filters.end && new Date(filters.start) > new Date(filters.end);

  function download() {
    const url = URL.createObjectURL(new Blob([exportRecords(records)], { type: "application/json" }));
    const link = document.createElement("a");
    link.href = url;
    link.download = "matching-records.json";
    link.click();
    setTimeout(() => URL.revokeObjectURL(url), 1000);
  }
</script>

<section class="filters" aria-label="Dashboard filters">
  <div class="heading"><h2>Dashboard filters</h2><span aria-live="polite">{records.length.toLocaleString()} of {data.extraction.length.toLocaleString()} records</span></div>
  <div class="controls">
    <label>From<input type="datetime-local" step="1" bind:value={filters.start} /></label>
    <label>To<input type="datetime-local" step="1" bind:value={filters.end} /></label>
    <label>Token<select bind:value={filters.token} onchange={() => filters.value = null}><option value="">All tokens</option>{#each tokens as token}<option value={token}>{token}</option>{/each}</select></label>
    <label>Value<select bind:value={filters.value} disabled={!filters.token}><option value={null}>Any value</option>{#each values as value}<option {value}>{value === "" ? "(empty string)" : value}</option>{/each}</select></label>
    <label>Source file<select bind:value={filters.source}><option value={null}>All sources</option>{#each sources as source}<option value={source}>{source || "Unknown source"}</option>{/each}</select></label>
    <label>Raw text<input bind:value={filters.query} placeholder="Filter raw log lines" /></label>
    <label class="checkbox"><input type="checkbox" bind:checked={filters.failuresOnly} /> Unmatched only</label>
  </div>
  <p>Filters apply to all charts and events. Times use your local timezone and include both bounds. Records without a valid timestamp are excluded when a time bound is set.</p>
  {#if invalidRange}<p role="alert">From must be earlier than or equal to To.</p>{/if}
  <div class="actions"><button onclick={() => filters = emptyFilters()}>Clear filters</button><button onclick={download} disabled={records.length === 0}>Export {records.length.toLocaleString()} matching records (JSON)</button></div>
</section>

<style>
  .filters { border: 1px solid var(--border); border-radius: var(--radius); background: var(--card); padding: 1.25rem 1.5rem; margin: 1.5rem 0; }
  .heading, .actions { display: flex; flex-wrap: wrap; gap: .75rem; align-items: center; justify-content: space-between; }
  h2 { margin: 0; font-size: 1rem; }
  span, p { color: var(--muted); font-size: .8rem; }
  .controls { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: .8rem; margin: 1rem 0; }
  label { display: flex; flex-direction: column; gap: .35rem; font-size: .85rem; min-width: 0; }
  input, select, button { background: #181818; color: var(--text); border: 1px solid #333; border-radius: 4px; padding: .45rem .6rem; min-width: 0; }
  .checkbox { flex-direction: row; align-items: center; }
  button { cursor: pointer; } :disabled { opacity: .5; } [role="alert"] { color: #ff8a8a; }
</style>
