import type { DrawflowNode } from "drawflow";

type Endpoint = {
  methods: string[];
  public: boolean;
};

const getEndpoints = (exported: { [nodeKey: string]: DrawflowNode }) => {
  const values = {} as Record<string, Endpoint>;
  for (const v of Object.values(exported)) {
    if (v.name == "endpoint") {
      console.log(v.data);
      values[v.data["endpoint"]] = {
        methods: (v.data["methods"] as string).replaceAll(" ", "").toUpperCase().split(","),
        public: v.data["public"],
      };
    }
  }

  return values;
};

export { getEndpoints };
