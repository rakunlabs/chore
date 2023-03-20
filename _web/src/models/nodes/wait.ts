import type { node } from "@/models/node";

export type waitData = {
  tags: string
};

export const wait: node = {
  name: "wait",
  html: `
  <div>
    <div class="title-box">Wait</div>
  </div>
  `,
  data: {
    tags: "",
  },
  input: 2,
  output: 1,
  class: "node-wait",
};
