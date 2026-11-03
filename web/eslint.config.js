// Safety-only rules: catch code-injection footguns, no style nits.
module.exports = [
  {
    files: ["static/js/**/*.js"],
    rules: {
      "no-eval": "error",
      "no-implied-eval": "error",
      "no-new-func": "error",
      "no-script-url": "error",
    },
  },
];
