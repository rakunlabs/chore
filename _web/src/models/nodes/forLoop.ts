import type { node } from "@/models/node";

export type forLoopData = {
  for: string
  tags: string
};

export const forLoop: node = {
  name: "for",
  html: `
  <div>
    <div class="title-box">For Loop</div>
    <div class="box">
      <p>Expression</p>
      <input type="text" placeholder="data" name="for" readonly disabled df-for>
    </div>
  </div>
  `,
  data: {
    for: "data",
    tags: "",
  } as forLoopData,
  input: 1,
  output: 1,
  class: "node-for",
};
