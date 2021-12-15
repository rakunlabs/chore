<script lang="ts">
  export let url = "/";

  let paths = [] as Array<string>;

  const changeInner = (url: string) => {
    if (!url) {
      paths = [];
      return;
    }

    url = url[0] != "/" ? url : url.slice(1);
    paths = url.split("/");
  };

  $: changeInner(url);
</script>

<div>
  {#each paths as path, i}
    {" / "}
    {#if paths.length - 1 == i || (paths.length - 2 == i && paths[i + 1] == "")}
      {path}
    {:else}
      <a
        class="text-blue-700"
        href={"#" + "/" + paths.slice(0, i + 1).join("/") + "/"}>{path}</a
      >
    {/if}
  {/each}
</div>
