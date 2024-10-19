const reactRecommended = require('eslint-plugin-react/configs/recommended');
const reactHooks = require('eslint-plugin-react-hooks');
const globals = require('globals');

const GLOBALS_BROWSER_FIX = Object.assign({}, globals.browser, {
    AudioWorkletGlobalScope: globals.browser['AudioWorkletGlobalScope ']
});

delete GLOBALS_BROWSER_FIX['AudioWorkletGlobalScope '];

module.exports = [
  {
    files: ['**/*.{js,mjs,cjs,jsx,mjsx,ts,tsx,mtsx}'],
    ...reactRecommended,
  },
  {
    files: ['**/*.{js,mjs,cjs,jsx,mjsx,ts,tsx,mtsx}'],
    plugins: {
      reactHooks,
    },
    languageOptions: {
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
      globals: {
        ...globals.serviceworker,
        ...GLOBALS_BROWSER_FIX,
      },
    },
    rules: {
      "react/prop-types": "off",
      "react/react-in-jsx-scope": "off",
      "react/jsx-uses-react": "off",
      "reactHooks/rules-of-hooks": "error",
      "reactHooks/exhaustive-deps": "warn"
    }
  },
];
