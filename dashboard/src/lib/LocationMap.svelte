<script lang="ts">
  import { onDestroy } from "svelte";
  import Plotly from "plotly.js-dist-min";
  import { chartLayout, plotConfig } from "./plot";

  function locationSummary(data: Data) {
    const locationCount: { [location: string]: number } = {};
    let total = 0;
    let mapped = 0;
    for (const extraction of data.extraction) {
      const param = extraction.params[token];
      if (param?.type === "ip") {
        total++;
        const location = data.locations[String(param.value)];
        if (location) {
          locationCount[location] = (locationCount[location] ?? 0) + 1;
          mapped++;
        }
      }
    }
    return { locationCount, total, mapped };
  }

  function unpackObject(obj: Object) {
    const locations: string[] = [];
    const z: number[] = [];
    for (const [k, v] of Object.entries(obj)) {
      locations.push(k);
      z.push(v);
    }
    return [locations, z];
  }

  function buildPlot() {
    if (!plotDiv) return;
    const { locationCount } = locationSummary(data);
    const [locations, z] = unpackObject(locationCount);

    plotData = [
      {
        type: "choropleth",
        locations: locations,
        z: z,
        text: locations,
        locationmode: "country names",
        colorscale: [
          [0, hex],
          [0.4, hex],
          [1, "#101010"],
        ],
        autocolorscale: false,
        reversescale: true,
        marker: {
          line: {
            color: "rgb(90,90,90)",
            width: 0.1,
          },
        },
        tick0: 0,
        zmin: 0,
        colorbar: {
          autotic: false,
        },
      },
    ];

    Plotly.react(
      plotDiv,
      plotData,
      {
        ...chartLayout({ height: 360, title: "Geographic distribution" }),
        hovermode: "closest",
        margin: { t: 50, l: 0, b: 0, r: 0 },
        geo: {
          showframe: false,
          countrycolor: "rgb(255, 255, 255)",
          showland: true,
          landcolor: "#111",
          showlakes: true,
          lakecolor: "black",
          lonaxis: {},
          lataxis: {},
          bgcolor: "transparent",
        },
      },
      plotConfig
    );
  }

  let plotDiv = $state<HTMLDivElement>();
  let plotData: any;
  let { data, token, hex }: { data: Data; token: string; hex: string } = $props();
  let coverage = $derived(locationSummary(data));
  $effect(() => { if (plotDiv) buildPlot(); });
  onDestroy(() => { if (plotDiv) Plotly.purge(plotDiv); });
</script>

<div class="container" class:hidden={coverage.mapped === 0}>
  <div id="plotDiv" bind:this={plotDiv}></div>
  <div class="coverage">Mapped {coverage.mapped.toLocaleString()} of {coverage.total.toLocaleString()} IP values ({coverage.total === 0 ? "0" : ((coverage.mapped / coverage.total) * 100).toFixed(1)}%)</div>
</div>

<style scoped>
  .container {
    display: flex;
    margin-bottom: 1.4em;
  }
  #plotDiv {
    width: 100%;
  }
  .coverage { color: var(--muted); font-size: .75rem; margin-top: -.5rem; }
</style>
