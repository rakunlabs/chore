import type { node } from "@/models/node";

export type endpointData = {
  endpoint: string
  methods: string
  public: boolean
  tags: string
};

export const endpoint: node = {
  name: "endpoint",
  html: `
  <div>
    <div class="title-box">Endpoint</div>
    <div class="box">
      <input type="text" placeholder="create" name="info" readonly disabled df-endpoint>
    </div>
  </div>
  `,
  data: {
    endpoint: "",
    methods: "POST",
    public: false,
    tags: "",
  } as endpointData,
  input: 0,
  output: 1,
  class: "node-endpoint",
};
