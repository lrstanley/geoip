/* eslint no-undef: "off" */
module.exports = {
    extends: [
        // add more generic rulesets here, such as:
        'eslint:recommended',
        // 'plugin:vue/strongly-recommended',
        'plugin:vue/recommended'
    ],
    rules: {
        // override/add rules settings here, such as:
        // 'vue/no-unused-vars': 'error'
        'vue/html-indent': 'off',
        'vue/max-attributes-per-line': 'off',
        'vue/singleline-html-element-content-newline': 'off',
        'vue/component-definition-name-casing': 'off',
        'newline-per-chained-call': 'off'
    }
}
