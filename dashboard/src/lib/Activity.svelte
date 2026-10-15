<script lang="ts">
  import { onDestroy } from "svelte";
  import Plotly from "plotly.js-dist-min";
  import { timestampRange, makeTimeSlots, nearestSlotIndex } from "./timeSlots";
  import { chartLayout, plotConfig } from "./plot";

  function buildPlot() {
    if (!plotDiv) return;
    const range = timestampRange(data, token);
    if (!range) { Plotly.purge(plotDiv); return; }
    const [min, max] = range;

    const timeSlots = makeTimeSlots(min, max, 20);
    const slotMs = timeSlots.map((s) => s.valueOf());

    const y: number[] = Array(timeSlots.length).fill(0);
    for (const extraction of data.extraction) {
      if (!(token in extraction.params)) continue;
      const ts = new Date(extraction.params[token].value as string).getTime();
      y[nearestSlotIndex(ts, slotMs)]++;
    }

    const x = timeSlots.map((s) => s.toDate());

    Plotly.react(
      plotDiv,
      [
        {
          x: x,
          y: y,
          type: "bar",
          marker: {
            color: hex,
          },
        },
      ],
      chartLayout({ height: 240, title: "Events over time" }),
      plotConfig
    );
  }

  let plotDiv = $state<HTMLDivElement>();
  let { data, token, hex }: { data: Data; token: string; hex: string } =
    $props();
  $effect(() => { if (plotDiv) buildPlot(); });
  onDestroy(() => { if (plotDiv) Plotly.purge(plotDiv); });
</script>

<div class="container">
  <div id="plotDiv" bind:this={plotDiv}></div>
</div>

<style scoped>
  .container {
    display: flex;
  }
  #plotDiv {
    width: 100%;
  }
</style>
