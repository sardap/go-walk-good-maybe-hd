<template>
  <div class="hello">
    <canvas id="canvas" :width="width" height="120" v-on:click="canvasClick" />
    <img class="icon" id="AIcon" src="@/assets/XboxSeriesX_A.png" />
    <img class="icon" id="BIcon" src="@/assets/XboxSeriesX_B.png" />
    <img class="icon" id="XIcon" src="@/assets/XboxSeriesX_X.png" />
    <img class="icon" id="YIcon" src="@/assets/XboxSeriesX_Y.png" />
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { Input, karaokeFile } from "../App.vue";

@Options({
  props: {
    width: Number,
    currentTime: Number,
    duration: Number,
    tempo: Number,
    inputs: {
      type: Array,
      default: [],
    },
  },
  data() {
    return {
      canvas: null,
    };
  },

  mounted() {
    var c = document.getElementById("canvas") as any;
    this.canvas = c.getContext("2d");

    this.draw();
  },
  methods: {
    drawLine(x1: number, y1: number, x2: number, y2: number, color: string) {
      let ctx = this.canvas;
      ctx.beginPath();
      ctx.strokeStyle = color;
      ctx.lineWidth = 4;
      ctx.moveTo(x1, y1);
      ctx.lineTo(x2, y2);
      ctx.stroke();
      ctx.closePath();
    },
    darwImage(image: any, x: number, y: number, w: number, h: number) {
      let ctx = this.canvas;
      ctx.drawImage(image, x, y, w, h);
    },
    canvasClick(event: any) {
      if (!karaokeFile) {
        return;
      }

      const time = (event.offsetX / this.width) * this.duration * 1000;
      console.log("Fucking time" + time);

      let bestIdx = -1;
      let bestDiff = this.duration;
      for (let i = 0; i < karaokeFile.inputs.length; i++) {
        const top = karaokeFile.inputs[i];
        const delta = Math.abs(top.hit_time - time);
        if (delta < bestDiff && delta < 500) {
          bestDiff = delta;
          bestIdx = i;
        }
      }

      if (bestIdx > -1) {
        karaokeFile.inputs.splice(bestIdx, 1);
      }
    },
    draw() {
      const step = this.tempo / 4 / 60;
      console.log("step " + step);
      console.log("duration " + this.duration);
      for (let time = 0; time < this.duration; time += step) {
        const x = (time / this.duration) * this.width;
        this.drawLine(x, 0, x, 120, "grey");
      }

      const x = (this.width * this.currentTime) / this.duration;
      this.drawLine(x, 0, x, 100, "red");

      for (let i of this.inputs) {
        let imgKey = i.sound + "Icon";
        let img = document.getElementById(imgKey);
        const x = (this.width * (i.hit_time / 1000)) / this.duration;
        let y = 0;
        switch (i.sound) {
          case "A":
            y = 75;
            break;
          case "B":
            y = 50;
            break;
          case "X":
            y = 25;
            break;
          case "Y":
            y = 0;
            break;
        }
        this.darwImage(img, x - 15, y, 30, 30);
      }
    },
  },
})
export default class RenderWindow extends Vue {
  width!: number;
  currentTime!: number;
  duration!: number;
  inputs!: Input[];
  tempo!: number;
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
#myCanvas {
  border: 1px solid grey;
}
.icon {
  width: 0px;
  height: 0px;
}
</style>
