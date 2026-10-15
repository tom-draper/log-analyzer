<script lang="ts">
  import Plotly from "plotly.js-dist-min";
  import { chartLayout, plotConfig } from "./plot";

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
    Plotly.react(
      div,
      [{ x, y, type: "bar", marker: { color: hex } }],
      chartLayout({ height }),
      plotConfig
    );
  }

  function buildHourChart(div: HTMLDivElement, counts: number[]) {
    // Each hour = 15° (360/24). Rotation 90 puts midnight at top, clockwise matches a clock face.
    Plotly.react(
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
        ...chartLayout({ height: 240 }),
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

  let counts = $derived.by(() => {
    const months = Array(12).fill(0);
    const days = Array(7).fill(0);
    const hours = Array(24).fill(0);
    for (const ts of collectTimestamps()) {
      months[ts.getMonth()]++;
      days[ts.getDay()]++;
      hours[ts.getHours()]++;
    }
    return { months, days, hours };
  });
  let showMonths = $derived(counts.months.filter(Boolean).length >= 2);
  let showDays = $derived(counts.days.filter(Boolean).length >= 2);
  let showHours = $derived(counts.hours.filter(Boolean).length >= 2);
  let monthDiv = $state<HTMLDivElement>();
  let dayDiv = $state<HTMLDivElement>();
  let hourDiv = $state<HTMLDivElement>();

  $effect(() => {
    if (monthDiv) buildChart(monthDiv, MONTHS, counts.months);
  });
  $effect(() => {
    if (dayDiv) buildChart(dayDiv, DAYS, counts.days, 280);
  });
  $effect(() => {
    if (hourDiv) buildHourChart(hourDiv, counts.hours);
  });
  function cleanupChart(div: HTMLDivElement) {
    return { destroy: () => Plotly.purge(div) };
  }

  let { data, token, hex }: { data: Data; token: string; hex: string } = $props();
</script>

{#if showMonths || showDays || showHours}
  <div class="breakdown">
    {#if showMonths}
      <div class="chart">
        <div class="label">Month</div>
        <div bind:this={monthDiv} use:cleanupChart></div>
      </div>
    {/if}
    {#if showDays}
      <div class="chart">
        <div class="label">Day of week</div>
        <div bind:this={dayDiv} use:cleanupChart></div>
      </div>
    {/if}
    {#if showHours}
      <div class="chart">
        <div class="label">Hour of day</div>
        <div bind:this={hourDiv} use:cleanupChart></div>
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
