<script lang="ts">
  /* eslint-disable no-unused-vars */

  enum IconTypes {
    folder = "M1.979 8l-.979-6h6.972c1.265 1.583 1.327 2 3.306 2h11.722l-.96 4h-2.04l.469-2h-9.191c-2.326 0-3.166-.612-4.278-2h-3.694l.694 4h-2.021zm22.021 2l-2 12h-20l-2-12h6.278l2.308-3 3.9 3h1.97l-1.429-1.097 1.465-1.903 3.9 3h5.608zm-2.331 2h-19.338l1.385 8h16.568l1.385-8z",
    file = "M22 24h-20v-24h14l6 6v18zm-7-23h-12v22h18v-16h-6v-6zm3 15v1h-12v-1h12zm0-3v1h-12v-1h12zm0-3v1h-12v-1h12zm-2-4h4.586l-4.586-4.586v4.586z",
    plus = "M24 10h-10v-10h-4v10h-10v4h10v10h4v-10h10z",
    minus = "M0 10h24v4h-24z",
    right = "M7.33 24l-2.83-2.829 9.339-9.175-9.339-9.167 2.83-2.829 12.17 11.996z",
    left = "M16.67 0l2.83 2.829-9.339 9.175 9.339 9.167-2.83 2.829-12.17-11.996z",
    ok = "M9 21.035l-9-8.638 2.791-2.87 6.156 5.874 12.21-12.436 2.843 2.817z",
    close = "M23.954 21.03l-9.184-9.095 9.092-9.174-2.832-2.807-9.09 9.179-9.176-9.088-2.81 2.81 9.186 9.105-9.095 9.184 2.81 2.81 9.112-9.192 9.18 9.1z",
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
