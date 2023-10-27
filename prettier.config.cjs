module.exports = {
  plugins: ["prettier-plugin-go-template"],
  overrides: [
    {
      files: ["*.gohtml", "*.tmpl", "*.html.tmpl"],
      options: {
        parser: "go-template",
      },
    },
  ],
};
