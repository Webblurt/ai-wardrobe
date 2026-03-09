export type JobResponse = {
  job_id: string
  status: string
  result_url?: string
}

const API_BASE = "/api/v1/wardrobe/try-on"

export async function createJob(
  person: File,
  garment: File,
  category: string,
  fit: string
): Promise<JobResponse> {

  const form = new FormData()

  form.append("person", person)
  form.append("garment", garment)
  form.append("category", category)
  form.append("fit", fit)

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