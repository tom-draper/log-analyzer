<script lang="ts">
  import ValueFreqGraph from "./ValueFreqGraph.svelte";
  import OverTimeGraph from "./OverTimeGraph.svelte";
  import Statistics from "./Statistics.svelte";
  import Distribution from "./Distribution.svelte";
  import Activity from "./Activity.svelte";
  import TimeBreakdown from "./TimeBreakdown.svelte";
  import LocationMap from "./LocationMap.svelte";
  import LatencyOverTime from "./LatencyOverTime.svelte";
  import ValueFreqGraphDependent from "./ValueFreqGraphDependent.svelte";
  import TopValuesTable from "./TopValuesTable.svelte";
  import DependencyHeatmap from "./DependencyHeatmap.svelte";

  export let data: Data,
    token: string,
    dependentToken: string | null,
    lineCount: number,
    timestampToken: string | null,
    hex: string;

  $: isLatencyToken = /(?:latency|duration|elapsed|response[_-]?time)/i.test(token);
  let showHeatmap = false;
  let showBreakdown = false;
  let showDistribution = false;
  let showExamples = false;
</script>

<div class="card">
  <div class="header">
    <div class="title">
      {#if dependentToken}
        {token} & {dependentToken}
      {:else}
        {token}
      {/if}
    </div>
    <div class="line-count">
      {lineCount.toLocaleString()} lines
      <span class="coverage">({((lineCount / data.extraction.length) * 100).toFixed(1)}%)</span>
    </div>
  </div>
  {#if dependentToken}
    <ValueFreqGraphDependent {data} {token} {dependentToken} />
    <details ontoggle={(event) => showHeatmap = event.currentTarget.open}>
      <summary>Correlation heatmap</summary>
      {#if showHeatmap}<DependencyHeatmap {data} {token} {dependentToken} />{/if}
    </details>
  {:else if token === timestampToken}
    <Activity {data} {token} {hex} />
    <details ontoggle={(event) => showBreakdown = event.currentTarget.open}>
      <summary>Time breakdown</summary>
      {#if showBreakdown}<TimeBreakdown {data} {token} {hex} />{/if}
    </details>
  {:else}
    <Statistics {data} {token} />
    <ValueFreqGraph {data} {token} />
    <details ontoggle={(event) => showDistribution = event.currentTarget.open}>
      <summary>Distribution & trends</summary>
      {#if showDistribution}<Distribution {data} {token} {hex} />
      {#if isLatencyToken && timestampToken}
        <LatencyOverTime {data} {token} {timestampToken} {hex} />
      {/if}
      <OverTimeGraph {data} {token} {timestampToken} />
      <LocationMap {data} {token} {hex} />{/if}
    </details>
    <details ontoggle={(event) => showExamples = event.currentTarget.open}>
      <summary>Values & examples</summary>
      {#if showExamples}<TopValuesTable {data} {token} />{/if}
    </details>
  {/if}
</div>

<style scoped>
  .header {
    margin-bottom: 20px;
    /* color: ; */
    display: flex;
  }
  .line-count {
    margin-left: auto;
    color: #888;
  }
  .coverage {
    color: #555;
    font-size: 0.9em;
    margin-left: 4px;
  }
  .title,
  .line-count {
    font-weight: 500;
  }
  .card {
    border: 1px solid var(--border);
    border-radius: var(--radius);
    margin: 1.5rem 0;
    padding: 1.5rem;
    background: var(--card);
  }
  details {
    border-top: 1px solid var(--border);
    margin-top: 1rem;
    padding-top: .75rem;
  }
  summary {
    color: var(--muted);
    cursor: pointer;
    font-size: .8125rem;
    font-weight: 500;
    list-style-position: inside;
  }
  summary:hover { color: var(--text); }

  @media screen and (max-width: 800px) {
    .card {
      padding: 1.25rem;
    }
  }
</style>
