import type { node } from "@/models/node";

export type respondData = {
  status: string,
  headers: string
  get: boolean,
  tags: string
};

export const respond: node = {
  name: "respond",
  html: `
  <div>
    <div class="title-box">Respond</div>
  </div>
  `,
  data: {
    status: "200",
    headers: "",
    get: false,
    tags: "",
  },
  input: 1,
  output: 0,
  class: "node-respond",
};
