<script lang="ts">
  import { onMount } from "svelte";
  import Moment from "moment";
  import { extendMoment } from "moment-range";

  const moment = extendMoment(Moment);

  type SortableValueCount = {
    value: string;
    total: number;
  }

  function valueByTimeSlot(
    data: Data,
    token: string,
    timestampToken: string | null,
    timeSlots: Moment.Moment[]
  ): ValueCounts {
    if (timestampToken === null) {
      return {};
    }

    const timeSlotTimestamps = timeSlots.map((timeSlot) => {
      return new Date(timeSlot).getTime();
    });

    const days: ValueCounts = {};
    for (let i = 0; i < data.extraction.params.length; i++) {
      const params = data.extraction.params;
      if (!(token in params[i]) || !(timestampToken in params[i])) {
        continue;
      }

      const timestamp = new Date(params[i][timestampToken]).getTime();

      // Find timeslot index
      const best = {
        index: -1,
        diff: Number.MAX_VALUE,
      };
      for (let j = 0; j < timeSlotTimestamps.length; j++) {
        const diff = Math.abs(timeSlotTimestamps[j] - timestamp);
        if (diff < best.diff) {
          best.index = j
          best.diff = diff
        } else {
          break;
        }
      }

      const value = params[i][token];
      if (!(value in days)) {
        days[value] = Array(timeSlots.length).fill(0);
      }

      days[value][best.index] += 1;
    }

    return days;
  }

  function sortedValueByTimeSlot(valueByTimeSlot: ValueCounts): SortableValueCount[] {
    let sortedValues: { value: string; total: number }[] = [];
    for (let [value, counts] of Object.entries(valueByTimeSlot)) {
      const total = counts.reduce((partialSum, a) => partialSum + a, 0);
      sortedValues.push({ value, total });
    }

    sortedValues.sort((a, b) => {
      return b.total - a.total;
    });

    return sortedValues;
  }

  function valueCountMax(days: any): ValueCount {
    const valueMax: ValueCount = {};
    for (let value of Object.keys(days)) {
      valueMax[value] = Math.max(...days[value]);
    }

    return valueMax;
  }

  function timestampRange(
    data: Data,
    timestampToken: string | null
  ): [Date | null, Date | null] {
    if (timestampToken === null) {
      return [null, null];
    }

    let maxDate = new Date(-8640000000000000);
    let minDate = new Date(8640000000000000);

    for (let i = 0; i < data.extraction.params.length; i++) {
      if (!(timestampToken in data.extraction.params[i])) {
        continue;
      }
      const timestamp = new Date(data.extraction.params[i][timestampToken]);
      if (timestamp > maxDate) {
        maxDate = timestamp;
      }
      if (timestamp < minDate) {
        minDate = timestamp;
      }
    }

    return [minDate, maxDate];
  }

  const numTiles = 50;
  let sortedValueCounts: SortableValueCount[];
  let valueCounts: ValueCounts;
  let valueMax: ValueCount;
  let timeSlots: Moment.Moment[];
  onMount(() => {
    const [minDate, maxDate] = timestampRange(data, timestampToken);
    if (minDate === null || maxDate === null) {
      return;
    }
    const dateRange = moment.range(minDate, maxDate);
    timeSlots = Array.from(dateRange.by("minutes", { step: numTiles }));

    valueCounts = valueByTimeSlot(data, token, timestampToken, timeSlots);
    sortedValueCounts = sortedValueByTimeSlot(valueCounts);
    valueMax = valueCountMax(valueCounts);
  });

  export let data: Data, token: string, timestampToken: string | null;
</script>

{#if valueMax != undefined}
  <div class="freq-graph">
    {#each sortedValueCounts.slice(0, 10) as value}
      <div class="value-name">{value.value}</div>
      <div class="days">
        {#each timeSlots as slot, i}
          <div
            class="day"
            title={`${slot.toLocaleString()}\n${
              valueCounts[value.value][i]
            } occurances`}
            style={valueCounts[value.value][i] == 0
              ? `background: #111`
              : `opacity: ${
                  (valueCounts[value.value][i] / valueMax[value.value]) * 100
                }%`}
          />
        {/each}
      </div>
    {/each}
  </div>
  <div class="time-range">
    <div class="time">{timeSlots[0].toLocaleString()}</div>
    <div class="time max-time">
      {timeSlots[timeSlots.length - 1].toLocaleString()}
    </div>
  </div>
{/if}

<style scoped>
  .freq-graph {
    overflow: auto;
    margin: 2em 0 1em;
  }
  .days {
    display: flex;
    height: 20px;
    margin: 2px 0;
  }
  .day {
    flex: 1;
    margin: 0 1px;
    border-radius: 2px;
    background: #0070f3;
  }
  .value-name {
    font-size: 0.85em;
    margin-top: 8px;
    text-wrap: nowrap;
  }

  .time-range {
    display: flex;
    font-size: 0.75em;
    color: #a1a1a1;
  }
  .time {
    flex: 1;
  }
  .max-time {
    text-align: right;
  }
</style>
