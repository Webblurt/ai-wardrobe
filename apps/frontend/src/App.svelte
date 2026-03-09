<script lang="ts">
import { createJob, getJob } from "./api"

let person: File | null = null
let garment: File | null = null

let personPreview = ""
let garmentPreview = ""

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

  const job = await createJob(person, garment)

  jobId = job.job_id
  status = job.status

  poll()
}

async function poll() {

  if (!jobId) return

  const interval = setInterval(async () => {

    const job = await getJob(jobId!)

    status = job.status

    if (job.status === "completed") {
      resultUrl = job.result_url
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

button{
  padding:10px 20px;
  font-size:16px;
}

img{
  margin-top:10px;
}

</style>