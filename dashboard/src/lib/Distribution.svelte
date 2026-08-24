<script lang="ts">
  import { onMount } from "svelte";
  import Plotly from "plotly.js-dist-min";
  import { chartLayout, plotConfig } from "./plot";

  function numericValueCounts(data: Data) {
    const values: ValueCount = {};
    for (let i = 0; i < data.extraction.length; i++) {
      if (!(token in data.extraction[i].params)) {
        continue;
      }
      const value = data.extraction[i].params[token];
      if (typeof value.value !== "number") {
        continue;
      }
      values[value.value] ||= 0;
      values[value.value] += 1;
    }

    return values;
  }

  function buildPlot(values: ValueCount) {
    const sortedValues = Object.keys(values).sort();
    const y = Array(sortedValues.length).fill(0);
    for (const value of sortedValues) {
      y[sortedValues.indexOf(value)] = values[value];
    }

    Plotly.newPlot(
      plotDiv,
      [
        {
          x: sortedValues,
          y: y,
          type: "bar",
          marker: {
            // color: "#0070f3",
            color: hex,
          },
        },
      ],
      chartLayout({ title: "Distribution" }),
      plotConfig
    );
  }

  $effect(() => { updatePlotColour(); });

  function updatePlotColour() {
    if (plotDiv) {
      try {
        Plotly.restyle(plotDiv, {
          "marker.color": hex,
        });
      } catch (e) {
        // console.error(e);
      }
    }
  }

  let isNumeric = $state(false);
  let plotDiv: HTMLDivElement;
  let { data, token, hex }: { data: Data; token: string; hex: string } = $props();
  onMount(() => {
    const values = numericValueCounts(data);
    if (Object.keys(values).length === 0) {
      return;
    }
    isNumeric = true;
    buildPlot(values);
    // Resize window to snap graph to window width
    window.dispatchEvent(new Event("resize"));
    return () => Plotly.purge(plotDiv);
  });
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
