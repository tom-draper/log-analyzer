<script lang="ts">
  import { onMount } from "svelte";

  function getLocationCount(data: Data) {
    const locationCount: { [location: string]: number } = {};
    for (let i = 0; i < data.extraction.length; i++) {
      for (let [_token, value] of Object.entries(data.extraction[i].params)) {
        if (_token !== token) {
          continue;
        }
        const location = data.locations[value];
        if (!(location in locationCount)) {
          locationCount[location] = 0;
        }
        locationCount[location] += 1;
      }
    }
    return locationCount;
  }

  function unpackObject(obj: Object) {
    const locations: string[] = [];
    const z: number[] = [];
    for (let [k, v] of Object.entries(obj)) {
      locations.push(k);
      z.push(v);
    }
    return [locations, z];
  }

  function buildPlot() {
    const locationCount = getLocationCount(data);
    const [locations, z] = unpackObject(locationCount);

    let d = [
      {
        type: "choropleth",
        locations: locations,
        z: z,
        text: locations,
        locationmode: "country names",
        colorscale: [
          [0, '#0070f3'], [0.4, '#0070f3'], [1, '#111111']],
        autocolorscale: false,
        reversescale: true,
        marker: {
            line: {
                color: 'rgb(90,90,90)',
                width: 0.1
            }
        },
        tick0: 0,
        zmin: 0,
        colorbar: {
            autotic: false,
        }
      },
    ];

    console.log(locations, z);
    Plotly.newPlot(
      plotDiv,
      d,
      {
        title: false,
        hovermode: "closest",
        plot_bgcolor: "transparent",
        paper_bgcolor: "transparent",
        margin: { t: 50, l: 0, b: 0, r: 0 },
        autosize: true,
        height: 600,
        dragmode: false,
        yaxis: {
          gridcolor: "gray",
          showgrid: false,
        },
        xaxis: {
          fixedrange: true,
        },
        geo:{
              showframe: false,
                      countrycolor: 'rgb(255, 255, 255)',
        showland: true,
        landcolor: '#111',
        showlakes: true,
        lakecolor: 'rgb(255, 255, 255)',
        subunitcolor: 'rgb(255, 255, 255)',
        lonaxis: {},
        lataxis: {},
        bgcolor: 'transparent',
      },
      mapbox: {
        style: "dark"
      }
      },
      { responsive: true, showSendToCloud: false, displayModeBar: false }
    );
  }

  let plotDiv: HTMLDivElement;
  let Plotly;
  onMount(async () => {
    Plotly = await import("plotly.js-dist-min");
    buildPlot();
  });
  export let data: Data, token: string;
</script>

<div class="container">
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