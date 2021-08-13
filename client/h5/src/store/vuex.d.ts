import { ComponentCustomProperties } from "vue";
import { Store } from "vuex";
import { AllState } from "./index.d";

declare module "@vue/runtime-core" {
  // provide typings for `this.$store`
  interface ComponentCustomProperties {
    $store: Store<AllState>;
  }
}
