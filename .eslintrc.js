module.exports = {
  extends: ["plugin:@einride/default"],
  rules: {
    "jest/no-deprecated-functions": "off", // we're not using Jest
    "prettier/prettier": "off", // we're not concerned with code style
    "@typescript-eslint/ban-ts-comment": "off", // we need ts comment in generated files to disable type checking of files for consumers
  },
};
