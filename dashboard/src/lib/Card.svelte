<script lang="ts">
  import { onMount } from "svelte";
  import ValueFreqGraph from "./ValueFreqGraph.svelte";

  function tokenValueFrequency(data: Data, token: string): ValueFreq {
    let freq: ValueFreq = {};
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

  let freq: ValueFreq;
  onMount(() => {
    freq = tokenValueFrequency(data, token);
    console.log(freq);
  });

  // your script goes here
  export let data: Data, token: string, timestampToken: string | null;
</script>

<div class="card">
  <div class="title">{token}</div>
  {#if freq !== undefined}
    <ValueFreqGraph {freq} />
  {/if}
</div>

<style scoped>
  .title {
    margin-bottom: 20px;
  }
  .card {
    border: 1px solid #ffffff24;
    border-radius: 5px;
    margin: 3em 0;
  }
</style>
