import type { node } from "@/models/node";

export type noteData = {
  note: string
  backgroundColor: string
  textColor: string
  width: number
  height: number
};

export const note: node = {
  name: "note",
  html: `
  <div>
    <div class="box">
      <textarea df-note data-name="note"></textarea>
    </div>
  </div>
  `,
  data: {
    note: "",
    backgroundColor: "#FEF9C3",
    textColor: "#000000",
    width: 262,
    height: 32,
  },
  input: 0,
  output: 0,
  class: "node-note",
};
