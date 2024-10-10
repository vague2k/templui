const baseConfig = require("./tailwind.config.js");

module.exports = {
  ...baseConfig,
  content: ["./pkg/**/*.templ", "./pkg/**/*.go"],
};
