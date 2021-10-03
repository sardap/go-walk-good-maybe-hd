import { createApp } from "vue";
import App from "./App.vue";
import AudioVisual from "vue-audio-visual";

const app = createApp(App);
app.use(AudioVisual);
app.mount("#app");
