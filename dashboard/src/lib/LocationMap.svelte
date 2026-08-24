<script lang="ts">
  import { onMount } from "svelte";
  import Plotly from "plotly.js-dist-min";
  import { chartLayout, plotConfig } from "./plot";

  function isIPAddressToken(data: Data, token: string) {
    for (let i = 0; i < data.extraction.length; i++) {
      const params = data.extraction[i].params;
      if (
        token in params &&
        params[token].type === "ip" &&
        params[token].value in data.locations
      ) {
        return true;
      }
    }
    return false;
  }

  function getLocationCount(data: Data) {
    const locationCount: { [location: string]: number } = {};
    for (let i = 0; i < data.extraction.length; i++) {
      const params = data.extraction[i].params;
      if (token in params && params[token].type === "ip") {
        const location = data.locations[params[token].value];
        locationCount[location] ||= 0;
        locationCount[location] += 1;
      }
    }
    return locationCount;
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
    const locationCount = getLocationCount(data);
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

    Plotly.newPlot(
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
        mapbox: {
          style: "dark",
        },
      },
      plotConfig
    );
  }

  $effect(() => {
    const color = hex;
    if (plotDiv && plotData) {
      Plotly.restyle(plotDiv, {
        colorscale: [[[0, color], [0.4, color], [1, "#101010"]]],
      } as any);
    }
  });

  let plotDiv: HTMLDivElement;
  let plotData: any;
  let { data, token, hex }: { data: Data; token: string; hex: string } = $props();
  onMount(() => {
    buildPlot();
    return () => Plotly.purge(plotDiv);
  });
</script>

<div class="container" class:hidden={!isIPAddressToken(data, token)}>
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
