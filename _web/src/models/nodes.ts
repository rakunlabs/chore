export type node = {
  name: string;
  html: string,
  data: Record<string, string>
  input: number,
  output: number,
  optionalInput?: boolean,
  class?: string
};

const endpoint: node = {
  name: "endpoint",
  html: `
  <div>
    <div class="title-box">Endpoint</div>
    <div class="box">
      <p>Enter endpoint name</p>
      <input type="text" placeholder="create" name="endpoint" df-endpoint>
    </div
  </div>
  `,
  data: {
    "endpoint": "",
  },
  input: 0,
  output: 1,
};

const template = {
  name: "template",
  html: `
  <div>
    <div class="title-box">Template</div>
    <div class="box">
      <p>Enter template name</p>
      <input type="text" placeholder="deepcore/create-issue" name="template" df-template>
    </div
  </div>
  `,
  data: {
    "template": "",
  },
  input: 1,
  output: 1,
};

const request = {
  name: "request",
  html: `
  <div>
    <div class="title-box">Request</div>
    <div class="box">
      <p>Enter request url</p>
      <input type="url" placeholder="https://createmyissue.com" name="url" df-url>
      <p>Enter method</p>
      <input type="text" placeholder="POST" name="method" df-method>
      <p>Enter additional headers</p>
      <textarea df-headers placeholder="json/yaml key:value"></textarea>
      <p>Enter auth</p>
      <input type="text" placeholder="myauth" name="auth" df-auth>
    </div
  </div>
  `,
  data: {},
  input: 2,
  output: 1,
  class: "input1-value",
};

const script = {
  name: "script",
  html: `
  <div>
    <div class="title-box">Script</div>
    <div class="box">
      <button data-action="editor">Open Editor</button>
      <textarea df-script></textarea>
    </div
  </div>
  `,
  data: {
    "script": `function main(data) {
  return data;
}
  `,
  },
  input: 1,
  output: 1,
  optionalInput: true,
};

// const control = {
//   name: "control",
//   html: `
//   <div>
//     <div class="title-box">Control</div>
//     <div class="box">
//       <p>Enter control name</p>
//       <input type="text" placeholder="mycontrol" name="url" df-control>
//       <p>Enter endpoint name</p>
//       <input type="text" placeholder="create" name="endpoint" df-endpoint>
//     </div
//   </div>
//   `,
//   data: {
//     "control": "",
//     "endpoint": "",
//   },
//   input: 1,
//   output: 1,
// };

const respond = {
  name: "respond",
  html: `
  <div>
    <div class="title-box">Respond</div>
  </div>
  `,
  data: {},
  input: 1,
  output: 0,
};

const log = {
  name: "log",
  html: `
  <div>
    <div class="title-box">Log</div>
    <div class="box">
      <p>Message</p>
      <input type="text" placeholder="awesome log message" name="url" df-message>
      <p>Log Level</p>
      <select df-level>
        <option value="debug">Debug</option>
        <option value="info">Info</option>
        <option value="warn">Warn</option>
        <option value="error">Error</option>
      </select>
      <hr>
      <label>
        <span>Print data</span>
        <input type="checkbox" name="data" data-action="checkbox" df-data>
      </label>
    </div
  </div>
  `,
  data: {
    message: "",
    level: "debug",
    data: "false",
  },
  input: 1,
  output: 1,
};

const email = {
  name: "email",
  html: `
  <div>
    <div class="title-box">Email</div>
    <div class="box">
      <button data-action="editor">Open Editor</button>
      <textarea df-script></textarea>
    </div
  </div>
  `,
  data: {
    "script":
`To:
  - user@example.com
Cc:
  - user@example.com
Subject:
  - my subject
`,
  },
  input: 1,
  output: 0,
};

export const nodes = { endpoint, template, request, script, respond, log, email } as Record<string, node>;
