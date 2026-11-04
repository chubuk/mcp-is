/** @type {import('stylelint').Config} */
module.exports = {
  extends: ["stylelint-config-recommended"],
  ignoreFiles: ["static/vendor/**"],
  rules: {
    // конвенция порядка, не баг: переупорядочивание легаси-каскада может изменить UI
    "no-descending-specificity": null,
  },
};
