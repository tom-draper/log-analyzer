<script lang="ts">
  import demo from "./lib/demo.json";
  import { onMount } from "svelte";
  import moment from "moment";
  import Card from "./lib/Card.svelte";
  import FailedLines from "./lib/FailedLines.svelte";

  function isDate(date: string): boolean {
    return moment(date, moment.ISO_8601, true).isValid();
  }

  function sortedTokenCounts(data: Data): {token: string, count: number}[] {
    let tokenCount: {[token: string]: number} = {}
    for (let i = 0; i < data.extraction.length; i++) {
      for (let token of Object.keys(data.extraction[i].params)) {
        if (!(token in tokenCount)) {
          tokenCount[token] = 0
        }
        tokenCount[token] += 1
      }
    }
    
    let tokens: {token: string, count: number}[] = []
    for (let [token, count] of Object.entries(tokenCount)) {
      tokens.push({token, count})
    }

    tokens.sort((a, b) => {
      return b.count - a.count
    })

    return tokens
  }

  function identifyTimestampToken(data: Data): string | null {
    let dateCount: { [token: string]: number } = {};
    for (let i = 0; i < data.extraction.length; i++) {
      for (let [token, value] of Object.entries(data.extraction[i].params)) {
        if (!(token in dateCount)) {
          dateCount[token] = 0;
        }
        if (typeof value === "string" && isDate(value)) {
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
    if (best.total < data.extraction.length * 0.5) {
      return null;
    }

    return best.token;
  }

  function performConversions(data: Data) {
    if (data.config.conversions === undefined) {
      return
    }

    for (let i = 0; i < data.extraction.length; i++) {
      const params = data.extraction[i].params;
      for (let [token, conversion] of Object.entries(data.config.conversions)) {
        if (!(token in params) || typeof params[token] !== "number") {
          continue
        }
        const value = params[token]
        if (typeof value !== "number") {
          continue
        }
        params[conversion.token] = conversion.multiplier * value
        delete params[token]
      }
    }
  }

  let data: Data;
  let tokenCounts: {token: string, count: number}[]
  let timestampToken: string | null;
  const production = import.meta.env.MODE === "production";
  onMount(async () => {
    if (production) {
      const response = await fetch("/data");
      data = await response.json();
    } else {
      //@ts-ignore
      data = demo;
    }
    performConversions(data);
    console.log(data);
    timestampToken = identifyTimestampToken(data);

    tokenCounts = sortedTokenCounts(data)
  });
</script>

<main>
  {#if tokenCounts != undefined}
    <div class="content">
      <div class="header">
        <div class="title">{data.extraction.length.toLocaleString()} lines</div>
      </div>
      {#each tokenCounts as token}
        {#if token.token !== timestampToken}
          <Card {data} token={token.token} lineCount={token.count} {timestampToken} />
        {/if}
      {/each}
      <FailedLines {data} />
    </div>
  {/if}
</main>

<style>
  .content {
    margin: 4em;
  }

  .title {
    margin: 0 0 20px;
  }
  .header {
    font-size: 2em;
    font-weight: 500;
  }
</style>
