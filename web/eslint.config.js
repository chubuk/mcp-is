// Safety-only rules: catch code-injection footguns, no style nits.
const noUnsanitized = require("eslint-plugin-no-unsanitized");

module.exports = [
  {
    files: ["static/js/**/*.js"],
    plugins: {
      "no-unsanitized": noUnsanitized,
    },
    rules: {
      "no-eval": "error",
      "no-implied-eval": "error",
      "no-new-func": "error",
      "no-script-url": "error",
      // escape.methods: esc() из api.js помечает значение экранированным;
      // осознанные innerHTML-присваивания несут eslint-disable-next-line с причиной
      "no-unsanitized/method": ["error", { escape: { methods: ["esc"] } }],
      "no-unsanitized/property": ["error", { escape: { methods: ["esc"] } }],
    },
  },
];
