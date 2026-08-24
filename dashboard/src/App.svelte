<script lang="ts">
  import demo from "./assets/demo/demo.json";
  import { onMount } from "svelte";
  import Header from "./lib/Header.svelte";
  import Card from "./lib/Card.svelte";
  import Failed from "./lib/Failed.svelte";
  import TypeWarnings from "./lib/TypeWarnings.svelte";
  import Footer from "./lib/Footer.svelte";
  import RequestFunnel from "./lib/RequestFunnel.svelte";
  import EventExplorer from "./lib/EventExplorer.svelte";
  import PatternCoverage from "./lib/PatternCoverage.svelte";

  type TokenCount = { token: string; dependentToken?: string; count: number };

  function sortedTokenCounts(data: Data) {
    const tokenCount: { [token: string]: number } = {};
    for (let i = 0; i < data.extraction.length; i++) {
      for (const token in data.extraction[i].params) {
        tokenCount[token] ||= 0;
        tokenCount[token] += 1;
      }
    }

    const tokens: TokenCount[] = [];
    for (const [token, count] of Object.entries(tokenCount)) {
      tokens.push({ token, count });
    }

    tokens.sort((a, b) => {
      // Force timestamp token to top of list
      if (a.token === timestampToken) return -1;
      if (b.token === timestampToken) return 1;
      return b.count - a.count;
    });

    return tokens;
  }

  function sortedTokenDependencyCounts(data: Data) {
    if (data.config.dependencies === undefined) {
      return [];
    }

    const tokenCount: Map<readonly [string, string], number> = new Map();
    for (const token in data.config.dependencies) {
      for (const dependentToken of data.config.dependencies[token]) {
        const key = [token, dependentToken] as const;
        for (let i = 0; i < data.extraction.length; i++) {
          if (
            token in data.extraction[i].params &&
            dependentToken in data.extraction[i].params
          ) {
            tokenCount.has(key) || tokenCount.set(key, 0);
            tokenCount.set(key, (tokenCount.get(key) ?? 0) + 1);
          }
        }
      }
    }

    const tokens: TokenCount[] = [];
    for (const [[token, dependentToken], count] of tokenCount.entries()) {
      tokens.push({ token, dependentToken, count });
    }

    tokens.sort((a, b) => {
      // Force timestamp token to top of list
      if (a.token === timestampToken) return -1;
      if (b.token === timestampToken) return 1;
      return b.count - a.count;
    });

    return tokens;
  }

  function identifyTimestampToken(data: Data) {
    const dateCount: { [token: string]: number } = {};
    for (let i = 0; i < data.extraction.length; i++) {
      for (const [token, value] of Object.entries(data.extraction[i].params)) {
        dateCount[token] ||= 0;
        if (value.type === "time") {
          dateCount[token] += 1;
        }
      }
    }

    const best: { token: string | null; total: number } = {
      token: null,
      total: 0,
    };
    for (const [token, count] of Object.entries(dateCount)) {
      if (count > best.total) {
        best.token = token;
        best.total = count;
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
      return;
    }

    for (let i = 0; i < data.extraction.length; i++) {
      const params = data.extraction[i].params;
      for (let [token, conversion] of Object.entries(data.config.conversions)) {
        if (!(token in params)) {
          continue;
        }
        const param = params[token];
        if (typeof param.value !== "number") {
          continue;
        }
        const converted = conversion.multiplier * param.value;
        // Group the converted value under the target token, creating it if the
        // line doesn't already have it.
        if (conversion.token in params) {
          params[conversion.token].value = converted;
        } else {
          params[conversion.token] = { value: converted, type: param.type };
        }
        delete params[token];
      }
    }
  }

  function getFailedLines(data: Data) {
    const failedLines: FailedLines = {};
    for (let i = 0; i < data.extraction.length; i++) {
      if (data.extraction[i].pattern === "") {
        failedLines[data.extraction[i].lineNumber] = data.extraction[i].line;
      }
    }
    return failedLines;
  }

  function getDataTypeBreakdown(data: Data) {
    const dataTypes: DataTypes = {};

    for (let i = 0; i < data.extraction.length; i++) {
      for (const [token, params] of Object.entries(data.extraction[i].params)) {
        dataTypes[token] ||= {};
        dataTypes[token][params.type] ||= 0;
        dataTypes[token][params.type] += 1;
      }
    }

    return dataTypes;
  }

  function getMultiTypeTokens(dataTypes: DataTypes) {
    const multiTypes: DataTypes = {};
    for (const [token, types] of Object.entries(dataTypes)) {
      if (Object.keys(types).length > 1) {
        multiTypes[token] = types;
      }
    }
    return multiTypes;
  }

  let data: Data;
  let failedLines: FailedLines;
  let multiTypeTokens: DataTypes;
  let tokenCounts: TokenCount[];
  let timestampToken: string | null;
  let dataVersion = 0;
  let loadError: string | null = null;
  let hex: string = "#55ff5c";
  const production = import.meta.env.MODE === "production";

  function applyData(d: Data) {
    data = d;
    performConversions(data);

    timestampToken = identifyTimestampToken(data);

    failedLines = getFailedLines(data);

    const dataTypes = getDataTypeBreakdown(data);
    multiTypeTokens = getMultiTypeTokens(dataTypes);

    const counts = sortedTokenCounts(data);
    counts.push(...sortedTokenDependencyCounts(data));
    tokenCounts = counts;
    dataVersion += 1;
  }

  async function loadData() {
    try {
      const response = await fetch("/data");
      if (!response.ok) {
        throw new Error(`server returned ${response.status}`);
      }
      applyData(await response.json());
      loadError = null;
    } catch (error) {
      loadError = error instanceof Error ? error.message : "unknown error";
    }
  }

  onMount(() => {
    if (production) {
      void loadData();
      // Subscribe to live updates when the server is following a log source.
      const events = new EventSource("/events");
      events.addEventListener("update", () => {
        void loadData();
      });
      return () => events.close();
    } else {
      //@ts-ignore
      applyData(demo);
    }
  });
</script>

<main>
  {#if loadError}
    <p class="error">Unable to load dashboard data: {loadError}</p>
  {/if}
  {#if tokenCounts !== undefined}
    <div class="content">
      <Header
        {failedLines}
        {multiTypeTokens}
        lineCount={data.extraction.length}
      />
      <RequestFunnel {data} />
      <PatternCoverage {data} />
      <EventExplorer {data} />
      {#each tokenCounts as token (`${token.token}:${token.dependentToken ?? ""}:${dataVersion}`)}
        <Card
          {data}
          token={token.token}
          dependentToken={token.dependentToken ?? null}
          lineCount={token.count}
          {timestampToken}
          {hex}
        />
      {/each}
      <TypeWarnings {data} multiTypes={multiTypeTokens} />
      <Failed {failedLines} />
      <Footer />
    </div>
  {/if}
</main>

<style>
  .content {
    max-width: 1280px;
    margin: 2.5rem auto 1rem;
    padding: 0 1.5rem;
  }

  .error {
    color: #ff8a8a;
    margin: 2rem 4rem;
  }

  @media screen and (max-width: 800px) {
    .content {
      margin: 1.5rem auto 1rem;
      padding: 0 1rem;
    }
  }
</style>
