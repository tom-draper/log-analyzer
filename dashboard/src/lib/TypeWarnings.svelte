<script lang="ts">
  import { onMount } from "svelte";

  type DataTypes = {
    [token: string]: {
      [type: string]: number;
    };
  };

  function getDataTypeBreakdown(data: Data) {
    const dataTypes: DataTypes = {};

    for (let i = 0; i < data.extraction.length; i++) {
      for (let [token, params] of Object.entries(data.extraction[i].params)) {
        if (!(token in dataTypes)) {
          dataTypes[token] = {};
        }
        if (!(params.type in dataTypes[token])) {
          dataTypes[token][params.type] = 0;
        }
        dataTypes[token][params.type] += 1;
      }
    }

    return dataTypes;
  }

  function getMultiTypeTokens(dataTypes: DataTypes) {
    const multiTypes: DataTypes = {};
    for (let [token, types] of Object.entries(dataTypes)) {
      if (Object.keys(types).length > 1) {
        multiTypes[token] = types;
      }
    }
    return multiTypes;
  }

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

  let multiTypes: DataTypes;
  let examples: { [token: string]: { [type: string]: string } };
  onMount(() => {
    const dataTypes = getDataTypeBreakdown(data);
    const tokens = getMultiTypeTokens(dataTypes);
    if (Object.keys(tokens).length > 0) {
      examples = getExamples(tokens);
      multiTypes = tokens;
    }
  });

  export let data: Data;
</script>

{#if multiTypes !== undefined}
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
