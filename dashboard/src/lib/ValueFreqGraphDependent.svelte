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
    const sortedFreq: { [tokenValue: string]: Bar[] } = {};
    for (const tokenCounts of Object.values(freq)) {
      for (const valueCounts of Object.values(tokenCounts)) {
        for (const [tokenValue, dependentTokenValueCounts] of Object.entries(
          valueCounts
        )) {
          for (const [dependentTokenValue, count] of Object.entries(
            dependentTokenValueCounts
          )) {
            tokenValue in sortedFreq || (sortedFreq[tokenValue] = []);
            sortedFreq[tokenValue].push({
              dependentTokenValue,
              count,
              width: 0,
            });
          }
        }
      }
    }

    // Sort descending
    for (const token in sortedFreq) {
      sortedFreq[token].sort((a, b) => {
        return b.count - a.count;
      });
    }

    let maxBar = 0;
    for (const token in sortedFreq) {
      for (let i = 0; i < sortedFreq[token].length; i++) {
        if (sortedFreq[token][i].count > maxBar) {
          maxBar = sortedFreq[token][i].count;
        }
      }
    }

    // Set widths
    for (const token in sortedFreq) {
      for (let i = 0; i < sortedFreq[token].length; i++) {
        sortedFreq[token][i].width =
          (sortedFreq[token][i].count / maxBar) * 100;
      }
      // Cap at 10 values
      sortedFreq[token] = sortedFreq[token].slice(0, numBars);
    }

    const tokenBarsCounts: { [token: string]: number } = {};
    for (const token in sortedFreq) {
      let tokenBarsTotal = 0;
      for (const bar of sortedFreq[token]) {
        tokenBarsTotal += bar.count;
      }
      tokenBarsCounts[token] = tokenBarsTotal;
    }

    const bars: { [token: string]: Bar[] }[] = [];
    for (const token in sortedFreq) {
      const tokenBars: (typeof bars)[number] = {};
      tokenBars[token] = sortedFreq[token];
      bars.push(tokenBars);
    }

    bars.sort((tokenBars1, tokenBars2) => {
      let aTotal = 0;
      for (const token in tokenBars1) {
        aTotal += tokenBarsCounts[token];
      }

      let bTotal = 0;
      for (const token in tokenBars2) {
        bTotal += tokenBarsCounts[token];
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
        token in freq || (freq[token] = {});
        dependentToken in freq[token] || (freq[token][dependentToken] = {});
        const tokenValue = data.extraction[i].params[token].value;
        tokenValue in freq[token][dependentToken] ||
          (freq[token][dependentToken][tokenValue] = {});
        const dependentTokenValue =
          data.extraction[i].params[dependentToken].value;
        dependentTokenValue in freq[token][dependentToken][tokenValue] ||
          (freq[token][dependentToken][tokenValue][dependentTokenValue] = 0);
        freq[token][dependentToken][tokenValue][dependentTokenValue] += 1;
      }
    }

    return freq;
  }

  const numBars = 10;
  let bars: {
    [token: string]: Bar[];
  }[] = [];
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
        <div class="value-name">
          {tokenValue}
        </div>
        {#each group[tokenValue] as bar}
          <div class="token-frequency" title={bar.count.toString()}>
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
    background: #e2b269;
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
  .value-name {
    font-size: 0.85em;
    margin-top: 5px;
    margin-bottom: -2px;
  }
</style>
