<script lang="ts">
  import { onMount } from "svelte";

  function typeString(types: { [type: string]: number }) {
    let string = "";
    for (const [type, count] of Object.entries(types)) {
      string += type + ": " + count.toLocaleString() + ", ";
    }
    // Remove final comma and space
    string = string.slice(0, -2);
    return string;
  }

  type Examples = { [token: string]: Example };
  type Example = { [type: string]: { line: string; pattern: string } };

  function getDataTypeExamples(
    token: string,
    types: { [type: string]: number }
  ) {
    const examples: Example = {};
    for (const type in types) {
      if (type in examples) {
        continue
      }
      for (const e of data.extraction) {
        if (!(token in e.params) || (e.params[token].type != type)) {
          continue;
        }
        examples[e.params[token].type] = {
          line: e.line,
          pattern: e.pattern,
        };
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
      <div class="lineNumber">
        <b>{token}</b> has {Object.keys(dataTypes).length} data types - {typeString(
          dataTypes
        )}
      </div>
      {#if examples !== undefined}
        {#each Object.entries(examples[token]) as [type, example]}
          <div class="examples">
            <div class="token"><b>{type}</b></div>
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
  .examples,
  .lineNumber {
    padding: 5px 15px;
  }
  .examples {
    margin-left: 20px;
  }
  .lineNumber {
    margin-bottom: 10px;
  }
  .example-line {
    overflow-wrap: break-word;
  }
  .example-pattern {
    margin-top: 5px;
  }
</style>
