<script lang="ts">
  import { onMount } from "svelte";
  import Plotly from "plotly.js-dist-min";
  import { timestampRange, makeTimeSlots, nearestSlotIndex } from "./timeSlots";

  function buildPlot() {
    const range = timestampRange(data, token);
    if (!range) return;
    const [min, max] = range;

    const timeSlots = makeTimeSlots(min, max, 20);
    const slotMs = timeSlots.map((s) => new Date(s).getTime());

    const y: number[] = Array(timeSlots.length).fill(0);
    for (const extraction of data.extraction) {
      if (!(token in extraction.params)) continue;
      const ts = new Date(extraction.params[token].value as string).getTime();
      y[nearestSlotIndex(ts, slotMs)]++;
    }

    const x = timeSlots.map((s) => new Date(s));

    Plotly.newPlot(
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
      {
        title: false,
        hovermode: "closest",
        plot_bgcolor: "transparent",
        paper_bgcolor: "transparent",
        margin: { t: 0, l: 50, b: 30, r: 10 },
        autosize: true,
        dragmode: false,
        height: 300,
        yaxis: {
          gridcolor: "gray",
          showgrid: false,
        },
        xaxis: {
          fixedrange: true,
        },
      },
      { responsive: true, showSendToCloud: false, displayModeBar: false }
    );
  }

  $effect(() => {
    updatePlotColour();
  });

  function updatePlotColour() {
    if (plotDiv) {
      try {
        Plotly.restyle(plotDiv, {
          "marker.color": hex,
        });
      } catch (e) {}
    }
  }

  let plotDiv: HTMLDivElement;
  let { data, token, hex }: { data: Data; token: string; hex: string } =
    $props();
  onMount(async () => {
    buildPlot();
  });
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
