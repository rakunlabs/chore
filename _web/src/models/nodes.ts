export type node = {
  name: string
  html: string
  data: Record<string, string>
  input: number
  output: number
  optionalInput?: boolean
  class?: string
  style?: string
};

const endpoint: node = {
  name: "endpoint",
  html: `
  <div>
    <div class="title-box">Endpoint</div>
    <div class="box">
      <p>Enter endpoint name</p>
      <input type="text" placeholder="create" name="endpoint" df-endpoint>
      <p>Methods</p>
      <input type="text" placeholder="POST, GET" name="methods" df-methods>
      <label>
        <span>Public</span>
        <input type="checkbox" data-parent="5" name="public" data-action="checkbox" df-public>
      </label>
    </div>
  </div>
  `,
  data: {
    endpoint: "",
    methods: "POST",
    public: "false",
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
    </div>
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
      <p>Enter auth</p>
      <input type="text" placeholder="myauth" name="auth" df-auth>
      <label>
        <span>Payload set to nil</span>
        <input type="checkbox" name="payload_nil" data-action="checkbox" df-payload_nil>
      </label>
      <label>
        <span>Skip verify certificate</span>
        <input type="checkbox" name="skip_verify" data-action="checkbox" df-skip_verify>
      </label>
      <label>
        <span>Use pooled client</span>
        <input type="checkbox" name="pool_client" data-action="checkbox" df-pool_client>
      </label>
      <details><summary>Enter additional headers</summary>
        <textarea df-headers placeholder="json/yaml key:value"></textarea>
      </details>
      <details><summary>Retry with status codes</summary>
        <p>Enabled Status Codes</p>
        <input type="text" placeholder="Ex: 401, 403" name="retry_codes" df-retry_codes>
        <p>Disabled Status Codes</p>
        <input type="text" placeholder="Ex: 500" name="retry_decodes" df-retry_decodes>
      </details>
    </div>
  </div>
  `,
  data: {
    skip_verify: "false",
    payload_nil: "false",
    pool_client: "false",
  },
  input: 2,
  output: 3,
  class: "node-request",
};

const script = {
  name: "script",
  html: `
  <div>
    <div class="title-box">Script</div>
    <div class="box">
      <button data-action="editor">Open Editor</button>
      <textarea df-script></textarea>
    </div>
  </div>
  `,
  data: {
    "inputs": null,
    "script": `function main(data) {
  return data;
}
`,
  },
  input: 1,
  output: 3,
  optionalInput: true,
  class: "node-script",
};

const forLoop = {
  name: "for",
  html: `
  <div>
    <div class="title-box">For Loop</div>
    <div class="box">
      <p>Return an array</p>
      <input type="text" placeholder="data" name="for" df-for>
    </div>
  </div>
  `,
  data: {
    "for": "data",
  },
  input: 1,
  output: 1,
};

const ifCase = {
  name: "if",
  html: `
  <div>
    <div class="title-box">IF</div>
    <div class="box">
      <p>Expression</p>
      <input type="text" placeholder="write expression" name="if" df-if>
    </div>
  </div>
  `,
  data: {
    "if": "data > 0",
  },
  input: 1,
  output: 2,
  class: "node-if",
};

const control = {
  name: "control",
  html: `
  <div>
    <div class="title-box">Control</div>
    <div class="box">
      <p>Enter control name</p>
      <input type="text" name="control" df-control>
      <p>Enter endpoint name</p>
      <input type="text" name="endpoint" df-endpoint>
      <p>Methods</p>
      <input type="text" placeholder="POST" name="method" df-method>
    </div>
  </div>
  `,
  data: {
    "control": "",
    "endpoint": "",
    "method": "POST",
  },
  input: 1,
  output: 1,
};

const respond = {
  name: "respond",
  html: `
  <div>
    <div class="title-box flex items-center gap-2">Respond</div>
    <div class="box">
      <p>Enter respond status code</p>
      <input class="mr-2" type="number" name="status" df-status>
      <p>Enter headers</p>
      <textarea df-headers placeholder="json/yaml key:value"></textarea>
      <hr>
      <label>
        <span>Get respond in data</span>
        <input type="checkbox" name="get" data-action="checkbox" df-get>
      </label>
    </div>
  </div>
  `,
  data: {
    status: "200",
    get: "false",
  },
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
      <input type="text" placeholder="awesome log message" df-message>
      <p>Log Level</p>
      <select df-level>
        <option value="debug">Debug</option>
        <option value="info">Info</option>
        <option value="warn">Warn</option>
        <option value="error">Error</option>
        <option value="">NoLevel</option>
      </select>
      <hr>
      <label>
        <span>Print data</span>
        <input type="checkbox" name="data" data-action="checkbox" df-data>
      </label>
    </div>
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
      <p>From</p>
      <input type="text" name="email-from" df-from>
      <p>To</p>
      <input type="text" name="email-to" df-to>
      <p>CC</p>
      <input type="text" name="email-cc" df-cc>
      <p>BCC</p>
      <input type="text" name="email-bcc" df-bcc>
      <p>Subject</p>
      <input type="text" name="email-subject" df-subject>
    </div>
  </div>
  `,
  data: {
    "from": "",
    "to": "",
    "cc": "",
    "bcc": "",
    "subject": "",
  },
  input: 2,
  output: 0,
  class: "node-email",
};

const note = {
  name: "note",
  html: `
  <div>
    <div class="box">
      <textarea df-note style="width: 250px; height: 100px"></textarea>
    </div>
  </div>
  `,
  data: {},
  input: 0,
  output: 0,
  class: "node-note",
};

export const nodes = { endpoint, template, request, script, forLoop, ifCase, control, respond, log, email, note } as Record<string, node>;
