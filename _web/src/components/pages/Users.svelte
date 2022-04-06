<script lang="ts">
  import { requestSender } from "@/helper/api";

  import { storeHead } from "@/store/store";
  import Pagination from "@/components/ui/Pagination.svelte";
  import { addToast } from "@/store/toast";
  import axios from "axios";
  import { onDestroy, onMount } from "svelte";
  import { formToObject } from "@/helper/codec";
  import NoData from "@/components/ui/NoData.svelte";

  storeHead.set("Users");

  let listenElement: HTMLDivElement;

  let datas = [];
  let meta = {} as { limit: number; count: number; offset: number };
  let editMode = false;
  let error = "";

  let editID = "";

  let viewPass = false;
  let viewPassR = false;

  let formEdit: HTMLFormElement;

  const setEditMode = (v: boolean) => {
    editID = "";
    editMode = v;
  };

  const listUsers = async (offset: number, limit = 20) => {
    try {
      const l = await requestSender(
        "users",
        { offset, limit },
        "GET",
        null,
        true,
        {
          noAlert: true,
        }
      );
      // console.log(l);
      const items = l ? l.data : {};
      modify(items);
    } catch (reason: unknown) {
      let msg = reason;
      if (axios.isAxiosError(reason)) {
        msg = reason.response.data.error ?? reason.message;
      }
      addToast(msg as string, "warn");
    }
  };

  const modify = (i: Record<string, any>) => {
    datas = i.data;
    meta = i.meta;
  };

  const deleteUser = async (id: string) => {
    try {
      await requestSender("user", { id }, "DELETE", null, true);

      datas = datas.filter((d) => d.id != id);
    } catch (reason: unknown) {
      if (axios.isAxiosError(reason)) {
        const msg = reason.response.data.error ?? reason.message;
        addToast(msg, "alert");
      }
    }
  };

  const createUser = async (
    e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
  ) => {
    const submitterValue = (e.submitter as HTMLButtonElement).value;
    const data = formToObject(e.currentTarget);

    // check password are same
    if (data["password"] != data["passwordr"]) {
      error = "miss matched password";
      return;
    }

    delete data["passwordr"];

    // delete unused fields
    for (const key of ["password", "groups", "email", "name"]) {
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

    try {
      const response = await requestSender(
        "user",
        null,
        submitterValue == "create" ? "POST" : "PATCH",
        data,
        true
      );
      // console.log(l);

      // filter and add
      datas = datas.filter((d) => d["id"] != response.data.data.id);

      try {
        const responseGet = await requestSender(
          "user",
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

  const clickListen = (e: Event) => {
    const action = (e.target as HTMLElement).dataset["action"];
    if (action == "delete") {
      e.preventDefault();
      e.stopPropagation();

      const id = (e.target as HTMLElement).dataset["id"];

      if (confirm("Are you sure to delete?")) {
        deleteUser(id);
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

  onMount(() => {
    listenElement.addEventListener("click", clickListen);
    listUsers(0);
  });

  onDestroy(() => {
    listenElement.removeEventListener("click", clickListen);
  });
</script>

<div class="bg-slate-50 p-5 mb-3">
  <div class="flex flex-row flex-wrap gap-4">
    <div class="flex-1">
      <div class="flex justify-between">
        <span class="font-bold block">{editMode ? "Edit" : "Create"} User</span>
        <div>
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={() => setEditMode(!editMode)}
            >{editMode ? "Create" : "Edit"} Mode</button
          >
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={() => listUsers(0)}>Reload</button
          >
        </div>
      </div>

      <hr class="mb-4" />

      <form
        on:submit|preventDefault|stopPropagation={createUser}
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
            placeholder="userX"
            autocomplete="off"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Email</span>
          <input
            type="text"
            name="email"
            placeholder="user@worldline.com"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Groups</span>
          <input
            type="text"
            name="groups"
            placeholder="admin, deepcore"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">Password</span>
          <input
            type={viewPass ? "text" : "password"}
            name="password"
            placeholder="supersecretpass"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
          <button
            type="button"
            class={`px-5 ml-1 border border-gray-300 hover:bg-green-500 hover:text-white ${
              viewPass ? "bg-indigo-200" : "bg-yellow-200"
            }`}
            on:click|stopPropagation={() => (viewPass = !viewPass)}>view</button
          >
        </label>
        <label class="mb-1 flex">
          <span class="w-20 inline-block">PasswordR</span>
          <input
            type={viewPassR ? "text" : "password"}
            name="passwordr"
            placeholder="supersecretpass"
            class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
          />
          <button
            type="button"
            class={`px-5 ml-1 border border-gray-300 hover:bg-green-500 hover:text-white ${
              viewPassR ? "bg-indigo-200" : "bg-yellow-200"
            }`}
            on:click|stopPropagation={() => (viewPassR = !viewPassR)}
            >view</button
          >
        </label>
        <button
          type="submit"
          name="action"
          value={editMode ? "edit" : "create"}
          class="w-full inline-flex items-center justify-center px-4 py-1 text-black bg-yellow-200 font-semibold capitalize hover:text-white hover:bg-red-500 active:bg-red-500 focus:outline-none focus:border-red-500 focus:ring focus:ring-red-200 disabled:opacity-25 transition"
        >
          {editMode ? "Edit" : "Create"}
        </button>
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
          <th style="width:15%">name</th>
          <th>email</th>
          <th>groups</th>
          <th style="width:20%" />
        </tr>
      </thead>
      <tbody>
        {#each datas as d, i (d.id)}
          <tr class={editID == d.id ? "!bg-indigo-200" : ""}>
            <th>{i + 1}</th>
            <th>{d.name}</th>
            <th>{d.email ?? ""}</th>
            <th>{d.groups ?? ""}</th>
            <th>
              <button
                data-id={d.id}
                data-action="edit"
                class="bg-yellow-200 text-black hover:bg-green-500 hover:text-white px-2 rounded-sm"
              >
                edit
              </button>
              <button
                data-id={d.id}
                data-action="delete"
                class="bg-yellow-200 text-black hover:bg-red-500 hover:text-white px-2 rounded-sm"
              >
                delete
              </button>
            </th>
          </tr>
        {/each}
      </tbody>
    </table>
    <Pagination {meta} listF={listUsers} />
    <NoData hide={!!datas.length} />
  </div>
</div>
