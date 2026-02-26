/**
 * @type {import('prettier').Options}
 */
export default {
  useTabs: true,
  semi: false,
  overrides: [
    {
      files: "*.json",
      options: {
        useTabs: false,
      },
    },
  ],
}
