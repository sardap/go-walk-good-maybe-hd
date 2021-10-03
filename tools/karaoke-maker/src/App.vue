<template>
  <div>
    <div v-if="uploaded">
      <Maker
        :src="musicPath"
        :tempo="tempo"
        :aSrc="aSrc"
        :bSrc="bSrc"
        :xSrc="xSrc"
        :ySrc="ySrc"
      />
    </div>
    <div v-else>
      <div>
        <button class="btn btn-info" @click="onPickMusicFile">
          Upload music
        </button>
        <input
          type="file"
          style="display: none"
          ref="fileInputMusic"
          accept="audio/mp3"
          @change="onAudioFilePicked"
        />
      </div>
      <div>
        <button class="btn btn-info" @click="onPickKarokeFile">
          Upload karoke
        </button>
        <input
          type="file"
          style="display: none"
          ref="fileInputKaroke"
          accept="application/json"
          @change="onKaraokeFilePicked"
        />
      </div>
      <div>
        <button class="btn btn-info" @click="onSoundAPick">
          Upload Button A Sound
        </button>
        <input
          type="file"
          style="display: none"
          ref="fileInputASound"
          accept="audio/wav"
          @change="onSoundAFilePicked"
        />
      </div>
      <div>
        <button class="btn btn-info" @click="onSoundBPick">
          Upload Button B Sound
        </button>
        <input
          type="file"
          style="display: none"
          ref="fileInputBSound"
          accept="audio/wav"
          @change="onSoundBFilePicked"
        />
      </div>
      <div>
        <button class="btn btn-info" @click="onSoundXPick">
          Upload Button X Sound
        </button>
        <input
          type="file"
          style="display: none"
          ref="fileInputXSound"
          accept="audio/wav"
          @change="onSoundXFilePicked"
        />
      </div>
      <div>
        <button class="btn btn-info" @click="onSoundYPick">
          Upload Button Y Sound
        </button>
        <input
          type="file"
          style="display: none"
          ref="fileInputYSound"
          accept="audio/wav"
          @change="onSoundYFilePicked"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Maker from "./components/Maker.vue";
import { analyze } from "web-audio-beat-detector";
import { dataUriToBuffer } from "data-uri-to-buffer";

export interface Input {
  hit_time: number;
  sound: string;
}

export interface KaraokeFile {
  inputs: Array<Input>;
}

export let karaokeFile: KaraokeFile | null = null;

@Options({
  components: {
    Maker,
  },
  data() {
    return {
      uploaded: false,
      musicPath: null,
      tempo: null,
      aSrc: null,
      bSrc: null,
      xSrc: null,
      ySrc: null,
    };
  },
  methods: {
    onPickMusicFile() {
      this.$refs.fileInputMusic.click();
    },
    onPickKarokeFile() {
      this.$refs.fileInputKaroke.click();
    },
    onAudioFilePicked(event: any) {
      const files = event.target.files;
      const fileReader = new FileReader();
      fileReader.addEventListener("load", () => {
        this.musicPath = fileReader.result;
        console.log("loaded music");

        const context = new AudioContext();
        const buffer = dataUriToBuffer(fileReader.result as string);

        context.decodeAudioData(buffer.buffer).then((buffer: AudioBuffer) => {
          analyze(buffer)
            .then((tempo) => {
              this.tempo = tempo;
              this.uploaded = this.uploadedComplete();
            })
            .catch((err) => {
              console.log("Error getting tempo " + err);
            });
        });

        this.uploaded = this.uploadedComplete();
      });
      fileReader.readAsDataURL(files[0]);
    },
    onKaraokeFilePicked(event: any) {
      const files = event.target.files;
      const fileReader = new FileReader();
      fileReader.addEventListener("load", () => {
        karaokeFile = JSON.parse(fileReader.result as any);
        console.log("Loaded karoke file");
        this.uploaded = this.uploadedComplete();
      });
      fileReader.readAsText(files[0]);
    },
    onSoundAPick() {
      this.$refs.fileInputASound.click();
    },
    onSoundAFilePicked(event: any) {
      const files = event.target.files;
      const fileReader = new FileReader();
      fileReader.addEventListener("load", () => {
        this.aSrc = fileReader.result;
        console.log("loaded aSrc");
        this.uploaded = this.uploadedComplete();
      });
      fileReader.readAsDataURL(files[0]);
    },
    onSoundBPick() {
      this.$refs.fileInputBSound.click();
    },
    onSoundBFilePicked(event: any) {
      const files = event.target.files;
      const fileReader = new FileReader();
      fileReader.addEventListener("load", () => {
        this.bSrc = fileReader.result;
        console.log("loaded bSrc");
        this.uploaded = this.uploadedComplete();
      });
      fileReader.readAsDataURL(files[0]);
    },
    onSoundXPick() {
      this.$refs.fileInputXSound.click();
    },
    onSoundXFilePicked(event: any) {
      const files = event.target.files;
      const fileReader = new FileReader();
      fileReader.addEventListener("load", () => {
        this.xSrc = fileReader.result;
        console.log("loaded xSrc");
        this.uploaded = this.uploadedComplete();
      });
      fileReader.readAsDataURL(files[0]);
    },
    onSoundYPick() {
      this.$refs.fileInputYSound.click();
    },
    onSoundYFilePicked(event: any) {
      const files = event.target.files;
      const fileReader = new FileReader();
      fileReader.addEventListener("load", () => {
        this.ySrc = fileReader.result;
        console.log("loaded ySrc");
        this.uploaded = this.uploadedComplete();
      });
      fileReader.readAsDataURL(files[0]);
    },
    uploadedComplete(): boolean {
      return (
        this.musicPath != null &&
        karaokeFile != null &&
        this.tempo != null &&
        this.aSrc != null &&
        this.bSrc != null &&
        this.xSrc != null &&
        this.ySrc != null
      );
    },
  },
})
export default class App extends Vue {}
</script>

<style lang="scss">
.btn {
  margin: 10px;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
