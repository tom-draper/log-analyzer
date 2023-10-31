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
    const tokenCount: {[token: string]: number} = {}
    for (let i = 0; i < data.extraction.length; i++) {
      for (const token of Object.keys(data.extraction[i].params)) {
        if (!(token in tokenCount)) {
          tokenCount[token] = 0
        }
        tokenCount[token] += 1
      }
    }
    
    const tokens: {token: string, count: number}[] = []
    for (const [token, count] of Object.entries(tokenCount)) {
      tokens.push({token, count})
    }

    tokens.sort((a, b) => {
      // Force timestamp token to top of list
      if (a.token === timestampToken) {
        return Number.MIN_SAFE_INTEGER
      }
      return b.count - a.count
    })

    return tokens
  }

  function identifyTimestampToken(data: Data): string | null {
    const dateCount: { [token: string]: number } = {};
    for (let i = 0; i < data.extraction.length; i++) {
      for (const [token, value] of Object.entries(data.extraction[i].params)) {
        if (!(token in dateCount)) {
          dateCount[token] = 0;
        }
        if (typeof value === "string" && isDate(value)) {
          dateCount[token] += 1;
        }
      }
    }

    const best: { token: string | null; total: number } = {
      token: null,
      total: 0,
    };
    for (let [token, count] of Object.entries(dateCount)) {
      if (count > best.total) {
        best.token = token
        best.total = count
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

  function scrollToBottom() {
    window.scrollTo(0, document.body.scrollHeight);
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
        {#if Object.keys(data.failed).length >= 1}
          <button on:click={scrollToBottom}>{Object.keys(data.failed).length} {Object.keys(data.failed).length == 1 ? 'error' : 'errors'}</button>
        {/if}
      </div>
      {#each tokenCounts as token}
        <Card {data} token={token.token} lineCount={token.count} {timestampToken} />
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
    display: flex;
  }
  button {
    margin-top: 10px;
    margin-left: auto;
    background: #271515;
    color: #dd7178;
    border: 1px solid #dd71787d;
    padding: 5px 10px;
    border-radius: 3px;
    height: min-content;
    font-size: 0.9rem;
    outline: none
  }

@media screen and (max-width: 800px) {
  .content {
    margin: 3em 2em 2em;
  }
}
</style>
