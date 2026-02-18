<script lang="ts">
  import { onMount } from "svelte";

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

<div class="card" class:hidden={Object.keys(multiTypes).length === 0}>
  {#each Object.entries(multiTypes) as [token, dataTypes]}
    <div class="line-container">
      <div class="token-header">
        <b>{token}</b> has {Object.keys(dataTypes).length} data types
        <span class="type-counts">
          {#each Object.entries(dataTypes) as [type, count], i}
            <span class="type-pill"
              >{type} <span class="pill-count">× {count.toLocaleString()}</span
              ></span
            >{#if i < Object.keys(dataTypes).length - 1}{" "}{/if}
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
</div>

<style scoped>
  .card {
    border: 1px solid #ffffff24;
    border-radius: 5px;
    margin: 3em 0;
    padding: 2rem;
  }
  .line-container {
    background: #5e4c1589;
    color: #ddd871;
    border-radius: 5px;
    margin: 5px;
    font-size: 0.9em;
    padding: 20px 20px;
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
  .type-pill {
    background: #ffffff18;
    border-radius: 4px;
    padding: 1px 8px;
    font-size: 0.9em;
  }
  .pill-count {
    color: #b8b44a;
  }
  .examples {
    padding: 5px 15px;
    margin-left: 20px;
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
