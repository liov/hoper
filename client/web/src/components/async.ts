import {
  defineAsyncComponent,
  defineComponent,
  h,
  resolveDynamicComponent,
  Suspense,
} from "vue";
import type { Component, AsyncComponentLoader } from "vue";

const AsyncComponent = (component: Component) => (props, context) => {
  console.log(component);
  console.log(props);
  console.log(context);
  context.slots = component;
  return h(Suspense, context.attrs, context.slots);
};

export default AsyncComponent;

const AsyncComponent5: Component = {
  render() {
    const component = resolveDynamicComponent(this.name);
    return h(component as Component);
  },
  props: {
    level: {
      type: Number,
      required: true,
    },
  },
};

const AsyncComponent6: Component = {
  props: ["modelValue"],
  emits: ["update:modelValue"],
  render() {
    return h(
      "div",
      {
        onClick: () => {
          this.$emit("update:modelValue", "123");
        },
      },
      this.modelValue
    );
  },
  template: `
      <input v-model="value">
    `,
  computed: {
    value: {
      get() {
        return this.modelValue;
      },
      set(value) {
        this.$emit("update:modelValue", value);
      },
    },
  },
};

const Component = defineComponent({
  // 已启用类型推断
});

const AsyncComponent1 = (component: Component) =>
  defineComponent({
    // 已启用类型推断
  });

const AsyncComponent2 = (component: AsyncComponentLoader) =>
  defineAsyncComponent({
    delay: 0,
    errorComponent: undefined,
    loader: component,
    loadingComponent: undefined,
    onError: (error, retry, fail, attempts) => {
      //
    },
    suspensible: false,
    timeout: 0,
    // 已启用类型推断
  });

const AsyncComponent3 = (path: string) =>
  defineAsyncComponent(() => import(`./src/${path}.vue`));
