export type JobResponse = {
  job_id: string
  status: string
  result_url?: string
}

export type CreateJobParams = {
  provider: string

  person: File
  garment: File

  description?: string
  category?: string

  steps?: number
  seed?: number
  autocrop?: boolean
  upscale?: number
  upscaler?: string

  // fashn
  garmentPhotoType?: string
  numSamples?: number
  numTimesteps?: number
  guidanceScale?: number
  segmentationFree?: boolean
}

const API_BASE = "/api/v1/wardrobe/try-on"


export async function createJob(
  params: CreateJobParams
): Promise<JobResponse> {

  const form = new FormData()

  form.append("provider", params.provider)

  form.append("person", params.person)
  form.append("garment", params.garment)


  // ---------- common ----------

  if (params.description)
    form.append("description", params.description)

  if (params.category)
    form.append("category", params.category)

  if (params.steps !== undefined)
    form.append("steps", params.steps.toString())

  if (params.seed !== undefined)
    form.append("seed", params.seed.toString())

  if (params.autocrop !== undefined)
    form.append("autocrop", params.autocrop ? "true" : "false")

  if (params.upscale !== undefined)
    form.append("upscale", params.upscale.toString())

  if (params.upscaler)
    form.append("upscaler", params.upscaler)


  // ---------- fashn ----------

  if (params.garmentPhotoType)
    form.append("garmentPhotoType", params.garmentPhotoType)

  if (params.numSamples !== undefined)
    form.append("numSamples", params.numSamples.toString())

  if (params.numTimesteps !== undefined)
    form.append("numTimesteps", params.numTimesteps.toString())

  if (params.guidanceScale !== undefined)
    form.append("guidanceScale", params.guidanceScale.toString())

  if (params.segmentationFree !== undefined)
    form.append(
      "segmentationFree",
      params.segmentationFree ? "true" : "false"
    )


  const res = await fetch(API_BASE, {
    method: "POST",
    body: form
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text)
  }

  return res.json()
}


export async function getJob(
  jobId: string
): Promise<JobResponse> {

  const res = await fetch(`${API_BASE}/${jobId}`)

  if (!res.ok) {
    throw new Error("job fetch failed")
  }

  return res.json()
}