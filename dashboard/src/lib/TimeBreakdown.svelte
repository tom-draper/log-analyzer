<script lang="ts">
  import { onMount, tick } from "svelte";
  import Plotly from "plotly.js-dist-min";

  const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
  const DAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

  function collectTimestamps(): Date[] {
    const values: Date[] = [];
    for (const extraction of data.extraction) {
      if (token in extraction.params) {
        values.push(new Date(extraction.params[token].value as string));
      }
    }
    return values;
  }

  function buildChart(div: HTMLDivElement, x: string[], y: number[], height = 200) {
    Plotly.newPlot(
      div,
      [{ x, y, type: "bar", marker: { color: hex } }],
      {
        plot_bgcolor: "transparent",
        paper_bgcolor: "transparent",
        margin: { t: 10, l: 40, b: 50, r: 10 },
        autosize: true,
        dragmode: false,
        height,
        yaxis: { gridcolor: "gray", showgrid: false },
        xaxis: { fixedrange: true },
      },
      { responsive: true, showSendToCloud: false, displayModeBar: false }
    );
  }

  function buildHourChart(div: HTMLDivElement, counts: number[]) {
    // Each hour = 15° (360/24). Rotation 90 puts midnight at top, clockwise matches a clock face.
    Plotly.newPlot(
      div,
      [
        {
          type: "barpolar",
          r: counts,
          theta: counts.map((_, i) => i * 15),
          width: Array(24).fill(15),
          marker: { color: hex },
        },
      ],
      {
        paper_bgcolor: "transparent",
        dragmode: false,
        height: 280,
        polar: {
          bgcolor: "transparent",
          angularaxis: {
            rotation: 90,
            direction: "clockwise",
            tickmode: "array",
            tickvals: [0, 90, 180, 270],
            ticktext: ["00:00", "06:00", "12:00", "18:00"],
            gridcolor: "#444",
            linecolor: "#444",
          },
          radialaxis: {
            showticklabels: false,
            showline: false,
            showgrid: false,
          },
        },
        margin: { t: 20, l: 20, b: 20, r: 20 },
      },
      { responsive: true, showSendToCloud: false, displayModeBar: false }
    );
  }

  let showMonths = $state(false);
  let showDays = $state(false);
  let showHours = $state(false);
  let monthCounts: number[];
  let dayCounts: number[];
  let hourCounts: number[];
  let monthDiv: HTMLDivElement = $state();
  let dayDiv: HTMLDivElement = $state();
  let hourDiv: HTMLDivElement = $state();

  onMount(async () => {
    const timestamps = collectTimestamps();

    monthCounts = Array(12).fill(0);
    dayCounts = Array(7).fill(0);
    hourCounts = Array(24).fill(0);

    for (const ts of timestamps) {
      monthCounts[ts.getMonth()]++;
      dayCounts[ts.getDay()]++;
      hourCounts[ts.getHours()]++;
    }

    showMonths = monthCounts.filter(Boolean).length >= 2;
    showDays = dayCounts.filter(Boolean).length >= 2;
    showHours = hourCounts.filter(Boolean).length >= 2;

    await tick();

    if (showMonths && monthDiv) buildChart(monthDiv, MONTHS, monthCounts);
    if (showDays && dayDiv) buildChart(dayDiv, DAYS, dayCounts, 280);
    if (showHours && hourDiv) buildHourChart(hourDiv, hourCounts);
  });

  $effect(() => {
    updateColours();
  });

  function updateColours() {
    for (const div of [monthDiv, dayDiv, hourDiv]) {
      if (div) {
        try {
          Plotly.restyle(div, { "marker.color": hex });
        } catch (e) {}
      }
    }
  }

  let { data, token, hex }: { data: Data; token: string; hex: string } = $props();
</script>

{#if showMonths || showDays || showHours}
  <div class="breakdown">
    {#if showMonths}
      <div class="chart">
        <div class="label">Month</div>
        <div bind:this={monthDiv}></div>
      </div>
    {/if}
    {#if showDays}
      <div class="chart">
        <div class="label">Day of week</div>
        <div bind:this={dayDiv}></div>
      </div>
    {/if}
    {#if showHours}
      <div class="chart">
        <div class="label">Hour of day</div>
        <div bind:this={hourDiv}></div>
      </div>
    {/if}
  </div>
{/if}

<style scoped>
  .breakdown {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    margin-top: 1.5rem;
  }

  .label {
    font-size: 0.8em;
    color: #888;
    margin-bottom: 0.25rem;
  }
</style>
