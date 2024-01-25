<script lang="ts">
  import { onMount } from "svelte";
  import Moment from "moment";
  import { extendMoment } from "moment-range";
  import Plotly from "plotly.js-dist-min";

  const moment = extendMoment(Moment);

  function timestampValues(data: Data) {
    const values: Date[] = [];
    for (let i = 0; i < data.extraction.length; i++) {
      if (!(token in data.extraction[i].params)) {
        continue;
      }
      const value = data.extraction[i].params[token].value;
      values.push(new Date(value));
    }

    return values;
  }

  function bars(timestamps: Date[], timeSlots: Moment.Moment[]) {
    const timeSlotTimestamps = timeSlots.map((timeSlot) => {
      return new Date(timeSlot).getTime();
    });

    const y: number[] = Array(timeSlots.length).fill(0);
    for (let i = 0; i < timestamps.length; i++) {
      // Find timeslot index
      const best = {
        index: -1,
        diff: Number.MAX_VALUE,
      };
      const timestamp = timestamps[i].getTime();
      for (let j = 0; j < timeSlotTimestamps.length; j++) {
        const diff = Math.abs(timeSlotTimestamps[j] - timestamp);
        if (diff < best.diff) {
          best.index = j;
          best.diff = diff;
        } else {
          // Assuming timestamps sorted, we've found the closest time slot time
          break;
        }
      }
      y[best.index] += 1;
    }

    const x = timeSlots.map((timeSlot) => {
      return new Date(timeSlot);
    });

    return [x, y] as const;
  }

  function buildPlot() {
    const values = timestampValues(data).sort((a, b) => {
      return a.getTime() - b.getTime();
    });
    const dateRange = moment.range(values[0], values[values.length - 1]);
    const timeSlots = Array.from(dateRange.by("minutes", { step: 20 }));
    const [x, y] = bars(values, timeSlots);

    Plotly.newPlot(
      plotDiv,
      [
        {
          x: x,
          y: y,
          type: "bar",
          marker: {
            // color: "#0070f3",
            color: "#ffdfaf",
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

  let plotDiv: HTMLDivElement;
  onMount(async () => {
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
    /* margin-bottom: 0.5em; */
  }
  #plotDiv {
    width: 100%;
  }
</style>
