<script lang="ts">
  import { storeView, storeViewReload } from "@/store/store";

  let endpoint = storeView.view;

  const submit = (e: Event) => {
    const form = e.target as HTMLFormElement;
    const formData = new FormData(form);

    const endpoint = formData.get("endpoint") as string;

    storeView.view = endpoint;
    storeViewReload.update((n) => !n);
  };

  const reset = () => {
    storeView.view = "";
    storeViewReload.update((n) => !n);
  };
</script>

<form on:submit|preventDefault={submit} on:reset|preventDefault={reset}>
  <p>View Endpoint</p>
  <input
    type="text"
    placeholder="endpoint"
    name="endpoint"
    bind:value={endpoint}
  />

  <div class="border-t border-black p-1 my-1 flex justify-between gap-1">
    <button
      type="submit"
      class="flex-1 px-2 py-1 hover:bg-red-400 hover:text-white border border-gray-400 bg-yellow-50 w-full"
      >View</button
    >
    <button
      type="reset"
      class="flex-1 px-2 py-1 hover:bg-red-400 hover:text-white border border-gray-400 bg-yellow-50 w-full"
      >Reset</button
    >
  </div>
</form>
