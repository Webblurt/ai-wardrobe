<script lang="ts">
import { createJob, getJob } from "./api"

let provider = "fedjaz"

let person: File | null = null
let garment: File | null = null

let personPreview = ""
let garmentPreview = ""

let description = ""
let category = ""

let steps = 30
let seed = 1
let autocrop = false
let upscale = 1
let upscaler = "ultrasharp"

// fashn
let garmentPhotoType = "model"
let numSamples = 1
let numTimesteps = 30
let guidanceScale = 5
let segmentationFree = false

let jobId: string | null = null
let status = ""
let resultUrl: string | null = null
let loading = false


// ---------- computed ----------

$: showDefault = provider === "fedjaz" || provider === "replicate"
$: showFashn = provider === "fedjazfashnv15"


// ---------- files ----------

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


// ---------- tryon ----------

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
      upscaler,

      garmentPhotoType,
      numSamples,
      numTimesteps,
      guidanceScale,
      segmentationFree,
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
      <option value="fedjaz">Fedjaz Default</option>
      <option value="fedjazfashnv15">Fedjaz FASHN v1.5</option>
      <option value="replicate">Replicate</option>
    </select>
  </div>


  <div>
    <label>Description</label>
    <input type="text" bind:value={description}/>
  </div>


  <!-- CATEGORY -->

  {#if showDefault}

    <div>
      <label>Category</label>
      <select bind:value={category}>
        <option value="upper_body">upper_body</option>
        <option value="upper_body_open">upper_body_open</option>
        <option value="lower_body">lower_body</option>
        <option value="dresses">dresses</option>
        <option value="skirt">skirt</option>
        <option value="skirt_short">skirt_short</option>
        <option value="skirt_mini">skirt_mini</option>
        <option value="shoes">shoes</option>
        <option value="socks">socks</option>
        <option value="stockings">stockings</option>
      </select>
    </div>

  {/if}


  {#if showFashn}

    <div>
      <label>Category</label>
      <select bind:value={category}>
        <option value="tops">tops</option>
        <option value="bottoms">bottoms</option>
        <option value="one_pieces">one_pieces</option>
      </select>
    </div>

  {/if}


  <!-- DEFAULT PARAMS -->

  {#if showDefault}

    <div>
      <label>Steps</label>
      <input type="number" bind:value={steps}/>
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
      <select bind:value={upscale}>
        <option value={1}>1</option>
        <option value={2}>2</option>
        <option value={4}>4</option>
      </select>
    </div>

    <div>
      <label>Upscaler</label>
      <select bind:value={upscaler}>
        <option value="ultrasharp">ultrasharp</option>
        <option value="realesrgan">realesrgan</option>
      </select>
    </div>

  {/if}


  <!-- FASHN PARAMS -->

  {#if showFashn}

    <div>
      <label>Garment photo type</label>
      <select bind:value={garmentPhotoType}>
        <option value="model">model</option>
        <option value="flat_lay">flat_lay</option>
      </select>
    </div>

    <div>
      <label>Num samples</label>
      <input type="number" bind:value={numSamples} min="1" max="4"/>
    </div>

    <div>
      <label>Num timesteps</label>
      <input type="number" bind:value={numTimesteps} min="1" max="100"/>
    </div>

    <div>
      <label>Guidance scale</label>
      <input type="number" step="0.1" bind:value={guidanceScale} min="0.1" max="10"/>
    </div>

    <div>
      <label>Segmentation free</label>
      <input type="checkbox" bind:checked={segmentationFree}/>
    </div>

  {/if}

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
  grid-template-columns:repeat(2,220px);
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