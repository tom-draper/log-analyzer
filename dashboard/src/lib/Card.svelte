<script lang="ts">
  import { onMount } from "svelte";
  import ValueFreqGraph from "./ValueFreqGraph.svelte";
  import OverTimeGraph from "./OverTimeGraph.svelte";
  import Statistics from "./Statistics.svelte";

  function tokenValueFrequency(data: Data, token: string): ValueCount {
    let freq: ValueCount = {};
    for (let i = 0; i < data.extraction.params.length; i++) {
      for (let [_token, value] of Object.entries(data.extraction.params[i])) {
        if (_token !== token) {
          continue;
        }
        if (!(value in freq)) {
          freq[value] = 0;
        }
        freq[value] += 1;
      }
    }

    return freq;
  }

  let freq: ValueCount;
  onMount(() => {
    freq = tokenValueFrequency(data, token);
  });

  export let data: Data,
    token: string,
    lineCount: number,
    timestampToken: string | null;
</script>

<div class="card">
  <div class="header">
    <div class="title">{token}</div>
    <div class="line-count">{lineCount.toLocaleString()} lines</div>
  </div>
  <Statistics {data} {token}/>
  {#if freq !== undefined}
    <ValueFreqGraph {freq} />
  {/if}
  <OverTimeGraph {data} {token} {timestampToken} />
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
    /* color: #ededed; */
  }
  .card {
    border: 1px solid #ffffff24;
    border-radius: 5px;
    margin: 3em 0;
    padding: 3rem;
  }
</style>
