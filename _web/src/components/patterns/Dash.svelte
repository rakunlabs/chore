<script lang="ts">
  let className = "";
  export { className as class };
</script>

<svg height="100%" width="100%" class={className}>
  <defs>
    <pattern
      id="doodad"
      width="55"
      height="55"
      viewBox="0 0 40 40"
      patternUnits="userSpaceOnUse"
      patternTransform="rotate(45)"
    >
      <rect width="100%" height="100%" fill="#e9ecef" />
      <path d="M-10 1h60v2h-60z" fill="#868e96" />
      <path d="M-10 41h60v2h-60z" fill="#868e96" />
      <path d="M-10 29h60v2h-60z" fill="#ffffff" />
      <path d="M29-10v60h2v-60z" fill="#ffffff" />
      <path d="M9-10v60h2v-60z" fill="#868e96" />
    </pattern>
  </defs>
  <rect fill="url(#doodad)" height="200%" width="200%" />
</svg>
