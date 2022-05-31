import type { DrawflowNode } from "drawflow";

const getPublicEndpoints = (exported: {[nodeKey: string]: DrawflowNode}) => {
  const values = [] as string[];
  for (const v of Object.values(exported)) {
    if (v.name == "endpoint" && v.data["public"] == "true") {
      values.push(v.data["endpoint"]);
    }
  }

  return values;
};

export { getPublicEndpoints };
