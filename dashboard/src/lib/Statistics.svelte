<script lang="ts">
  import { onMount } from "svelte";

  function asc(arr: number[]) {
    return arr.sort((a, b) => a - b);
  }

  function quantile(arr: number[], q: number) {
    const sorted = asc(arr);
    const pos = (sorted.length - 1) * q;
    const base = Math.floor(pos);
    const rest = pos - base;
    if (sorted[base + 1] != undefined) {
      return sorted[base] + rest * (sorted[base + 1] - sorted[base]);
    } else if (sorted[base] != undefined) {
      return sorted[base];
    } else {
      return 0;
    }
  }

  function isNumericField(data: Data): boolean {
    for (let i = 0; i < data.extraction.length; i++) {
      if (!(token in data.extraction[i].params)) {
        continue;
      }
      const value = data.extraction[i].params[token].value;
      if (typeof value === "number") {
        return true;
      }
    }
    return false;
  }

  function sortedValues(data: Data): number[] {
    const values: number[] = [];
    for (let i = 0; i < data.extraction.length; i++) {
      if (!(token in data.extraction[i].params)) {
        continue;
      }
      const value = data.extraction[i].params[token].value;
      if (typeof value === "number") {
        values.push(value);
      }
    }

    return values.sort();
  }

  type Statistics = {
    median: number;
    lq: number;
    uq: number;
    p95: number;
    p99: number;
  };

  let isNumeric = false;
  let statistics: Statistics;
  onMount(() => {
    isNumeric = isNumericField(data);
    const values = sortedValues(data);
    statistics = {
      median: values[Math.floor(values.length / 2)],
      lq: quantile(values, 0.25),
      uq: quantile(values, 0.75),
      p95: quantile(values, 0.95),
      p99: quantile(values, 0.99),
    };
  });
  export let data: Data, token: string;
</script>

{#if isNumeric}
  <div class="container">
    <div class="statistic">
      <div class="value">{statistics.lq}</div>
      <div class="label">LQ</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.median}</div>
      <div class="label">Median</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.uq}</div>
      <div class="label">UQ</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.p95}</div>
      <div class="label">p95</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.p99}</div>
      <div class="label">p99</div>
    </div>
  </div>
{/if}

<style scoped>
  .container {
    border: 1px solid #ffffff24;
    border-radius: 5px;
    margin: 3em 0 2em;
    padding: 2rem;
    display: flex;
  }
  .statistic {
    flex: 1;
    font-size: 1.6em;
    display: grid;
    place-items: center;
  }
  .value {
    font-weight: 600;
  }
  .label {
    font-size: 0.62em;
    color: #888;
  }
</style>
