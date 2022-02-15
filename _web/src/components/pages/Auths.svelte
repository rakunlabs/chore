<script lang="ts">
  import { requestSender } from "@/helper/api";
  import { formToObjectMulti } from "@/helper/codec";

  import { storeHead } from "@/store/store";
  import { addToast } from "@/store/toast";
  import axios from "axios";
  import { onDestroy, onMount } from "svelte";
  import Icon from "@/components/ui/Icon.svelte";
  import Pagination from "@/components/ui/Pagination.svelte";
  import NoData from "@/components/ui/NoData.svelte";

  storeHead.set("Authentications");

  let listenElement: HTMLDivElement;

  let editMode = false;

  let formEdit: HTMLFormElement;
  let editID = "";

  const setEditMode = (v: boolean) => {
    editID = "";
    editMode = v;
  };

  const listAuth = async (offset: number, limit = 20) => {
    try {
      const l = await requestSender(
        "auths",
        { offset, limit },
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

  let datas = [];
  let meta = {} as { limit: number; count: number; offset: number };

  const modify = (i: Record<string, any>) => {
    datas = i.data;
    meta = i.meta;
  };

  const deleteAuth = async (id: string) => {
    try {
      await requestSender("auth", { id }, "DELETE", null, true);

      datas = datas.filter((d) => d.id != id);
    } catch (reason: unknown) {
      if (axios.isAxiosError(reason)) {
        const msg = reason.response.data.error ?? reason.message;
        addToast(msg, "alert");
      }
    }
  };

  const clickListen = (e: Event) => {
    const action = (e.target as HTMLElement).dataset["action"];
    if (action == "delete") {
      e.preventDefault();
      e.stopPropagation();

      const id = (e.target as HTMLElement).dataset["id"];

      if (confirm("Are you sure to delete?")) {
        deleteAuth(id);
      }
    }

    if (action == "edit") {
      e.preventDefault();
      e.stopPropagation();

      editMode = true;

      formEdit.reset();

      const dataset = (e.target as HTMLElement).dataset;

      editID = dataset["id"];
    }
  };

  const createAuth = async (
    e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
  ) => {
    const data = formToObjectMulti(e.currentTarget);
    const submitter = (e.submitter as HTMLButtonElement).value;

    // delete unused fields
    for (const key of ["id", "name", "groups", "headers"]) {
      if (data[key] == "") {
        delete data[key];
      }
    }

    // console.log(data);

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

    try {
      const response = await requestSender(
        "auth",
        null,
        submitter == "create" ? "POST" : "PATCH",
        data,
        true
      );

      // console.log(l);
      // filter and add
      datas = datas.filter((d) => d["id"] != response.data.data.id);

      try {
        const responseGet = await requestSender(
          "auth",
          {
            id: response.data.data.id,
          },
          "GET",
          null,
          true
        );

        datas.unshift(responseGet.data.data);
        datas = datas;
      } catch (reason: unknown) {
        let msg = reason;
        if (axios.isAxiosError(reason)) {
          msg = reason.response.data.error ?? reason.message;
        }
        addToast(msg as string, "warn");
      }

      error = "";
    } catch (reason: unknown) {
      if (axios.isAxiosError(reason)) {
        error = reason.response.data.error ?? reason.message;
      } else {
        error = reason as any;
      }
    }
  };

  let error = "";

  let headerCount = [];

  onMount(() => {
    listenElement.addEventListener("click", clickListen);
    listAuth(0);
  });

  onDestroy(() => {
    listenElement.removeEventListener("click", clickListen);
  });
</script>

<div class="bg-slate-50 p-5 mb-3">
  <div class="flex flex-row flex-wrap gap-4">
    <div class="flex-1">
      <div class="flex justify-between">
        <span class="font-bold block">{editMode ? "Edit" : "Create"} Auth</span>
        <div>
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={() => setEditMode(!editMode)}
            >{editMode ? "Create" : "Edit"} mode</button
          >
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={() => listAuth(0)}>Reload</button
          >
        </div>
      </div>

      <hr class="mb-4" />

      <form
        on:submit|preventDefault|stopPropagation={createAuth}
        bind:this={formEdit}
      >
        <label class="mb-1 flex">
          <span class="w-20 inline-block">ID</span>
          <input
            type="text"
            name="id"
            placeholder="----"
            disabled={!editMode}
            bind:value={editID}
            autocomplete="off"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Name</span>
          <input
            type="text"
            name="name"
            placeholder="jira1"
            autocomplete="off"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Groups</span>
          <input
            type="text"
            name="groups"
            placeholder="admin, deepcore"
            autocomplete="off"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <span class="mb-1 flex align-middle">
          <span class="w-20 inline-block">Headers</span>
          <button
            type="button"
            class="h-full"
            on:click|preventDefault|stopPropagation={() => {
              if (headerCount.length > 0) {
                headerCount.push(headerCount[headerCount.length - 1] + 1);
              } else {
                headerCount.push(1);
              }
              headerCount = headerCount;
            }}
          >
            <Icon icon="plus" class="p-1" />
          </button>
          <div class="flex flex-grow flex-col">
            {#each headerCount as h (h)}
              <div class="flex">
                <button
                  type="button"
                  on:click|preventDefault|stopPropagation={() =>
                    (headerCount = headerCount.filter((v) => v != h))}
                >
                  <Icon icon="minus" class="p-1" />
                </button>
                <input
                  type="text"
                  name={`headers-key-${h}`}
                  placeholder={`key-${h}`}
                  class="w-full px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
                />
                <input
                  type="text"
                  name={`headers-value-${h}`}
                  placeholder={`value-${h}`}
                  autocomplete="off"
                  class="w-full px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
                />
              </div>
            {/each}
          </div>
        </span>
        <button
          type="submit"
          value={editMode ? "edit" : "create"}
          class="w-full inline-flex items-center justify-center px-4 py-1 text-black bg-yellow-200 font-semibold capitalize hover:text-white hover:bg-red-500 active:bg-red-500 focus:outline-none focus:border-red-500 focus:ring focus:ring-red-200 disabled:opacity-25 transition"
          >{editMode ? "Edit" : "Create"}</button
        >
        <div
          class={`mt-2 bg-red-200 w-full h-6 ${error != "" ? "" : "invisible"}`}
        >
          <span class="break-all">{error}</span>
        </div>
      </form>
    </div>
  </div>
</div>

<div class="bg-slate-50 p-5" bind:this={listenElement}>
  <div class="overflow-x-auto rounded-none bg-white">
    <table class="w-full table-custom">
      <thead>
        <tr>
          <th style="width:5%" />
          <th style="width:10%">name</th>
          <th>headers</th>
          <th>groups</th>
          <th style="width:20%" />
        </tr>
      </thead>
      <tbody>
        {#each datas as d, i (d.id)}
          <tr class={editID == d.id ? "!bg-indigo-200" : ""}>
            <th>{i + 1}</th>
            <th>{d.name}</th>
            <th class="text-left">
              <input
                type="text"
                readonly
                class="w-full"
                value={d.headers ? JSON.stringify(d.headers) : ""}
              />
            </th>
            <th>{d.groups ?? ""}</th>
            <th>
              <button
                data-id={d.id}
                data-action="edit"
                class="bg-yellow-300 text-black hover:bg-green-500 hover:text-white px-2 rounded-sm"
              >
                edit
              </button>
              <button
                data-id={d.id}
                data-action="delete"
                class="bg-yellow-300 text-black hover:bg-red-500 hover:text-white px-2 rounded-sm"
              >
                delete
              </button>
            </th>
          </tr>
        {/each}
      </tbody>
    </table>
    <Pagination {meta} listF={listAuth} />
    <NoData hide={!!datas.length} />
  </div>
</div>
