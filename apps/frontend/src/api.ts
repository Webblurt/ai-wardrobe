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
}

const API_BASE = "/api/v1/wardrobe/try-on"

export async function createJob(
  params: CreateJobParams
): Promise<JobResponse> {

  const form = new FormData()

  form.append("provider", params.provider)
  form.append("person", params.person)
  form.append("garment", params.garment)

  if (params.description)
    form.append("description", params.description)

  if (params.category)
    form.append("category", params.category)

  if (params.steps)
    form.append("steps", params.steps.toString())

  if (params.seed)
    form.append("seed", params.seed.toString())

  if (params.autocrop !== undefined)
    form.append("autocrop", params.autocrop ? "true" : "false")

  if (params.upscale)
    form.append("upscale", params.upscale.toString())

  if (params.upscaler)
    form.append("upscaler", params.upscaler)

  const res = await fetch(API_BASE, {
    method: "POST",
    body: form
  })

  if (!res.ok) {
    throw new Error("create job failed")
  }

  return res.json()
}

export async function getJob(jobId: string): Promise<JobResponse> {

  const res = await fetch(`${API_BASE}/${jobId}`)

  if (!res.ok) {
    throw new Error("job fetch failed")
  }

  return res.json()
}