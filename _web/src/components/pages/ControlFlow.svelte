<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { requestSender } from "@/helper/api";
  import { storeViewReload, storeHead, storeView } from "@/store/store";
  import { addToast } from "@/store/toast";
  import { formToObject } from "@/helper/codec";
  import { nodes } from "@/models/node";
  import type { node } from "@/models/node";
  import Pagination from "@/components/ui/Pagination.svelte";
  import Icon from "@/components/ui/Icon.svelte";
  import NoData from "@/components/ui/NoData.svelte";

  import { Base64 } from "js-base64";
  import axios from "axios";
  import Drawflow from "drawflow";
  import type { DrawflowNode } from "drawflow";
  import { getEndpoints } from "@/helper/nodes";
  import Search from "@/components/ui/Search.svelte";
  import CodeEditor from "@/components/ui/CodeEditor.svelte";
  import Nodes from "../ui/Nodes.svelte";
  import type { noteData } from "@/models/nodes/note";
  import { viewEndpoint, viewEndpointClear } from "@/helper/view";

  storeHead.set("ControlFlow");

  let drawDiv: HTMLDivElement;
  let listenElement: HTMLDivElement;
  let listenNodes: HTMLDivElement;
  let formEdit: HTMLFormElement;
  let editor: Drawflow;

  let datas = [];
  let meta = {} as { limit: number; count: number; offset: number };
  let selected = "table";
  let error = "";
  let editID = "";

  let nodeSelected = "endpoint";
  let lastSelectedID = 0;
  let lastSelectedNode: DrawflowNode | null = null;

  let inputCount = 1;
  let showEditor = false;
  let codeEditorSave: (script: string, inputs: string) => void;
  let setCodeEditorValue: (
    script: string,
    inputs: string,
    info: string
  ) => void;
  let showEditorChange = (v: boolean) => {
    showEditor = v;
  };

  let currentGroups = "";
  let currentName = "";

  let fullScreen = false;

  let search = "";

  let drawMousePos = {
    x: 0,
    y: 0,
  };

  const handleMouseMove = (e: MouseEvent) => {
    drawMousePos.x = e.clientX;
    drawMousePos.y = e.clientY;
  };

  const sanitizeMove = (
    clientX: null | number = null,
    clientY: null | number = null
  ) => {
    const posX = clientX ?? drawMousePos.x;
    const posY = clientY ?? drawMousePos.y;
    const x =
      posX * (drawDiv.clientWidth / (drawDiv.clientWidth * editor.zoom)) -
      (drawDiv.firstChild as HTMLDivElement).getBoundingClientRect().x *
        (drawDiv.clientWidth / (drawDiv.clientWidth * editor.zoom));
    const y =
      posY * (drawDiv.clientHeight / (drawDiv.clientHeight * editor.zoom)) -
      (drawDiv.firstChild as HTMLDivElement).getBoundingClientRect().y *
        (drawDiv.clientHeight / (drawDiv.clientHeight * editor.zoom));

    return [x, y];
  };

  const setSelected = (v: string) => {
    storeView.view = "";

    formEdit.reset();
    currentGroups = "";
    currentName = "";

    editID = "";
    selected = v;
    error = "";

    editor.editor_mode = "edit";
    editor.zoom_reset();

    editor.canvas_x = 0;
    editor.canvas_y = 0;

    // style.transform unset
    (drawDiv.firstChild as HTMLDivElement).style.transform = "";

    editor.import({
      drawflow: {
        Home: {
          data: {},
        },
      },
    });
  };

  const listControlsSearch = async (
    search: string,
    offset: number,
    limit = 20
  ) => {
    try {
      const l = await requestSender(
        "controls",
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

  const listControls = async (offset: number, limit = 20) => {
    listControlsSearch(search, offset, limit);
  };

  const modify = (i: Record<string, any>) => {
    datas = i.data;
    meta = i.meta;
  };

  const deleteControl = async (id: string) => {
    try {
      await requestSender("control", { id }, "DELETE", null, true);

      datas = datas.filter((d) => d.id != id);
    } catch (reason: unknown) {
      if (axios.isAxiosError(reason)) {
        const msg = reason.response.data.error ?? reason.message;
        addToast(msg, "alert");
      }
    }
  };

  const createControl = async (submitterValue: string) => {
    const data = formToObject(formEdit);

    // delete unused fields
    for (const key of ["id", "name", "groups"]) {
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

    const exportedData = editor.export().drawflow.Home.data;
    const content = JSON.stringify(exportedData);
    data["content"] = Base64.encode(content);

    data["endpoints"] = getEndpoints(exportedData);

    try {
      const response = await requestSender(
        "control",
        null,
        submitterValue == "create" ? "POST" : "PATCH",
        data,
        true
      );
      // console.log(l);
      addToast("saved controlflow", "info");

      // filter and add
      datas = datas.filter((d) => d["id"] != response.data.data.id);

      if (submitterValue == "create") {
        formEdit.reset();
        editID = response.data.data.id;
        selected = "edit";
      }

      try {
        const responseGet = await requestSender(
          "control",
          {
            id: response.data.data.id,
            nodata: true,
          },
          "GET",
          null,
          true
        );

        currentGroups = responseGet.data.data.groups;
        currentName = responseGet.data.data.name;

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
      }
    }
  };

  const clickListenDraw = (e: Event) => {
    const action = (e.target as HTMLElement).dataset["action"];

    switch (action) {
      case "editor":
        {
          const modifyElement = (e.target as HTMLElement)
            .nextElementSibling as HTMLTextAreaElement;

          const nodeName = (e.target as HTMLElement).parentElement.parentElement
            .parentElement.parentElement.id;

          const id = nodeName.slice(nodeName.indexOf("-") + 1);
          const data = editor.getNodeFromId(id).data;

          const info = `Control=${currentName} NodeID=${id}`;
          setCodeEditorValue(data.script, data.inputs, info);
          codeEditorSave = (script, inputs) => {
            modifyElement.value = script;
            editor.updateNodeDataFromId(
              nodeName.slice(nodeName.indexOf("-") + 1),
              {
                script: script,
                inputs: inputs,
              }
            );
          };

          showEditor = true;
        }
        break;
      case "checkbox":
        {
          const modifyElement = e.target as HTMLInputElement;

          let nodeName = "";
          if (modifyElement.dataset["parent"]) {
            const count = parseInt(modifyElement.dataset["parent"]);
            let e = modifyElement as HTMLElement;
            for (let i = 0; i < count; i++) {
              e = e.parentElement;
            }
            nodeName = e.id;
          } else {
            nodeName =
              modifyElement.parentElement.parentElement.parentElement
                .parentElement.parentElement.id;
          }

          const nodeId = nodeName.slice(nodeName.indexOf("-") + 1);
          editor.updateNodeDataFromId(nodeId, {
            ...editor.getNodeFromId(nodeId).data,
            [modifyElement.name]: modifyElement.checked,
          });
        }
        break;

      default:
        break;
    }
  };

  const clickListen = (e: Event) => {
    const action = (e.target as HTMLElement).dataset["action"];
    if (action == "delete") {
      e.preventDefault();
      e.stopPropagation();

      const id = (e.target as HTMLElement).dataset["id"];

      if (confirm("Are you sure to delete?")) {
        deleteControl(id);
      }
    }

    if (action == "edit") {
      e.preventDefault();
      e.stopPropagation();

      formEdit.reset();
      selected = "edit";

      const dataset = (e.target as HTMLElement).dataset;

      // console.log(dataset);

      editID = dataset["id"];
      currentGroups = dataset["groups"];
      currentName = dataset["name"];

      (async () => {
        try {
          const responseGet = await requestSender(
            "control",
            {
              id: editID,
            },
            "GET",
            null,
            true
          );

          const rawContent = Base64.decode(responseGet.data.data.content);
          const content = JSON.parse(rawContent);

          // console.log(content);

          editor.import({
            drawflow: {
              Home: {
                data: content,
              },
            },
          });
        } catch (reason: unknown) {
          if (axios.isAxiosError(reason)) {
            error = reason.response.data.error ?? reason.message;
          } else {
            error = reason as any;
          }

          return;
        }

        let lastId = 1;

        // add id number to html
        drawDiv
          ?.querySelectorAll(".drawflow-node")
          .forEach((node: HTMLDivElement) => {
            // get id number
            let id = node.id.slice(node.id.indexOf("-") + 1);

            if (+id > lastId) {
              lastId = +id;
            }

            // add id number to title-box
            let selectedTitle = node.querySelector(".title-box");
            if (selectedTitle) {
              selectedTitle.innerHTML += ` (${id})`;
            }
          });

        // set lastId
        (editor as any).nodeId = lastId + 1;
      })();
    }
  };

  const nodeUpdatedListener = (e: CustomEvent) => {
    e.preventDefault();
    e.stopPropagation();

    nodeUpdated(e.detail);
  };

  const nodeUpdated = (id: number) => {
    const node = editor.getNodeFromId(id);
    if (node == null) {
      return;
    }

    switch (node.name) {
      case "note":
        updateNote(id, node.data);
        break;
      default:
        break;
    }
  };

  const dragNode = (e: DragEvent) => {
    e.dataTransfer.setData("node", nodeSelected);
  };

  const dropNode = (e: DragEvent) => {
    const nodeS = e.dataTransfer.getData("node");
    addNodeToDrawFlow(nodeS, e.clientX, e.clientY);
  };

  const addNodeToDrawFlow = (name: string, posX: number, posY: number) => {
    if (editor.editor_mode == "fixed") {
      return false;
    }

    [posX, posY] = sanitizeMove(posX, posY);

    const node = JSON.parse(JSON.stringify(nodes[name])) as node;
    if (node.optionalInput) {
      node.input = inputCount;
    }

    if (node) {
      editor.addNode(
        name,
        node.input,
        node.output,
        posX,
        posY,
        node.class ?? "",
        node.data,
        node.html,
        false
      );
    }
  };

  const listenKeys = (event: KeyboardEvent) => {
    if (event.altKey) {
      switch (event.key) {
        case "c":
          navigator.clipboard.writeText(
            JSON.stringify(editor.getNodeFromId(lastSelectedID))
          );
          break;
        case "v":
          (async () => {
            let readNode = await navigator.clipboard.readText();
            let readNodeObj: any;
            try {
              readNodeObj = JSON.parse(readNode);
            } catch {}

            if (
              typeof readNodeObj === "object" &&
              typeof readNodeObj.inputs === "object" &&
              typeof readNodeObj.outputs === "object"
            ) {
              const [x, y] = sanitizeMove();
              editor.addNode(
                readNodeObj.name,
                Object.keys(readNodeObj.inputs).length,
                Object.keys(readNodeObj.outputs).length,
                x,
                y,
                readNodeObj.class,
                readNodeObj.data,
                readNodeObj.html,
                readNodeObj.typenode
              );
            }
          })();

          break;
      }
    }

    if (event.ctrlKey || event.metaKey) {
      switch (event.key) {
        // log editor output to console
        case "l":
          event.preventDefault();

          if (lastSelectedID != 0) {
            console.log(editor.getNodeFromId(lastSelectedID));
            break;
          }

          console.log(editor.export().drawflow.Home.data);
          break;
      }
    }
  };

  const searchFn = (s: string) => {
    listControlsSearch(s, 0);
  };

  // const updateCheckbox = (id: number) => {
  //   let nodeSelectedQuery = drawDiv.querySelector(`#node-${id}`);
  //   updateCheckboxAll(nodeSelectedQuery);
  // };

  const updateCheckboxAll = (e: Element | null = null) => {
    if (e == null) {
      e = drawDiv;
    }
    e.querySelectorAll('input[type="checkbox"]').forEach(
      (input: HTMLInputElement) => {
        if (input.getAttribute("value") == "true") {
          input.checked = true;
        }
      }
    );
  };

  const updateNote = (id: number, data: noteData) => {
    let nodeSelectedQuery = drawDiv.querySelector(`#node-${id}`);
    nodeSelectedQuery
      .querySelectorAll("textarea")
      .forEach((textArea: HTMLTextAreaElement) => {
        textArea.style.backgroundColor = data.backgroundColor;
        textArea.style.color = data.textColor;
        textArea.style.width = `${data.width}px`;
        textArea.style.height = `${data.height}px`;
      });
  };

  const updateNoteAll = (e: Element | null = null) => {
    if (e == null) {
      e = drawDiv;
    }

    // query selector all for class 'drawflow-node' and 'node-note'
    e.querySelectorAll(".node-note").forEach((node: HTMLDivElement) => {
      let id = node.id.slice(node.id.indexOf("-") + 1);
      let nodeData = editor.getNodeFromId(+id).data;
      updateNote(+id, nodeData);
    });
  };

  const viewEndpointCaller = (_: boolean) => {
    if (drawDiv == null) {
      return;
    }

    if (storeView.view == "") {
      viewEndpointClear(drawDiv);
    } else {
      viewEndpointClear(drawDiv);
      viewEndpoint(storeView.view, drawDiv, editor);
    }
  };

  $: viewEndpointCaller($storeViewReload);

  const nodeUnselected = () => {
    lastSelectedID = 0;
    lastSelectedNode = null;
  };

  onMount(() => {
    editor = new Drawflow(drawDiv, null);
    editor.reroute = true;
    editor.reroute_fix_curvature = true;
    editor.force_first_input = false;

    editor.zoom_min = 0.4;

    // editor.curvature = 0;
    // editor.reroute_curvature = 0;
    // editor.reroute_curvature_start_end = 0;

    editor.start();

    // listen buttons
    listenElement.addEventListener("click", clickListen);
    listControls(0);

    drawDiv.addEventListener("click", clickListenDraw);

    editor.on("nodeCreated", (id: number) => {
      let nodeSelectedQuery = drawDiv.querySelector(`#node-${id}`);
      if (nodeSelectedQuery) {
        let selectedTitle = nodeSelectedQuery.querySelector(".title-box");
        if (selectedTitle) {
          selectedTitle.innerHTML += ` (${id})`;
        }
      }

      nodeUpdated(id);
      // updateCheckbox(id);
    });

    editor.on("nodeSelected", (id: number) => {
      lastSelectedID = id;
      const node = editor.getNodeFromId(id);
      lastSelectedNode = node;
    });
    editor.on("import", (e) => {
      updateNoteAll();
      updateCheckboxAll();
    });

    editor.on("nodeUnselected", nodeUnselected);

    listenNodes.addEventListener("nodeUpdated", nodeUpdatedListener);
  });

  onDestroy(() => {
    listenElement.removeEventListener("click", clickListen);
    drawDiv.removeEventListener("click", clickListenDraw);
    listenNodes.addEventListener("nodeUpdated", nodeUpdatedListener);
  });
</script>

<CodeEditor
  bind:setCodeEditorValue
  {codeEditorSave}
  {showEditor}
  {showEditorChange}
/>

<div class="grid h-full grid-rows-[auto_1fr]">
  <div class="bg-slate-50 p-2 mb-3">
    <div class="flex flex-row flex-wrap justify-between gap-4 items-start">
      <div class="flex-1">
        <div class={selected == "table" ? "hidden" : ""}>
          <details class="border" open={!!!currentName}>
            <summary class="fill-slate-300 p-1">{currentName}</summary>
            <form bind:this={formEdit} class="px-1">
              <label class="mb-1 flex">
                <span class="w-20 inline-block">ID</span>
                <input
                  type="text"
                  name="id"
                  placeholder="----"
                  disabled={selected == "create"}
                  value={editID}
                  readonly
                  class="bg-gray-100 flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
                />
              </label>
              <label class="mb-1 flex">
                <span class="w-20 inline-block">Name</span>
                <input
                  type="text"
                  name="name"
                  placeholder="uniquename"
                  autocomplete="off"
                  class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
                />
                <span class="w-40 pl-1">{currentName}</span>
              </label>
              <label class="mb-1 flex">
                <span class="w-20 inline-block">Groups</span>
                <input
                  type="text"
                  name="groups"
                  placeholder="admin, deepcore"
                  class="flex-grow px-2 border border-gray-300 focus:border-red-300 focus:outline-none focus:ring focus:ring-red-200 focus:ring-opacity-50 disabled:bg-gray-100"
                />
                <span class="w-40 pl-1">{currentGroups ?? ""}</span>
              </label>
            </form>
          </details>
        </div>
      </div>
      <div class="flex-1">
        <div class="flex justify-end gap-2">
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={() =>
              selected == "table"
                ? setSelected("create")
                : createControl(selected)}
            >{selected == "table" ? "Create" : "Save"}</button
          >
          <button
            class="bg-gray-200 p-1 font-bold inline-block hover:bg-yellow-200 w-40"
            on:click={() =>
              selected == "table" ? listControls(0) : setSelected("table")}
            >{selected == "table" ? "Reload" : "Cancel"}</button
          >
        </div>
        <div
          class={`bg-red-200 px-1 w-full ${!!error ? "mt-2" : "invisible"} ${
            selected == "table" ? "hidden" : ""
          }`}
        >
          <span class="break-all">{error}</span>
        </div>
      </div>
    </div>
  </div>

  <div>
    <div
      class={`bg-slate-50 p-5 ${selected == "table" ? "" : "hidden"}`}
      bind:this={listenElement}
    >
      <div class="flex items-center justify-end mb-1">
        <Search {searchFn} bind:search />
      </div>
      <div class="overflow-x-auto rounded-none bg-white">
        <table class="w-full table-custom">
          <thead>
            <tr>
              <th style="width:5%" />
              <th style="width:35%">name</th>
              <th>groups</th>
              <th style="width:20%" />
            </tr>
          </thead>
          <tbody>
            {#each datas as d, i (d.id)}
              <tr>
                <th>{i + 1}</th>
                <th>{d.name}</th>
                <th>{d.groups ? d.groups : ""}</th>
                <th>
                  <button
                    data-id={d.id}
                    data-name={d.name}
                    data-groups={d.groups}
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
        <Pagination {meta} listF={listControls} />
        <NoData hide={!!datas.length} />
      </div>
    </div>

    <div
      class={`h-full border border-gray-600 relative ${
        selected == "table" ? "hidden" : ""
      } ${fullScreen ? "fullscreen" : ""}`}
      on:keydown={listenKeys}
    >
      <div
        class="absolute z-10 bg-slate-200 flex items-stretch border-b border-r border-gray-600"
      >
        <span
          class="hover:bg-yellow-200 hover:cursor-move select-none inline-block px-2 py-1"
          draggable="true"
          on:dragstart={dragNode}><Icon icon="plus" class="py-1" /></span
        >
        <select
          name="nodes"
          bind:value={nodeSelected}
          class="px-1 border-l border-gray-600 bg-white"
        >
          {#each Object.entries(nodes) as n (n)}
            <option value={n[0]}>{n[1].name}</option>
          {/each}
        </select>
        {#if nodes[nodeSelected].optionalInput}
          <input
            type="number"
            class="py-1 px-2 border-l border-gray-600 max-w-[4rem]"
            min="1"
            max="99"
            bind:value={inputCount}
          />
        {/if}
        {#if editor}
          {#if editor.editor_mode == "fixed"}
            <button
              class="hover:bg-yellow-200 px-2 py-1 border-l border-gray-600"
              on:click={() => (editor.editor_mode = "edit")}
            >
              <Icon icon="lock" />
            </button>
          {:else if (editor.editor_mode = "edit")}
            <button
              class="hover:bg-yellow-200 px-2 py-1 border-l border-gray-600"
              on:click={() => (editor.editor_mode = "fixed")}
            >
              <Icon icon="unlock" />
            </button>
          {/if}
        {/if}
        <button
          class="border-l border-gray-600 hover:bg-yellow-200 px-2 py-1"
          on:click={() => (fullScreen = !fullScreen)}
        >
          {#if fullScreen}
            <Icon icon="unfull" />
          {:else}
            <Icon icon="full" />
          {/if}
        </button>
        <!-- <button
          class="border-l border-gray-600 hover:bg-yellow-200 px-2 py-1"
          on:click={() => console.log(editor.export().drawflow.Home.data)}
        >
          log
        </button> -->
      </div>
      <div
        class="grid grid-cols-[1fr_auto] h-full absolute top-0 left-0 w-full"
        bind:this={listenNodes}
      >
        <div
          bind:this={drawDiv}
          on:drop|preventDefault={dropNode}
          on:dragover|preventDefault={() => void {}}
          on:mousemove={handleMouseMove}
          class="parent-drawflow parent-drawflow-style h-full w-full"
        />
        <div
          class="w-56 bg-gray-100 border-l border-black p-1 overflow-auto properties"
        >
          <Nodes node={lastSelectedNode} {editor} {nodeUnselected} />
        </div>
      </div>
    </div>
  </div>
</div>
