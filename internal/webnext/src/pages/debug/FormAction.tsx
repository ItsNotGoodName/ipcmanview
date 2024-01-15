import { action, redirect, useSubmission } from "@solidjs/router";

// anywhere
const myAction = action(async (formData: FormData) => {
  const start = formData.get("start") as string
  console.log(Date.parse(start))
  await new Promise(r => setTimeout(r, 1000));
  throw redirect("debug");
});

export function FormAction() {
  const submission = useSubmission(myAction)

  return (
    <form method="post" action={myAction}>
      <input id="start" name="start" type="datetime-local" />
      <button type="submit" disabled={submission.pending} >Do</button>
    </form>
  )
}
