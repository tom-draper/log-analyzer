<script lang="ts">
  import { onMount } from "svelte";

  function typeString(types: { [type: string]: number }) {
    let total = 0;
    for (let count of Object.values(types)) {
      total += count;
    }

    let string = "";
    for (let [type, count] of Object.entries(types)) {
      string += type + ": " + count.toLocaleString() + ", ";
    }
    // Remove final comma and space
    string = string.slice(0, -2);
    return string;
  }

  function getDataTypeExamples(
    token: string,
    types: { [type: string]: number }
  ) {
    const examples: { [type: string]: string } = {};
    for (let type of Object.keys(types)) {
      if (type in examples) {
        continue;
      }
      for (let e of data.extraction) {
        if (!(token in e.params)) {
          continue;
        }
        if (e.params[token].type != type) {
          continue;
        }
        examples[e.params[token].type] = e.line;
      }
    }

    return examples;
  }

  function getExamples(dataTypes: DataTypes) {
    const examples: { [token: string]: { [type: string]: string } } = {};
    for (let [token, types] of Object.entries(dataTypes)) {
      examples[token] = getDataTypeExamples(token, types);
    }
    return examples;
  }

  let examples: { [token: string]: { [type: string]: string } };
  onMount(() => {
    examples = getExamples(multiTypes);
  });

  export let data: Data, multiTypes: DataTypes;
</script>

{#if Object.keys(multiTypes).length > 0}
  <div class="card">
    {#each Object.entries(multiTypes) as [token, dataTypes]}
      <div class="line-container">
        <div class="lineNumber">
          {token} has {Object.keys(dataTypes).length} data types
        </div>
        <div class="line">{typeString(dataTypes)}</div>
      </div>
    {/each}
  </div>
{/if}

<style scoped>
  .card {
    border: 1px solid #ffffff24;
    border-radius: 5px;
    margin: 3em 0;
    padding: 2rem;
  }
  .line-container {
    display: flex;
    background: #271515;
    color: #ddd871;
    border-radius: 5px;
    margin: 5px;
    font-size: 0.9em;
  }
  .lineNumber {
    margin: 6px 0;
    padding: 4px 20px;
    border-right: 1px solid #5e4c1e;
  }
  .line {
    margin: 10px 20px;
  }
</style>
