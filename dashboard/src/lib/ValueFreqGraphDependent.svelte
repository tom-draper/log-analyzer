<script lang="ts">
  import { onMount } from "svelte";

  type DependentTokenValueCounts = {
    [token: string]: {
      [dependentToken: string]: {
        [tokenValue: string]: { [dependentTokenValue: string]: number };
      };
    };
  };

  type Bar = { dependentTokenValue: string; count: number; width: number };

  function sortedBars(freq: DependentTokenValueCounts) {
    let sortedFreq: { [tokenValue: string]: Bar[] } = {};
    for (let tokenCounts of Object.values(freq)) {
      for (let valueCounts of Object.values(tokenCounts)) {
        for (let [tokenValue, dependentTokenValueCounts] of Object.entries(
          valueCounts
        )) {
          for (let [dependentTokenValue, count] of Object.entries(
            dependentTokenValueCounts
          )) {
            if (!(tokenValue in sortedFreq)) {
              sortedFreq[tokenValue] = [];
            }
            sortedFreq[tokenValue].push({
              dependentTokenValue,
              count,
              width: 0,
            });
          }
        }
      }
    }

    console.log(sortedFreq);

    // Sort descending
    for (const token of Object.keys(sortedFreq)) {
      sortedFreq[token].sort((a, b) => {
        return b.count - a.count;
      });
      // Cap at 10 values
      sortedFreq[token] = sortedFreq[token].slice(0, numBars);
    }

    let maxBar = 0;
    for (const token of Object.keys(sortedFreq)) {
      for (let i = 0; i < sortedFreq[token].length; i++) {
        if (sortedFreq[token][i].count > maxBar) {
          maxBar = sortedFreq[token][i].count;
        }
      }
    }

    // Set widths
    for (const token of Object.keys(sortedFreq)) {
      for (let i = 0; i < sortedFreq[token].length; i++) {
        sortedFreq[token][i].width =
          (sortedFreq[token][i].count / maxBar) * 100;
      }
    }
    let bars: { [token: string]: Bar[] }[] = [];
    for (const token of Object.keys(sortedFreq)) {
      const bar = {};
      bar[token] = sortedFreq[token];
      bars.push(bar);
    }
    bars.sort((a, b) => {
      let aTotal = 0;
      for (const key of Object.keys(a)) {
        for (let i = 0; i < a[key].length; i++) {
          aTotal += a[key][i].count;
        }
      }

      let bTotal = 0;
      for (const key of Object.keys(b)) {
        for (let i = 0; i < b[key].length; i++) {
          bTotal += b[key][i].count;
        }
      }

      return bTotal - aTotal;
    });

    return bars;
  }

  function tokenValueFrequency(
    data: Data,
    token: string,
    dependentToken: string
  ) {
    const freq: DependentTokenValueCounts = {};
    for (let i = 0; i < data.extraction.length; i++) {
      if (
        token in data.extraction[i].params &&
        dependentToken in data.extraction[i].params
      ) {
        if (!(token in freq)) {
          freq[token] = {};
        }
        if (!(dependentToken in freq[token])) {
          freq[token][dependentToken] = {};
        }

        const tokenValue = data.extraction[i].params[token].value;
        if (!(tokenValue in freq[token][dependentToken])) {
          freq[token][dependentToken][tokenValue] = {};
        }
        const dependentTokenValue =
          data.extraction[i].params[dependentToken].value;
        if (!(dependentTokenValue in freq[token][dependentToken][tokenValue])) {
          freq[token][dependentToken][tokenValue][dependentTokenValue] = 0;
        }
        freq[token][dependentToken][tokenValue][dependentTokenValue] += 1;
      }
    }

    return freq;
  }

  const numBars = 10;
  let bars: {}[];
  function setBars(freq: DependentTokenValueCounts) {
    bars = sortedBars(freq);
  }

  onMount(() => {
    const freq = tokenValueFrequency(data, token, dependentToken);
    setBars(freq);
  });

  export let data: Data, token: string, dependentToken: string;
</script>

{#if bars !== undefined}
  <div class="freq-graph">
    {#each bars.slice(0, 10) as group}
      {#each Object.keys(group) as tokenValue}
        {tokenValue}
        {#each group[tokenValue] as bar}
          <div class="token-frequency" title={bar.count}>
            <div class="bar" style="width: {bar.width}%">
              {bar.dependentTokenValue}
            </div>
          </div>
        {/each}
      {/each}
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
