import type { node } from "@/models/node";

export type ifCaseData = {
  if: string
  tags: string
};

export const ifCase: node = {
  name: "if",
  html: `
  <div>
    <div class="title-box">IF</div>
    <div class="box">
      <p>Expression</p>
      <input type="text" placeholder="write expression" name="if" readonly disabled df-if>
    </div>
  </div>
  `,
  data: {
    if: "data > 0",
    tags: "",
  } as ifCaseData,
  input: 1,
  output: 2,
  class: "node-if",
};
