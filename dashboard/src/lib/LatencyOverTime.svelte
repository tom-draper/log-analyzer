<script lang="ts">
  import { onMount } from "svelte";
  import Plotly from "plotly.js-dist-min";
  import { makeTimeSlots, nearestSlotIndex, timestampRange } from "./timeSlots";
  import { chartLayout, plotConfig } from "./plot";

  const quantile = (values: number[], q: number) => {
    const sorted = values.slice().sort((a, b) => a - b);
    const position = (sorted.length - 1) * q;
    const lower = Math.floor(position);
    const upper = Math.ceil(position);
    return sorted[lower] + (sorted[upper] - sorted[lower]) * (position - lower);
  };

  function buildPlot() {
    const range = timestampRange(data, timestampToken);
    if (!range || !timestampToken) return;
    const slots = makeTimeSlots(range[0], range[1], 20);
    const slotMillis = slots.map((slot) => slot.valueOf());
    const values = Array.from({ length: slots.length }, () => [] as number[]);
    for (const extraction of data.extraction) {
      const latency = extraction.params[token]?.value;
      const timestamp = extraction.params[timestampToken]?.value;
      if (typeof latency !== "number" || !timestamp) continue;
      values[nearestSlotIndex(new Date(timestamp as string).getTime(), slotMillis)].push(latency);
    }
    if (!values.some((slot) => slot.length > 0)) return;
    const x = slots.map((slot) => slot.toDate());
    const percentile = (q: number) => values.map((slot) => slot.length ? quantile(slot, q) : null);
    Plotly.newPlot(plotDiv, [
      { x, y: percentile(.5), name: "p50", mode: "lines", line: { color: "#aaa" } },
      { x, y: percentile(.95), name: "p95", mode: "lines", line: { color: hex } },
      { x, y: percentile(.99), name: "p99", mode: "lines", line: { color: "#dd7178" } },
    ], { ...chartLayout({ title: "Latency percentiles over time", height: 240 }), yaxis: { ...chartLayout().yaxis, title: { text: token } } }, plotConfig);
  }

  let plotDiv: HTMLDivElement;
  let { data, token, timestampToken, hex }: { data: Data; token: string; timestampToken: string | null; hex: string } = $props();
  onMount(() => { buildPlot(); return () => Plotly.purge(plotDiv); });
</script>

<div class="chart" bind:this={plotDiv}></div>

<style scoped>.chart { width: 100%; margin: 1.5rem 0; }</style>
