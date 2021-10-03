<template>
  <div>
    <av-waveform
      id="player"
      :line-width="2"
      line-color="lime"
      :audio-src="src"
      :key="windowWidth"
      :canv-width="windowWidth - 100"
      :canv-height="150"
      :playtime-line-width="100"
      :playtime-font-size="14"
    ></av-waveform>
    <RenderWindow
      :width="windowWidth - 100"
      :key="[currentTime, inputsCount, duration, windowWidth]"
      :currentTime="currentTime"
      :duration="duration"
      :inputs="inputs"
      :tempo="tempo"
    />
    <img
      class="icon"
      id="A"
      src="@/assets/XboxSeriesX_A.png"
      v-on:click="addInput"
    />
    <img
      class="icon"
      id="B"
      src="@/assets/XboxSeriesX_B.png"
      v-on:click="addInput"
    />
    <img
      class="icon"
      id="X"
      src="@/assets/XboxSeriesX_X.png"
      v-on:click="addInput"
    />
    <img
      class="icon"
      id="Y"
      src="@/assets/XboxSeriesX_Y.png"
      v-on:click="addInput"
    />
    <br />
    <button class="btn" @click="dump">DUMP</button>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import RenderWindow from "./RenderWindow.vue";
import { Input, karaokeFile } from "../App.vue";

@Options({
  components: {
    RenderWindow,
  },
  props: {
    src: String,
    tempo: Number,
    aSrc: String,
    bSrc: String,
    xSrc: String,
    ySrc: String,
  },

  data() {
    return {
      windowWidth: window.innerWidth,
      audio: null,
      currentTime: 0,
      duration: 0,
      inputsCount: karaokeFile ? karaokeFile.inputs.length : 0,
      inputs: karaokeFile ? karaokeFile.inputs : [],
      played: [],
      lastTime: 0,
    };
  },
  methods: {
    dump(event: any) {
      const data = JSON.stringify(karaokeFile, null, 2);
      const blob = new Blob([data], { type: "text/plain" });
      const e = document.createEvent("MouseEvents"),
        a = document.createElement("a");
      a.download = "test.json";
      a.href = window.URL.createObjectURL(blob);
      a.dataset.downloadurl = ["text/json", a.download, a.href].join(":");
      e.initEvent("click", true, false);
      a.dispatchEvent(e);
    },
    playSound(btn: string) {
      let path = "";
      switch (btn) {
        case "A":
          path = this.aSrc;
          break;
        case "B":
          path = this.bSrc;
          break;
        case "X":
          path = this.xSrc;
          break;
        case "Y":
          path = this.ySrc;
          break;
      }

      const sound = new Audio(path);
      sound.play();
    },
    addInput: function (event: any) {
      let input: Input = {
        sound: event.target.id,
        hit_time: Math.floor(this.currentTime * 1000),
      };

      if (karaokeFile == null) {
        return;
      }

      this.played.push(karaokeFile.inputs.length);
      karaokeFile.inputs.push(input);

      this.playSound(event.target.id);
    },
  },
  mounted() {
    window.addEventListener("resize", () => {
      this.windowWidth = window.innerWidth;
    });

    window.addEventListener("keydown", (e: KeyboardEvent) => {
      console.log(`Key down ${e.key}`);

      const audio = document
        .getElementById("player")!
        .getElementsByTagName("audio")!
        .item(0);

      if (!audio) {
        return;
      }

      switch (e.key) {
        case "p":
          if (audio.paused) {
            audio.play();
          } else {
            audio.pause();
          }
          break;
        case "ArrowLeft":
          audio.pause();
          audio.currentTime -= 0.1;
          break;
        case "ArrowRight":
          audio.pause();
          audio.currentTime += 0.1;
          break;
      }
    });

    window.setInterval(() => {
      this.audio = document
        .getElementById("player")!
        .getElementsByTagName("audio")!
        .item(0);

      if (this.audio == null) {
        return;
      }

      this.duration = this.audio.duration;
      this.lastTime = this.currentTime;
      this.currentTime = this.audio.currentTime;

      if (karaokeFile == null) {
        return;
      }

      if (this.lastTime > this.currentTime) {
        this.played = [];
      }

      const audio = document
        .getElementById("player")!
        .getElementsByTagName("audio")!
        .item(0);

      this.inputsCount = karaokeFile.inputs.length;

      if (!audio || audio.paused) {
        return;
      }

      for (let i = 0; i < karaokeFile.inputs.length; i++) {
        let top = karaokeFile.inputs[i];
        const diff = this.currentTime - top.hit_time / 1000;
        if (
          diff < 2 &&
          this.currentTime > top.hit_time / 1000 &&
          !this.played.includes(i)
        ) {
          this.playSound(top.sound);
          this.played.push(i);
        }
      }
    }, 10);
  },
})
export default class Maker extends Vue {
  src!: string;
  number!: number;
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
.btn {
  width: 100px;
  height: 30px;
}
#myCanvas {
  border: 1px solid grey;
}
</style>
