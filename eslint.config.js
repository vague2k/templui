const htmlPlugin = require("@html-eslint/eslint-plugin");
const eslintHTMLParser = require("@html-eslint/parser");
const tailwind = require("eslint-plugin-html-tailwind");

module.exports = [
  {
    files: ["out/**/*.html"],
    languageOptions: {
      parser: eslintHTMLParser,
    },
    plugins: {
      "@html-eslint": htmlPlugin,
      "html-tailwind": tailwind,
    },
    rules: {
      ...htmlPlugin.configs.recommended.rules,
      ...tailwind.configs.recommended.rules,
      "html-tailwind/classname-order": "off",
      "html-tailwind/no-style-attribute": "warn",
      "html-tailwind/no-contradicting-classnames": "warn",
    },
  },
];
