export default {
  mounted: (el: HTMLElement) =>
    setTimeout(() => el.querySelector<HTMLElement>("input, textarea")?.focus(), 300),
}
