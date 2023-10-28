<script lang="ts">
  import { onMount } from "svelte";

  function isNumericField(data: Data): boolean {
    for (let i = 0; i < data.extraction.params.length; i++) {
      if (!(token in data.extraction.params[i])) {
        continue;
      }
      const value = data.extraction.params[i][token];
      if (typeof value === "number") {
        return true;
      }
    }
    return false;
  }

  function valueCounts(data: Data): ValueCount {
    let values: ValueCount = {};
    for (let i = 0; i < data.extraction.params.length; i++) {
      if (!(token in data.extraction.params[i])) {
        continue;
      }
      const value = data.extraction.params[i][token];
      if (typeof value !== "number") {
        continue;
      }
      if (!(value in values)) {
        values[value] = 0;
      }
      values[value] += 1;
    }

    return values;
  }

  function buildPlot() {
    const sortedValues = Object.keys(values).sort();
    const y = Array(sortedValues.length).fill(0);
    for (let value of sortedValues) {
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
            color: "#0070f3",
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
  let values: ValueCount;
  let plotDiv: HTMLDivElement;
  let Plotly;
  onMount(async () => {
    Plotly = await import("plotly.js-dist-min");
    isNumeric = isNumericField(data);
    values = valueCounts(data);
    setTimeout(buildPlot, 100);
  });
  export let data: Data, token: string;
</script>

{#if isNumeric}
  <div class="container">
    <div id="plotDiv" bind:this={plotDiv} />
  </div>
{/if}

<style scoped>
  .container {
    display: flex;
    margin-bottom: 1.4em;
  }
  #plotDiv {
    width: 100%;
  }
</style>
