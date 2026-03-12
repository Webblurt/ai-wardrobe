<script lang="ts">
import { createJob, getJob } from "./api"

let provider = "fedjaz"

let person: File | null = null
let garment: File | null = null

let personPreview = ""
let garmentPreview = ""

let description = ""
let category = "upper_body"
let steps = 30
let seed = 0
let autocrop = false
let upscale = 1
let upscaler = "ultrasharp"

let jobId: string | null = null
let status = ""
let resultUrl: string | null = null
let loading = false

function onPerson(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return

  person = file
  personPreview = URL.createObjectURL(file)
}

function onGarment(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return

  garment = file
  garmentPreview = URL.createObjectURL(file)
}

async function startTryOn() {

  if (!person || !garment) {
    alert("Upload both images")
    return
  }

  loading = true

  try {

    const job = await createJob({
      provider,
      person,
      garment,
      description,
      category,
      steps,
      seed,
      autocrop,
      upscale,
      upscaler
    })

    jobId = job.job_id
    status = job.status

    poll()

  } catch (err) {
    alert("Failed to start job")
    loading = false
  }

}

async function poll() {

  if (!jobId) return

  const interval = setInterval(async () => {

    const job = await getJob(jobId!)

    status = job.status

    if (job.status === "completed") {
      resultUrl = job.result_url ?? null
      loading = false
      clearInterval(interval)
    }

    if (job.status === "failed") {
      loading = false
      clearInterval(interval)
      alert("Generation failed")
    }

  }, 2000)

}
</script>

<h1>AI Wardrobe</h1>

<div class="upload">

  <div>
    <p>Person</p>
    <input type="file" accept="image/*" on:change={onPerson} />

    {#if personPreview}
      <img src={personPreview} width="200"/>
    {/if}
  </div>

  <div>
    <p>Garment</p>
    <input type="file" accept="image/*" on:change={onGarment} />

    {#if garmentPreview}
      <img src={garmentPreview} width="200"/>
    {/if}
  </div>

</div>

<div class="options">

  <div>
    <label>Provider</label>
    <select bind:value={provider}>
      <option value="fedjaz">Fedjaz VTON</option>
      <option value="replicate">Replicate</option>
    </select>
  </div>

  <div>
    <label>Description</label>
    <input type="text" bind:value={description} />
  </div>

  <div>
    <label>Category</label>
    <select bind:value={category}>
      <option value="upper_body">Upper body</option>
      <option value="upper_body_open">Upper body open</option>
      <option value="lower_body">Lower body</option>
      <option value="dresses">Dresses</option>
      <option value="skirt">Skirt</option>
      <option value="skirt_short">Skirt short</option>
      <option value="skirt_mini">Skirt mini</option>
      <option value="shoes">Shoes</option>
      <option value="socks">Socks</option>
      <option value="stockings">Stockings</option>
    </select>
  </div>

  <div>
    <label>Steps</label>
    <input type="number" bind:value={steps} min="1" max="100"/>
  </div>

  <div>
    <label>Seed</label>
    <input type="number" bind:value={seed}/>
  </div>

  <div>
    <label>Autocrop</label>
    <input type="checkbox" bind:checked={autocrop}/>
  </div>

  <div>
    <label>Upscale</label>
    <select bind:value={upscaler}>
      <option value="1">1</option>
      <option value="2">2</option>
      <option value="3">4</option>
    </select>
  </div>

  <div>
    <label>Upscaler</label>
    <select bind:value={upscaler}>
      <option value="ultrasharp">UltraSharp</option>
      <option value="realesrgan">RealESRGAN</option>
    </select>
  </div>

</div>

<button on:click={startTryOn} disabled={loading}>
  Try On
</button>

{#if status}
<p>Status: {status}</p>
{/if}

{#if resultUrl}
<h2>Result</h2>
<img src={resultUrl} width="400"/>
{/if}

<style>

h1{
  font-family:sans-serif
}

.upload{
  display:flex;
  gap:40px;
  margin-bottom:20px;
}

.options{
  display:grid;
  grid-template-columns:repeat(2,200px);
  gap:12px;
  margin-bottom:20px;
  font-family:sans-serif;
}

label{
  display:block;
  font-size:14px;
  margin-bottom:4px;
}

input, select{
  width:100%;
  padding:6px;
}

button{
  padding:10px 20px;
  font-size:16px;
}

img{
  margin-top:10px;
}

</style>