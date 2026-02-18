<script lang="ts">
  import { onMount } from "svelte";

  function quantile(arr: number[], q: number) {
    const sorted = arr.slice().sort((a, b) => a - b);
    const pos = (sorted.length - 1) * q;
    const base = Math.floor(pos);
    const rest = pos - base;
    if (sorted[base + 1] !== undefined) {
      return sorted[base] + rest * (sorted[base + 1] - sorted[base]);
    } else if (sorted[base] !== undefined) {
      return sorted[base];
    }
    return 0;
  }

  function numericValues(data: Data, token: string) {
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

    return values;
  }

  type Statistics = {
    min: number;
    lq: number;
    median: number;
    uq: number;
    max: number;
    p95: number;
    p99: number;
  };

  let statistics: Statistics;
  onMount(() => {
    const values = numericValues(data, token);
    if (values.length == 0) return;
    values.sort((a, b) => a - b);
    statistics = {
      min: values[0],
      lq: quantile(values, 0.25),
      median: quantile(values, 0.5),
      uq: quantile(values, 0.75),
      max: values[values.length - 1],
      p95: quantile(values, 0.95),
      p99: quantile(values, 0.99),
    };
  });
  export let data: Data, token: string;
</script>

{#if statistics}
  <div class="container">
    <div class="statistic">
      <div class="value">{statistics.min.toLocaleString()}</div>
      <div class="label">Min</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.lq.toLocaleString()}</div>
      <div class="label">LQ</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.median.toLocaleString()}</div>
      <div class="label">Median</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.uq.toLocaleString()}</div>
      <div class="label">UQ</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.max.toLocaleString()}</div>
      <div class="label">Max</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.p95.toLocaleString()}</div>
      <div class="label">p95</div>
    </div>
    <div class="statistic">
      <div class="value">{statistics.p99.toLocaleString()}</div>
      <div class="label">p99</div>
    </div>
  </div>
{/if}

<style scoped>
  .container {
    border: 1px solid #ffffff24;
    border: 1px solid rgb(42, 42, 42);
    border-radius: 5px;
    margin: 3em 0 2em;
    padding: 2rem;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem 0;
  }
  .statistic {
    flex: 1;
    min-width: 80px;
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
