<script lang="ts">
  // import { getItem, setItem } from "@/helper/api";
  // import { b64ToUtf8 } from "@/helper/codec";
  // import YAML from "yaml";

  // import { onMount } from "svelte";

  // import Table from "@/components/ui/Table.svelte";
  // import TableAdd from "../ui/TableAdd.svelte";

  // const head = ["id", "headers", "URL", "method"];

  // let data = [] as Array<object>;

  // const getData = async () => {
  //   try {
  //     const respData = await getItem("auths", "");
  //     data = (respData.data as Array<string>).map((v) => {
  //       return YAML.parse(b64ToUtf8(v));
  //     });
  //   } catch (error) {
  //     console.log(error);
  //   }
  // };

  // onMount(() => {
  //   getData();
  // });

  // let addToggle = false;
  // const addNewItem = () => {
  //   addToggle = !addToggle;
  // };

  // const saveData = async (id: string, v: string) => {
  //   try {
  //     await setItem("auths", id, v);
  //   } catch (error) {
  //     console.log("error");
  //   }
  // };
</script>

<div class="pb-1 flex justify-between items-center">
  <span>Authentications</span>
  <!-- <button
    on:click|stopPropagation={addNewItem}
    class="w-20 bg-transparent border-2 border-gray-500 text-gray-500 text-sm hover:bg-gray-500 hover:text-gray-100"
  >
    Add
  </button> -->
</div>
<!-- 
{#if addToggle}
  <TableAdd {head} call={saveData} />
{/if}

<Table {head} {data} /> -->
