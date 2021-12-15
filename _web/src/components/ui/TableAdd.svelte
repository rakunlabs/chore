<script lang="ts">
  export let head = [] as Array<string>;
  export let call: (id: string, data: string) => void;

  let form: HTMLFormElement;
  const save = () => {
    const formData = new FormData(form);
    const data: Record<string, any> = {};
    for (const field of formData) {
      const [key, value] = field;
      data[key] = value;
    }
    call(data["id"], JSON.stringify(data));
  };
</script>

<div class="block p-6 bg-white">
  <form bind:this={form}>
    <div class="form-group mb-6">
      {#each head as h}
        <input type="text" name={h} placeholder={h} />
      {/each}
    </div>
    <button type="button" on:click|preventDefault|stopPropagation={save}>
      Save
    </button>
  </form>
</div>
