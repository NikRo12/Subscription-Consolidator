/**
 * @type {import('prettier').Options}
 */
export default {
  useTabs: true,
  semi: false,
  singleQuote: false,
  jsxSingleQuote: false,
  overrides: [
    {
      files: "*.json",
      options: {
        useTabs: false,
      },
    },
  ],
}
