<script lang="ts">
  import { onDestroy } from "svelte";
  import Plotly from "plotly.js-dist-min";
  import { chartLayout, plotConfig } from "./plot";

  function numericValues(data: Data) {
    const values: number[] = [];
    for (const extraction of data.extraction) {
      const value = extraction.params[token]?.value;
      if (typeof value === "number" && Number.isFinite(value)) values.push(value);
    }
    return values;
  }

  function buildPlot(values: number[]) {
    if (!plotDiv) return;
    const min = Math.min(...values);
    const max = Math.max(...values);
    const binCount = Math.max(5, Math.min(40, Math.ceil(Math.sqrt(values.length))));
    const size = max === min ? 1 : (max - min) / binCount;

    Plotly.react(
      plotDiv,
      [
        {
          x: values,
          type: "histogram",
          xbins: { start: min, end: max, size },
          marker: { color: hex },
          hovertemplate: "%{x}<br>%{y} records<extra></extra>",
        },
      ],
      { ...chartLayout({ title: "Distribution" }), bargap: 0.08, xaxis: { ...chartLayout().xaxis, title: { text: token } } },
      plotConfig
    );
  }

  let plotDiv = $state<HTMLDivElement>();
  let { data, token, hex }: { data: Data; token: string; hex: string } = $props();
  let values = $derived(numericValues(data));
  let isNumeric = $derived(values.length > 0);
  $effect(() => {
    if (plotDiv && isNumeric) buildPlot(values);
    else if (plotDiv) Plotly.purge(plotDiv);
  });
  onDestroy(() => { if (plotDiv) Plotly.purge(plotDiv); });
</script>

<div class="container" class:hidden={!isNumeric}>
  <div id="plotDiv" bind:this={plotDiv}></div>
</div>

<style scoped>
  .container {
    display: flex;
    margin-bottom: 1.4em;
  }

  #plotDiv {
    width: 100%;
  }
</style>
