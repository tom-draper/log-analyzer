<script lang="ts">
  import { onMount } from "svelte";
  import Moment from "moment";
  import { timestampRange, makeTimeSlots, nearestSlotIndex } from "./timeSlots";

  type SortableValueCount = {
    value: string;
    total: number;
  };

  function valueByTimeSlot(
    data: Data,
    token: string,
    timestampToken: string | null,
    timeSlots: Moment.Moment[]
  ) {
    if (timestampToken === null) {
      return {};
    }

    const slotMs = timeSlots.map((s) => new Date(s).getTime());

    const days: ValueCounts = {};
    for (const extraction of data.extraction) {
      if (!(token in extraction.params)) continue;

      const timestamp = new Date(
        extraction.params[timestampToken].value as string
      ).getTime();
      const slotIndex = nearestSlotIndex(timestamp, slotMs);

      const value = extraction.params[token].value;
      days[value] ||= Array(timeSlots.length).fill(0);
      days[value][slotIndex]++;
    }

    return days;
  }

  function sortedValueByTimeSlot(valueByTimeSlot: ValueCounts) {
    const sortedValues: SortableValueCount[] = [];
    for (let [value, counts] of Object.entries(valueByTimeSlot)) {
      const total = counts.reduce((partialSum, a) => partialSum + a, 0);
      sortedValues.push({ value, total });
    }

    sortedValues.sort((a, b) => {
      return b.total - a.total;
    });

    return sortedValues;
  }

  function valueCountMax(days: ValueCounts) {
    const valueMax: ValueCount = {};
    for (const value in days) {
      valueMax[value] = Math.max(...days[value]);
    }

    return valueMax;
  }

  let sortedValueCounts: SortableValueCount[];
  let valueCounts: ValueCounts;
  let valueMax: ValueCount;
  let timeSlots: Moment.Moment[];
  onMount(() => {
    const range = timestampRange(data, timestampToken);
    if (!range) return;
    const [min, max] = range;
    timeSlots = makeTimeSlots(min, max, 50);

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
            } occurrences`}
            style={valueCounts[value.value][i] == 0
              ? `background: #101010`
              : `opacity: ${
                  (valueCounts[value.value][i] / valueMax[value.value]) * 100
                }%`}
          ></div>
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
    background: var(--highlight);
  }
  .value-name {
    font-size: 0.85em;
    margin-top: 8px;
    text-wrap: nowrap;
    color: #888;
  }

  .time-range {
    display: flex;
    font-size: 0.75em;
    color: rgb(68, 68, 68);
  }
  .time {
    flex: 1;
  }
  .max-time {
    text-align: right;
  }
</style>
