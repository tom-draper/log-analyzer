<script lang="ts">
  import { onMount } from "svelte";

  type Bar = {token: string, count: number, width: number}

  function sortedBars(freq: ValueCount): Bar[] {
    let sortedFreq: Bar[] = [];
    for (let [token, count] of Object.entries(freq)) {
      sortedFreq.push({token, count, width: 0});
    }

    // Sort descending
    sortedFreq.sort((a, b) => {
      return b.count - a.count;
    });

    // Create "others" bar
    let othersTotal = 0;
    for (let i = numBars; i < sortedFreq.length; i++) {
      othersTotal += sortedFreq[i].count
    }

    // Cap at 10 values
    sortedFreq = sortedFreq.slice(0, numBars)
    if (othersTotal > 0) {
      sortedFreq.push({token: 'Others', count: othersTotal, width: 0})
    }

    let maxBar = 0;
    for (let i = 0; i < sortedFreq.length; i++) {
      if (sortedFreq[i].count > maxBar) {
        maxBar = sortedFreq[i].count
      }
    }
    
    // Set widths
    for (let i = 0; i < sortedFreq.length; i++) {
      sortedFreq[i].width = (sortedFreq[i].count / maxBar) * 100
    }

    return sortedFreq;
  }

  const numBars = 10
  let bars: Bar[];
  onMount(() => {
    bars = sortedBars(freq);
  });

  export let freq: ValueCount;
</script>

{#if bars !== undefined}
  <div class="freq-graph">
    {#each bars as bar}
      <div class="token-frequency" title="{bar.count.toLocaleString()}">
        <!-- <div class="value-name"></div>   -->
        <div class="bar" style="width: {bar.width}%">{bar.token}</div>
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
    background: #0070f3;
    border-radius: 4px;
    margin: 5px 0;
    padding: 1px 10px;
    font-size: 0.85em;
    font-weight: 500;
    text-wrap: nowrap;
    box-sizing: border-box;
  }
  .token-frequency {
    position: relative;
  }
</style>
