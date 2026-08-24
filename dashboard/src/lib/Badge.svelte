<script lang="ts">
  let { tone = "neutral", children }: { tone?: "neutral" | "warning" | "danger"; children: import("svelte").Snippet } = $props();
</script>

<span class={`badge ${tone}`}>{@render children()}</span>
