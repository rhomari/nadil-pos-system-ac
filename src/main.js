import { createApp } from "vue";

import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import { loadFonts } from "./plugins/webfontloader";

loadFonts();

const app = createApp(App).use(vuetify);
//app.config.errorHandler = () => null;
app.config.warnHandler = () => null;
history.pushState(null, null, location.href);
window.onpopstate = function () {
  history.go(1);
};

app.mount("#app");
