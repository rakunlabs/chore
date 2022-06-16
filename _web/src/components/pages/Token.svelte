<script lang="ts">
  import { requestSender } from "@/helper/api";
  import { formToObject } from "@/helper/codec";
  import { dayCount } from "@/helper/date";
  import Pagination from "@/components/ui/Pagination.svelte";
  import { storeHead } from "@/store/store";
  import axios from "axios";
  import { onDestroy, onMount } from "svelte";
  import NoData from "@/components/ui/NoData.svelte";
  import Search from "@/components/ui/Search.svelte";

  storeHead.set("Personal Access Token");

  let listenElement: HTMLDivElement;

  let showID = "";

  let search = "";

  const listTokenSearch = async (
    search: string,
    offset: number,
    limit = 20
  ) => {
    try {
      const l = await requestSender(
        "tokens",
        { offset, limit, search },
        "GET",
        null,
        true
      );
      // console.log(l);
      const items = l ? l.data : {};
      modify(items);
    } catch (error) {
      console.error(error);
    }
  };

  const listToken = async (offset: number, limit = 20) => {
    listTokenSearch(search, offset, limit);
  };

  let datas = [];
  let meta = {} as { limit: number; count: number; offset: number };
  let showData = null as Record<string, any>;

  const modify = (i: Record<string, any>) => {
    datas = i.data;
    meta = i.meta;
  };

  const revoke = async (id: string) => {
    try {
      await requestSender("token", { id }, "DELETE", null, true);

      datas = datas.filter((d) => d.id != id);
    } catch (error) {
      console.error(error);
    }
  };

  const show = async (id: string) => {
    try {
      const d = await requestSender("token", { id }, "GET", null, true);
      showData = d.data.data;
      // console.log(d);
    } catch (error) {
      console.error(error);
    }
  };

  const clickListen = (e: Event) => {
    const action = (e.target as HTMLElement).dataset["action"];
    if (action == "revoke") {
      e.preventDefault();
      e.stopPropagation();

      const id = (e.target as HTMLElement).dataset["id"];

      if (confirm("Are you sure to delete?")) {
        revoke(id);
      }
    }

    if (action == "show") {
      e.preventDefault();
      e.stopPropagation();

      const id = (e.target as HTMLElement).dataset["id"];
      show(id);
      showID = id;
    }
  };

  const createToken = async (
    e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
  ) => {
    const data = formToObject(e.currentTarget);

    // delete unused fields
    for (const key of ["id", "name", "date", "groups"]) {
      if (data[key] == "") {
        delete data[key];
      }
    }

    // fix groups
    if (data["groups"]) {
      if (data["groups"].replaceAll(" ", "") == "") {
        data["groups"] = null;
      } else {
        data["groups"] = (data["groups"] as string)
          .replaceAll(" ", "")
          .split(",");
      }
    }

    if (data.date) {
      data.date = new Date(data.date).toISOString();
    }

    try {
      const response = await requestSender("token", null, "POST", data, true);
      datas.unshift(response.data.data);
      datas = datas;
      error = "";
    } catch (reason: unknown) {
      if (axios.isAxiosError(reason)) {
        error = reason.response.data.error ?? reason.message;
      } else {
        error = reason as any;
      }
    }
  };

  let infinity = false;
  let error = "";

  const searchFn = (s: string) => {
    listTokenSearch(s, 0);
  };

  onMount(() => {
    listenElement.addEventListener("click", clickListen);
    listToken(0);
  });

  onDestroy(() => {
    listenElement.removeEventListener("click", clickListen);
  });
</script>

<div class="bg-slate-50 p-5 mb-3">
  <div class="text-right">
    <button
      class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
      on:click={() => listToken(0)}>Reload</button
    >
  </div>

  <hr class="mb-4" />
  <div class="flex flex-row flex-wrap gap-4">
    <div class="flex-1">
      <span class="font-bold block">Create Token</span>

      <form on:submit|preventDefault|stopPropagation={createToken}>
        <label class="mb-1 flex align-middle">
          <span class="w-20 inline-block">Infinity</span>
          <input
            type="checkbox"
            name="infinity"
            bind:checked={infinity}
            class="px-2 self-center border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          /></label
        >
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Last date</span>
          <input
            type="datetime-local"
            name="date"
            disabled={infinity}
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Groups</span>
          <input
            type="text"
            name="groups"
            placeholder="admin,group1"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Name</span>
          <input
            type="text"
            name="name"
            placeholder="mytoken"
            autocomplete="off"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <button
          type="submit"
          class="w-full inline-flex items-center justify-center px-4 py-1 text-black bg-yellow-200 font-semibold capitalize hover:text-white hover:bg-red-500 active:bg-red-500 focus:outline-none focus:border-red-500 focus:ring focus:ring-red-200 disabled:opacity-25 transition"
          >Create</button
        >
        <div
          class={`mt-2 bg-red-200 w-full h-6 ${error != "" ? "" : "invisible"}`}
        >
          <span class="break-all">{error}</span>
        </div>
      </form>
    </div>
    <div class="flex-1">
      <div class="flex justify-between">
        <span class="font-bold block">Show Token</span>
        {#if showData}
          <button
            type="button"
            class="px-4 bg-gray-200 hover:bg-yellow-200"
            on:click|stopPropagation={() => {
              showData = null;
              showID = null;
            }}>Close</button
          >
        {/if}
      </div>
      {#if showData}
        {#each Object.entries(showData) as [key, d]}
          <label class="flex">
            <span class="w-20 inline-block">
              {key}
            </span>
            <input
              type="text"
              class="flex-grow  border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
              readonly
              value={d}
            />
          </label>
        {/each}
      {/if}
    </div>
  </div>
</div>

<div class="bg-slate-50 p-5" bind:this={listenElement}>
  <div class="flex items-center justify-end mb-1">
    <Search {searchFn} bind:search />
  </div>
  <div class="overflow-x-auto rounded-none bg-white">
    <table class="w-full table-custom">
      <thead>
        <tr>
          <th style="width:5%" />
          <th style="width:20%">name</th>
          <th>last date</th>
          <th>groups</th>
          <th style="width:20%" />
        </tr>
      </thead>
      <tbody>
        {#each datas as d, i (d.id)}
          <tr class={showID == d.id ? "!bg-indigo-200" : ""}>
            <th>{i + 1}</th>
            <th>{d.name}</th>
            <th class="text-left pl-1"
              >{d.date ?? ""}{#if d.date}
                - <b>{dayCount(new Date(d.date))} days</b>
              {/if}
            </th>
            <th>{d.groups ?? ""}</th>
            <th>
              <button
                data-id={d.id}
                data-action="show"
                class="bg-yellow-300 text-black hover:bg-green-500 hover:text-white px-2 rounded-sm"
              >
                show
              </button>
              <button
                data-id={d.id}
                data-action="revoke"
                class="bg-yellow-300 text-black hover:bg-red-500 hover:text-white px-2 rounded-sm"
              >
                revoke
              </button>
            </th>
          </tr>
        {/each}
      </tbody>
    </table>
    <Pagination {meta} listF={listToken} />
    <NoData hide={!!datas.length} />
  </div>
</div>
