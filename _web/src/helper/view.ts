import type Drawflow from "drawflow";

let visit: Record<number, boolean> = {};
let tags = "";

const compareTags = (tagsCompare: string) => {
  if (tags == "" || tagsCompare == "") {
    return true;
  }

  const tags1List = tags.replaceAll(",", " ").split(RegExp("\\s+"));
  const tags2List = tagsCompare.replaceAll(",", " ").split(RegExp("\\s+"));

  // check if tags1List has one of the elements in tags2List
  return tags1List.some((tag: string) => tags2List.includes(tag));
}

export const viewEndpoint = (endpoint: string, drawDiv: HTMLDivElement, editor: Drawflow) => {
  visit = {};
  tags = "";

  editor.getNodesFromName("endpoint").forEach((id: number) => {
    let node = editor.getNodeFromId(id);

    if (node.data.endpoint != endpoint) {
      return;
    }

    addClassToNodeId(node.id, drawDiv);

    visit[node.id] = true;
    tags = node.data.tags ?? "";

    visitor(id, drawDiv, editor);
  });

  Object.keys(editor.drawflow.drawflow.Home.data).forEach((id: any) => {
    if (visit[id]) {
      return;
    }

    if (editor.drawflow.drawflow.Home.data[id].name == "note") {
      return;
    }

    let nodeSelectedQuery = drawDiv.querySelector(`#node-${id}`);
    if (nodeSelectedQuery) {
      nodeSelectedQuery.classList.add("view-endpoint-disabled");
    }
  });

  drawDiv.classList.add("view-disabled");
}

const addClassToNodeId = (id: number, drawDiv: HTMLDivElement) => {
  let nodeSelectedQuery = drawDiv.querySelector(`#node-${id}`);
  if (nodeSelectedQuery) {
    nodeSelectedQuery.classList.add("view-endpoint");
  }
}

const addClassToConnection = (idOut: number, idIn: number, drawDiv: HTMLDivElement) => {
  let nodeSelectedQuery = drawDiv.querySelector(`.connection.node_in_node-${idIn}.node_out_node-${idOut}`);
  if (nodeSelectedQuery) {
    nodeSelectedQuery.classList.add("view-connection");
  }
}

const addOtherClassToConnection = (drawDiv: HTMLDivElement) => {
  drawDiv.querySelectorAll(`.connection:not(.view-connection)`).forEach((node: any) => {
    node.classList.add("view-connection-disabled");
  });
}

const visitor = (id: number, drawDiv: HTMLDivElement, editor: Drawflow) => {
  visit[id] = true;
  addClassToNodeId(id, drawDiv);

  let node = editor.getNodeFromId(id);
  Object.keys(node.outputs).forEach((output: any) => {
    node.outputs[output].connections.forEach((connection: any) => {
      const nodeNext = editor.getNodeFromId(connection.node)
      if (!compareTags(nodeNext.data.tags ?? "")) {
        return;
      }

      addClassToConnection(id, connection.node, drawDiv);

      if (visit[connection.node]) {
        return;
      }

      visitor(connection.node, drawDiv, editor);
    });
  });

  addOtherClassToConnection(drawDiv);
}

export const viewEndpointClear = (drawDiv: HTMLDivElement) => {
  drawDiv.querySelectorAll(".view-endpoint").forEach((node: any) => {
    node.classList.remove("view-endpoint");
  });

  drawDiv.querySelectorAll(".view-endpoint-disabled").forEach((node: any) => {
    node.classList.remove("view-endpoint-disabled");
  });

  drawDiv.classList.remove("view-disabled");

  drawDiv.querySelectorAll(".view-connection").forEach((node: any) => {
    node.classList.remove("view-connection");
  });

  drawDiv.querySelectorAll(".view-connection-disabled").forEach((node: any) => {
    node.classList.remove("view-connection-disabled");
  });
}
