<script lang="ts">
  import { onMount } from "svelte";
  import Plotly from "plotly.js-dist-min";

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
            color: '#ffdfaf'
          },
        },
      ],
      {
        title: false,
        hovermode: "closest",
        plot_bgcolor: "transparent",
        paper_bgcolor: "transparent",
        margin: { t: 0, l: 40, b: 40, r: 20 },
        autosize: true,
        dragmode: false,
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

  let isNumeric = false;
  let plotDiv: HTMLDivElement;
  onMount(async () => {
    const values = numericValueCounts(data);
    if (Object.keys(values).length === 0) {
      return;
    }
    isNumeric = true;
    buildPlot(values);
    // Resize window to snap graph to window width
    window.dispatchEvent(new Event("resize"));
  });
  export let data: Data, token: string;
</script>

<div class="container" class:hidden={!isNumeric}>
  <div id="plotDiv" bind:this={plotDiv} />
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
