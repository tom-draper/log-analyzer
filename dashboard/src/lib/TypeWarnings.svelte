<script lang="ts">
  import { onMount } from "svelte";
  import Badge from "./Badge.svelte";
  import Panel from "./Panel.svelte";

  type Examples = { [token: string]: Example };
  type Example = { [type: string]: { line: string; pattern: string } };

  function getDataTypeExamples(
    token: string,
    types: { [type: string]: number }
  ) {
    const examples: Example = {};
    for (const type in types) {
      if (type in examples) {
        continue;
      }
      for (const e of data.extraction) {
        if (!(token in e.params) || e.params[token].type != type) {
          continue;
        }
        examples[e.params[token].type] = {
          line: e.line,
          pattern: e.pattern,
        };
        break;
      }
    }

    return examples;
  }

  function getExamples(dataTypes: DataTypes) {
    const examples: Examples = {};
    for (const [token, types] of Object.entries(dataTypes)) {
      examples[token] = getDataTypeExamples(token, types);
    }
    return examples;
  }

  let examples: Examples;
  onMount(() => {
    examples = getExamples(multiTypes);
  });

  export let data: Data, multiTypes: DataTypes;
</script>

{#if Object.keys(multiTypes).length > 0}
<Panel class="card">
  {#each Object.entries(multiTypes) as [token, dataTypes]}
    <div class="line-container">
      <div class="token-header">
        <b>{token}</b> has {Object.keys(dataTypes).length} data types
        <span class="type-counts">
          {#each Object.entries(dataTypes) as [type, count], i}
            <Badge tone="warning">{type} × {count.toLocaleString()}</Badge>{#if i < Object.keys(dataTypes).length - 1}{" "}{/if}
          {/each}
        </span>
      </div>
      {#if examples !== undefined}
        {#each Object.entries(examples[token]) as [type, example]}
          <div class="examples">
            <div class="type-label">
              <b>{type}</b>
              <span class="type-count">× {dataTypes[type].toLocaleString()}</span>
            </div>
            <div class="example-line">{example.line}</div>
            <div class="example-pattern">Pattern: {example.pattern}</div>
          </div>
        {/each}
      {/if}
    </div>
  {/each}
</Panel>
{/if}

<style scoped>
  .line-container {
    background: #5e4c1589;
    color: #ddd871;
    border-radius: calc(var(--radius) - 2px);
    margin: .5rem 0;
    font-size: 0.9em;
    padding: 1rem;
  }
  .token-header {
    padding: 5px 15px;
    margin-bottom: 10px;
  }
  .type-counts {
    display: inline-flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-left: 8px;
    vertical-align: middle;
  }
  .examples {
    padding: 5px 15px;
    margin-left: 1rem;
  }
  .type-label {
    margin-bottom: 4px;
  }
  .type-count {
    color: #b8b44a;
    font-size: 0.9em;
    margin-left: 4px;
  }
  .example-line {
    overflow-wrap: break-word;
  }
  .example-pattern {
    margin-top: 5px;
  }
</style>
