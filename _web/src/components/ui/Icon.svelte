<script lang="ts">
  /* eslint-disable no-unused-vars */

  enum IconTypes {
    folder = "M1.979 8l-.979-6h6.972c1.265 1.583 1.327 2 3.306 2h11.722l-.96 4h-2.04l.469-2h-9.191c-2.326 0-3.166-.612-4.278-2h-3.694l.694 4h-2.021zm22.021 2l-2 12h-20l-2-12h6.278l2.308-3 3.9 3h1.97l-1.429-1.097 1.465-1.903 3.9 3h5.608zm-2.331 2h-19.338l1.385 8h16.568l1.385-8z",
    file = "M22 24h-20v-24h14l6 6v18zm-7-23h-12v22h18v-16h-6v-6zm3 15v1h-12v-1h12zm0-3v1h-12v-1h12zm0-3v1h-12v-1h12zm-2-4h4.586l-4.586-4.586v4.586z",
  }
  type IconTypesString = keyof typeof IconTypes;

  type PosTypes = "left" | "down" | "up" | undefined;

  export let icon: IconTypesString;

  export let height = "24";
  export let width = "24";

  let className: string = undefined;
  export { className as class };
  export let position: PosTypes = undefined;

  const totalClass = [className, position].join(" ");
</script>

<svg
  xmlns="http://www.w3.org/2000/svg"
  class={totalClass ? totalClass : undefined}
  {width}
  {height}
  viewBox="0 0 24 24"
>
  <path d={IconTypes[icon]} />
</svg>

<style lang="scss">
  .left {
    transform: rotate(180deg);
  }

  .down {
    transform: rotate(90deg);
  }

  .up {
    transform: rotate(-90deg);
  }
</style>
