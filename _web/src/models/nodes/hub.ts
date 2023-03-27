import type { node } from "@/models/node";

export type hubData = {
  tags: string
};

export const hub: node = {
  name: "hub",
  html: `
  <div>
    <div class="title-box">Hub</div>
  </div>
  `,
  data: {
    tags: "",
  },
  input: 1,
  output: 1,
  class: "node-hub title-box-alone",
};
