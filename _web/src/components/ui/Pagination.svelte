<script lang="ts">
  // offset, limit, total
  type Meta = {
    limit: number;
    count: number;
    offset: number;
  };

  export let meta: Meta;
  export let listF: (offset: number, limit: number) => void;

  let totalPage = 0;
  let currentPage = 0;

  let prevOffset = 0;
  let nextOffset = 0;
  let nextDisabled = false;
  let prevDisabled = false;

  const calculate = (meta: Meta) => {
    totalPage = Math.ceil(meta.count / meta.limit);
    currentPage = meta.offset / meta.limit + 1;

    nextOffset = meta.offset + meta.limit;
    prevOffset = meta.offset - meta.limit;

    nextDisabled = nextOffset >= meta.count;
    prevDisabled = prevOffset < 0;
  };

  $: calculate(meta);
</script>

<div class="flex justify-between bg-indigo-200 border-t-2 border-indigo-400">
  <div>
    {#if meta.offset != undefined}
      <span class="pl-2">
        Showing {meta.offset + 1} to {nextOffset} of {meta.count}
      </span>
    {/if}
  </div>
  <ul class="flex">
    <button
      disabled={prevDisabled}
      class="w-20 border-r hover:bg-indigo-400 hover:disabled:bg-gray-300"
      on:click|stopPropagation={() => listF(prevOffset, meta.limit)}
      >Prev</button
    >
    <button
      disabled={nextDisabled}
      class="w-20 hover:bg-indigo-400 hover:disabled:bg-gray-300"
      on:click|stopPropagation={() => listF(nextOffset, meta.limit)}
      >Next</button
    >
  </ul>
</div>
