import plugin from "tailwindcss/plugin";

/** @type {import("tailwindcss").Config} */
const config = {
  content: { files: ["internal/web/**/*.templ"] },
  plugins: [
    plugin(({ addVariant }) => {
      addVariant("hocus", ["&:hover", "&:focus"]);
      addVariant("tab-active", "&[data-active='true']");
    })
  ],
  theme: {
    extend: {
      colors: {
        border: "#1c1c23",
        inactive: "#86868c"
      }
    }
  }
};

export default config;