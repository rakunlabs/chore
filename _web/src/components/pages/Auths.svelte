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
  import Search from "@/components/ui/Search.svelte";

  storeHead.set("Authentications");

  let listenElement: HTMLDivElement;

  let editMode = false;

  let formEdit: HTMLFormElement;
  // let editID = "";
  let editData = {
    id: "",
    name: "",
    groups: "",
    headers: {},
    data: "",
  };

  let error = "";
  let headerCount = [];

  const setEditMode = (v: boolean) => {
    editData = {
      id: "",
      name: "",
      groups: "",
      headers: {},
      data: "",
    };

    headerCount = [];

    editMode = v;
  };

  let search = "";

  const listAuthSearch = async (search: string, offset: number, limit = 20) => {
    try {
      const l = await requestSender(
        "auths",
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

  const listAuth = async (offset: number, limit = 20) => {
    listAuthSearch(search, offset, limit);
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

      if (
        confirm(
          `Are you sure to delete ${datas.find((v) => v.id == id)?.name}?`
        )
      ) {
        deleteAuth(id);
      }

      return;
    }

    if (action == "edit") {
      e.preventDefault();
      e.stopPropagation();

      if (!editMode) {
        editMode = true;

        formEdit.reset();
      }

      const dataset = (e.target as HTMLElement).dataset;

      let editID = dataset["id"];

      let data = datas.find((d) => d.id == editID);

      editData = {
        id: data["id"],
        name: data["name"],
        groups: data["groups"],
        headers: data["headers"],
        data: data["data"],
      };

      return;
    }
  };

  const createAuth = async (
    e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
  ) => {
    const data = formToObjectMulti(e.currentTarget);
    const submitter = (e.submitter as HTMLButtonElement).value;

    // delete unused fields
    for (const key of ["id", "name", "groups", "headers", "data"]) {
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
        submitter == "create" ? "POST" : "PUT",
        data,
        true
      );

      // console.log(l);
      // filter and add
      datas = datas.filter((d) => d["id"] != data.id);

      try {
        const responseGet = await requestSender(
          "auth",
          {
            id: submitter == "create" ? response.data.data.id : data.id,
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

  const searchFn = (s: string) => {
    listAuthSearch(s, 0);
  };

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
            readonly={true}
            bind:value={editData.id}
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
            value={editData.name}
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
            value={editData.groups}
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <span class="mb-1 flex align-middle">
          <span class="w-20 inline-block">Headers</span>
          <button
            type="button"
            class="h-full"
            on:click|preventDefault|stopPropagation={() => {
              const data = formToObjectMulti(formEdit);

              editData.headers = {
                ...data.headers,
                "": "",
              };
            }}
          >
            <Icon icon="plus" class="p-1" />
          </button>
          <div class="flex flex-grow flex-col">
            {#each Object.keys(editData?.headers ?? {}) as key}
              <div class="flex">
                <button
                  type="button"
                  on:click|preventDefault|stopPropagation={() => {
                    const data = formToObjectMulti(formEdit);

                    delete data.headers[key];
                    editData.headers = {
                      ...data.headers,
                    };
                  }}
                >
                  <Icon icon="minus" class="p-1" />
                </button>
                <input
                  type="text"
                  name={`headers-key-${key}`}
                  placeholder={`key-${key}`}
                  value={key}
                  class="w-full px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
                />
                <input
                  type="text"
                  name={`headers-value-${key}`}
                  placeholder={`value-${key}`}
                  autocomplete="off"
                  value={editData.headers[key]}
                  class="w-full px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
                />
              </div>
            {/each}
          </div>
        </span>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Data</span>
          <textarea
            name="data"
            placeholder="any data"
            autocomplete="off"
            rows="5"
            value={editData.data}
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
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
  <div class="flex items-center justify-end mb-1">
    <Search {searchFn} bind:search />
  </div>
  <div class="overflow-x-auto rounded-none bg-white">
    <table class="w-full table-custom">
      <thead>
        <tr>
          <th style="width:5%" />
          <th style="width:10%">name</th>
          <th>headers</th>
          <th>data</th>
          <th>groups</th>
          <th style="width:20%" />
        </tr>
      </thead>
      <tbody>
        {#each datas as d, i (d.id)}
          <tr class={editData.id == d.id ? "!bg-indigo-200" : ""}>
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
            <th class="text-left">
              <input type="text" readonly class="w-full" value={d.data ?? ""} />
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
