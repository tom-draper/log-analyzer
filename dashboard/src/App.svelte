<script lang="ts">
  import demo from "./lib/demo.json";
  import { onMount } from "svelte";
  import moment from "moment";

  type LineParams = {
    [token: string]: string;
  };

  type TokenValueFreq = {
    [token: string]: { [value: string]: number };
  };

  function valueFreq(data: Data): TokenValueFreq {
    let freq: TokenValueFreq = {};
    for (let i = 0; i < data.extraction.params.length; i++) {
      for (let [token, value] of Object.entries(data.extraction.params[i])) {
        if (!(token in freq)) {
          freq[token] = {};
        }
        if (!(value in freq[token])) {
          freq[token][value] = 0;
        }
        freq[token][value] += 1;
      }
    }

    return freq;
  }

  function isDate(date: string): boolean {
    return moment(date, moment.ISO_8601, true).isValid();
  }

  function identifyTimestampToken(data: Data): string | null {
    let dateCount: { [token: string]: number } = {};
    for (let i = 0; i < data.extraction.params.length; i++) {
      for (let [token, value] of Object.entries(data.extraction.params[i])) {
        if (!(token in dateCount)) {
          dateCount[token] = 0;
        }
        if (isDate(value)) {
          dateCount[token] += 1;
        }
      }
    }

    let best: { token: string | null; total: number } = {
      token: null,
      total: 0,
    };

    for (let [token, count] of Object.entries(dateCount)) {
      if (count > best.total) {
        best = {
          token,
          total: count,
        };
      }
    }

    // Disqualify timestamp if only on less than 50% of log lines
    // Assume there is no timestamp token
    if (best.total < data.extraction.params.length * 0.5) {
      return null;
    }

    return best.token;
  }

  type Data = {
    extraction: {
      params: LineParams[];
      patterns: string[];
    };
    config: {
      tokens: string[];
      patterns: string[];
    };
  };

  let data: Data;
  const production = import.meta.env.MODE === "production";
  onMount(async () => {
    if (production) {
      const response = await fetch("/data");
      data = await response.json();
    } else {
      //@ts-ignore
      data = demo;
    }
    console.log(data);
    const timestampToken = identifyTimestampToken(data);
    console.log(timestampToken);
    const freq = valueFreq(data);
    console.log(freq);
  });
</script>

<main>
  {#if data != undefined}
    <div class="content">
      <div class="header">
        <div class="title">{data.extraction.params.length} lines</div>
      </div>
    </div>
  {/if}
</main>

<style>
  .header {
    margin: 2em 8%;
    font-size: 2em;
    font-weight: 500;
  }
</style>
