const js = require('@eslint/js')
const pluginVue = require('eslint-plugin-vue')
const parserVue = require('vue-eslint-parser')
const configPrettier = require('eslint-config-prettier')
const pluginPrettier = require('eslint-plugin-prettier')
const parserTypeScript = require('@typescript-eslint/parser')
const pluginTypeScript = require('@typescript-eslint/eslint-plugin')

;(module.exports = {
  ...js.configs.recommended,
  ignores: ['**/.*', 'dist/*', '*.d.ts', 'public/*', 'src/assets/**', 'src/**/iconfont/**'],
  languageOptions: {
    globals: {
      // index.d.ts
      RefType: 'readonly',
      EmitType: 'readonly',
      TargetContext: 'readonly',
      ComponentRef: 'readonly',
      ElRef: 'readonly',
      ForDataType: 'readonly',
      AnyFunction: 'readonly',
      PropType: 'readonly',
      Writable: 'readonly',
      Nullable: 'readonly',
      NonNullable: 'readonly',
      Recordable: 'readonly',
      ReadonlyRecordable: 'readonly',
      Indexable: 'readonly',
      DeepPartial: 'readonly',
      Without: 'readonly',
      Exclusive: 'readonly',
      TimeoutHandle: 'readonly',
      IntervalHandle: 'readonly',
      Effect: 'readonly',
      ChangeEvent: 'readonly',
      WheelEvent: 'readonly',
      ImportMetaEnv: 'readonly',
      Fn: 'readonly',
      PromiseFn: 'readonly',
      ComponentElRef: 'readonly',
      parseInt: 'readonly',
      parseFloat: 'readonly',
    },
  },
  plugins: {
    prettier: pluginPrettier,
  },
  rules: {
    ...configPrettier.rules,
    ...pluginPrettier.configs.recommended.rules,
    'no-debugger': 'off',
    'no-unused-vars': [
      'error',
      {
        argsIgnorePattern: '^_',
        varsIgnorePattern: '^_',
      },
    ],
    'prettier/prettier': [
      'error',
      {
        endOfLine: 'auto',
      },
    ],
  },
}),
  {
    files: ['**/*.?([cm])ts', '**/*.?([cm])tsx'],
    languageOptions: {
      parser: parserTypeScript,
      parserOptions: {
        sourceType: 'module',
        warnOnUnsupportedTypeScriptVersion: false,
      },
    },
    plugins: {
      '@typescript-eslint': pluginTypeScript,
    },
    rules: {
      ...pluginTypeScript.configs.strict.rules,
      '@typescript-eslint/ban-types': 'off',
      '@typescript-eslint/no-redeclare': 'error',
      '@typescript-eslint/ban-ts-comment': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/prefer-as-const': 'warn',
      '@typescript-eslint/no-empty-function': 'off',
      '@typescript-eslint/no-non-null-assertion': 'off',
      '@typescript-eslint/no-unused-expressions': 'off',
      '@typescript-eslint/no-unsafe-function-type': 'off',
      '@typescript-eslint/no-import-type-side-effects': 'error',
      '@typescript-eslint/explicit-module-boundary-types': 'off',
      '@typescript-eslint/consistent-type-imports': [
        'error',
        { disallowTypeAnnotations: false, fixStyle: 'inline-type-imports' },
      ],
      '@typescript-eslint/prefer-literal-enum-member': ['error', { allowBitwiseExpressions: true }],
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
        },
      ],
    },
  },
  {
    files: ['**/*.d.ts'],
    rules: {
      'eslint-comments/no-unlimited-disable': 'off',
      'import/no-duplicates': 'off',
      'unused-imports/no-unused-vars': 'off',
    },
  },
  {
    files: ['**/*.?([cm])js'],
    rules: {
      '@typescript-eslint/no-require-imports': 'off',
      '@typescript-eslint/no-var-requires': 'off',
    },
  },
  {
    files: ['**/*.vue'],
    languageOptions: {
      globals: {
        $: 'readonly',
        $$: 'readonly',
        $computed: 'readonly',
        $customRef: 'readonly',
        $ref: 'readonly',
        $shallowRef: 'readonly',
        $toRef: 'readonly',
      },
      parser: parserVue,
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
        extraFileExtensions: ['.vue'],
        parser: '@typescript-eslint/parser',
        sourceType: 'module',
      },
    },
    plugins: {
      vue: pluginVue,
    },
    processor: pluginVue.processors['.vue'],
    rules: {
      ...pluginVue.configs.base.rules,
      ...pluginVue.configs['vue3-essential'].rules,
      ...pluginVue.configs['vue3-recommended'].rules,
      'no-undef': 'off',
      'no-unused-vars': 'off',
      'vue/no-v-html': 'off',
      'vue/require-default-prop': 'off',
      'vue/require-explicit-emits': 'off',
      'vue/multi-word-component-names': 'off',
      'vue/no-setup-props-reactivity-loss': 'off',
      'vue/html-self-closing': [
        'error',
        {
          html: {
            void: 'always',
            normal: 'always',
            component: 'always',
          },
          svg: 'always',
          math: 'always',
        },
      ],
    },
  }