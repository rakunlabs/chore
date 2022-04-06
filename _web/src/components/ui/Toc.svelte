<script lang="ts">
  import type { tocT } from "@/helper/toc";

  export let articles: tocT[] = [];
  let className = "";
  export { className as class };
</script>

<div class={`z-10 ${className}`}>
  <span class="font-bold border-b border-gray-400 pb-1">IN THIS ARTICLE</span>
  <ul class="pt-2">
    {#each articles as article}
      <li
        class={`border-l border-orange-200 pl-4 hover:border-red-500 hover:bg-red-100 toc-indent-${article.indent}`}
      >
        <a href={article.href} class="block">{article.name}</a>
      </li>
    {/each}
  </ul>
</div>

<style lang="scss">
  .toc-indent-1 {
    @apply border-l;
  }
  .toc-indent-2 {
    @apply border-l-4 pl-3;
  }
  .toc-indent-3 {
    @apply border-l-8 pl-2;
  }
</style>
