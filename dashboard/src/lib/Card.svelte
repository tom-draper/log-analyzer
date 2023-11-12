<script lang="ts">
  import ValueFreqGraph from "./ValueFreqGraph.svelte";
  import OverTimeGraph from "./OverTimeGraph.svelte";
  import Statistics from "./Statistics.svelte";
  import Distribution from "./Distribution.svelte";
  import Activity from "./Activity.svelte";
  import LocationMap from "./LocationMap.svelte";
  import ValueFreqGraphDependent from "./ValueFreqGraphDependent.svelte";

  export let data: Data,
    token: string,
    dependentToken: string | null,
    lineCount: number,
    timestampToken: string | null;
</script>

<div class="card">
  <div class="header">
    <div class="title">{dependentToken ? `${token} & ${dependentToken}` : token}</div>
    <div class="line-count">{lineCount.toLocaleString()} lines</div>
  </div>
  {#if dependentToken} 
    <ValueFreqGraphDependent {data} {token} {dependentToken} />
  {:else if token === timestampToken}
    <Activity {data} {token} />
  {:else}
    <Statistics {data} {token} />
    <Distribution {data} {token} />
    <ValueFreqGraph {data} {token} />
    <OverTimeGraph {data} {token} {timestampToken} />
    <LocationMap {data} {token} />
  {/if}
</div>

<style scoped>
  .header {
    margin-bottom: 20px;
    color: white;
    display: flex;
  }
  .line-count {
    margin-left: auto;
    color: #888;
  }
  .card {
    border: 1px solid #ffffff24;
    border-radius: 5px;
    margin: 3em 0;
    padding: 3rem;
  }

  @media screen and (max-width: 800px) {
    .card {
      padding: 2em;
    }
  }
</style>
