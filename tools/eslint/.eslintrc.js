module.exports = {
  extends: ["plugin:@einride/default"],
  rules: {
    "jest/no-deprecated-functions": "off", // we're not using Jest
    "prettier/prettier": "off", // we're not concerned with code style
  },
};
