<script lang="ts">
  import { onMount } from "svelte";
  import ValueFreqGraph from "./ValueFreqGraph.svelte";
  import OverTimeGraph from "./OverTimeGraph.svelte";
  import Statistics from "./Statistics.svelte";
  import Distribution from "./Distribution.svelte";
  import Activity from "./Activity.svelte";
  import LocationMap from "./LocationMap.svelte";

  function tokenValueFrequency(data: Data, token: string): ValueCount {
    let freq: ValueCount = {};
    for (let i = 0; i < data.extraction.length; i++) {
      for (let [_token, value] of Object.entries(data.extraction[i].params)) {
        if (_token !== token) {
          continue;
        }
        if (!(value.value in freq)) {
          freq[value.value] = 0;
        }
        freq[value.value] += 1;
      }
    }

    return freq;
  }

  function isIPAddressToken(data: Data, token: string): boolean {
    for (let i = 0; i < data.extraction.length; i++) {
      for (let [_token, value] of Object.entries(data.extraction[i].params)) {
        if (_token !== token) {
          continue;
        }
        if (value.type == "net.IP" && value.value in data.locations) {
          return true;
        }
      }
    }
    return false;
  }

  let freq: ValueCount;
  let isIPAddress = false;
  onMount(() => {
    freq = tokenValueFrequency(data, token);
    isIPAddress = isIPAddressToken(data, token);
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
  {#if token === timestampToken}
    <Activity {data} {token} />
  {:else}
    <Statistics {data} {token} />
    <Distribution {data} {token} />
    {#if freq !== undefined}
      <ValueFreqGraph {freq} />
    {/if}
    <OverTimeGraph {data} {token} {timestampToken} />
    {#if isIPAddress}
      <LocationMap {data} {token} />
    {/if}
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
