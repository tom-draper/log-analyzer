<script lang="ts">
  import { onMount } from "svelte";

  function sortedBars(freq: ValueFreq): [string, number, number][] {
    let sortedFreq: [string, number, number][] = [];
    for (let [token, count] of Object.entries(freq)) {
      sortedFreq.push([token, count, 0]);
    }

    // Sort descending
    sortedFreq.sort((a, b) => {
      return b[1] - a[1];
    });

    // Create "others" bar
    let othersTotal = 0;
    for (let i = numBars; i < sortedFreq.length; i++) {
      othersTotal += sortedFreq[i][1]
    }

    // Cap at 10 values
    sortedFreq = sortedFreq.slice(0, numBars)
    if (othersTotal > 0) {
      sortedFreq.push(['Others', othersTotal, 0])
    }

    let maxBar = 0;
    for (let i = 0; i < sortedFreq.length; i++) {
      if (sortedFreq[i][1] > maxBar) {
        maxBar = sortedFreq[i][1]
      }
    }
    
    // Set widths
    for (let i = 0; i < sortedFreq.length; i++) {
      sortedFreq[i][2] = (sortedFreq[i][1] / maxBar) * 97
    }

    return sortedFreq;
  }

  const numBars = 10
  let bars: [string, number, number][] = [];
  onMount(() => {
    bars = sortedBars(freq);
    console.log(bars)
  });

  export let freq: ValueFreq;
</script>

{#if bars.length > 0}
  <div class="freq-graph">
    {#each bars as bar}
      <div class="token-frequency" title="{bar[1]}">
        <!-- <div class="value-name"></div>   -->
        <div class="bar" style="width: {bar[2]}%">{bar[0]}</div>
      </div>
    {/each}
  </div>
{/if}

<style scoped>
  .token-frequency {
    margin: 2px 0;
  }
  .freq-graph {
    overflow: auto;
  }
  .bar {
    /* height: px; */
    background: #0070f3;
    border-radius: 4px;
    margin: 5px 0;
    padding: 1px 10px;
    font-size: 0.9em;
    font-weight: 500;
    text-wrap: nowrap;
  }
  .token-frequency {
    position: relative;
  }
  .value-name {
    color: white;
    position: absolute;
  }
</style>
